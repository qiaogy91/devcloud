package impl

import (
	"context"
	kafkaClient "github.com/qiaogy91/devcloud/maudit/apps/kafka"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type Client struct {
	ioc.ObjectImpl
	log       *zerolog.Logger
	Username  string   `json:"username" yaml:"username"`
	Password  string   `json:"password" yaml:"password"`
	Brokers   []string `json:"brokers" yaml:"brokers"`
	Async     bool     `json:"async" yaml:"async"`
	Offset    int64    `json:"offset" yaml:"offset"`
	producers []*kafka.Writer
	consumers []*kafka.Reader
	mechanism sasl.Mechanism
}

func (c *Client) Name() string  { return kafkaClient.AppName }
func (c *Client) Priority() int { return 201 }
func (c *Client) Init() {
	c.log = log.Sub(kafkaClient.AppName)
	if c.Username == "" {
		return
	}

	mechanism, err := scram.Mechanism(scram.SHA512, c.Username, c.Password)
	if err != nil {
		panic(err)
	}
	c.mechanism = mechanism
}

func (c *Client) Close(ctx context.Context) error {
	for _, item := range c.producers {
		if err := item.Close(); err != nil {
			c.log.Error().Msgf("close producer failed, %s", err)
		}
	}

	for _, item := range c.consumers {
		if err := item.Close(); err != nil {
			c.log.Error().Msgf("close consumer failed, %s", err)
		}
	}

	return nil
}
func init() {
	ioc.Default().Registry(&Client{})
}
