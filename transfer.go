package dwolla

import (
	"errors"
	"fmt"
)

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
	ID              string           `json:"id"`
	Status          TransferStatus   `json:"status"`
	Amount          Amount           `json:"amount"`
	Created         string           `json:"created"`
	MetaData        MetaData         `json:"metadata"`
	Clearing        TransferClearing `json:"clearing"`
	CorrelationID   string           `json:"correlationId"`
	IndividualACHID string           `json:"individualAchId"`
}

// Transfers is a collection of dwolla transfers
type Transfers struct {
	Collection
	Embedded map[string][]Transfer `json:"_embedded"`
}

// TransferACHDetails contains data sent to the bank account
type TransferACHDetails struct {
	Destination TransferAddenda `json:"destination,omitempty"`
	Source      TransferAddenda `json:"source,omitempty"`
}

// TransferAddenda is a transfer addenda
type TransferAddenda struct {
	Addenda TransferAddendaValues `json:"addenda,omitempty"`
}

// TransferAddendaValues is the addenda values
type TransferAddendaValues struct {
	Values []string `json:"values,omitempty"`
}

// TransferClearing is a transfer clearing schedule
type TransferClearing struct {
	Destination string `json:"destination,omitempty"`
	Source      string `json:"source,omitempty"`
}

// TransferFailureReason contains details about a failed transfer
type TransferFailureReason struct {
	Resource
	Code        string `json:"code"`
	Description string `json:"description"`
}

// TransferFee is a transfer fee
type TransferFee struct {
	Resource
	Amount Amount `json:"amount"`
}

// TransferFees contains fees related to a transfer
type TransferFees struct {
	Transactions []Transfer `json:"transactions"`
	Total        int        `json:"total"`
}

// TransferRequest is a transfer request
type TransferRequest struct {
	Resource
	Status        TransferStatus     `json:"status,omitempty"`
	Amount        Amount             `json:"amount,omitempty"`
	MetaData      MetaData           `json:"metadata,omitempty"`
	Fees          []TransferFee      `json:"fees,omitempty"`
	Clearing      TransferClearing   `json:"clearing,omitempty"`
	CorrelationID string             `json:"correlationId"`
	ACHDetails    TransferACHDetails `json:"achDetails,omitempty"`
}

// Create initiates a transfer
// see: https://docsv2.dwolla.com/#initiate-a-transfer
func (t *TransferServiceOp) Create(body *TransferRequest) (*Transfer, error) {
	var transfer Transfer

	if err := t.client.Post("transfers", body, nil, &transfer); err != nil {
		return nil, err
	}

	transfer.client = t.client

	return &transfer, nil
}

// Retrieve returns the transfer matching the id
// see: https://docsv2.dwolla.com/#retrieve-a-transfer
func (t *TransferServiceOp) Retrieve(id string) (*Transfer, error) {
	var transfer Transfer

	if err := t.client.Get(fmt.Sprintf("transfers/%s", id), nil, nil, &transfer); err != nil {
		return nil, err
	}

	transfer.client = t.client

	return &transfer, nil
}

// Cancel cancels the transfer
// see: https://docsv2.dwolla.com/#cancel-a-transfer
func (t *Transfer) Cancel() error {
	if _, ok := t.Links["cancel"]; !ok {
		return errors.New("No cancel resource link")
	}

	body := &TransferRequest{Status: TransferStatusCancelled}

	return t.client.Post(t.Links["cancel"].Href, body, nil, t)
}

// Destination returns the customer transfer destination
func (t *Transfer) Destination() (*Customer, error) {
	if _, ok := t.Links["destination"]; !ok {
		return nil, errors.New("No destination resource link")
	}

	return t.client.Customer.Retrieve(t.Links["destination"].Href)
}

// DestinationFundingSource returns the transfer funding destination
func (t *Transfer) DestinationFundingSource() (*FundingSource, error) {
	if _, ok := t.Links["destination-funding-source"]; !ok {
		return nil, errors.New("No destination funding source resource link")
	}

	return t.client.FundingSource.Retrieve(t.Links["destination-funding-source"].Href)
}

// ListFees returns the fees associated with the transfer
// see: https://docsv2.dwolla.com/#list-fees-for-a-transfer
func (t *Transfer) ListFees() (*TransferFees, error) {
	var fees TransferFees

	if _, ok := t.Links["fees"]; !ok {
		return nil, errors.New("No fees resource link")
	}

	if err := t.client.Get(t.Links["fees"].Href, nil, nil, &fees); err != nil {
		return nil, err
	}

	for i := range fees.Transactions {
		fees.Transactions[i].client = t.client
	}

	return &fees, nil
}

// Source returns the customer transfer source
func (t *Transfer) Source() (*Customer, error) {
	if _, ok := t.Links["source"]; !ok {
		return nil, errors.New("No source resource link")
	}

	return t.client.Customer.Retrieve(t.Links["source"].Href)
}

// SourceFundingSource returns the transfer funding source
func (t *Transfer) SourceFundingSource() (*FundingSource, error) {
	if _, ok := t.Links["source-funding-source"]; !ok {
		return nil, errors.New("No source funding source resource link")
	}

	return t.client.FundingSource.Retrieve(t.Links["source-funding-souce"].Href)
}

// RetrieveFailureReason returns the transfer's failure reason
// see: https://docsv2.dwolla.com/#retrieve-a-transfer-failure-reason
func (t *Transfer) RetrieveFailureReason() (*TransferFailureReason, error) {
	var reason TransferFailureReason

	if _, ok := t.Links["failure"]; !ok {
		return nil, errors.New("No failure resource link")
	}

	if err := t.client.Get(t.Links["failure"].Href, nil, nil, &reason); err != nil {
		return nil, err
	}

	return &reason, nil
}
