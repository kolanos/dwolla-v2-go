package dwolla

import (
	"fmt"
)

// FundingSourceService is the funding source service interface
// see: https://docsv2.dwolla.com/#funding-sources
type FundingSourceService interface {
	Retrieve(string) (*FundingSource, error)
	Update(string, *FundingSourceRequest) (*FundingSource, error)
	Remove(string) error
}

// FundingSourceServiceOp is an implementation of the funding source interface
type FundingSourceServiceOp struct {
	client *Client
}

// FundingSourceStatus is a funding source's status
type FundingSourceStatus string

const (
	// FundingSourceStatusUnverified is when the funding source is unverified
	FundingSourceStatusUnverified FundingSourceStatus = "unverified"
	// FundingSourceStatusVerified is when the funding source is verified
	FundingSourceStatusVerified FundingSourceStatus = "verified"
)

// FundingSourceType is the funding source type
type FundingSourceType string

const (
	// FundingSourceTypeBank is when the funding source is a bank account
	FundingSourceTypeBank FundingSourceType = "bank"
	// FundingSourceTypeBalance is when the funding source is a dwolla balance
	FundingSourceTypeBalance FundingSourceType = "balance"
)

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

const (
	// FundingSourceBankAccountTypeChecking is a checking bank account
	FundingSourceBankAccountTypeChecking FundingSourceBankAccountType = "checking"
	// FundingSourceBankAccountTypeSavings is a savings bank account
	FundingSourceBankAccountTypeSavings FundingSourceBankAccountType = "savings"
)

// FundingSourceBalance is a funding source balance
type FundingSourceBalance struct {
	Resource
	Balance     Amount `json:"balance"`
	LastUpdated string `json:"lastUpdated"`
}

// FundingSourceRequest is a funding source request
type FundingSourceRequest struct {
	Resource
	RoutingNumber   string                       `json:"routingNumber,omitempty"`
	AccountNumber   string                       `json:"accountNumber,omitempty"`
	BankAccountType FundingSourceBankAccountType `json:"bankAccountType,omitempty"`
	Name            string                       `json:"name,omitempty"`
	Channels        []string                     `json:"channels,omitempty"`
	Removed         bool                         `json:"removed,omitempty"`
}

// Retrieve retrieves a funding source with the matching id
// see: https://docsv2.dwolla.com/#retrieve-a-funding-source
func (f *FundingSourceServiceOp) Retrieve(id string) (*FundingSource, error) {
	var source FundingSource

	if err := f.client.Get(fmt.Sprintf("funding-sources/%s", id), nil, nil, &source); err != nil {
		return nil, err
	}

	source.client = f.client

	return &source, nil
}

// Update updates the funding source with matching id
// see: https://docsv2.dwolla.com/#update-a-funding-source
func (f *FundingSourceServiceOp) Update(id string, body *FundingSourceRequest) (*FundingSource, error) {
	var source FundingSource

	if err := f.client.Post(fmt.Sprintf("funding-sources/%s", id), body, nil, &source); err != nil {
		return nil, err
	}

	source.client = f.client

	return &source, nil
}

// Remove removes a funding source matching the id
// see: https://docsv2.dwolla.com/#remove-a-funding-source
func (f *FundingSourceServiceOp) Remove(id string) error {
	body := &FundingSourceRequest{Removed: true}

	return f.client.Post(fmt.Sprintf("funding-sources/%s", id), body, nil, f)
}

// Customer returns the funding source's customer
func (f *FundingSource) Customer() (*Customer, error) {
	if _, ok := f.Links["customer"]; !ok {
		return nil, fmt.Errorf("No customer resource link")
	}

	return f.client.Customer.Retrieve(f.Links["customer"].Href)
}

// FailedVerificationMicroDeposits returns true if micro deposit
// verificationfailed
func (f *FundingSource) FailedVerificationMicroDeposits() bool {
	_, ok := f.Links["failed-verification-micro-deposits"]
	return ok
}

// InitiateMicroDeposits initiates micro deposit verification
// see: https://docsv2.dwolla.com/#initiate-micro-deposits
func (f *FundingSource) InitiateMicroDeposits() (*MicroDeposit, error) {
	var deposit MicroDeposit

	if _, ok := f.Links["initiate-micro-deposits"]; !ok {
		return nil, fmt.Errorf("No initiate micro deposits resource link")
	}

	if err := f.client.Post(f.Links["initiate-micro-deposits"].Href, nil, nil, &deposit); err != nil {
		return nil, err
	}

	deposit.client = f.client

	return &deposit, nil
}

// Remove removes the funding source
// see: https://docsv2.dwolla.com/#remove-a-funding-source
func (f *FundingSource) Remove() error {
	if _, ok := f.Links["remove"]; !ok {
		return fmt.Errorf("No remove resource link")
	}

	request := &FundingSourceRequest{Removed: true}

	return f.client.Post(f.Links["remove"].Href, request, nil, f)
}

// RetrieveBalance retrieves the funding source balance
// see: https://docsv2.dwolla.com/#retrieve-a-funding-source-balance
func (f *FundingSource) RetrieveBalance() (*FundingSourceBalance, error) {
	var balance FundingSourceBalance

	if _, ok := f.Links["balance"]; !ok {
		return nil, fmt.Errorf("No balance resource link")
	}

	if err := f.client.Get(f.Links["balance"].Href, nil, nil, &balance); err != nil {
		return nil, err
	}

	balance.client = f.client

	return &balance, nil
}

// RetrieveMicroDeposits retrieves funding source micro deposits
// see: https://docsv2.dwolla.com/#retrieve-micro-deposits-details
func (f *FundingSource) RetrieveMicroDeposits() (*MicroDeposit, error) {
	var deposit MicroDeposit

	if _, ok := f.Links["verify-micro-deposits"]; !ok {
		return nil, fmt.Errorf("No verify micro deposits resource link")
	}

	if err := f.client.Get(f.Links["verify-micro-deposits"].Href, nil, nil, &deposit); err != nil {
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
// see: https://docsv2.dwolla.com/#update-a-funding-source
func (f *FundingSource) Update(body *FundingSourceRequest) error {
	if _, ok := f.Links["self"]; !ok {
		return fmt.Errorf("No self resource link")
	}

	return f.client.Post(f.Links["self"].Href, body, nil, f)
}

// VerifyMicroDeposits verifies micro deposit amounts
// see: https://docsv2.dwolla.com/#verify-micro-deposits
func (f *FundingSource) VerifyMicroDeposits(body *MicroDepositRequest) error {
	if _, ok := f.Links["verify-micro-deposits"]; !ok {
		return fmt.Errorf("No verify micro deposits resource link")
	}

	return f.client.Post(f.Links["verify-micro-deposits"].Href, body, nil, nil)
}
