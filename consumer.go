package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gitscm.cisco.com/mcmp/bus"
)

var Handler = func(string, []byte) {}

func main() {
	// initialize bus options
	opts := bus.Options{}

	// set options
	opts.Hosts = []string{
		"b-1.runonmskpocnoiam.cjf0rv.c18.kafka.us-east-1.amazonaws.com:9092",
		"b-2.runonmskpocnoiam.cjf0rv.c18.kafka.us-east-1.amazonaws.com:9092",
		"b-3.runonmskpocnoiam.cjf0rv.c18.kafka.us-east-1.amazonaws.com:9092",
	}

	opts.Topic = "runon-msk-poc-noiam-topic"
	opts.Options.Logger = logrus.New()

	fmt.Println("Creating Consumer")

	// create new producer
	consumer, err := bus.NewConsumer(opts, Handler, "Runon")
	if err != nil {
		panic(err)
	}

	defer consumer.Close()
	channel := make(chan bool)

	consumer.Start(channel)
}
