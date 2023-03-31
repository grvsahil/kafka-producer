package kafkaproducer

import (
	payment "github.com/grvsahil/golang-kafka/kafka-producer/api/v1/payments"
)

func (a *Application) SetupHandlers() {
	payment.RegisterHandlers(
		a.router,
		a.producer,
		a.log,
	)
}
