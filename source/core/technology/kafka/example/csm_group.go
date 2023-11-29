package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"source/core/technology/kafka"
	"sync"
)

type exampleConsumerGroupHandler struct {
	GroupName string
}

func (exampleConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (exampleConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (c exampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d value:%s group:%s ---- \n", msg.Topic, msg.Partition, msg.Offset, string(msg.Value), c.GroupName)
		fmt.Println("===========")
		//time.Sleep(1 * time.Second)
		//r := mysql.AdsetTable{}
		//mysql.Client.Table("adset").Last(&r)
		sess.MarkMessage(msg, "")
	}
	return nil
}
func (c exampleConsumerGroupHandler) Setting() (info kafka.Info) {
	info.Assignor = kafka.Range
	info.Oldest = false
	info.Verbose = false
	info.Topics = []string{"tung_topic"}
	info.Group = c.GroupName
	return
}

var wg sync.WaitGroup

func main() {

	fmt.Println("Main start")
	for i := 0; i <= 2; i++ {
		wg.Add(1)
		go func() {
			kafka.SetupConsume(exampleConsumerGroupHandler{GroupName: "my-group"})
			wg.Done()
		}()
	}
	wg.Wait()
}
