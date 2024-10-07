package impl_test

import (
	"context"
	"github.com/qiaogy91/ioc"

	kafkaClient "github.com/qiaogy91/devcloud/maudit/apps/kafka"
	"github.com/segmentio/kafka-go"
	"testing"
)

var (
	ctx = context.Background()
	i   kafkaClient.Service
	r   *kafka.Reader
	w   *kafka.Writer
)

func init() {
	confPath := "/Users/qiaogy/GolandProjects/projects/github/devcloud/maudit/etc/application.yaml"
	if err := ioc.ConfigIocObject(confPath); err != nil {
		panic(err)
	}

	i = kafkaClient.GetClient()
	r = i.Consumer("maudit", "group01")
	w = i.Producer("maudit")
}

func TestClient_Producer(t *testing.T) {
	if err := w.WriteMessages(ctx, kafka.Message{Value: []byte("密文消息01")}); err != nil {
		t.Fatalf("send message err, %s", err)
	}
}

func TestClient_Consumer(t *testing.T) {
	msg, err := r.ReadMessage(ctx)
	if err != nil {
		t.Fatalf("recive err, %s", err)
	}
	t.Logf("recive data, %s", msg.Value)
}
