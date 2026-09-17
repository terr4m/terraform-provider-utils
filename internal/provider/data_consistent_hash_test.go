package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccConsistentHashDataSource(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "utils_consistent_hash" "test" {
  members = ["member1", "member2", "member3"]
  keys    = ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"]
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.utils_consistent_hash.test", tfjsonpath.New("mapping"), knownvalue.ObjectExact(map[string]knownvalue.Check{
							"keys":    knownvalue.MapSizeExact(12),
							"members": knownvalue.NotNull(),
						})),
					},
				},
			},
		})
	})

	t.Run("add_keys", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "utils_consistent_hash" "test" {
  members = ["member1", "member2", "member3"]
  keys    = ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"]
}
`,

					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.utils_consistent_hash.test", tfjsonpath.New("mapping"), knownvalue.ObjectExact(map[string]knownvalue.Check{
							"keys":    knownvalue.MapSizeExact(12),
							"members": knownvalue.NotNull(),
						})),
					},
				},
				{
					Config: `
data "utils_consistent_hash" "test" {
  members = ["member1", "member2", "member3"]
  keys    = ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13"]
}
`,

					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.utils_consistent_hash.test", tfjsonpath.New("mapping"), knownvalue.ObjectExact(map[string]knownvalue.Check{
							"keys":    knownvalue.MapSizeExact(13),
							"members": knownvalue.NotNull(),
						})),
					},
				},
			},
		})
	})

	t.Run("add_members", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `
data "utils_consistent_hash" "test" {
  members = ["member1", "member2", "member3"]
  keys    = ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"]
}
`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.utils_consistent_hash.test", tfjsonpath.New("mapping"), knownvalue.ObjectExact(map[string]knownvalue.Check{
							"keys":    knownvalue.MapSizeExact(12),
							"members": knownvalue.NotNull(),
						})),
					},
				},
				{
					Config: `
data "utils_consistent_hash" "test" {
  members = ["member1", "member2", "member3", "members4"]
  keys    = ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"]
}
`,

					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.utils_consistent_hash.test", tfjsonpath.New("mapping"), knownvalue.ObjectExact(map[string]knownvalue.Check{
							"keys":    knownvalue.MapSizeExact(12),
							"members": knownvalue.NotNull(),
						})),
					},
				},
			},
		})
	})
}
