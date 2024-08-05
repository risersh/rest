package util

import (
	"log"

	"github.com/mateothegreat/go-rabbitmq/producer"
	"github.com/risersh/rest/conf"
)

var Producer *producer.Producer

func ConnectRabbitMQ() {
	Producer = &producer.Producer{}
	err := Producer.Connect(conf.Config.RabbitMQ.URI)
	if err != nil {
		log.Fatalf("could not connect to rabbitmq: %v", err)
	}
}
