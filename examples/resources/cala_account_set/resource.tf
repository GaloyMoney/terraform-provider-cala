resource "random_uuid" "journal_id" { }

resource "cala_journal" "journal" {
  id   = random_uuid.journal_id.result
  name = "Default"
}

resource "random_uuid" "account_set_id" { }

resource "cala_account_set" "set" {
  id   = random_uuid.account_set_id.result
  name = "Assets"
  journal_id = cala_journal.journal.id
}
