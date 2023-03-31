package kafkaproducer

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/grvsahil/golang-kafka/kafka-producer/internal/log"
	"github.com/grvsahil/golang-kafka/kafka-producer/internal/producer"
)

const (
	server = "localhost:8088"
)

type Application struct {
	log      log.Logger
	router   *gin.Engine
	producer *producer.Producer
}

func (a *Application) Init(ctx context.Context) {
	log := log.New().With(ctx)
	a.log = log

	producer := a.BuildProducer()
	a.producer = producer

	router := gin.Default()
	a.router = router
	a.SetupHandlers()

	err := a.router.Run(server)
	if err != nil {
		a.log.Fatalf("failed to start server: %s ", err)
		return
	}
}
