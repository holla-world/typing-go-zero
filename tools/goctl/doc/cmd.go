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
	cmdFlags.StringVar(&docgen.AdminFolder, "internal_path")
	cmdFlags.StringVar(&docgen.InternalFolder, "admin_path")
	cmdFlags.StringVar(&docgen.ApiProjectId, "api_project_id")
	cmdFlags.StringVar(&docgen.AdminProjectId, "admin_project_id")
	cmdFlags.StringVar(&docgen.InternalProjectId, "api_project_id")

}
