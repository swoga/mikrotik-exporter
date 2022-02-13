package connection

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/swoga/go-routeros"
)

type Connection struct {
	mu      sync.RWMutex
	Client  *routeros.Client
	inUse   bool
	healthy bool
	lastUse time.Time
}

func (c *Connection) check(log zerolog.Logger) bool {
	log.Trace().Msg("run healthcheck")
	_, err := c.Client.Run("/system/identity/print")
	c.healthy = err == nil
	if err != nil {
		log.Warn().Err(err).Msg("error during healthcheck")
	} else {
		log.Trace().Msg("healthcheck successful")
	}
	return c.healthy
}

func (c *Connection) freeInternal(log zerolog.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()

	log.Trace().Msg("free connection")
	c.inUse = false
	c.lastUse = time.Now()
	c.check(log)
}

func (c *Connection) Free(log zerolog.Logger) {
	if c == nil {
		return
	}
	// do not block caller
	go c.freeInternal(log)
}

// Check if the connection is usable, if yes mark as used (blocks during healthcheck)
func (c *Connection) Use(log zerolog.Logger) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.inUse {
		log.Trace().Msg("skip connection in use")
		return false
	}
	if !c.healthy {
		log.Trace().Msg("skip unhealthy connection")
		return false
	}
	if !c.check(log) {
		return false
	}
	log.Trace().Msg("return existing connection")
	c.inUse = true
	return true
}

func (c *Connection) IsInUse() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.inUse
}

func (c *Connection) IsHealthy() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.healthy
}

func (c *Connection) GetLastUse() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.lastUse
}
