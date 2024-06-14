resource "cala_account" "name" {
  name       = "name"
  code       = "slghsg"
  account_id = "1ce41f76-1e86-4879-b326-0c4a8501ded3"
}

terraform {
  required_providers {
    cala = {
      source  = "galoymoney/cala"
      version = "0.0.1"
    }
  }
}
