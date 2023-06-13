package producer

import (
	"time"

	"github.com/grvsahil/golang-kafka/kafka-producer/internal/log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	log           log.Logger
	producers     []*kafka.Producer
	producerIndex int
}

func NewProducer(log log.Logger, bootstrapServer string) (*Producer, error) {
	prod1, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"client.id":         "my-client",
		"acks":              "all"})
	if err != nil {
		log.Errorf("Failed to create producer1: %s\n", err)
		return nil, err
	}

	prod2, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"client.id":         "my-client",
		"acks":              "all"})
	if err != nil {
		log.Errorf("Failed to create producer2: %s\n", err)
		return nil, err
	}

	prod3, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"client.id":         "my-client",
		"acks":              "all"})
	if err != nil {
		log.Errorf("Failed to create producer3: %s\n", err)
		return nil, err
	}

	producers := []*kafka.Producer{prod1, prod2, prod3}
	
	return &Producer{
		log:           log,
		producers:     producers,
		producerIndex: 0,
	}, nil
}

func (prod Producer) ProduceWithRetry(topic string, value []byte) error {
	// retry 3 times after 2 seconds
	var err error
	for i := 0; i < 3; i++ {
		err = prod.produce(topic, value)
		if err == nil {
			prod.log.Infof("Produced message after retry %d", i)
			return nil
		}
		prod.log.Infof("Failed to produce message, retrying in 2 seconds...")
		time.Sleep(time.Second * 2)
	}

	prod.log.Error("Failed to produce message after retrying 3 times")
	return err
}

func (prod Producer) produce(topic string, value []byte) error {
	deliveryCh := make(chan kafka.Event)
	err := prod.producers[prod.producerIndex].Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value},
		deliveryCh,
	)
	if err != nil {
		prod.log.Errorf("Error producing message ", err)
		return err
	}

	e := <-deliveryCh
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		prod.log.Errorf("Delivery failed: %v using producer%d", m.TopicPartition.Error, prod.producerIndex)
	} else {
		prod.log.Infof("Delivered message to topic %s [%d] at offset %v using producer%d",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset, prod.producerIndex)
	}

	close(deliveryCh)
	prod.producerIndex++
	if prod.producerIndex == len(prod.producers) {
		prod.producerIndex = 0
	}

	return nil
}
