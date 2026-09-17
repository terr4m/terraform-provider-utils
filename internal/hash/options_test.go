package hash

import (
	"regexp"
	"testing"
)

func TestWithMembers(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name    string
		members []string
	}{
		{
			name:    "has_members",
			members: []string{"member1", "member2"},
		},
		{
			name:    "empty_members",
			members: []string{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opts := &options{}
			err := WithMembers(tt.members)(opts)
			if err != nil {
				t.Errorf("WithMembers() error = %v", err)
				return
			}

			if len(opts.members) != len(tt.members) {
				t.Errorf("WithMembers() got %d members, want %d", len(opts.members), len(tt.members))
			}
		})
	}
}

func TestWithPartitionCount(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name           string
		partitionCount int
		wantErr        *string
	}{
		{
			name:           "empty_partition_count",
			partitionCount: 0,
		},
		{
			name:           "valid_partition_count",
			partitionCount: 10,
		},
		{
			name:           "invalid_partition_count",
			partitionCount: -1,
			wantErr:        new("invalid partition count"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opts := &options{}
			err := WithPartitionCount(tt.partitionCount)(opts)
			if err != nil {
				if tt.wantErr == nil || !regexp.MustCompile(regexp.QuoteMeta(*tt.wantErr)).MatchString(err.Error()) {
					t.Errorf("WithPartitionCount() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if tt.wantErr != nil {
				t.Errorf("WithPartitionCount() expected error %v, got nil", *tt.wantErr)
			}

			if opts.partitionCount != tt.partitionCount {
				t.Errorf("WithPartitionCount() got %d, want %d", opts.partitionCount, tt.partitionCount)
			}
		})
	}
}

func TestWithReplicationFactor(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name              string
		replicationFactor int
		wantErr           *string
	}{
		{
			name:              "empty_replication_factor",
			replicationFactor: 0,
		},
		{
			name:              "valid_replication_factor",
			replicationFactor: 10,
		},
		{
			name:              "invalid_replication_factor",
			replicationFactor: -1,
			wantErr:           new("invalid replication factor"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opts := &options{}
			err := WithReplicationFactor(tt.replicationFactor)(opts)
			if err != nil {
				if tt.wantErr == nil || !regexp.MustCompile(regexp.QuoteMeta(*tt.wantErr)).MatchString(err.Error()) {
					t.Errorf("WithReplicationFactor() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if tt.wantErr != nil {
				t.Errorf("WithReplicationFactor() expected error %v, got nil", *tt.wantErr)
			}

			if opts.replicationFactor != tt.replicationFactor {
				t.Errorf("WithReplicationFactor() got %d, want %d", opts.replicationFactor, tt.replicationFactor)
			}
		})
	}
}

func TestWithLoad(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name    string
		load    float64
		wantErr *string
	}{
		{
			name:    "empty_load",
			load:    0,
			wantErr: nil,
		},
		{
			name:    "valid_load",
			load:    1.05,
			wantErr: nil,
		},
		{
			name:    "invalid_load",
			load:    -1.0,
			wantErr: new("invalid load factor"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			opts := &options{}
			err := WithLoad(tt.load)(opts)
			if err != nil {
				if tt.wantErr == nil || !regexp.MustCompile(regexp.QuoteMeta(*tt.wantErr)).MatchString(err.Error()) {
					t.Errorf("WithLoad() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if tt.wantErr != nil {
				t.Errorf("WithLoad() expected error %v, got nil", *tt.wantErr)
			}

			if opts.load != tt.load {
				t.Errorf("WithLoad() got %f, want %f", opts.load, tt.load)
			}
		})
	}
}
