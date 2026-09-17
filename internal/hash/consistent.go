package hash

import (
	"fmt"

	"github.com/buraksezer/consistent"
)

const (
	defaultPartitionCount    int     = 1021
	defaultReplicationFactor int     = 100
	defaultLoad              float64 = 1.05
)

// newConsistent creates a new consistent hash implementation and catches package panics.
func newConsistent(members []consistent.Member, cfg consistent.Config) (c *consistent.Consistent, err error) {
	// Required as package panics on error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("invalid consistent hash configuration: %v", r)
		}
	}()

	return consistent.New(members, cfg), nil
}
