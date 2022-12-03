package main

import (
	"go-log/log-transfer/es"
	"go-log/log-transfer/kafka"
)

const (
	TOPIC         = "go-log"
	WORKERS       = 2
	CHAN_MAX_SIZE = 16
)

func main() {
	k := kafka.NewConsumer(
		[]string{"kafka:9092"},
		TOPIC,
		es.New("http://es:9200", WORKERS, CHAN_MAX_SIZE),
	)
	k.Start()
}
