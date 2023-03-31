package kafkaproducer

import (
	"github.com/grvsahil/golang-kafka/kafka-producer/internal/producer"
)

func (a *Application) BuildProducer() *producer.Producer {
	producer := producer.NewProducer(a.log)
	return producer
}
