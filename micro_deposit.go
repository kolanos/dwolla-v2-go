package dwolla

// MicroDeposit is a dwolla micro deposit
type MicroDeposit struct {
	Resource
	Created string      `json:"created"`
	Status  string      `json:"status"`
	Failure interface{} `json:"failure"`
}
