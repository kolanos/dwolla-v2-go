package dwolla

// MassPayment is a dwolla mass payment
type MassPayment struct {
	Resource
	ID            string            `json:"id"`
	Status        string            `json:"status"`
	Created       string            `json:"created"`
	MetaData      map[string]string `json:"metadata"`
	Total         Amount            `json:"total"`
	TotalFees     Amount            `json:"totalFees"`
	CorrelationID string            `json:"correlationId"`
}

// MassPayments is a collection of mass payments
type MassPayments struct {
	Collection
	Embedded map[string][]MassPayment `json:"_embedded"`
}

// MassPaymentItem is a dwolla mass payment item
type MassPaymentItem struct {
	Resource
	Amount        Amount      `json:"amount"`
	MetaData      interface{} `json:"metadata"`
	CorrelationID string      `json:"correlationId"`
}
