provider "cala" {
  endpoint = "http://localhost:2252/graphql"
}

module "account" {
  source = "./resources/cala_account"
}

terraform {
  required_providers {
    cala = {
      source  = "registry.terraform.io/galoymoney/cala"
      version = "0.0.10"
    }
  }
}
