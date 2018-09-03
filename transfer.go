package dwolla

// TransferService is the transfer service interface
// see: https://docsv2.dwolla.com/#transfers
type TransferService interface {
	Create(*TransferRequest) (*Transfer, error)
	Retrieve(string) (*Transfer, error)
}

// TransferServiceOp is an implementation of the transfer service interface
type TransferServiceOp struct {
	client *Client
}

// TransferStatus is a transfer's status
type TransferStatus string

const (
	// TransferStatusCancelled is when the transfer has been cancelled
	TransferStatusCancelled TransferStatus = "cancelled"
	// TransferStatusFailed is when the transfer failed
	TransferStatusFailed TransferStatus = "failed"
	// TransferStatusPending is when the transfer is pending
	TransferStatusPending TransferStatus = "pending"
	// TransferStatusProcessed is when the transfer is processed
	TransferStatusProcessed TransferStatus = "processed"
	// TransferStatusReclaimed is when the transfer is reclaimed
	TransferStatusReclaimed TransferStatus = "reclaimed"
)

// Transfer is a dwolla transfer
type Transfer struct {
	Resource
	ID              string         `json:"id"`
	Status          TransferStatus `json:"status"`
	Amount          Amount         `json:"amount"`
	Created         string         `json:"created"`
	Clearing        interface{}    `json:"clearing"`
	IndividualACHID string         `json:"individualAchId"`
}

// Transfers is a collection of dwolla transfers
type Transfers struct {
	Collection
	Embedded map[string][]Transfer `json:"_embedded"`
}

// TransferRequest is a transfer request
type TransferRequest struct{}

// Create initiates a transfer
// see: https://docsv2.dwolla.com/#initiate-a-transfer
func (t *TransferServiceOp) Create(body *TransferRequest) (*Transfer, error) {
	var transfer Transfer
	return &transfer, nil
}

// Retrieve returns the transfer matching the id
// see: https://docsv2.dwolla.com/#retrieve-a-transfer
func (t *TransferServiceOp) Retrieve(id string) (*Transfer, error) {
	var transfer Transfer
	return &transfer, nil
}
