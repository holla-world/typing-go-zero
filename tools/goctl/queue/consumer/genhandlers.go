package consumer

import (
	_ "embed"
	"fmt"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

const defaultLogicPackage = "logic"

//go:embed handler.tpl
var handlerTemplate string

func genHandler(dir, rootPkg, name, style string, cfg Consumer) error {
	pkgName := "handler"

	filename, err := format.FileNamingFormat(style, name)
	if err != nil {
		return err
	}

	imports := genHandlerImports(rootPkg, cfg.XMsgStructPkg)
	return genFile(fileGenConfig{
		dir:          dir,
		subdir:       "handler",
		filename:     filename + "handler.go",
		templateName: "handlerTemplate",
		// category:        "api",
		// templateFile:    "handler.tpl",
		builtinTemplate: handlerTemplate,
		data: map[string]any{
			"PkgName":         pkgName,
			"ImportPackages":  imports,
			"HandlerName":     name,
			"LHandlerName":    ToLowerFirst(name),
			"MsgMetaPkgShort": GetPackageName(cfg.XMsgStructPkg),
			"MsgMeta":         cfg.XMsgStruct,
			"Timeout":         cfg.XHandleTimeout,
		},
	})
}

func genHandlers(dir, style, rootPkg string, cfg ConsumerCfg) error {
	for name, c := range cfg.Consumers {
		err := genHandler(dir, rootPkg, name, style, c)
		if err != nil {
			return err
		}
	}

	return nil
}

func genHandlerImports(parentPkg string, pkgs ...string) string {
	imports := []string{
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, "svc")),
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, "logic")),
	}
	for _, pkg := range pkgs {
		imports = append(imports, fmt.Sprintf("\"%s\"", pkg))
	}
	return strings.Join(imports, "\n\t")
}

func GetPackageName(pkgPath string) string {
	return path.Base(pkgPath)
}
