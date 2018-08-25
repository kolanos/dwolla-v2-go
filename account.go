package dwolla

import (
	"fmt"
	"net/url"
)

// AccountService is the account service interface
// see: https://docsv2.dwolla.com/#accounts
type AccountService interface {
	Get() (*Account, error)
}

// AccountServiceOp is an implementation of the account service interface
type AccountServiceOp struct {
	client *Client
}

// Account is a dwolla account
type Account struct {
	Resource
	client         *Client
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	TimezoneOffset float32 `json:"timezoneOffset"`
	Type           string  `json:"type"`
}

// Get returns the dwolla account
// see: https://docsv2.dwolla.com/#retrieve-account-details
func (a *AccountServiceOp) Get() (*Account, error) {
	root, err := a.client.Root()
	if err != nil {
		return nil, err
	}

	var account Account

	if _, ok := root.Links["account"]; !ok {
		return nil, fmt.Errorf("No account resource link")
	}

	if err := a.client.Get(root.Links["account"].HREF, nil, nil, &account); err != nil {
		return nil, err
	}

	account.client = a.client
	return &account, nil
}

// CreateFundingSource creates a funding source for the account
// see: https://docsv2.dwolla.com/#create-a-funding-source-for-an-account
func (a *Account) CreateFundingSource(body *FundingSourceCreate) (*FundingSource, error) {
	var source FundingSource

	if _, ok := a.Links["funding-sources"]; !ok {
		return nil, fmt.Errorf("No funding sources resource link")
	}

	if err := a.client.Post(a.Links["funding-sources"].HREF, body, nil, &source); err != nil {
		return nil, err
	}

	source.client = a.client
	return &source, nil
}

// ListFundingSources returns the account's funding sources
// see: https://docsv2.dwolla.com/#list-funding-sources-for-an-account
func (a *Account) ListFundingSources(params *url.Values) (*FundingSources, error) {
	var sources FundingSources

	if _, ok := a.Links["funding-sources"]; !ok {
		return nil, fmt.Errorf("No funding sources resource link")
	}

	if err := a.client.Get(a.Links["funding-sources"].HREF, params, nil, &sources); err != nil {
		return nil, err
	}

	sources.client = a.client

	for i := range sources.Embedded["funding-sources"] {
		sources.Embedded["funding-sources"][i].client = a.client
	}

	return &sources, nil
}

// ListMassPayments returns mass payments for the account
// see: https://docsv2.dwolla.com/#list-mass-payments-for-an-account
func (a *Account) ListMassPayments(params *url.Values) (*MassPayments, error) {
	var payments MassPayments

	if _, ok := a.Links["mass-payments"]; !ok {
		return nil, fmt.Errorf("No mass payments resource link")
	}

	if err := a.client.Get(a.Links["mass-payments"].HREF, params, nil, &payments); err != nil {
		return nil, err
	}

	payments.client = a.client

	for i := range payments.Embedded["mass-payments"] {
		payments.Embedded["mass-payments"][i].client = a.client
	}

	return &payments, nil
}

// ListTransfers returns the account's transfers
// see: https://docsv2.dwolla.com/#list-and-search-transfers-for-an-account
func (a *Account) ListTransfers(params *url.Values) (*Transfers, error) {
	var transfers Transfers

	if _, ok := a.Links["transfers"]; !ok {
		return nil, fmt.Errorf("No transfers resource link")
	}

	if err := a.client.Get(a.Links["transfers"].HREF, params, nil, &transfers); err != nil {
		return nil, err
	}

	transfers.client = a.client

	for i := range transfers.Embedded["transfers"] {
		transfers.Embedded["transfers"][i].client = a.client
	}

	return &transfers, nil
}
