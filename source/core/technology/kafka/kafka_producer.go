package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
	"source/pkg/logger"
	"time"
)

var producer sarama.AsyncProducer

func init() {
	//producer, _ = setupProducer()
	producer = NewAsyncProducer(Connector.Brokers)
}

type ProducerMessage interface {
	GetTopicName() string
}

// Push message to kafka producer
//
// param: message
func Push(message ProducerMessage) {
	jsonData, err := json.Marshal(message)
	if err != nil {
		logger.Debug(err.Error())
		return
	}
	msg := &sarama.ProducerMessage{
		Topic: message.GetTopicName(),
		Key:   nil,
		Value: sarama.StringEncoder(jsonData),
	}
	select {
	case producer.Input() <- msg:
	case err := <-producer.Errors():
		log.Println(err)
	}
}

// Start Async Producer
//
// param: brokerList
// return:
func NewAsyncProducer(brokerList []string) (producer sarama.AsyncProducer) {
	// For the access log, we are looking for AP semantics, with high throughput.
	// By creating batches of compressed messages, we reduce network I/O at a cost of more latency.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	config.Version = sarama.V2_0_0_0

	var err error
	producer, err = sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		logger.Log.Fatal("Failed to start Sarama producer: " + err.Error())
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			logger.Info("Failed to write access log entry: " + err.Error())
		}
	}()

	return
}

// setupProducer will create a AsyncProducer and returns it
func setupProducer() (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	return sarama.NewAsyncProducer(Connector.Brokers, config)
}
