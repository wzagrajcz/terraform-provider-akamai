provider "akamai" {
  edgerc = "~/.edgerc"
}

data "akamai_cloudlets_policy" "test" {
  policy_id = 1234
  version   = 3
}