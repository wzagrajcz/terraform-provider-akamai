provider "akamai" {
  edgerc = "../../test/edgerc"
}

resource "akamai_edgeworker" "edgeworker" {
  name             = "example"
  group_id         = "12345"
  resource_tier_id = 54321
}