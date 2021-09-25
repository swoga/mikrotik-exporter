package connection

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/swoga/mikrotik-exporter/config"
)

func CreateConnectionManager() *ConnectionManager {
	cm := ConnectionManager{
		targetConnections: make(map[string]*targetConnections),
		stopCleanup:       make(chan bool),
	}
	cm.StartCleanup()
	return &cm
}

type ConnectionManager struct {
	mu                sync.Mutex
	targetConnections map[string]*targetConnections
	stopCleanup       chan (bool)
}

func (cm *ConnectionManager) Get(log zerolog.Logger, target *config.Target) (*Connection, error) {
	cm.mu.Lock()

	tc, found := cm.targetConnections[target.Name]
	if !found {
		log.Trace().Msg("first connection to this target")
		tc = createTargetConnections(target.Name)
		cm.targetConnections[target.Name] = tc
	} else {
		log.Trace().Msg("target found in connection cache")
	}

	cm.mu.Unlock()

	return tc.get(log, target)
}

func (cm *ConnectionManager) cleanup() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	log.Logger.Trace().Msg("run cleanup")
	for _, tc := range cm.targetConnections {
		go tc.cleanup()
	}
}

func (cm *ConnectionManager) StartCleanup() {
	log.Logger.Debug().Msg("start cleanup job")
	ticker := time.NewTicker(time.Minute)

	go func() {
		for {
			select {
			case <-cm.stopCleanup:
				ticker.Stop()
				return
			case <-ticker.C:
				cm.cleanup()
				continue
			}
		}
	}()
}

func (cm *ConnectionManager) StopCleanup() {
	log.Logger.Debug().Msg("stop cleanup job")

	select {
	case cm.stopCleanup <- true:
		break
	default:
		break
	}
}
