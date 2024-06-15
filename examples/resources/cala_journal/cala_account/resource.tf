resource "random_uuid" "journal_id" { }

resource "cala_journal" "default" {
  id   = random_uuid.journal_id.result
  name = "Default"
}
