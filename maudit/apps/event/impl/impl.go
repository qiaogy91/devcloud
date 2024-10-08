package impl

import (
	"context"
	"github.com/qiaogy91/devcloud/maudit/apps/event"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/datasource"
	"github.com/qiaogy91/ioc/config/log"
	kafkaClient "github.com/qiaogy91/ioc/default/kafka"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

var _ event.Service = &Impl{}

type Impl struct {
	ioc.ObjectImpl
	log      *zerolog.Logger
	db       *gorm.DB
	kafkaSvc kafkaClient.Service
	Topic    string `json:"topic" yaml:"topic"`
}

func (i *Impl) Name() string  { return event.AppName }
func (i *Impl) Priority() int { return 301 }
func (i *Impl) Init() {
	i.log = log.Sub(event.AppName)
	i.db = datasource.DB()
	i.kafkaSvc = kafkaClient.GetClient()

	// 启动一个协程去执行Sync 操作
	i.log.Info().Msgf("kafka consumer at %s/topic %s/groupId", i.Topic, event.AppName)
	go i.Sync(context.Background(), &event.SyncReq{})
}

func init() {
	ins := &Impl{}
	ioc.Controller().Registry(ins)
}
