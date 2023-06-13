package payment

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/grvsahil/golang-kafka/kafka-producer/internal/log"
	"github.com/grvsahil/golang-kafka/kafka-producer/internal/payment"
	"github.com/grvsahil/golang-kafka/kafka-producer/internal/producer"
)

const (
	topic = "payment"
)

func RegisterHandlers(router *gin.Engine, producer producer.Produce, logger log.Logger) {
	res := resource{producer, logger}
	router.POST("/entry", res.Entry)
}

type resource struct {
	producer producer.Produce
	log      log.Logger
}

func (res resource) Entry(c *gin.Context) {
	var payment payment.Payment
	if err := c.BindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	paymentBytes, err := json.Marshal(payment)
	if err != nil {
		fmt.Println("Error marshalling payment:", err)
		return
	}

	err = res.producer.ProduceWithRetry(topic, paymentBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Payment received!"})
}
