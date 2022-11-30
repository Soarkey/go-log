package main

import "go-log/log-transfer/kafka"

const TOPIC = "log"

func main() {
	k := kafka.NewConsumer([]string{"kafka:9092"}, TOPIC)
	k.Start()
}
