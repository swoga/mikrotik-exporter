package connection

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/swoga/mikrotik-exporter/config"
)

type targetConnections struct {
	targetName  string
	connections map[*Connection]struct{}
	mu          sync.Mutex
}

func createTargetConnections(targetName string) *targetConnections {
	tc := targetConnections{
		targetName:  targetName,
		connections: make(map[*Connection]struct{}),
	}
	return &tc
}

func (tc *targetConnections) get(log zerolog.Logger, target *config.Target) (*Connection, error) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	log.Trace().Msg("try to find existing connection")
	for c := range tc.connections {
		if c.inUse {
			log.Trace().Msg("skip connection in use")
			continue
		}
		if !c.healthy {
			log.Trace().Msg("skip unhealthy connection")
			continue
		}
		if !c.check() {
			continue
		}
		log.Trace().Msg("return existing connection")
		c.inUse = true
		return c, nil
	}

	log.Info().Msg("connect to target")
	client, err := target.Dial()
	if err != nil {
		return nil, err
	}
	connection := Connection{
		targetConnections: tc,
		Client:            client,
	}
	tc.connections[&connection] = struct{}{}

	return &connection, nil
}

func (tc *targetConnections) cleanup(useTimeout time.Duration) {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	for c := range tc.connections {
		if c.inUse {
			continue
		}

		expired := time.Since(c.lastUse) > useTimeout

		if !c.healthy || expired {
			log.Logger.Info().Str("target", tc.targetName).Bool("healthy", c.healthy).Bool("expired", expired).Time("lastUse", c.lastUse).Msg("close and cleanup connection")
			c.Client.Close()
			delete(tc.connections, c)
		}
	}
}
