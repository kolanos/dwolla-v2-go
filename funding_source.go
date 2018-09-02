package dwolla

import (
	"fmt"
	"net/url"
)

// FundingSourceService is the funding source service interface
// see: https://docsv2.dwolla.com/#funding-sources
type FundingSourceService interface {
	Retrieve(string) (*FundingSource, error)
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
func (f *FundingSourceServiceOp) Retrieve(id string) (*FundingSource, error) {
	var fundingSource FundingSource
	return &fundingSource, nil
}

// List returns a collection of funding sources
func (f *FundingSourceServiceOp) List(params *url.Values) (*FundingSources, error) {
	var fundingSources FundingSources
	return &fundingSources, nil
}

// Remove removes a funding source from the account/customer
// see: https://docsv2.dwolla.com/#remove-a-funding-source
func (f *FundingSource) Remove() error {
	if _, ok := f.Links["self"]; !ok {
		return fmt.Errorf("No self resource link")
	}

	request := &FundingSourceRequest{Removed: true}

	return f.client.Post(f.Links["self"].Href, request, nil, f)
}
