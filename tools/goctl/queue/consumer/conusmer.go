package consumer

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zeromicro/go-zero/tools/goctl/pkg/golang"
	"github.com/zeromicro/go-zero/tools/goctl/util/pathx"
)

var (
	VarStringDir string
	VarStringCfg string
)

type ConsumerCfg struct {
	Consumers map[string]Consumer
}

type Consumer struct {
	kq.KqConf
	XMsgStructPkg  string
	XMsgStruct     string
	XGroup         string `json:",optional"`
	XHandleTimeout int    `json:",default=5"`
}

// GoCommand gen go project files from command line
func GoCommand(_ *cobra.Command, _ []string) error {
	file := VarStringCfg
	if len(file) == 0 {
		return errors.New("missing --cfg")
	}
	dir := VarStringDir
	if len(file) == 0 {
		return errors.New("missing --dir")
	}
	cf := ConsumerCfg{}
	conf.MustLoad(file, &cf)
	return Gen(cf, dir, "gozero")
}

func Gen(cfg ConsumerCfg, abs, style string) error {
	if len(cfg.Consumers) == 0 {
		return errors.New("no consumers found")
	}
	logx.Must(pathx.MkdirIfNotExist(abs))
	rootPkg, err := golang.GetParentPackage(abs)
	if err != nil {
		return err
	}
	logx.Must(genHandlers(abs, style, rootPkg, cfg))
	logx.Must(genLogics(abs, style, rootPkg, cfg))
	logx.Must(genGroups(abs, rootPkg, cfg))
	println("Done!")
	return nil
}
