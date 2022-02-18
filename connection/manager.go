package connection

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/swoga/mikrotik-exporter/config"
)

func CreateConnectionManager(cleanupInterval time.Duration, useTimeout time.Duration) *ConnectionManager {
	cm := ConnectionManager{
		targetConnections: make(map[string]*targetConnections),
		cleanupInterval:   cleanupInterval,
		useTimeout:        useTimeout,
	}
	return &cm
}

type ConnectionManager struct {
	mu                sync.Mutex
	targetConnections map[string]*targetConnections
	cleanupInterval   time.Duration
	useTimeout        time.Duration
}

func (cm *ConnectionManager) Get(log zerolog.Logger, target *config.Target) (*Connection, zerolog.Logger, error) {
	cm.mu.Lock()

	tc, found := cm.targetConnections[target.Name]
	if !found {
		log.Trace().Msg("first connection to this target")
		tc = createTargetConnections(log, target.Name, cm.cleanupInterval, cm.useTimeout)
		cm.targetConnections[target.Name] = tc
	} else {
		log.Trace().Msg("target found in connection cache")
	}

	cm.mu.Unlock()

	return tc.get(log, target)
}
