package provider

import "fmt"

// toDebitOrCredit converts a string to the DebitOrCredit enum type.
func toDebitOrCredit(value string) (DebitOrCredit, error) {
	switch value {
	case "DEBIT":
		return DebitOrCreditDebit, nil
	case "debit":
		return DebitOrCreditDebit, nil
	case "CREDIT":
		return DebitOrCreditCredit, nil
	case "credit":
		return DebitOrCreditCredit, nil
	default:
		return DebitOrCredit(""), fmt.Errorf("invalid value for DebitOrCredit: %s", value)
	}
}

// toStatus converts a string to the Status enum type.
func toStatus(value string) (Status, error) {
	switch value {
	case "ACTIVE":
		return StatusActive, nil
	case "active":
		return StatusActive, nil
	case "INACTIVE":
		return StatusLocked, nil
	case "inactive":
		return StatusLocked, nil
	default:
		return Status(""), fmt.Errorf("invalid value for Status: %s", value)
	}
}
