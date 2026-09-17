locals {
  targets = ["target1", "target2", "target3"]
  payload = ["1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"]
}

data "utils_consistent_hash" "example" {
  members = local.targets
  keys    = local.payload
}

resource "terraform_data" "example" {
  for_each = toset(local.targets)

  input = data.utils_consistent_hash.example.mapping.members[each.key]
}
