package provider

import "fmt"

// toDebitOrCredit converts a string to the DebitOrCredit enum type.
func toDebitOrCredit(value string) (DebitOrCredit, error) {
	switch value {
	case "DEBIT":
		return DebitOrCreditDebit, nil
	case "CREDIT":
		return DebitOrCreditCredit, nil
	default:
		return DebitOrCreditCredit, nil
	}
}

// toStatus converts a string to the Status enum type.
func toStatus(value string) (Status, error) {
	switch value {
	case "ACTIVE":
		return StatusActive, nil
	case "INACTIVE":
		return StatusLocked, nil
	default:
		return Status(""), fmt.Errorf("invalid value for Status: %s", value)
	}
}
