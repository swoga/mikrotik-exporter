package connection

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/swoga/go-routeros"
)

type Connection struct {
	mu       sync.RWMutex
	Client   *routeros.Client
	username string
	address  string
	inUse    bool
	healthy  bool
	lastUse  time.Time
	id       int
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
func (c *Connection) Use(log zerolog.Logger, address string, username string) (bool, zerolog.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()

	useLog := log.With().Int("connection_no", c.id).Logger()

	if c.inUse {
		useLog.Trace().Msg("skip connection in use")
		return false, log
	}
	if !c.healthy {
		useLog.Trace().Msg("skip unhealthy connection")
		return false, log
	}
	if c.address != address {
		useLog.Trace().Msg("target address has changed")
		return false, log
	}
	if c.username != username {
		useLog.Trace().Msg("target username has changed")
		return false, log
	}
	if !c.check(useLog) {
		return false, log
	}
	useLog.Trace().Msg("return existing connection")
	c.inUse = true
	return true, useLog
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
