package kafkaproducer

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/grvsahil/golang-kafka/kafka-producer/internal/log"
	"github.com/grvsahil/golang-kafka/kafka-producer/internal/producer"
)

const (
	server          = "localhost:8088"
	bootstrapServer = "localhost:9092"
	timeout         = 5 * time.Second
)

type Application struct {
	log        log.Logger
	router     *gin.Engine
	httpServer *http.Server
	producer   *producer.Producer
}

func (a *Application) Init(ctx context.Context) {
	log := log.New().With(ctx)
	a.log = log

	producer, err := a.BuildProducer()
	if err != nil {
		a.log.Fatalf("failed to initialize producer: %s ", err)
		return
	}
	a.producer = producer

	router := gin.Default()
	a.router = router
	a.SetupHandlers()
}

func (a *Application) Start(ctx context.Context) {
	a.httpServer = &http.Server{
		Addr:              server,
		Handler:           a.router,
		ReadHeaderTimeout: timeout,
	}
	go func() {
		defer a.log.Infof("server stopped listening")
		if err := a.httpServer.ListenAndServe(); err != nil {
			a.log.Errorf("failed to listen and serve: %v ", err)
			return
		}
		return
	}()
	a.log.Infof("http server started on %s ...", server)
}

func (a *Application) Stop(ctx context.Context) {
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		a.log.Error(err)
	}
	a.CloseProducer()
	a.log.Info("shutting down....")
	return
}
