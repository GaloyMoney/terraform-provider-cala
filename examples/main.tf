provider "cala" {
  endpoint = "http://localhost:2252/graphql"
}

module "account" {
  source = "./resources/cala_account"
}

module "journal" {
  source = "./resources/cala_journal"
}

module "account_set" {
  source = "./resources/cala_account_set"
}

module "account_set_member_account" {
  source = "./resources/cala_account_set_member_account"
}

terraform {
  required_providers {
    cala = {
      source  = "registry.terraform.io/galoymoney/cala"
      version = "0.0.12"
    }
  }
}
