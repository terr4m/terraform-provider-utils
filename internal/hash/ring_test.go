package hash

import (
	"regexp"
	"testing"

	"github.com/buraksezer/consistent"
)

func TestNewRing(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name    string
		fn      []OptionsFunc
		wantErr *string
	}{
		{
			name: "empty",
		},
		{
			name: "with_members",
			fn: []OptionsFunc{func(o *options) error {
				o.members = []consistent.Member{member("member1"), member("member2")}
				return nil
			}},
		},
		{
			name: "invalid",
			fn: []OptionsFunc{func(o *options) error {
				o.members = []consistent.Member{member("member1"), member("member2")}
				o.partitionCount = 10
				o.replicationFactor = 2
				o.load = 0.5
				return nil
			}},
			wantErr: new("invalid consistent hash configuration"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ring, err := NewRing(tt.fn...)
			if err != nil {
				if tt.wantErr == nil || !regexp.MustCompile(regexp.QuoteMeta(*tt.wantErr)).MatchString(err.Error()) {
					t.Errorf("NewRing() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if tt.wantErr != nil {
				t.Errorf("NewRing() expected error %v, got nil", *tt.wantErr)
			}

			if ring == nil {
				t.Errorf("NewRing() got = nil, want non-nil")
			}
		})
	}
}
