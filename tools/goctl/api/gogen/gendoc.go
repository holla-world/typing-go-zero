package gogen

import (
	_ "embed"
	"path"
)

const (
	docDir = "doc"
)

// 核心思路，如果配置了远程仓库地址，那么会将远程仓库拉取到本地，然后找到预设的模板文件，并执行渲染
func genDoc(abs string) error {
	const filename = "client.api"
	serviceName := path.Base(abs)
	return genFile(fileGenConfig{
		dir:             abs,
		subdir:          docDir,
		filename:        filename,
		templateName:    "client_api",
		category:        category,
		templateFile:    apiIdlTemplateFile,
		builtinTemplate: apiIdlTemplate,
		data: map[string]string{
			"name":    serviceName,
			"handler": "Hello",
		},
	})
}
