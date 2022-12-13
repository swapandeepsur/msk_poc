/*
Package kafka provides an event bus implementation using kafka.

The implementation is a thin wrapper for github.com/Shopify/sarama that handles configuring

	- Consumer
	- Producer (uses a pool)
	- AsyncProducer
*/
package kafka
