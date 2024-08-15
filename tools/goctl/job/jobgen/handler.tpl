package {{.PkgName}}

import (
	"context"
	"github.com/holla-world/typing-golib/xzero/xcmd/engine"

	{{.ImportPackages}}
)
// {{.HandlerName}}
// 本地运行命令: go run main.go {{.JobType}} {{.Action}}
func {{.HandlerName}}(csvc *svc.ServiceContext) engine.HandlerFunc {
	return func(ctx context.Context) error {
		// TODO business logic
		return nil
	}
}
