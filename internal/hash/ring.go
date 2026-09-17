package hash

import (
	"github.com/buraksezer/consistent"
)

// Ring is a consistent hash ring.
type Ring struct {
	c *consistent.Consistent
}

// NewRing creates a new consistent hash ring with the given options.
func NewRing(f ...OptionsFunc) (*Ring, error) {
	opts := options{
		partitionCount:    defaultPartitionCount,
		replicationFactor: defaultReplicationFactor,
		load:              defaultLoad,
	}

	for _, fn := range f {
		if err := fn(&opts); err != nil {
			return nil, err
		}
	}

	cfg := consistent.Config{
		Hasher:            &hasher{},
		PartitionCount:    opts.partitionCount,
		ReplicationFactor: opts.replicationFactor,
		Load:              opts.load,
	}

	c, err := newConsistent(opts.members, cfg)
	if err != nil {
		return nil, err
	}

	return &Ring{c: c}, nil
}

// Members returns the list of members in the consistent hash ring.
func (r *Ring) Members() []string {
	mem := r.c.GetMembers()

	s := make([]string, len(mem))
	for i, m := range mem {
		s[i] = m.String()
	}

	return s
}

// Calculate calculates the mapping of keys to members in the consistent hash ring.
func (r *Ring) Calculate(keys []string) (map[string]string, error) {
	res := make(map[string]string, len(keys))

	for _, k := range keys {
		m := r.c.LocateKey([]byte(k))
		res[k] = m.String()
	}

	return res, nil
}
