package dwolla

import "time"

const (
	// MicroDepositStatusPending is when micro-deposits initiated and are en route to their destination
	MicroDepositStatusPending MicroDepositStatus = "pending"
	// MicroDepositStatusProcessed is when micro-deposits have reached the destination account and are awaiting verification
	MicroDepositStatusProcessed MicroDepositStatus = "processed"
	// MicroDepositStatusFailed is when micro-deposits failed to clear successfully to the destination
	MicroDepositStatusFailed MicroDepositStatus = "failed"
)

// MicroDeposit is a dwolla micro deposit
type MicroDeposit struct {
	Resource
	Created string              `json:"created"`
	Status  MicroDepositStatus  `json:"status"`
	Failure MicroDepositFailure `json:"failure"` // Determines if micro-deposits fail to complete to a bank.
}

// MicroDepositFailure is detail about a micro deposit failure
type MicroDepositFailure struct {
	Code        string `json:"code"`        // ACH return code
	Description string `json:"description"` // description of the return
}

// MicroDepositRequest is a micro deposit verification request
type MicroDepositRequest struct {
	Amount1        Amount `json:"amount1,omitempty"`
	Amount2        Amount `json:"amount2,omitempty"`
	IdempotencyKey string `json:"-"`
}

// MicroDepositStatus is the status of the micro deposit
type MicroDepositStatus string

// CreatedTime returns the created value as time.Time
func (m *MicroDeposit) CreatedTime() time.Time {
	t, _ := time.Parse(time.RFC3339, m.Created)
	return t
}
