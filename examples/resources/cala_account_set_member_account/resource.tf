resource "random_uuid" "journal_id" {}

resource "cala_journal" "journal" {
  id   = random_uuid.journal_id.result
  name = "Default"
}

resource "random_uuid" "account_set_id" {}

resource "cala_account_set" "set" {
  id         = random_uuid.account_set_id.result
  name       = "Assets"
  journal_id = cala_journal.journal.id
}

resource "random_uuid" "bob_account_id" {}

resource "cala_account" "bob" {
  id   = random_uuid.bob_account_id.result
  name = "Bob Account"
  code = "USER.ACCOUNTS.${random_uuid.bob_account_id.result}"
}

resource "cala_account_set_member_account" "bob" {
  account_set_id    = cala_account_set.set.id
  member_account_id = cala_account.bob.id
}
