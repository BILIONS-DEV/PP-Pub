package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Assignor int

const (
	Sticky Assignor = iota + 1
	RoundRobin
	Range
)

func (a Assignor) String() string {
	switch a {
	case Sticky:
		return "sticky"
	case RoundRobin:
		return "roundrobin"
	case Range:
		return "range"
	default:
		return ""
	}
}

type Consumer interface {
	Setting() Info
	sarama.ConsumerGroupHandler
}

type Info struct {
	Verbose  bool
	Oldest   bool
	Assignor Assignor
	Group    string
	Topics   []string
}

func SetupConsume(consumer Consumer) {
	setting := consumer.Setting()
	log.Println("Starting a new Sarama consumer")
	if setting.Verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0 // specify appropriate version
	config.Consumer.Return.Errors = true
	switch setting.Assignor {
	case Sticky:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case RoundRobin:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case Range:
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", setting.Assignor)
	}
	if setting.Oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(Connector.Brokers, setting.Group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}
	defer func() { _ = client.Close() }()

	// Track errors
	go func() {
		for err := range client.Errors() {
			log.Println("ERROR", err)
		}
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, setting.Topics, consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
		}
	}()

	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}
