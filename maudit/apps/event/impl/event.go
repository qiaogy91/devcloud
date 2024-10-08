package impl

import (
	"context"
	"encoding/json"
	"github.com/qiaogy91/devcloud/maudit/apps/event"
	"time"
)

func (i *Impl) Sync(ctx context.Context, req *event.SyncReq) {
	consumer := i.kafkaSvc.Consumer(i.Topic, event.AppName)
	defer func() { _ = consumer.Close() }()
	for {
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			i.log.Error().Msgf("event read failed, %s", err)
			time.Sleep(5 * time.Second)
			continue
		}
		ins := &event.Event{}
		if err := json.Unmarshal(msg.Value, ins); err != nil {
			i.log.Error().Msgf("event unmarshal failed, %s", err)
			time.Sleep(5 * time.Second)
			continue
		}
	}
}

func (i *Impl) Query(ctx context.Context, req *event.QueryReq) (*event.EventsSet, error) {
	//TODO implement me
	panic("implement me")
}
