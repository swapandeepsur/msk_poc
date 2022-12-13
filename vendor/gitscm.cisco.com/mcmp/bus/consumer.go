package bus

import (
	"gitscm.cisco.com/mcmp/bus/kafka"
)

// Handler represents a generic function that accepts a message key
// and byte array containing the received message, and handles the message.
type Handler = kafka.Handler

// Consumer defines a minimal interface for an Message Bus Consumer.
type Consumer interface {
	// Start will start listening for messages and call the provided handler
	// for any event that was subcribed. The channel is used to stop listening
	// for messages.
	Start(stop <-chan bool)
	// Close closes any resources in use.
	Close()
}

// NewConsumer creates and configures a Consumer.
func NewConsumer(opts Options, h Handler, events ...string) (Consumer, error) {
	return kafka.NewConsumer(opts.Options, h, events...)
}
