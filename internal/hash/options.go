package hash

import (
	"fmt"

	"github.com/buraksezer/consistent"
)

// options represents the configuration options for creating a consistent hash ring.
type options struct {
	members           []consistent.Member
	partitionCount    int
	replicationFactor int
	load              float64
}

// OptionsFunc is a function type that modifies the options for creating a consistent hash ring.
type OptionsFunc func(*options) error

// WithMembers sets the initial members for the consistent hash ring.
func WithMembers(members []string) OptionsFunc {
	return func(o *options) error {
		l := len(members)
		if l == 0 {
			return nil
		}

		o.members = make([]consistent.Member, l)
		for i, m := range members {
			o.members[i] = member(m)
		}
		return nil
	}
}

// WithPartitionCount sets the partition count for the consistent hash ring.
func WithPartitionCount(partitionCount int) OptionsFunc {
	return func(o *options) error {
		if partitionCount < 0 {
			return fmt.Errorf("invalid partition count: %d", partitionCount)
		}

		o.partitionCount = partitionCount
		return nil
	}
}

// WithReplicationFactor sets the replication factor for the consistent hash ring.
func WithReplicationFactor(replicationFactor int) OptionsFunc {
	return func(o *options) error {
		if replicationFactor < 0 {
			return fmt.Errorf("invalid replication factor: %d", replicationFactor)
		}

		o.replicationFactor = replicationFactor
		return nil
	}
}

// WithLoad sets the load factor for the consistent hash ring.
func WithLoad(load float64) OptionsFunc {
	return func(o *options) error {
		if load < 0 {
			return fmt.Errorf("invalid load factor: %f", load)
		}

		o.load = load
		return nil
	}
}
