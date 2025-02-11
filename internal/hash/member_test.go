package hash

import (
	"testing"
)

func TestMember(t *testing.T) {
	t.Parallel()

	t.Run("String", func(t *testing.T) {
		t.Parallel()

		for _, tt := range []struct {
			name     string
			member   member
			expected string
		}{
			{
				name:     "empty",
				member:   "",
				expected: "",
			},
			{
				name:     "non_empty",
				member:   "member1",
				expected: "member1",
			},
		} {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				if actual := tt.member.String(); actual != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, actual)
				}
			})
		}
	})
}
