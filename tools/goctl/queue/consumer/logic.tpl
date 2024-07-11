package {{.PkgName}}

import (
    "context"
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

func (l *{{.LogicName}}Logic) {{.LogicName}}(key string, msg *{{.MsgMetaPkgShort}}.{{.MsgMeta}}) error {

	return nil
}