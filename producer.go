package main

import (
	"gitscm.cisco.com/mcmp/bus"
	logutil "gitscm.cisco.com/mcmp/utils/log"
)

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

	opts.Producer.InitCapacity = 1
	opts.Producer.MaxCapacity = 2
	opts.Logger = logutil.GetLogger()

	// create new producer
	producer, err := bus.NewProducer(opts)
	if err != nil {
		panic(err)
	}

	// close the producer once messages are published
	defer producer.Close()

	// publish the event with event key
	if err = producer.Publish("Runon", []byte("NEW MESSAGE11")); err != nil {
		panic(err)
	}
}
