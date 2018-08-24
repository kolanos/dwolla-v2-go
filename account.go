package dwolla

import (
	"net/url"
	"strconv"
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

	if err := a.client.Post(a.Links["funding-sources"].HREF, body, nil, &source); err != nil {
		return nil, err
	}

	source.client = a.client
	return &source, nil
}

// ListFundingSources returns the account's funding sources
// see: https://docsv2.dwolla.com/#list-funding-sources-for-an-account
func (a *Account) ListFundingSources(removed bool) (*FundingSources, error) {
	var sources FundingSources

	params := &url.Values{}
	params.Add("removed", strconv.FormatBool(removed))

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
func (a *Account) ListMassPayments() (*MassPayments, error) {
	var payments MassPayments

	if err := a.client.Get(a.Links["mass-payments"].HREF, nil, nil, &payments); err != nil {
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
func (a *Account) ListTransfers() (*Transfers, error) {
	var transfers Transfers

	if err := a.client.Get(a.Links["transfers"].HREF, nil, nil, &transfers); err != nil {
		return nil, err
	}

	transfers.client = a.client

	for i := range transfers.Embedded["transfers"] {
		transfers.Embedded["transfers"][i].client = a.client
	}

	return &transfers, nil
}
