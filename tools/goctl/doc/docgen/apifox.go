package docgen

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// APIDocConfig 用于传递配置参数
type APIDocConfig struct {
	FileName        string
	APIFolder       string
	APIPath         string
	ApifoxProjectID string
}

// GenerateAndImportSwagger 生成 swagger.json 并导入 Apifox
func GenerateAndImportSwagger(configs []APIDocConfig, apifoxAccessToken string) error {
	// 检查是否有 api 文件存在
	hasApi := false
	for _, cfg := range configs {
		if fileExists(cfg.APIPath) {
			hasApi = true
			break
		}
	}
	if !hasApi {
		return nil // 没有 api 文件则不执行
	}

	for _, config := range configs {
		if !fileExists(config.APIPath) {
			continue
		}
		apiFolder := config.APIFolder
		apiPath := config.APIPath
		fileName := config.FileName

		// 1. 生成 swagger.json
		swaggerPath := filepath.Join(apiFolder, fileName+".json")
		cmd := exec.Command("goctl", "api", "plugin",
			"-plugin", fmt.Sprintf(`goctl-swagger='swagger -filename %s.json'`, fileName),
			"-api", apiPath,
			"-dir", apiFolder,
		)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to generate swagger: %v, output: %s", err, string(out))
		}

		// 2. 读取 swagger.json
		swaggerData, err := os.ReadFile(swaggerPath)
		if err != nil {
			return fmt.Errorf("failed to read swagger file: %v", err)
		}
		var swagger map[string]interface{}
		if err := json.Unmarshal(swaggerData, &swagger); err != nil {
			return fmt.Errorf("failed to unmarshal swagger: %v", err)
		}
		if len(swagger) == 0 {
			continue
		}

		// 3. 自定义处理
		serviceName := getServiceName(fileName)
		// 3.1 添加 x-apifox-folder
		if paths, ok := swagger["paths"].(map[string]interface{}); ok {
			for _, pathValue := range paths {
				if methods, ok := pathValue.(map[string]interface{}); ok {
					for _, methodValue := range methods {
						if methodMap, ok := methodValue.(map[string]interface{}); ok {
							methodMap["x-apifox-folder"] = serviceName
						}
					}
				}
			}
		}

		// 3.2 schema 重命名
		definitions, ok := swagger["definitions"].(map[string]interface{})
		if ok {
			defCopy := make(map[string]interface{})
			allSchemaNames := make([]string, 0, len(definitions))
			for k, v := range definitions {
				defCopy[k] = v
				allSchemaNames = append(allSchemaNames, k)
			}
			for schemaName, schemaValue := range defCopy {
				if schemaMap, ok := schemaValue.(map[string]interface{}); ok {
					schemaMap["title"] = serviceName + schemaName
				}
				delete(definitions, schemaName)
				definitions[serviceName+schemaName] = schemaValue
			}
			// 替换所有 $ref
			swaggerStr, _ := json.Marshal(swagger)
			swaggerStrNew := string(swaggerStr)
			for _, schemaName := range allSchemaNames {
				oldRef := "#/definitions/" + schemaName
				newRef := "#/definitions/" + serviceName + schemaName
				swaggerStrNew = strings.ReplaceAll(swaggerStrNew, oldRef, newRef)
			}
			_ = json.Unmarshal([]byte(swaggerStrNew), &swagger)
		}

		// 3.3 responses 外层包裹
		if paths, ok := swagger["paths"].(map[string]interface{}); ok {
			for _, pathValue := range paths {
				if methods, ok := pathValue.(map[string]interface{}); ok {
					for _, methodValue := range methods {
						if methodMap, ok := methodValue.(map[string]interface{}); ok {
							if responses, ok := methodMap["responses"].(map[string]interface{}); ok {
								for _, resp := range responses {
									if respMap, ok := resp.(map[string]interface{}); ok {
										if schema, ok := respMap["schema"].(map[string]interface{}); ok {
											// 有 $ref
											if ref, ok := schema["$ref"].(string); ok {
												respMap["schema"] = map[string]interface{}{
													"type": "object",
													"properties": map[string]interface{}{
														"result":  map[string]interface{}{"type": "integer"},
														"message": map[string]interface{}{"type": "string"},
														"data":    map[string]interface{}{"$ref": ref},
													},
													"required": []string{"result", "message", "data"},
												}
											} else {
												// 空 map
												respMap["schema"] = map[string]interface{}{
													"type": "object",
													"properties": map[string]interface{}{
														"result":  map[string]interface{}{"type": "integer"},
														"message": map[string]interface{}{"type": "string"},
														"data":    map[string]interface{}{"type": "object"},
													},
													"required": []string{"result", "message", "data"},
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}

		// 4. 写回 swagger.json
		swaggerOut, err := json.MarshalIndent(swagger, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal swagger: %v", err)
		}
		if err := os.WriteFile(swaggerPath, swaggerOut, 0644); err != nil {
			return fmt.Errorf("failed to write swagger file: %v", err)
		}

		// 5. 导入到 Apifox
		url := fmt.Sprintf("https://api.apifox.com/v1/projects/%s/import-openapi", config.ApifoxProjectID)
		payload := map[string]interface{}{
			"input": string(swaggerOut),
			"options": map[string]interface{}{
				"endpointOverwriteBehavior": "OVERWRITE_EXISTING",
				"schemaOverwriteBehavior":   "OVERWRITE_EXISTING",
			},
		}
		payloadBytes, _ := json.Marshal(payload)
		curlCmd := exec.Command("curl", "-X", "POST", url,
			"-H", "Content-Type: application/json",
			"-H", "X-Apifox-Version: 2024-03-28",
			"-H", fmt.Sprintf("Authorization: Bearer %s", apifoxAccessToken),
			"-d", string(payloadBytes),
		)
		curlOut, err := curlCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to import to Apifox: %v, output: %s", err, string(curlOut))
		}
		fmt.Println(string(curlOut))
	}
	return nil
}

// fileExists 判断文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// getServiceName 获取服务名
func getServiceName(fileName string) string {
	wd, _ := os.Getwd()
	parts := strings.Split(wd, "/")
	last := parts[len(parts)-1]
	lastPart := last
	if idx := strings.LastIndex(last, "-"); idx != -1 {
		lastPart = last[idx+1:]
	}
	return strings.Title(lastPart) + "_" + fileName
}
