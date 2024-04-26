package gogen

import (
	_ "embed"
	"fmt"
	"path"
	"strings"

	"github.com/holla-world/typing-go-zero/tools/goctl/api/spec"
	"github.com/holla-world/typing-go-zero/tools/goctl/config"
	"github.com/holla-world/typing-go-zero/tools/goctl/util/format"
	"github.com/holla-world/typing-go-zero/tools/goctl/util/pathx"
	"github.com/holla-world/typing-go-zero/tools/goctl/vars"
)

const contextFilename = "service_context"

//go:embed svc.tpl
var contextTemplate string

func genServiceContext(dir, rootPkg string, cfg *config.Config, api *spec.ApiSpec) error {
	filename, err := format.FileNamingFormat(cfg.NamingFormat, contextFilename)
	if err != nil {
		return err
	}

	var middlewareStr string
	var middlewareAssignment string
	middlewares := getMiddleware(api)

	for _, item := range middlewares {
		middlewareStr += fmt.Sprintf("%s rest.Middleware\n", item)
		name := strings.TrimSuffix(item, "Middleware") + "Middleware"
		middlewareAssignment += fmt.Sprintf("%s: %s,\n", item,
			fmt.Sprintf("middleware.New%s().%s", strings.Title(name), "Handle"))
	}

	importList := "\"" + pathx.JoinPackages(path.Dir(rootPkg), xsvcDir) + "\""
	if len(middlewareStr) > 0 {
		importList += "\n\t\"" + pathx.JoinPackages(rootPkg, middlewareDir) + "\""
		importList += fmt.Sprintf("\n\t\"%s/rest\"", vars.ProjectOpenSourceURL)
	}

	return genFile(fileGenConfig{
		dir:             dir,
		subdir:          contextDir,
		filename:        filename + ".go",
		templateName:    "contextTemplate",
		category:        category,
		templateFile:    contextTemplateFile,
		builtinTemplate: contextTemplate,
		data: map[string]string{
			"importList":           importList,
			"xsvc":                 "*xsvc.ServiceContext",
			"middleware":           middlewareStr,
			"middlewareAssignment": middlewareAssignment,
		},
	})
}
