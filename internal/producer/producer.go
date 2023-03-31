package producer

import (
	"time"

	"github.com/grvsahil/golang-kafka/kafka-producer/internal/log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	log log.Logger
}

func NewProducer(log log.Logger) *Producer {
	return &Producer{
		log: log,
	}
}

func (prod Producer) ProduceWithRetry(bootstrapServer string, topic string, value []byte) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"client.id":         "my-client",
		"acks":              "all"})
	if err != nil {
		prod.log.Errorf("Failed to create producer: %s\n", err)
		return err
	}
	defer p.Close()

	// retry 3 times after 2 seconds
	for i := 0; i < 3; i++ {
		err = produce(p, topic, value, prod.log)
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

func produce(p *kafka.Producer, topic string, value []byte, log log.Logger) error {
	deliveryCh := make(chan kafka.Event, 10000)
	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value},
		deliveryCh,
	)
	if err != nil {
		log.Errorf("Error producing message ", err)
		return err
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Errorf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					log.Infof("Successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()

	return nil
}
