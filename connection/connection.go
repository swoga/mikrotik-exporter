package connection

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/swoga/go-routeros"
)

type Connection struct {
	targetConnections *targetConnections
	Client            *routeros.Client
	inUse             bool
	healthy           bool
	lastUse           time.Time
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
	c.targetConnections.mu.Lock()
	defer c.targetConnections.mu.Unlock()

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
