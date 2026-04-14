package event

type Producer interface {
	Send(topic, key string, payload []byte) error
}
