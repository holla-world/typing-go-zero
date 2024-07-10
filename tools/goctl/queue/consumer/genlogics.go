package consumer

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util/format"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

//go:embed logic.tpl
var logicTemplate string

func genLogic(dir, rootPkg, name, style string, cfg Consumer) error {
	pkgName := "logic"

	filename, err := format.FileNamingFormat(style, name)
	if err != nil {
		return err
	}

	imports := genLogicImports(rootPkg, cfg.XMsgStructPkg)
	c := fileGenConfig{
		dir:          dir,
		subdir:       "logic",
		filename:     filename + "logic.go",
		templateName: "logicTemplate",
		// category:        "api",
		// templateFile:    "logic.tpl",
		builtinTemplate: logicTemplate,
		data: map[string]any{
			"PkgName":         pkgName,
			"ImportPackages":  imports,
			"LogicName":       name,
			"MsgMetaPkgShort": GetPackageName(cfg.XMsgStructPkg),
			"MsgMeta":         cfg.XMsgStruct,
		},
	}
	return genFile(c)
}

func genLogics(dir, style, rootPkg string, cfg ConsumerCfg) error {
	for name, c := range cfg.Consumers {
		err := genLogic(dir, rootPkg, name, style, c)
		if err != nil {
			return err
		}
	}

	return nil
}

func genLogicImports(parentPkg string, pkgs ...string) string {
	imports := []string{
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, "svc")),
	}
	for _, pkg := range pkgs {
		imports = append(imports, fmt.Sprintf("\"%s\"", pkg))
	}
	return strings.Join(imports, "\n\t")
}
