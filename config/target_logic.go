package config

import (
	"context"

	"github.com/swoga/go-routeros"
)

func (t *Target) Dial() (*routeros.Client, error) {
	c, err := routeros.DialContext(context.Background(), t.Address, t.Credentials.Username, t.Credentials.Password, t.timeoutDuration)
	if err != nil {
		return nil, err
	}
	c.Queue = t.Queue
	return c, err
}
