terraform {
  required_providers {
    cala = {
      source = "registry.terraform.io/galoymoney/cala"
    }
  }
}

provider "cala" {
  endpoint = "http://localhost:2252/graphql"
}
