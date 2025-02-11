package hash

import (
	"fmt"

	"github.com/buraksezer/consistent"
)

const (
	DefaultPartitionCount    int     = 271
	DefaultReplicationFactor int     = 20
	DefaultLoad              float64 = 1.25
)

// ConsistentHash is a consistent hash implementation.
type ConsistentHash struct {
	Members           []string
	PartitionCount    int
	ReplicationFactor int
	Load              float64
}

// NewConsistentHash creates a new consistent hash implementation.
func NewConsistentHash(members []string, partitionCount, replicationFactor int, load float64) *ConsistentHash {
	return &ConsistentHash{
		Members:           members,
		PartitionCount:    partitionCount,
		ReplicationFactor: replicationFactor,
		Load:              load,
	}
}

// CalculateMapping calculates the mapping of keys to members.
func (ch *ConsistentHash) CalculateMapping(keys []string) (map[string][]string, error) {
	members := make([]consistent.Member, len(ch.Members))
	for i, m := range ch.Members {
		members[i] = member(fmt.Sprintf("%s-", m))
	}

	c, err := newConsistent(members, consistent.Config{PartitionCount: ch.PartitionCount, ReplicationFactor: ch.ReplicationFactor, Load: ch.Load, Hasher: &hasher{}})
	if err != nil {
		return nil, err
	}

	mapping := make(map[string][]string, len(ch.Members))
	for _, k := range keys {
		member := c.LocateKey([]byte(k))
		mm := member.String()
		m := mm[:len(mm)-1]

		mapping[m] = append(mapping[m], k)
	}

	return mapping, nil
}

// newConsistent creates a new consistent hash implementation and catches package panics.
func newConsistent(members []consistent.Member, config consistent.Config) (c *consistent.Consistent, err error) {
	// Required as package panics on error
	defer func() {
		if errp := recover(); errp != nil {
			err = fmt.Errorf("invalid consistent hash configuration: %v", errp)
		}
	}()

	c = consistent.New(members, config)
	return
}
