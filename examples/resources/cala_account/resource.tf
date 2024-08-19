resource "random_uuid" "alice_account_id" {}

resource "cala_account" "alice" {
  id   = random_uuid.alice_account_id.result
  name = "Alice Account"
  code = "USER.ACCOUNTS.${random_uuid.alice_account_id.result}"
}

resource "random_uuid" "bank" {}
resource "cala_account" "bank" {
  id                  = random_uuid.bank.result
  name                = "Bank cash"
  code                = "BANK.DEPOSITS.${random_uuid.bank.result}"
  normal_balance_type = "DEBIT"
}
