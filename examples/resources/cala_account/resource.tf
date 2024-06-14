resource "random_uuid" "alice_account_id" { }

resource "cala_account" "alice" {
  id   = random_uuid.alice_account_id.result
  name = "Alice Account"
  code = "USER.ACCOUNTS.${random_uuid.alice_account_id.result}"
}
