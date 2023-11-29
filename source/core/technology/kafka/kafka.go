package kafka

import (
	"source/pkg/utility"
)

type connector struct {
	User     string
	Password string
	Brokers  []string
}

type topic struct {
	Report string
}

var Connector connector
var Topics topic

func init() {
	/**
	Kafka Topic
	*/
	Topics.Report = "report"

	/**
	Config Kafka Login
	*/
	if utility.IsWindow() {
		Connector = connector{
			User:     "",
			Password: "",
			Brokers:  []string{"127.0.0.1:9092"},
		}
	} else {
		Connector = connector{
			User:     "",
			Password: "",
			Brokers:  []string{"66.206.12.106:2909"},
		}
	}
}
