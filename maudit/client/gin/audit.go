package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/qiaogy91/devcloud/maudit/apps/event"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/application"
	ioc_gin "github.com/qiaogy91/ioc/config/gin"
	"github.com/qiaogy91/ioc/config/log"
	kafkaClient "github.com/qiaogy91/ioc/default/kafka"
	"github.com/qiaogy91/ioc/labels"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"time"
)

const AppName = "Audit"

type Audit struct {
	ioc.ObjectImpl
	log      *zerolog.Logger
	producer *kafka.Writer
	Topic    string `json:"topic" yaml:"topic"`
}

func (a *Audit) Name() string  { return AppName }
func (a *Audit) Priority() int { return 202 }
func (a *Audit) Init() {
	a.log = log.Sub(AppName)
	a.producer = kafkaClient.GetClient().Producer(a.Topic)

	// 初始化一个kafka producer client，然后注册到全局webservie
	r := ioc_gin.RootRouter()
	r.Use(a.Middleware)
}

func (a *Audit) Middleware(ctx *gin.Context) {
	// 没有开启审计功能，直接跳过
	if !ctx.GetBool(labels.AuditEnable) {
		ctx.Next()
		return
	}

	// 开启审计后，在Response 过程中进行处理
	ctx.Next()
	e := &event.Event{
		User:       "",
		Time:       time.Now().Unix(),
		SourceIP:   ctx.ClientIP(),
		UserAgent:  ctx.Request.UserAgent(),
		Service:    ctx.GetString(application.Get().Name()),
		Resource:   ctx.GetString(labels.ResourceName),
		Action:     ctx.GetString(labels.ActionName),
		StatusCode: ctx.Writer.Status(),
	}
	bs, err := e.GetBS()
	if err != nil {
		a.log.Error().Msgf("marshal event failed, %s", err)
	}
	if err := a.producer.WriteMessages(ctx, kafka.Message{Value: bs}); err != nil {
		a.log.Error().Msgf("maudit middleware write message failed, %s", err)
	}
}

func init() {
	ioc.Default().Registry(&Audit{})
}
