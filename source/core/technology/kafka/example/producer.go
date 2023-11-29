package main

import (
	"fmt"
	"source/core/technology/kafka"
	"time"
)

type Producer struct {
	Id   int
	Name string
}

func (Producer) GetTopicName() string {
	return "tung_topic"
}

func main() {
	i := 1
	for {
		if i == 100 {
			break
		}
		time.Sleep(1 * time.Second)
		kafka.Push(Producer{
			Id:   i,
			Name: "Dang Thanh Tung",
		})
		i++
	}
	fmt.Println("Finish")
}
