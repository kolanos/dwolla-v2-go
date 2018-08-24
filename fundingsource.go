package dwolla

import (
	"net/url"
)

// FundingSourceService is the funding source service interface
// see: https://docsv2.dwolla.com/#funding-sources
type FundingSourceService interface {
	Get(string) (*FundingSource, error)
	List(*url.Values) (*FundingSources, error)
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

// FundingSourceCreate is a funding source create request
type FundingSourceCreate struct {
	RoutingNumber   string                       `json:"routingNumber"`
	AccountNumber   string                       `json:"accountNumber"`
	BankAccountType FundingSourceBankAccountType `json:"bankAccountType"`
	Name            string                       `json:"name"`
	Channels        []string                     `json:"channels,omitempty"`
}

// Get returns a funding source with the matching the id
func (f *FundingSourceServiceOp) Get(id string) (*FundingSource, error) {
	var fundingSource FundingSource
	return &fundingSource, nil
}

// List returns a collection of funding sources
func (f *FundingSourceServiceOp) List(params *url.Values) (*FundingSources, error) {
	var fundingSources FundingSources
	return &fundingSources, nil
}
