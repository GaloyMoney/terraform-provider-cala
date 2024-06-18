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

module "account_set_member_account_set" {
  source = "./resources/cala_account_set_member_account_set"
}

module "balance_sheet" {
  source = "./resources/cala_balance_sheet"
}

module "bitfinex_integration" {
  source = "./resources/cala_bitfinex_integration"
}

terraform {
  required_providers {
    cala = {
      source  = "registry.terraform.io/galoymoney/cala"
      version = "0.0.14"
    }
  }
}
