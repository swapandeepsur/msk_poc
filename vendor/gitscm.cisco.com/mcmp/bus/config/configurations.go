/*
Package config defines the supported configuration options for event bus.

Example Configuration file (config.yaml)

	bus.producer:
		capacity:
			initial: 3
			maximum: 10
		flush.frequency: 500ms
*/
package config

import (
	"github.com/spf13/viper"
)

// all configuration keys.
const (
	// Environment Variable: "BUS_HOSTS".
	BusHosts = "bus.hosts"
	// Environment Variable: "EVENT_TOPIC".
	BusTopicEvent = "bus.topic.event"
	// Environment Variable: "LOGS_TOPIC".
	BusTopicLogs = "bus.topic.logs"

	// Environment Variable: "BUS_PRODUCER_INIT_CAP"		Default: 3.
	ProducerInitCap = "bus.producer.capacity.initial"
	// Environment Variable: "BUS_PRODUCER_MAX_CAP"			Default: 10.
	ProducerMaxCap = "bus.producer.capacity.maximum"
	// Default: 10.
	ProducerMaxRetry = "bus.producer.retry.maximum"
	// Default: 500ms.
	ProducerFlushFrequency = "bus.producer.flush.frequency"

	// Environment Variable: "KAFKA_CLIENT_CERT".
	KafkaClientCertLocation = "kafka.certs.client.certificate.location"
	// Environment Variable: "KAFKA_CLIENT_KEY".
	KafkaClientKeyLocation = "kafka.certs.client.key.location"
	// Environment Variable: "KAFKA_CACERT".
	KafkaCACertLocation = "kafka.certs.ca.certificate.location"
)

func init() {
	viper.SetDefault(ProducerInitCap, 3)
	viper.SetDefault(ProducerMaxCap, 10)
	viper.SetDefault(ProducerMaxRetry, 10)
	viper.SetDefault(ProducerFlushFrequency, "500ms")

	_ = viper.BindEnv(BusHosts, "BUS_HOSTS")
	_ = viper.BindEnv(BusTopicEvent, "EVENT_TOPIC")
	_ = viper.BindEnv(BusTopicLogs, "LOGS_TOPIC")

	_ = viper.BindEnv(ProducerInitCap, "BUS_PRODUCER_INIT_CAP")
	_ = viper.BindEnv(ProducerMaxCap, "BUS_PRODUCER_MAX_CAP")

	_ = viper.BindEnv(KafkaClientCertLocation, "KAFKA_CLIENT_CERT")
	_ = viper.BindEnv(KafkaClientKeyLocation, "KAFKA_CLIENT_KEY")
	_ = viper.BindEnv(KafkaCACertLocation, "KAFKA_CACERT")
}
