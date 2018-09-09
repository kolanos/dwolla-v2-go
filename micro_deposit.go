package dwolla

const (
	// MicroDepositStatusPending is when the micro deposit is pending
	// processing
	MicroDepositStatusPending MicroDepositStatus = "pending"
	// MicroDepositStatusProcessed is when the micro deposit is processed
	MicroDepositStatusProcessed MicroDepositStatus = "processed"
)

// MicroDeposit is a dwolla micro deposit
type MicroDeposit struct {
	Resource
	Created string              `json:"created"`
	Status  MicroDepositStatus  `json:"status"`
	Failure MicroDepositFailure `json:"failure"`
}

// MicroDepositFailure is detail about a micro deposit failure
type MicroDepositFailure struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// MicroDepositRequest is a micro deposit verification request
type MicroDepositRequest struct {
	Amount1 Amount `json:"amount1,omitempty"`
	Amount2 Amount `json:"amount2,omitempty"`
}

// MicroDepositStatus is the status of the micro deposit
type MicroDepositStatus string
