package {{.PkgName}}

import (
"context"
"encoding/json"
"time"

"github.com/holla-world/typing-golib/xzero/xlog"
"github.com/holla-world/typing-golib/xzero/xqueue"
"github.com/holla-world/typing-golib/xzero/xqueue/kafkaq"
"github.com/zeromicro/go-zero/core/logx"

	{{.ImportPackages}}
)

func New{{.HandlerName}}ConsumerHandler(cfg kafkaq.KqConf, svcCtx *svc.ServiceContext) xqueue.Consumer {
return xqueue.MustKqConsumer(cfg, new{{.HandlerName}}Consumer(cfg, svcCtx))
}

type {{.LHandlerName}}Consumer struct {
	svcCtx *svc.ServiceContext
cfg    *kafkaq.KqConf
}

func new{{.HandlerName}}Consumer(cfg kafkaq.KqConf, svcCtx *svc.ServiceContext) {{.LHandlerName}}Consumer {
	return {{.LHandlerName}}Consumer{
		svcCtx: svcCtx,
cfg:    &cfg,
	}
}

func (c {{.LHandlerName}}Consumer) Consume(ctx1 context.Context, msg xqueue.MsgOut) error {
timeout := time.Duration(c.cfg.Timeout) * time.Second
ctx, cancel := context.WithTimeout(ctx1, timeout)
defer cancel()

dst := {{.MsgMetaPkgShort}}.{{.MsgMeta}}{}
if err := json.Unmarshal(msg.Value(), &dst); err != nil {
xlog.Errorm(ctx,
"unmarshal queue msg err",
logx.Field("err", err),
logx.Field("msg", msg.Value()),
logx.Field("dst", "{{.MsgMeta}}"),
)
return err
	}

return logic.New{{.HandlerName}}Logic(ctx, c.svcCtx).{{.HandlerName}}(msg, &dst)
}