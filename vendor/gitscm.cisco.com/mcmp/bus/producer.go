package bus

import (
	"gitscm.cisco.com/mcmp/bus/kafka"
)

// Producer defines a minimal interface for an Message Bus Producer.
type Producer interface {
	// Publish writes a named event and message to the Message Bus.
	Publish(string, []byte) error
	// Close releases any resources in use.
	Close()
}

// NewProducer creates and configures a Producer.
func NewProducer(opts Options) (Producer, error) {
	return kafka.NewProducer(opts.Options)
}
