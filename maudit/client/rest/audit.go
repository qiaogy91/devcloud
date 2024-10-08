package rest

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/qiaogy91/devcloud/maudit/apps/event"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/application"
	"github.com/qiaogy91/ioc/config/gorestful"
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
	ws := gorestful.RootContainer()
	ws.Filter(a.Middleware)
}

func (a *Audit) Middleware(r *restful.Request, w *restful.Response, chain *restful.FilterChain) {
	// 获取不到路由、或元数据为空，则直接执行下一个中间件
	sr := r.SelectedRoute()
	if sr == nil || sr.Metadata() == nil {
		chain.ProcessFilter(r, w)
		return
	}

	// 未开启审计，则直接执行下一个中间件
	meta := Metadata{data: sr.Metadata()}
	if !meta.Bool("audit") {
		chain.ProcessFilter(r, w)
		return
	}

	chain.ProcessFilter(r, w)

	// response 过程中进行处理
	e := &event.Event{
		Time:       time.Now().Unix(),
		SourceIP:   r.Request.RemoteAddr,
		UserAgent:  r.Request.UserAgent(),
		Service:    application.Get().Name(),
		Resource:   meta.Str(labels.ResourceName),
		Action:     meta.Str(labels.ActionName),
		StatusCode: w.StatusCode(),
	}

	bs, err := e.GetBS()
	if err != nil {
		a.log.Error().Msgf("marshal event failed, %s", err)
		return
	}
	if err := a.producer.WriteMessages(r.Request.Context(), kafka.Message{Value: bs}); err != nil {
		a.log.Error().Msgf("maudit middleware write message failed, %s", err)
		return
	}
}

type Metadata struct {
	data map[string]any
}

func (m *Metadata) Bool(k string) bool {
	v, ok := m.data[k]
	if !ok {
		return false
	}
	return v.(bool)
}

func (m *Metadata) Str(k string) string {
	v, ok := m.data[k]
	if !ok {
		return ""
	}
	return v.(string)
}

func init() {
	ioc.Default().Registry(&Audit{})
}
