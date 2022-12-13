/*
Package pool implements a pool of kafka SyncProducer interfaces to manage and reuse them.
*/
package pool

import (
	"errors"

	"github.com/Shopify/sarama"
)

// ErrClosed is used to signal that the pool is closed.
var ErrClosed = errors.New("pool is closed")

// Pool interface describes a pool implementation that holds sarama.SyncProducer connections.
type Pool interface {
	// Get returns a new SyncProducer instance from the backing pool
	// Calling Close() on the object does not terminate connection to bus
	// It rather retutns it back to the pool to be re-used at a later
	// point in time
	Get() (sarama.SyncProducer, error)

	// Markunusable is used by the client to mark a connection unusable and hence
	// allowing the pool to close it rather than reclaiming it
	MarkUnusable(p sarama.SyncProducer)

	// Close closes the pool and terminates all connections
	Close()

	// Returns the number of active connections in the pool
	Len() int
}
