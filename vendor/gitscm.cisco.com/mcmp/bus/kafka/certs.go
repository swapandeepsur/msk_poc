package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/spf13/viper"

	"gitscm.cisco.com/mcmp/bus/config"
)

func createTLSConfiguration(certFile, keyFile, caFile string) (*tls.Config, error) {
	// if a cert was not specified by the environment variable
	// then return `nil`
	if certFile == "" || keyFile == "" || caFile == "" {
		return nil, nil
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,
	}, nil
}

// LoadClientCertificate loads a certficate from a file specified by
// environment variable `KAFKA_CLIENT_CERT` and creates a tls.Config instance.
// If the environment variable is not set or has no value, the tls.Config will be `nil`.
func LoadClientCertificate() (*tls.Config, error) {
	return createTLSConfiguration(viper.GetString(config.KafkaClientCertLocation), viper.GetString(config.KafkaClientKeyLocation), viper.GetString(config.KafkaCACertLocation))
}
