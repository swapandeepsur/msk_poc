package bus

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gitscm.cisco.com/mcmp/bus/config"
	"gitscm.cisco.com/mcmp/bus/kafka"
)

// Options provides the available configurations for Consumers and Producers.
type Options struct {
	kafka.Options
}

// DefaultOptions creates an instance of Options with default values for each
// of the Options attibutes.
func DefaultOptions() Options {
	opts := Options{}
	opts.Logger = defaultLogger()
	opts.Hosts = cleanHosts(viper.GetString(config.BusHosts))
	opts.Topic = viper.GetString(config.BusTopicEvent)
	opts.Producer.InitCapacity = viper.GetInt(config.ProducerInitCap)
	opts.Producer.MaxCapacity = viper.GetInt(config.ProducerMaxCap)

	return opts
}

// Validate verifies the values provided for Options are valid.
func (o Options) Validate() error {
	return o.Options.Validate()
}

func defaultLogger() logrus.FieldLogger {
	l := logrus.New()
	// configure the default logger to include timestamps and quote empty fields
	// to make visually seeing an empty Field easier.
	l.Formatter = &logrus.TextFormatter{
		FullTimestamp:    true,
		QuoteEmptyFields: true,
	}
	l.SetLevel(logrus.ErrorLevel)

	return l
}

func cleanHosts(val string) []string {
	hosts := make([]string, 0)

	for _, host := range strings.Split(val, ",") {
		if h := strings.TrimSpace(host); h != "" {
			hosts = append(hosts, h)
		}
	}

	return hosts
}
