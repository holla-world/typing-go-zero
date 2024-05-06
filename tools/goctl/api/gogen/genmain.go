package gogen

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
	"github.com/zeromicro/go-zero/tools/goctl/vars"
)

//go:embed main.tpl
var mainTemplate string

func genMain(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	// name := strings.ToLower(api.Service.Name)
	// filename, err := format.FileNamingFormat(cfg.NamingFormat, name)
	// if err != nil {
	// 	return err
	// }
	//
	// configName := filename
	// if strings.HasSuffix(filename, "-api") {
	// 	filename = strings.ReplaceAll(filename, "-api", "")
	// }

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          "",
		filename:        filepath.Base(dir) + ".go",
		templateName:    "mainTemplate",
		category:        category,
		templateFile:    mainTemplateFile,
		builtinTemplate: mainTemplate,
		data: map[string]string{
			"importPackages": genMainImports(rootPkg),
			// "serviceName":    configName,
		},
	})
}

func genMainImports(parentPkg string) string {
	var imports []string
	// imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, configDir)))
	imports = append(imports, xsvcImport(parentPkg))
	imports = append(imports, fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, handlerDir)))
	imports = append(imports, fmt.Sprintf("\"%s\"\n", pathx.JoinPackages(parentPkg, contextDir)))
	// imports = append(imports, fmt.Sprintf("\"%s/core/conf\"", vars.ProjectOpenSourceURL))
	imports = append(imports, fmt.Sprintf("\"%s/rest\"", vars.ProjectOpenSourceURL))
	return strings.Join(imports, "\n\t")
}

func xsvcImport(parentPkg string) string {
	projectDir := filepath.Dir(parentPkg)
	return fmt.Sprintf("\"%s\"", pathx.JoinPackages(projectDir, xsvcDir))
}
