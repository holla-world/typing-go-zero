package command

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/zeromicro/go-zero/tools/goctl/orm/gormgen"
)

var (
	// VarStringSrc describes the source file of sql.
	VarStringSrc string
	VarStringPkg string
)

// MysqlDDL generates model code from ddl
func GenOrm(_ *cobra.Command, _ []string) error {
	src := VarStringSrc
	if src == "" {
		return errors.New("请指定-s, ddl文件不能为空")
	}
	pkg := VarStringPkg
	if pkg == "" {
		pkg = "query"
	} else {
		pkg = "query/" + pkg
	}
	return gormgen.Gen(src, pkg)
}
