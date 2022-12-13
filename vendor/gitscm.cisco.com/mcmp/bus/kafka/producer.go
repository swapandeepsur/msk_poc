package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gitscm.cisco.com/mcmp/bus/config"
	"gitscm.cisco.com/mcmp/bus/kafka/pool"
)

// Producer provides the details for connecting to Kafka.
type Producer struct {
	pool  pool.Pool
	topic string
	log   logrus.FieldLogger
}

// NewProducer creates and configures a new Producer client as a connection pool.
func NewProducer(opts Options) (*Producer, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	p, err := pool.NewChannelPool(opts.Producer.InitCapacity, opts.Producer.MaxCapacity, makeFactory(opts))
	if err != nil {
		opts.Logger.Errorf("error while creating a new channel pool: %v", err)

		return nil, err
	}

	return &Producer{
		pool:  p,
		topic: opts.Topic,
		log:   opts.Logger,
	}, nil
}

func makeFactory(opts Options) pool.Factory {
	return func() (sarama.SyncProducer, error) {
		cfg := sarama.NewConfig()
		cfg.Producer.RequiredAcks = sarama.WaitForAll
		cfg.Producer.Retry.Max = viper.GetInt(config.ProducerMaxRetry)
		cfg.Producer.Return.Successes = true

		tlsConfig, err := LoadClientCertificate()
		if err != nil {
			opts.Logger.Errorf("error in loading client certificate: %v", err)

			return nil, err
		}

		if tlsConfig != nil {
			cfg.Net.TLS.Config = tlsConfig
			cfg.Net.TLS.Enable = true
		} else {
			opts.Logger.Warn("No client certificates provided; connection will be attempted without authentication")
		}

		return sarama.NewSyncProducer(opts.Hosts, cfg)
	}
}

// Publish writes a message on bus.
func (p *Producer) Publish(key string, msg []byte) error {
	client, err := p.pool.Get()
	if err != nil {
		p.log.Errorf("failed to get usable connection: %v", err)
		p.pool.MarkUnusable(client)

		return err
	}
	defer client.Close()

	_, _, err = client.SendMessage(&sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(msg),
	})

	return err
}

// Close will close the connection(s) to the bus.
func (p *Producer) Close() {
	p.pool.Close()
}
