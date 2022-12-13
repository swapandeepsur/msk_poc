package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"

	"gitscm.cisco.com/mcmp/bus/config"
	"gitscm.cisco.com/mcmp/bus/errors"
)

// NewAsyncProducer creates an AsyncProducer using the provided options.
func NewAsyncProducer(opts Options) (sarama.AsyncProducer, error) {
	if len(opts.Hosts) == 0 || opts.Hosts[0] == "" {
		return nil, errors.ConfigurationError("no host(s) provided")
	}

	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForLocal
	cfg.Producer.Compression = sarama.CompressionSnappy
	cfg.Producer.Flush.Frequency = viper.GetDuration(config.ProducerFlushFrequency)

	tlsConfig, err := LoadClientCertificate()
	if err != nil {
		return nil, err
	}

	if tlsConfig != nil {
		cfg.Net.TLS.Config = tlsConfig
		cfg.Net.TLS.Enable = true
	}

	return sarama.NewAsyncProducer(opts.Hosts, cfg)
}
