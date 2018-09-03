package dwolla

// MicroDeposit is a dwolla micro deposit
type MicroDeposit struct {
	Resource
	Created string      `json:"created"`
	Status  string      `json:"status"`
	Failure interface{} `json:"failure"`
}

// MicroDepositRequest is a micro deposit verification request
type MicroDepositRequest struct {
	Amount1 Amount `json:"amount1,omitempty"`
	Amount2 Amount `json:"amount2,omitempty"`
}
