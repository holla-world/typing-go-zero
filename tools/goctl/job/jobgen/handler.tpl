package {{.PkgName}}

import (
	"context"
	"github.com/holla-world/typing-golib/xzero/xcmd/engine"

	{{.ImportPackages}}
)

func {{.HandlerName}}(csvc *svc.ServiceContext) engine.HandlerFunc {
	return func(ctx context.Context) error {
		// TODO business logic
		return nil
	}
}
