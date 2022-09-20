package dwolla

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const (
	// FundingSourceBankAccountTypeChecking is a checking bank account
	FundingSourceBankAccountTypeChecking FundingSourceBankAccountType = "checking"
	// FundingSourceBankAccountTypeSavings is a savings bank account
	FundingSourceBankAccountTypeSavings FundingSourceBankAccountType = "savings"
)

const (
	// FundingSourceStatusUnverified is when the funding source is unverified
	FundingSourceStatusUnverified FundingSourceStatus = "unverified"
	// FundingSourceStatusVerified is when the funding source is verified
	FundingSourceStatusVerified FundingSourceStatus = "verified"
)

const (
	// FundingSourceTypeBank is when the funding source is a bank account
	FundingSourceTypeBank FundingSourceType = "bank"
	// FundingSourceTypeBalance is when the funding source is a dwolla balance
	FundingSourceTypeBalance FundingSourceType = "balance"
)

// FundingSourceService is the funding source service interface
//
// see: https://developers.dwolla.com/api-reference/funding-sources
type FundingSourceService interface {
	Retrieve(context.Context, string) (*FundingSource, error)
	Update(context.Context, string, *FundingSourceRequest) (*FundingSource, error)
	Remove(context.Context, string) error
}

// FundingSourceServiceOp is an implementation of the funding source interface
type FundingSourceServiceOp struct {
	client *Client
}

// FundingSourceStatus is a funding source's status
type FundingSourceStatus string

// FundingSourceType is the funding source type
type FundingSourceType string

// FundingSource is a dwolla funding source
type FundingSource struct {
	Resource
	ID              string                       `json:"id"`
	Status          FundingSourceStatus          `json:"status"`
	Type            FundingSourceType            `json:"type"`
	BankAccountType FundingSourceBankAccountType `json:"bankAccountType"`
	Name            string                       `json:"name"`
	Created         string                       `json:"created"`
	Balance         Amount                       `json:"balance"`
	Removed         bool                         `json:"removed"`
	Channels        []string                     `json:"channels"`
	BankName        string                       `json:"bankName"`
	Fingerprint     string                       `json:"fingerprint"`
}

// FundingSources is a collection of funding sources
type FundingSources struct {
	Collection
	Embedded map[string][]FundingSource `json:"_embedded"`
}

// FundingSourceBankAccountType is a dwolla bank account type enum
type FundingSourceBankAccountType string

// FundingSourceBalance is a funding source balance
type FundingSourceBalance struct {
	Resource
	Balance     Amount `json:"balance"`
	Total       Amount `json:"total"`
	LastUpdated string `json:"lastUpdated"`
}

// FundingSourceRequest is a funding source request
type FundingSourceRequest struct {
	Resource
	RoutingNumber   string                       `json:"routingNumber,omitempty"`
	AccountNumber   string                       `json:"accountNumber,omitempty"`
	BankAccountType FundingSourceBankAccountType `json:"bankAccountType,omitempty"`
	Name            string                       `json:"name,omitempty"` // Arbitrary nickname for the funding source. Must be 50 characters or less.
	Channels        []string                     `json:"channels,omitempty"`
	Removed         bool                         `json:"removed,omitempty"`
	PlaidToken      string                       `json:"plaidToken,omitempty"` // A processor token obtained from Plaid for adding and verifying a bank
	IdempotencyKey  string                       `json:"-"`
}

// FundingSourceToken is a funding source dwolla.js token
type FundingSourceToken struct {
	Resource
	Token string `json:"token"`
}

// Retrieve retrieves a funding source with the matching id
//
// see: https://docsv2.dwolla.com/#retrieve-a-funding-source
func (f *FundingSourceServiceOp) Retrieve(ctx context.Context, id string) (*FundingSource, error) {
	var source FundingSource

	if err := f.client.Get(ctx, fmt.Sprintf("funding-sources/%s", id), nil, nil, &source); err != nil {
		return nil, err
	}

	source.client = f.client

	return &source, nil
}

// Update updates the funding source with matching id
//
// see: https://docsv2.dwolla.com/#update-a-funding-source
func (f *FundingSourceServiceOp) Update(ctx context.Context, id string, body *FundingSourceRequest) (*FundingSource, error) {
	var source FundingSource

	var headers *http.Header
	if body.IdempotencyKey != "" {
		headers = &http.Header{}
		headers.Set(HeaderIdempotency, body.IdempotencyKey)
	}

	if err := f.client.Post(ctx, fmt.Sprintf("funding-sources/%s", id), body, headers, &source); err != nil {
		return nil, err
	}

	source.client = f.client

	return &source, nil
}

// Remove removes a funding source matching the id
//
// see: https://docsv2.dwolla.com/#remove-a-funding-source
func (f *FundingSourceServiceOp) Remove(ctx context.Context, id string) error {
	body := &FundingSourceRequest{Removed: true}

	return f.client.Post(ctx, fmt.Sprintf("funding-sources/%s", id), body, nil, f)
}

// Customer returns the funding source's customer
func (f *FundingSource) Customer(ctx context.Context) (*Customer, error) {
	var customer Customer

	if _, ok := f.Links["customer"]; !ok {
		return nil, errors.New("No customer resource link")
	}

	if err := f.client.Get(ctx, f.Links["customer"].Href, nil, nil, &customer); err != nil {
		return nil, err
	}

	customer.client = f.client

	return &customer, nil
}

// FailedVerificationMicroDeposits returns true if micro deposit
// verificationfailed
func (f *FundingSource) FailedVerificationMicroDeposits() bool {
	_, ok := f.Links["failed-verification-micro-deposits"]
	return ok
}

// InitiateMicroDeposits initiates micro deposit verification
//
// see: https://docsv2.dwolla.com/#initiate-micro-deposits
func (f *FundingSource) InitiateMicroDeposits(ctx context.Context, idempotencyKey string) (*MicroDeposit, error) {
	var deposit MicroDeposit

	if _, ok := f.Links["initiate-micro-deposits"]; !ok {
		return nil, errors.New("No initiate micro deposits resource link")
	}

	var headers *http.Header
	if idempotencyKey != "" {
		headers = &http.Header{}
		headers.Set(HeaderIdempotency, idempotencyKey)
	}

	if err := f.client.Post(ctx, f.Links["initiate-micro-deposits"].Href, nil, headers, &deposit); err != nil {
		return nil, err
	}

	deposit.client = f.client

	return &deposit, nil
}

// Remove removes the funding source
//
// see: https://docsv2.dwolla.com/#remove-a-funding-source
func (f *FundingSource) Remove(ctx context.Context) error {
	if _, ok := f.Links["remove"]; !ok {
		return errors.New("No remove resource link")
	}

	request := &FundingSourceRequest{Removed: true}

	return f.client.Post(ctx, f.Links["remove"].Href, request, nil, f)
}

// RetrieveBalance retrieves the funding source balance
//
// see: https://docsv2.dwolla.com/#retrieve-a-funding-source-balance
func (f *FundingSource) RetrieveBalance(ctx context.Context) (*FundingSourceBalance, error) {
	var balance FundingSourceBalance

	if _, ok := f.Links["balance"]; !ok {
		return nil, errors.New("No balance resource link")
	}

	if err := f.client.Get(ctx, f.Links["balance"].Href, nil, nil, &balance); err != nil {
		return nil, err
	}

	balance.client = f.client

	return &balance, nil
}

// RetrieveMicroDeposits retrieves funding source micro deposits
//
// see: https://docsv2.dwolla.com/#retrieve-micro-deposits-details
func (f *FundingSource) RetrieveMicroDeposits(ctx context.Context) (*MicroDeposit, error) {
	var deposit MicroDeposit

	if _, ok := f.Links["verify-micro-deposits"]; !ok {
		return nil, errors.New("No verify micro deposits resource link")
	}

	if err := f.client.Get(ctx, f.Links["verify-micro-deposits"].Href, nil, nil, &deposit); err != nil {
		return nil, err
	}

	deposit.client = f.client

	return &deposit, nil
}

// TransferFromBalance returns true if funding source can transfer from balance
func (f *FundingSource) TransferFromBalance() bool {
	_, ok := f.Links["transfer-from-balance"]
	return ok
}

// TransferToBalance returns true if funding source can transfer to balance
func (f *FundingSource) TransferToBalance() bool {
	_, ok := f.Links["transfer-to-balance"]
	return ok
}

// TransferReceive returns true if funding source can receive transfers
func (f *FundingSource) TransferReceive() bool {
	_, ok := f.Links["transfer-receive"]
	return ok
}

// TransferSend returns true if funding source can send transfers
func (f *FundingSource) TransferSend() bool {
	_, ok := f.Links["transfer-send"]
	return ok
}

// Update updates the funding source
//
// see: https://docsv2.dwolla.com/#update-a-funding-source
func (f *FundingSource) Update(ctx context.Context, body *FundingSourceRequest) error {
	if _, ok := f.Links["self"]; !ok {
		return errors.New("No self resource link")
	}

	var headers *http.Header
	if body.IdempotencyKey != "" {
		headers = &http.Header{}
		headers.Set(HeaderIdempotency, body.IdempotencyKey)
	}

	return f.client.Post(ctx, f.Links["self"].Href, body, headers, f)
}

// VerifyMicroDeposits verifies micro deposit amounts
//
// see: https://docsv2.dwolla.com/#verify-micro-deposits
func (f *FundingSource) VerifyMicroDeposits(ctx context.Context, body *MicroDepositRequest) error {
	if _, ok := f.Links["verify-micro-deposits"]; !ok {
		return errors.New("No verify micro deposits resource link")
	}

	var headers *http.Header
	if body.IdempotencyKey != "" {
		headers = &http.Header{}
		headers.Set(HeaderIdempotency, body.IdempotencyKey)
	}

	return f.client.Post(ctx, f.Links["verify-micro-deposits"].Href, body, headers, nil)
}
