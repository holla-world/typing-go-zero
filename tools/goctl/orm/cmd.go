package orm

import (
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
	"github.com/zeromicro/go-zero/tools/goctl/orm/command"
)

var (
	// Cmd describes a model command.
	Cmd = cobrax.NewCommand("orm", cobrax.WithRunE(command.GenOrm))
)

func init() {
	Cmd.Flags().StringVarP(&command.VarStringSrc, "src", "s")
}
