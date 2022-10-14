package dwolla

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

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

// TransferService is the transfer service interface
//
// see: https://docsv2.dwolla.com/#transfers
type TransferService interface {
	Create(context.Context, *TransferRequest) (*Transfer, error)
	Retrieve(context.Context, string) (*Transfer, error)
}

// TransferServiceOp is an implementation of the transfer service interface
type TransferServiceOp struct {
	client *Client
}

// TransferStatus is a transfer's status
type TransferStatus string

// Transfer is a dwolla transfer
type Transfer struct {
	Resource
	ID              string         `json:"id"`
	Status          TransferStatus `json:"status"`
	Amount          Amount         `json:"amount"`
	Created         string         `json:"created"`
	MetaData        MetaData       `json:"metadata"`
	Clearing        Clearing       `json:"clearing"`
	CorrelationID   string         `json:"correlationId"`
	IndividualACHID string         `json:"individualAchId"`
}

// Transfers is a collection of dwolla transfers
type Transfers struct {
	Collection
	Embedded map[string][]Transfer `json:"_embedded"`
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
	Status         TransferStatus `json:"status,omitempty"`
	Amount         Amount         `json:"amount,omitempty"`
	MetaData       MetaData       `json:"metadata,omitempty"`
	Fees           []TransferFee  `json:"fees,omitempty"`
	Clearing       Clearing       `json:"clearing,omitempty"`
	CorrelationID  string         `json:"correlationId,omitempty"`
	ACHDetails     *ACHDetails    `json:"achDetails,omitempty"`
	IdempotencyKey string         `json:"-"`
}

// Create initiates a transfer
//
// see: https://docsv2.dwolla.com/#initiate-a-transfer
func (t *TransferServiceOp) Create(ctx context.Context, body *TransferRequest) (*Transfer, error) {
	var transfer Transfer

	var headers *http.Header
	if body.IdempotencyKey != "" {
		headers = &http.Header{}
		headers.Set(HeaderIdempotency, body.IdempotencyKey)
	}

	if err := t.client.Post(ctx, "transfers", body, headers, &transfer); err != nil {
		return nil, err
	}

	transfer.client = t.client

	return &transfer, nil
}

// Retrieve returns the transfer matching the id
//
// see: https://docsv2.dwolla.com/#retrieve-a-transfer
func (t *TransferServiceOp) Retrieve(ctx context.Context, id string) (*Transfer, error) {
	var transfer Transfer

	if err := t.client.Get(ctx, fmt.Sprintf("transfers/%s", id), nil, nil, &transfer); err != nil {
		return nil, err
	}

	transfer.client = t.client

	return &transfer, nil
}

// Cancel cancels the transfer
//
// see: https://docsv2.dwolla.com/#cancel-a-transfer
func (t *Transfer) Cancel(ctx context.Context) error {
	if _, ok := t.Links["cancel"]; !ok {
		return errors.New("No cancel resource link")
	}

	body := &TransferRequest{Status: TransferStatusCancelled}

	return t.client.Post(ctx, t.Links["cancel"].Href, body, nil, t)
}

// Destination returns the customer transfer destination
func (t *Transfer) Destination(ctx context.Context) (*Customer, error) {
	if _, ok := t.Links["destination"]; !ok {
		return nil, errors.New("No destination resource link")
	}

	return t.client.Customer.Retrieve(ctx, t.Links["destination"].Href)
}

// DestinationString returns the customer transfer destination id
func (t Transfer) DestinationString() string {
	if _, ok := t.Links["destination"]; !ok {
		return ""
	}

	parts := strings.Split(t.Links["destination"].Href, "/")

	return parts[len(parts)-1]
}

// DestinationFundingSource returns the transfer funding source destination
func (t *Transfer) DestinationFundingSource(ctx context.Context) (*FundingSource, error) {
	if _, ok := t.Links["destination-funding-source"]; !ok {
		return nil, errors.New("No destination funding source resource link")
	}

	return t.client.FundingSource.Retrieve(ctx, t.Links["destination-funding-source"].Href)
}

// DestinationFundingSourceString returns the funding source destination id
func (t Transfer) DestinationFundingSourceString() string {
	if _, ok := t.Links["destination-funding-source"]; !ok {
		return ""
	}

	parts := strings.Split(t.Links["destination-funding-source"].Href, "/")

	return parts[len(parts)-1]
}

// ListFees returns the fees associated with the transfer
//
// see: https://docsv2.dwolla.com/#list-fees-for-a-transfer
func (t *Transfer) ListFees(ctx context.Context) (*TransferFees, error) {
	var fees TransferFees

	if _, ok := t.Links["fees"]; !ok {
		return nil, errors.New("No fees resource link")
	}

	if err := t.client.Get(ctx, t.Links["fees"].Href, nil, nil, &fees); err != nil {
		return nil, err
	}

	for i := range fees.Transactions {
		fees.Transactions[i].client = t.client
	}

	return &fees, nil
}

// Source returns the customer transfer source
func (t *Transfer) Source(ctx context.Context) (*Customer, error) {
	if _, ok := t.Links["source"]; !ok {
		return nil, errors.New("No source resource link")
	}

	return t.client.Customer.Retrieve(ctx, t.Links["source"].Href)
}

// SourceString returns the customer transfer source id
func (t Transfer) SourceString() string {
	if _, ok := t.Links["source"]; !ok {
		return ""
	}

	parts := strings.Split(t.Links["source"].Href, "/")

	return parts[len(parts)-1]
}

// SourceFundingSource returns the transfer funding source
func (t *Transfer) SourceFundingSource(ctx context.Context) (*FundingSource, error) {
	if _, ok := t.Links["source-funding-source"]; !ok {
		return nil, errors.New("No source funding source resource link")
	}

	return t.client.FundingSource.Retrieve(ctx, t.Links["source-funding-source"].Href)
}

// SourceFundingSourceString returns the transfer funding source
func (t Transfer) SourceFundingSourceString() string {
	if _, ok := t.Links["source-funding-source"]; !ok {
		return ""
	}

	parts := strings.Split(t.Links["source-funding-source"].Href, "/")

	return parts[len(parts)-1]
}

// RetrieveFailureReason returns the transfer's failure reason
//
// see: https://developers.dwolla.com/api-reference/transfers/retrieve-a-transfer-failure-reason
func (t *Transfer) RetrieveFailureReason(ctx context.Context) (*TransferFailure, error) {
	var transferFailure TransferFailure

	if _, ok := t.Links["failure"]; !ok {
		return nil, errors.New("No failure resource link")
	}

	if err := t.client.Get(ctx, t.Links["failure"].Href, nil, nil, &transferFailure); err != nil {
		return nil, err
	}

	return &transferFailure, nil
}

// CreatedTime returns the created value as time.Time
func (t *Transfer) CreatedTime() time.Time {
	createdTime, _ := time.Parse(time.RFC3339, t.Created)
	return createdTime
}
