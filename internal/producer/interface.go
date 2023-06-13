package producer

type Produce interface {
	ProduceWithRetry(topic string, value []byte) error
}
