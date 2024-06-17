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

resource "random_uuid" "member_account_set_id" {}

resource "cala_account_set" "member_set" {
  id         = random_uuid.member_account_set_id.result
  name       = "Cash"
  journal_id = cala_journal.journal.id
}

resource "cala_account_set_member_account_set" "member_account_set" {
  account_set_id        = cala_account_set.set.id
  member_account_set_id = cala_account_set.member_set.id
}
