package pool

import (
	stderrors "errors"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"

	"gitscm.cisco.com/mcmp/bus/errors"
)

// channelPool implements the Pool interface based on buffered channels.
type channelPool struct {
	// storage for our sarama.SyncProducer clients
	mu      sync.Mutex
	clients chan sarama.SyncProducer

	// sarama.SyncProducer factory function
	factory Factory
}

// Factory is used to create new SyncProducer clients.
type Factory func() (sarama.SyncProducer, error)

// NewChannelPool returns a new pool based on buffered channels with an initial
// capacity and maximum capacity. Factory is used when initial capacity is
// greater than zero to fill the pool. A zero initialCap doesn't fill the Pool
// until a new Get() is called. During a Get(), If there is no new client
// available in the pool, a new client will be created via the Factory() method.
func NewChannelPool(initialCap, maxCap int, factory Factory) (Pool, error) {
	if initialCap < 0 || maxCap <= 0 || initialCap > maxCap {
		return nil, errors.ConfigurationError("invalid capacity")
	}

	c := &channelPool{
		clients: make(chan sarama.SyncProducer, maxCap),
		factory: factory,
	}

	// create initial clients, if something goes wrong,
	// just close the pool error out.
	for i := 0; i < initialCap; i++ {
		client, err := factory()
		if err != nil {
			c.Close()

			return nil, stderrors.Unwrap(fmt.Errorf("factory is not able to fill the pool: %w", err))
		}

		c.clients <- client
	}

	return c, nil
}

func (c *channelPool) getConnsAndFactory() (chan sarama.SyncProducer, Factory) {
	c.mu.Lock()
	clients := c.clients
	factory := c.factory
	c.mu.Unlock()

	return clients, factory
}

// Get implements the Pool interfaces Get() method. If there is no new
// client available in the pool, a new client will be created via the
// Factory() method.
func (c *channelPool) Get() (sarama.SyncProducer, error) {
	clients, factory := c.getConnsAndFactory()
	if clients == nil {
		return nil, ErrClosed
	}

	select {
	case client := <-clients:
		if client == nil {
			return nil, ErrClosed
		}

		return c.wrapClient(client), nil
	default:
		p, err := factory()
		if err != nil {
			return nil, err
		}

		return c.wrapClient(p), nil
	}
}

func (c *channelPool) MarkUnusable(p sarama.SyncProducer) {
	if cp, ok := p.(*ClientPool); ok {
		cp.markUnusable()
	}
}

// put is used to return SyncProducer clients back to bus channel pool.
func (c *channelPool) put(p sarama.SyncProducer) {
	if p == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.clients == nil {
		// pool is closed, close passed client
		c.MarkUnusable(p)
		_ = p.Close()

		return
	}

	// put the resource back into the pool. If the pool is full, this will
	// block and the default case will be executed.
	select {
	case c.clients <- p:
		return
	default:
		// pool is full, close passed client
		c.MarkUnusable(p)
		_ = p.Close()

		return
	}
}

// Close is called to shutdown the pool and rendering all the clients unusable.
func (c *channelPool) Close() {
	c.mu.Lock()
	clients := c.clients
	c.clients = nil
	c.factory = nil
	c.mu.Unlock()

	if clients == nil {
		return
	}

	close(clients)

	for client := range clients {
		c.MarkUnusable(client)
		// Close has error as part of signature but it is always nil
		_ = client.Close()
	}
}

func (c *channelPool) Len() int {
	clients, _ := c.getConnsAndFactory()

	return len(clients)
}
