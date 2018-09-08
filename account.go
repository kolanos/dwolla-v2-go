package dwolla

import (
	"errors"
	"net/url"
)

// AccountService is the account service interface
// see: https://docsv2.dwolla.com/#accounts
type AccountService interface {
	Retrieve() (*Account, error)
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

// Retrieve retrieves the dwolla account
// see: https://docsv2.dwolla.com/#retrieve-account-details
func (a *AccountServiceOp) Retrieve() (*Account, error) {
	root, err := a.client.Root()
	if err != nil {
		return nil, err
	}

	var account Account

	if _, ok := root.Links["account"]; !ok {
		return nil, errors.New("No account resource link")
	}

	if err := a.client.Get(root.Links["account"].Href, nil, nil, &account); err != nil {
		return nil, err
	}

	account.client = a.client

	return &account, nil
}

// CreateFundingSource creates a funding source for the account
// see: https://docsv2.dwolla.com/#create-a-funding-source-for-an-account
func (a *Account) CreateFundingSource(body *FundingSourceRequest) (*FundingSource, error) {
	var source FundingSource

	if err := a.client.Post("funding-sources", body, nil, &source); err != nil {
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
		return nil, errors.New("No funding sources resource link")
	}

	if err := a.client.Get(a.Links["funding-sources"].Href, params, nil, &sources); err != nil {
		return nil, err
	}

	sources.client = a.client

	if _, ok := sources.Embedded["funding-sources"]; ok {
		for i := range sources.Embedded["funding-sources"] {
			sources.Embedded["funding-sources"][i].client = a.client
		}
	}

	return &sources, nil
}

// ListMassPayments returns mass payments for the account
// see: https://docsv2.dwolla.com/#list-mass-payments-for-an-account
func (a *Account) ListMassPayments(params *url.Values) (*MassPayments, error) {
	var payments MassPayments

	if _, ok := a.Links["mass-payments"]; !ok {
		return nil, errors.New("No mass payments resource link")
	}

	if err := a.client.Get(a.Links["mass-payments"].Href, params, nil, &payments); err != nil {
		return nil, err
	}

	payments.client = a.client

	if _, ok := payments.Embedded["mass-payments"]; ok {
		for i := range payments.Embedded["mass-payments"] {
			payments.Embedded["mass-payments"][i].client = a.client
		}
	}

	return &payments, nil
}

// ListTransfers returns the account's transfers
// see: https://docsv2.dwolla.com/#list-and-search-transfers-for-an-account
func (a *Account) ListTransfers(params *url.Values) (*Transfers, error) {
	var transfers Transfers

	if _, ok := a.Links["transfers"]; !ok {
		return nil, errors.New("No transfers resource link")
	}

	if err := a.client.Get(a.Links["transfers"].Href, params, nil, &transfers); err != nil {
		return nil, err
	}

	transfers.client = a.client

	if _, ok := transfers.Embedded["transfers"]; ok {
		for i := range transfers.Embedded["transfers"] {
			transfers.Embedded["transfers"][i].client = a.client
		}
	}

	return &transfers, nil
}
