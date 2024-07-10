package consumer

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

//go:embed consumergroup.tpl
var groupTemplate string

// removeFileIfExists 检查文件是否存在，如果存在则删除它
func removeFileIfExists(filePath string) error {
	if _, err := os.Stat(filePath); err == nil {
		// 文件存在，删除文件
		if err := os.Remove(filePath); err != nil {
			return fmt.Errorf("failed to delete file: %w", err)
		}
	} else if !os.IsNotExist(err) {
		// 发生错误，但不是文件不存在的错误
		return fmt.Errorf("failed to check file: %w", err)
	}
	return nil
}

func genGroups(dir, rootPkg string, cfg ConsumerCfg) error {
	pkgName := "cmd"
	imports := genCgroupImports(rootPkg)
	gs := groups(cfg)
	filename := "consumer.gen.go"
	err := removeFileIfExists(path.Join(dir, "cmd", filename))
	if err != nil {
		return err
	}
	c := fileGenConfig{
		dir:             dir,
		subdir:          "cmd",
		filename:        filename,
		templateName:    "groupTemplate",
		category:        "api",
		templateFile:    "group.tpl",
		builtinTemplate: groupTemplate,
		data: map[string]any{
			"PkgName":        pkgName,
			"ImportPackages": imports,
			"Groups":         gs,
		},
	}
	return genFile(c)
}

type Groups struct {
	GroupName string
	Handlers  []string
}

func groups(cfg ConsumerCfg) []Groups {
	const def = "defaultGroup"
	gs := make(map[string]Groups, 0)
	for k, c := range cfg.Consumers {
		if c.XGroup == "" {
			c.XGroup = def
		}
		g, ok := gs[c.XGroup]
		if !ok {
			g = Groups{
				GroupName: c.XGroup,
				Handlers:  []string{k},
			}
		} else {
			g.Handlers = append(g.Handlers, k)
		}
		gs[c.XGroup] = g
	}
	list := make([]Groups, 0)
	for _, v := range gs {
		list = append(list, Groups{
			GroupName: v.GroupName,
			Handlers:  v.Handlers,
		})
	}
	return list
}

func genCgroupImports(parentPkg string, pkgs ...string) string {
	imports := []string{
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, "handler")),
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(parentPkg, "svc")),
		fmt.Sprintf("\"%s\"", pathx.JoinPackages(projectPkg(parentPkg), "internal", "xsvc")),
	}

	for _, pkg := range pkgs {
		imports = append(imports, fmt.Sprintf("\"%s\"", pkg))
	}
	return strings.Join(imports, "\n\t")
}

func projectPkg(parentPkg string) string {
	return strings.TrimSuffix(parentPkg, "/job")
}
