package producer

type Produce interface {
	ProduceWithRetry(bootstrapServer string, topic string, value []byte) error
}
