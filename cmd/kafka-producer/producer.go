package kafkaproducer

import (
	"github.com/grvsahil/golang-kafka/kafka-producer/internal/producer"
)

func (a *Application) BuildProducer() (*producer.Producer, error) {
	producer, err := producer.NewProducer(a.log, bootstrapServer)
	return producer, err
}

func (a *Application) CloseProducer() {
	producers := a.producer.Producers
	for _, producer := range producers {
		producer.Close()
	}
}
