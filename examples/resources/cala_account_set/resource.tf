resource "random_uuid" "account_set_id" { }

resource "cala_account_set" "set" {
  id   = random_uuid.account_set_id.result
  name = "Assets"
}
