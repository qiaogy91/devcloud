package impl

import (
	"github.com/segmentio/kafka-go"
	"time"
)

func (c *Client) Producer(topic string) *kafka.Writer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(c.Brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		Transport:              &kafka.Transport{SASL: c.mechanism},
		AllowAutoTopicCreation: true,
	}
	c.producers = append(c.producers, w)
	return w
}

func (c *Client) Consumer(topic, groupId string) *kafka.Reader {
	conf := kafka.ReaderConfig{
		Brokers:     c.Brokers,
		Topic:       topic,
		GroupID:     groupId,
		MaxBytes:    10e6,
		StartOffset: c.Offset,
		Dialer:      &kafka.Dialer{Timeout: 10 * time.Second, SASLMechanism: c.mechanism},
	}

	r := kafka.NewReader(conf)
	c.consumers = append(c.consumers, r)
	return r
}
