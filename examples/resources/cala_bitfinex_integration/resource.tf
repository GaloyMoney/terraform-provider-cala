variable "bitfinex_key" {
  sensitive = true
  default = "dummy"
}

variable "bitfinex_secret" {
  sensitive = true
  default = "dummy"
}

resource "random_uuid" "journal_id" {}

resource "cala_journal" "journal" {
  id   = random_uuid.journal_id.result
  name = "Default"
}

resource "random_uuid" "integration_id" {}

resource "cala_bitfinex_integration" "bfx" {
  id         = random_uuid.integration_id.result
  name       = "Main account"
  journal_id = cala_journal.journal.id
  key        = var.bitfinex_key
  secret     = var.bitfinex_secret
}
