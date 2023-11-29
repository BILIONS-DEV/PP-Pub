package config

import (
	"source/pkg/helpers"
	"source/pkg/utility"
)

type Config struct {
	Kafka     *Kafka
	Mysql     *Mysql
	Aerospike *Aerospike
}

func NewConfig() *Config {
	return &Config{
		Mysql:     newConfigMysql(),
		Aerospike: newConfigAerospike(),
		Kafka:     newConfigKafka(),
	}
}

type Mysql struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	Encoding string
}

func newConfigMysql() *Mysql {
	mysqlConfig := &Mysql{
		Username: "apacadmin",
		Password: "iK29&6%!9XKjs@",
		//Host:     "23.92.74.98", //=> PublicIP server backend
		Host:     "192.168.9.10", //=> PrivateIP server backend
		Port:     "9682",
		Database: "apac_ss",
		Encoding: "utf8mb4",
	}
	if utility.IsWindow() {
		mysqlConfig.Host = "localhost"
		mysqlConfig.Port = "3306"
	}
	if !utility.IsWindow() && utility.IsDev() {
		mysqlConfig.Database = "apac_ss_dev"
	}
	if utility.IsDemo() {
		if !utility.IsWindow() {
			mysqlConfig.Database = "apac_ss_demo"
		} else {
			mysqlConfig.Database = "apac_ss"
		}
	}

	if helpers.IsHaiMode() {
		mysqlConfig = &Mysql{
			Username: "assyrian",
			Password: "xadwXLLKa*(as2",
			//Host:     "23.92.74.98", //=> PublicIP server backend
			Host:     "localhost", //=> PrivateIP server backend
			Port:     "3306",
			Database: "apacadmin",
			Encoding: "utf8mb4",
		}
	}
	return mysqlConfig
}

type Aerospike struct {
	Host      string
	Port      int
	Namespace string
}

func newConfigAerospike() *Aerospike {
	aerospikeConfig := &Aerospike{
		Host:      "192.168.9.13",
		Port:      3000,
		Namespace: "disk_cache",
	}
	if utility.IsWindow() {
		aerospikeConfig = &Aerospike{
			Host:      "127.0.0.1",
			Port:      3000,
			Namespace: "common",
		}
	}
	return aerospikeConfig
}

type Kafka struct {
	Brokers []string
	Topic   string
	GroupID string
}

func newConfigKafka() *Kafka {
	kafkaConfig := &Kafka{
		Brokers: []string{"127.0.0.1:2909"},
		Topic:   "demand-bid",
		GroupID: "demand-bid-group",
	}
	if utility.IsWindow() {
		kafkaConfig.Brokers = []string{"kafka:9092"}
	}
	return kafkaConfig
}
