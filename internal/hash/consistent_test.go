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

		for _, d := range []struct {
			name              string
			members           []string
			partitionCount    int
			replicationFactor int
			load              float64
			keys              []string
			errMsg            string
		}{
			{
				name:              "empty",
				members:           []string{},
				partitionCount:    DefaultPartitionCount,
				replicationFactor: DefaultReplicationFactor,
				load:              DefaultLoad,
				keys:              []string{},
				errMsg:            "invalid consistent hash configuration: not enough room to distribute partitions",
			},
			{
				name:              "single_member_single_key",
				members:           []string{"member1"},
				partitionCount:    DefaultPartitionCount,
				replicationFactor: DefaultReplicationFactor,
				load:              DefaultLoad,
				keys:              []string{"key1"},
				errMsg:            "",
			},
			{
				name:              "single_member_multiple_keys",
				members:           []string{"member1"},
				partitionCount:    DefaultPartitionCount,
				replicationFactor: DefaultReplicationFactor,
				load:              DefaultLoad,
				keys:              []string{"key1", "key2", "key3"},
				errMsg:            "",
			},
			{
				name:              "multiple_members_single_key",
				members:           []string{"member1", "member2"},
				partitionCount:    DefaultPartitionCount,
				replicationFactor: DefaultReplicationFactor,
				load:              DefaultLoad,
				keys:              []string{"key1"},
				errMsg:            "",
			},
			{
				name:              "multiple_members_multiple_keys",
				members:           []string{"member1", "member2"},
				partitionCount:    DefaultPartitionCount,
				replicationFactor: DefaultReplicationFactor,
				load:              DefaultLoad,
				keys:              []string{"key1", "key2", "key3"},
				errMsg:            "",
			},
		} {
			t.Run(d.name, func(t *testing.T) {
				t.Parallel()

				ch := NewConsistentHash(d.members, d.partitionCount, d.replicationFactor, d.load)
				actual, err := ch.CalculateMapping(d.keys)

				var errMsg string
				if err != nil {
					errMsg = err.Error()
				}

				if errMsg != d.errMsg {
					t.Errorf("expected error message %s, got %s", d.errMsg, errMsg)
				}

				actualMembers := make([]string, 0, len(actual))
				for k := range actual {
					actualMembers = append(actualMembers, k)
				}

				for _, v := range actualMembers {
					if !slices.Contains(d.members, v) {
						t.Errorf("expected member %s to have been configured", v)
					}
				}

				actualKeys := make([]string, 0, len(d.keys))
				for _, v := range actual {
					actualKeys = append(actualKeys, v...)
				}

				slices.Sort(actualKeys)
				slices.Sort(d.keys)
				if !reflect.DeepEqual(actualKeys, d.keys) {
					t.Errorf("expected keys to be %v, got %v", d.keys, actualKeys)
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
