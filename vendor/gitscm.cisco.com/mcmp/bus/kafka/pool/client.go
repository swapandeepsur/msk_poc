package pool

import (
	"sync"

	"github.com/Shopify/sarama"
)

// ClientPool is a wrapper around sarama.SyncProducer to modify the the behavior of
// sarama.SyncProducer's Close() method.
type ClientPool struct {
	sarama.SyncProducer
	mu       sync.RWMutex
	c        *channelPool
	unusable bool
}

// Close puts the given client back to the pool instead of closing it.
func (p *ClientPool) Close() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.unusable {
		if p.SyncProducer != nil {
			return p.SyncProducer.Close()
		}
	} else {
		p.c.put(p.SyncProducer)
	}

	return nil
}

// markUnusable marks the client not usable any more, to let the pool close it instead of returning it to pool.
func (p *ClientPool) markUnusable() {
	p.mu.Lock()
	p.unusable = true
	p.mu.Unlock()
}

// wrapClient wraps a standard sarama.SyncProducer to a clientPool of sarama.SyncProducer.
func (c *channelPool) wrapClient(p sarama.SyncProducer) sarama.SyncProducer {
	return &ClientPool{
		c:            c,
		SyncProducer: p,
	}
}
