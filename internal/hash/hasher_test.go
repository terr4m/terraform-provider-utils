package hash

import (
	"testing"

	"github.com/cespare/xxhash/v2"
)

func TestHasher(t *testing.T) {
	t.Parallel()

	t.Run("Sum64", func(t *testing.T) {
		t.Parallel()

		for _, tt := range []struct {
			name     string
			data     []byte
			expected uint64
		}{
			{
				name:     "empty",
				data:     []byte{},
				expected: xxhash.Sum64([]byte{}),
			},
			{
				name:     "non_empty",
				data:     []byte("data"),
				expected: xxhash.Sum64([]byte("data")),
			},
		} {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				h := &hasher{}
				if actual := h.Sum64(tt.data); actual != tt.expected {
					t.Errorf("expected %d, got %d", tt.expected, actual)
				}
			})
		}
	})
}
