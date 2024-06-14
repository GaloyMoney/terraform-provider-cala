resource "cala_account" "name" {
  name = "name"
  code = "slghsg"
  id   = "1ce41f76-1e86-4879-b326-0c4a8501ded3"
}

provider "cala" {
  endpoint = "http://localhost:2252"
}

terraform {
  required_providers {
    cala = {
      source  = "galoymoney/cala"
      version = "0.0.1"
    }
  }
}
