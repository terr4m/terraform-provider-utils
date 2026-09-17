package hash

import (
	"regexp"
	"testing"

	"github.com/buraksezer/consistent"
)

func Test_newConsistent(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name    string
		members []consistent.Member
		cfg     consistent.Config
		wantErr *string
	}{
		{
			name:    "valid_configuration",
			cfg:     consistent.Config{Hasher: &hasher{}, PartitionCount: 10, ReplicationFactor: 2, Load: 1.05},
			wantErr: nil,
		},
		{
			name:    "valid_configuration_with_members",
			members: []consistent.Member{member("member1"), member("member2")},
			cfg:     consistent.Config{Hasher: &hasher{}, PartitionCount: 10, ReplicationFactor: 2, Load: 1.05},
			wantErr: nil,
		},
		{
			name:    "invalid_configuration",
			wantErr: new("invalid consistent hash configuration"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := newConsistent(tt.members, tt.cfg)
			if err != nil {
				if tt.wantErr == nil || !regexp.MustCompile(regexp.QuoteMeta(*tt.wantErr)).MatchString(err.Error()) {
					t.Errorf("newConsistent() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if tt.wantErr != nil {
				t.Errorf("newConsistent() expected error %v, got nil", *tt.wantErr)
			}

			if got == nil {
				t.Errorf("newConsistent() got = nil, want non-nil")
			}
		})
	}
}
