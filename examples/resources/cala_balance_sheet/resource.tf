resource "random_uuid" "journal_id" {}

resource "cala_journal" "journal" {
  id   = random_uuid.journal_id.result
  name = "Default"
}

resource "cala_balance_sheet" "balance_sheet" {
  journal_id = cala_journal.journal.id
}
