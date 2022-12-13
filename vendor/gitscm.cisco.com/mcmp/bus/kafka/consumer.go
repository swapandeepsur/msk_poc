package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"gitscm.cisco.com/ccdev/go-common/sets"

	"gitscm.cisco.com/mcmp/bus/errors"
)

// Handler represents a generic function that accepts a message key
// and byte array containing the received message, and handles the message.
type Handler func(string, []byte)

// Consumer provides a basic Kafka Consumer client.
type Consumer struct {
	client   sarama.Consumer
	listener sarama.PartitionConsumer
	log      logrus.FieldLogger
	handler  Handler
	events   sets.String
}

// NewConsumer creates and configures new Consumer.
func NewConsumer(opts Options, h Handler, events ...string) (*Consumer, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	if h == nil {
		return nil, errors.ConfigurationError("no handler provided")
	}

	if len(events) == 0 {
		return nil, errors.ConfigurationError("no events subscribed")
	}

	c := &Consumer{
		log:     opts.Logger,
		handler: h,
		events:  sets.NewString(events...),
	}

	if err := c.configure(opts); err != nil {
		opts.Logger.Errorf("error in configuring consumer: %v", err)

		return nil, err
	}

	return c, nil
}

func (c *Consumer) configure(opts Options) error {
	tlsConfig, err := LoadClientCertificate()
	if err != nil {
		return err
	}

	var config *sarama.Config

	if tlsConfig != nil {
		config = sarama.NewConfig()
		config.Net.TLS.Config = tlsConfig
		config.Net.TLS.Enable = true
	} else {
		c.log.Warn("No client certificates provided; connection will be attempted without authentication")
	}

	c.client, err = sarama.NewConsumer(opts.Hosts, config)
	if err != nil {
		return err
	}

	c.listener, err = c.client.ConsumePartition(opts.Topic, 0, sarama.OffsetNewest)

	return err
}

// Close closes resources in use.
func (c *Consumer) Close() {
	if c.listener != nil {
		// always returned as nil within library
		_ = c.listener.Close()
		_ = c.client.Close()
	}

	if c.client != nil {
		// always returned as nil within library
		_ = c.client.Close()
	}

	c.log.Info("consumer has been closed")
}

// Start will start listening for messages and call the provided handler
// for any event that was subscribed. The channel is used to stop listening
// for messages.
func (c *Consumer) Start(stop <-chan bool) {
ConsumerLoop:
	for {
		select {
		case msg := <-c.listener.Messages():
			key := string(msg.Key)
			c.log.Infof("Received message on key: %s", key)
			if c.events.Has(key) {
				c.log.Debugf("invoking handler for message consumed %v with key: %s", msg.Value, key)
				c.handler(key, msg.Value)
			} else {
				c.log.Debugf("no subscription for key: %s", key)
			}
		case <-stop:
			break ConsumerLoop
		}
	}

	c.log.Info("consumer has stopped as requested")
}
