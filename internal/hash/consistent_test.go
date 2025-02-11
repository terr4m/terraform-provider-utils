package hash

import (
	"reflect"
	"slices"
	"testing"

	"github.com/buraksezer/consistent"
)

func TestConsistentHash(t *testing.T) {
	t.Parallel()

	t.Run("CalculateMapping", func(t *testing.T) {
		t.Parallel()

		for _, tt := range []struct {
			name              string
			members           []string
			partitionCount    int
			replicationFactor int
			load              float64
			keys              []string
			err               bool
		}{
			{
				name:              "empty",
				members:           []string{},
				partitionCount:    DefaultPartitionCount,
				replicationFactor: DefaultReplicationFactor,
				load:              DefaultLoad,
				keys:              []string{},
				err:               true,
			},
			{
				name:              "single_member_single_key",
				members:           []string{"member1"},
				partitionCount:    DefaultPartitionCount,
				replicationFactor: DefaultReplicationFactor,
				load:              DefaultLoad,
				keys:              []string{"key1"},
				err:               false,
			},
			{
				name:              "single_member_multiple_keys",
				members:           []string{"member1"},
				partitionCount:    DefaultPartitionCount,
				replicationFactor: DefaultReplicationFactor,
				load:              DefaultLoad,
				keys:              []string{"key1", "key2", "key3"},
				err:               false,
			},
			{
				name:              "multiple_members_single_key",
				members:           []string{"member1", "member2"},
				partitionCount:    DefaultPartitionCount,
				replicationFactor: DefaultReplicationFactor,
				load:              DefaultLoad,
				keys:              []string{"key1"},
				err:               false,
			},
			{
				name:              "multiple_members_multiple_keys",
				members:           []string{"member1", "member2"},
				partitionCount:    DefaultPartitionCount,
				replicationFactor: DefaultReplicationFactor,
				load:              DefaultLoad,
				keys:              []string{"key1", "key2", "key3"},
				err:               false,
			},
		} {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				ch := NewConsistentHash(tt.members, tt.partitionCount, tt.replicationFactor, tt.load)
				actual, err := ch.CalculateMapping(tt.keys)

				if !tt.err && err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if tt.err && err == nil {
					t.Errorf("expected an error")
				}

				actualMembers := make([]string, 0, len(actual))
				for k := range actual {
					actualMembers = append(actualMembers, k)
				}

				for _, v := range actualMembers {
					if !slices.Contains(tt.members, v) {
						t.Errorf("expected member %s to have been configured", v)
					}
				}

				actualKeys := make([]string, 0, len(tt.keys))
				for _, v := range actual {
					actualKeys = append(actualKeys, v...)
				}

				slices.Sort(actualKeys)
				slices.Sort(tt.keys)
				if !reflect.DeepEqual(actualKeys, tt.keys) {
					t.Errorf("expected keys to be %v, got %v", tt.keys, actualKeys)
				}
			})
		}
	})
}

func Test_newConsistent(t *testing.T) {
	t.Parallel()

	t.Run("invalid_configuration", func(t *testing.T) {
		t.Parallel()

		_, err := newConsistent([]consistent.Member{member("member1"), member("member2")}, consistent.Config{PartitionCount: 1, ReplicationFactor: 1, Load: 1, Hasher: &hasher{}})
		if err == nil {
			t.Errorf("expected an error")
		}
	})
}
