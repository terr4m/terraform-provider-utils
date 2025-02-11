package hash

import "github.com/cespare/xxhash/v2"

// hasher supports creating hashes.
type hasher struct{}

// Sum64 returns the hash of the given data using the hasher.
func (h *hasher) Sum64(data []byte) uint64 {
	return xxhash.Sum64(data)
}
