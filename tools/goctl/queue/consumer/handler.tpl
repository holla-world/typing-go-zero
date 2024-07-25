package {{.PkgName}}

import (
    "context"
    "encoding/json"
    "time"

	"github.com/zeromicro/go-queue/kq"
    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/core/queue"
	{{.ImportPackages}}
)

func New{{.HandlerName}}ConsumerHandler(kqConf kq.KqConf, svcCtx *svc.ServiceContext) queue.MessageQueue {
	return kq.MustNewQueue(kqConf, new{{.HandlerName}}Consumer(svcCtx))
}

type {{.LHandlerName}}Consumer struct {
	svcCtx *svc.ServiceContext
}

func new{{.HandlerName}}Consumer(svcCtx *svc.ServiceContext) {{.LHandlerName}}Consumer {
	return {{.LHandlerName}}Consumer{
		svcCtx: svcCtx,
	}
}

func (c {{.LHandlerName}}Consumer) Consume(key, value string) (err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), {{.Timeout}}*time.Second)
	defer cancelFunc()

	msg := {{.MsgMetaPkgShort}}.{{.MsgMeta}}{}
	if err = json.Unmarshal([]byte(value), &msg); err != nil {
		logx.
			WithContext(ctx).
			Errorw(
				"kafka msg unmarshal error",
				logx.Field("error", err),
				logx.Field("key", key),
				logx.Field("msg", value),
				logx.Field("struct", "{{.MsgMeta}}"),
			)
		return
	}

	l := logic.New{{.HandlerName}}Logic(ctx, c.svcCtx)
	return l.{{.HandlerName}}(key, &msg)
}