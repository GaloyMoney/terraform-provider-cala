resource "random_uuid" "journal_id" {}

resource "cala_journal" "journal" {
  id   = random_uuid.journal_id.result
  name = "Default"
}
