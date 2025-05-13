package doc

import (
	"github.com/zeromicro/go-zero/tools/goctl/doc/docgen"
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
)

var (
	Cmd = cobrax.NewCommand("doc", cobrax.WithRunE(docgen.GenDoc))
)

func init() {
	var (
		cmdFlags = Cmd.Flags()
	)

	cmdFlags.StringVar(&docgen.ApiFolder, "api_path")
	cmdFlags.StringVar(&docgen.AdminFolder, "admin_path")
	cmdFlags.StringVar(&docgen.InternalFolder, "internal_path")
	cmdFlags.StringVar(&docgen.ApiProjectId, "api_project_id")
	cmdFlags.StringVar(&docgen.AdminProjectId, "admin_project_id")
	cmdFlags.StringVar(&docgen.InternalProjectId, "internal_project_id")

}
