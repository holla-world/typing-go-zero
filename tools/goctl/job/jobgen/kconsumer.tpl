package {{.PkgName}}

import (
	"context"
	"github.com/holla-world/typing-golib/xzero/xkafka"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/queue"

	{{.ImportPackages}}
)

// 消费者名：{{.ConsumerName}}
// 消费TOPIC：{{.ConsumerTopic}}
// 消费者组：{{.ConsumerGroup}}

type {{.ConsumerHandle}} struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	// lg     logic.DemoLogicInterface 可以挂载某个logic实例
}

func Run{{.ConsumerHandle}}(ctx context.Context, svcCtx *svc.ServiceContext) error {
	q := New{{.ConsumerHandle}}(ctx, svcCtx)
	q.Start()
	defer q.Stop()
	return nil
}

func New{{.ConsumerHandle}}(ctx context.Context, svcCtx *svc.ServiceContext) queue.MessageQueue {
	consumerConf := xkafka.GetKqConsumerConfByName({{.ConsumerName}}, svcCtx.XvcCtx.Config.KqConsumer)
	return xkafka.NewConsumer(consumerConf, func() kq.ConsumeHandler {
		return &{{.ConsumerHandle}}{
			ctx:    ctx,
			svcCtx: svcCtx,
			// lg:     logic.NewDemoLogic(ctx, svcCtx),可以挂载某个logic实例
		}
	})
}

func (c *{{.ConsumerHandle}}) Consume(k, value string) error {
	// var msg ktypes.YourType
	// xlog.Debugf(c.ctx, "{{.ConsumerHandle}}:%s", value)
	// err := jsoniter.UnmarshalFromString(value, &msg)
	// if err != nil {
	// 	return err
	// }
	// TODO business logic here

	return nil
}