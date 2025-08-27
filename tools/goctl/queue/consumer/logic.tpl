package {{.PkgName}}

import (
    "context"
"github.com/holla-world/typing-golib/xzero/xqueue"

	{{.ImportPackages}}
)

type {{.LogicName}}Logic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func New{{.LogicName}}Logic(ctx context.Context, svcCtx *svc.ServiceContext) *{{.LogicName}}Logic {
	return &{{.LogicName}}Logic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *{{.LogicName}}Logic) {{.LogicName}}(raw xqueue.MsgOut, msg *{{.MsgMetaPkgShort}}.{{.MsgMeta}}) error {

	return nil
}