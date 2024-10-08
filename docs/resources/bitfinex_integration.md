---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cala_bitfinex_integration Resource - terraform-provider-cala"
subcategory: ""
description: |-
  Cala Bitfinex Integration.
---

# cala_bitfinex_integration (Resource)

Cala Bitfinex Integration.

## Example Usage

```terraform
variable "bitfinex_key" {
  sensitive = true
  default   = "dummy"
}

variable "bitfinex_secret" {
  sensitive = true
  default   = "dummy"
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) ID of the integration.
- `journal_id` (String) journal_id
- `key` (String, Sensitive) The bitfinex API key
- `name` (String) Name of the integration.
- `secret` (String, Sensitive) The bitfinex API secret

### Optional

- `description` (String) Description of the integration.

### Read-Only

- `omnibus_account_id` (String) The Account id for the omnibus Account
