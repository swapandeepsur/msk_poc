package kafka

import (
	"github.com/sirupsen/logrus"

	"gitscm.cisco.com/mcmp/bus/errors"
)

// Options is used runtime to send the needed config params.
type Options struct {
	Logger   logrus.FieldLogger
	Hosts    []string
	Topic    string
	Producer struct {
		InitCapacity int
		MaxCapacity  int
	}
}

// Validate verifies the values provided for Options are valid.
func (o Options) Validate() error {
	if len(o.Hosts) == 0 || o.Hosts[0] == "" {
		return errors.ConfigurationError("no host(s) provided")
	}

	if o.Topic == "" {
		return errors.ConfigurationError("no topic provided")
	}

	if o.Logger == nil {
		return errors.ConfigurationError("no logger provided")
	}

	return nil
}
