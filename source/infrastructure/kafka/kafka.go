package kafka

import "github.com/segmentio/kafka-go"

type Config struct {
	Brokers []string
	Topic   string
	GroupID string
}

type Client struct {
	Writer *kafka.Writer
	Reader *kafka.Reader
}

func NewKafka(config Config) *Client {
	return &Client{
		Writer: NewProducer(config),
		//Reader: NewConsumes(config),
	}
}

func NewProducer(config Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(config.Brokers...),
		Topic:    config.Topic,
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	}
}

func NewConsumes(config Config) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: config.Brokers,
		Topic:   config.Topic,
		GroupID: config.GroupID,
		//MinBytes: 10e3, // 10KB
		//MaxBytes: 10e6, // 10MB
		MinBytes: 1, // 10KB
		MaxBytes: 57671680,
	})
}
