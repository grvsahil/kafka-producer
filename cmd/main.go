package main

import (
	"context"

	kafkaProducer "github.com/grvsahil/golang-kafka/kafka-producer/cmd/kafka-producer"
)

const (
	server = "localhost:8088"
)

func main() {
	application := &kafkaProducer.Application{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	application.Init(ctx)
}
