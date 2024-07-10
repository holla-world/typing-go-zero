package queue

import (
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
	"github.com/zeromicro/go-zero/tools/goctl/queue/consumer"
)

var (
	// Cmd describes a model command.
	Cmd         = cobrax.NewCommand("queue")
	consumerCmd = cobrax.NewCommand("consumer", cobrax.WithRunE(consumer.GoCommand))
)

func init() {
	consumerCmd.Flags().StringVarP(&consumer.VarStringDir, "dir", "d")
	consumerCmd.Flags().StringVarP(&consumer.VarStringCfg, "cfg", "c")
	Cmd.AddCommand(consumerCmd)
}
