package job

import (
	"github.com/zeromicro/go-zero/tools/goctl/internal/cobrax"
	"github.com/zeromicro/go-zero/tools/goctl/job/jobgen"
)

var (
	Cmd = cobrax.NewCommand("job", cobrax.WithRunE(jobgen.GenJob))
)

func init() {
	var (
		cmdFlags = Cmd.Flags()
	)

	cmdFlags.StringVar(&jobgen.CronFile, "cron")
	cmdFlags.StringVar(&jobgen.DaemonFile, "daemon")

}
