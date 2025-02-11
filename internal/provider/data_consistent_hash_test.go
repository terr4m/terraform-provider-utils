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
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `data "utils_consistent_hash" "test" {
  members = ["member1", "member2", "member3"]
  keys    = ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"]
}

output "keys" {
  value = flatten([for k, v in data.utils_consistent_hash.test.mapping : v])
}`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.utils_consistent_hash.test", tfjsonpath.New("mapping"), knownvalue.MapSizeExact(3)),
						statecheck.ExpectKnownOutputValue("keys", knownvalue.ListSizeExact(12)),
					},
				},
			},
		})
	})

	t.Run("add_keys", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `data "utils_consistent_hash" "test" {
  members = ["member1", "member2", "member3"]
  keys    = ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"]
}

output "keys" {
  value = flatten([for k, v in data.utils_consistent_hash.test.mapping : v])
}`,

					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.utils_consistent_hash.test", tfjsonpath.New("mapping"), knownvalue.MapSizeExact(3)),
						statecheck.ExpectKnownOutputValue("keys", knownvalue.ListSizeExact(12)),
					},
				},
				{
					Config: `data "utils_consistent_hash" "test" {
  members = ["member1", "member2", "member3"]
  keys    = ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13"]
}

output "keys" {
  value = flatten([for k, v in data.utils_consistent_hash.test.mapping : v])
}`,

					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.utils_consistent_hash.test", tfjsonpath.New("mapping"), knownvalue.MapSizeExact(3)),
						statecheck.ExpectKnownOutputValue("keys", knownvalue.ListSizeExact(13)),
					},
				},
			},
		})
	})

	t.Run("add_members", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: `data "utils_consistent_hash" "test" {
  members = ["member1", "member2", "member3"]
  keys    = ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"]
}

output "keys" {
  value = flatten([for k, v in data.utils_consistent_hash.test.mapping : v])
}`,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.utils_consistent_hash.test", tfjsonpath.New("mapping"), knownvalue.MapSizeExact(3)),
						statecheck.ExpectKnownOutputValue("keys", knownvalue.ListSizeExact(12)),
					},
				},
				{
					Config: `data "utils_consistent_hash" "test" {
  members = ["member1", "member2", "member3", "members4"]
  keys    = ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"]
}

output "keys" {
  value = flatten([for k, v in data.utils_consistent_hash.test.mapping : v])
}`,

					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue("data.utils_consistent_hash.test", tfjsonpath.New("mapping"), knownvalue.MapSizeExact(4)),
						statecheck.ExpectKnownOutputValue("keys", knownvalue.ListSizeExact(12)),
					},
				},
			},
		})
	})
}
