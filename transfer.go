package dwolla

import (
	"net/url"
)

// TransferService is the transfer service interface
// see: https://docsv2.dwolla.com/#transfers
type TransferService interface {
	Get(string) (*Transfer, error)
	List(*url.Values) (*Transfers, error)
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
