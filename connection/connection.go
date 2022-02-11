package connection

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/swoga/go-routeros"
)

type Connection struct {
	targetConnections *targetConnections
	Client            *routeros.Client
	inUse             bool
	healthy           bool
	lastUse           time.Time
}

func (c *Connection) check() bool {
	checkLogger := log.Logger.With().Str("target", c.targetConnections.targetName).Logger()
	checkLogger.Trace().Msg("run healthcheck")
	_, err := c.Client.Run("/system/identity/print")
	c.healthy = err == nil
	if err != nil {
		checkLogger.Warn().Err(err).Msg("error during healthcheck")
	} else {
		checkLogger.Trace().Msg("healthcheck successful")
	}
	return c.healthy
}

func (c *Connection) freeInternal() {
	c.targetConnections.mu.Lock()
	defer c.targetConnections.mu.Unlock()

	log.Logger.Trace().Str("target", c.targetConnections.targetName).Msg("free connection")
	c.inUse = false
	c.lastUse = time.Now()
	c.check()
}

func (c *Connection) Free() {
	if c == nil {
		return
	}
	// do not block caller
	go c.freeInternal()
}
