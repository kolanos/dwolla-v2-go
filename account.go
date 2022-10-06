package dwolla

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// AccountService is the account service interface
//
// see: https://developers.dwolla.com/api-reference/accounts
type AccountService interface {
	Retrieve(context.Context) (*Account, error)
}

// AccountServiceOp is an implementation of the account service interface
type AccountServiceOp struct {
	client *Client
}

// Account is a dwolla account
type Account struct {
	Resource
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	TimezoneOffset float32 `json:"timezoneOffset"`
	Type           string  `json:"type"`
}

// Retrieve retrieves the dwolla account
//
// see: https://docsv2.dwolla.com/#retrieve-account-details
func (a *AccountServiceOp) Retrieve(ctx context.Context) (*Account, error) {
	root, err := a.client.Root(ctx)

	if err != nil {
		return nil, err
	}

	var account Account

	if _, ok := root.Links["account"]; !ok {
		return nil, errors.New("No account resource link")
	}

	if err := a.client.Get(ctx, root.Links["account"].Href, nil, nil, &account); err != nil {
		return nil, err
	}

	account.client = a.client

	return &account, nil
}

// CreateFundingSource creates a funding source for the account
//
// see: https://docsv2.dwolla.com/#create-a-funding-source-for-an-account
func (a *Account) CreateFundingSource(ctx context.Context, body *FundingSourceRequest) (*FundingSource, error) {
	var source FundingSource

	var headers *http.Header
	if body.IdempotencyKey != "" {
		headers = &http.Header{}
		headers.Set(HeaderIdempotency, body.IdempotencyKey)
	}

	if err := a.client.Post(ctx, "funding-sources", body, headers, &source); err != nil {
		return nil, err
	}

	source.client = a.client

	return &source, nil
}

// ListFundingSources returns the account's funding sources
//
// see: https://developers.dwolla.com/api-reference/accounts/list-funding-sources
func (a *Account) ListFundingSources(ctx context.Context, removed *bool) (*FundingSources, error) {
	var sources FundingSources

	if _, ok := a.Links["funding-sources"]; !ok {
		return nil, errors.New("No funding sources resource link")
	}

	var params *url.Values
	if removed != nil {
		params = &url.Values{}
		params.Add("removed", strconv.FormatBool(*removed))
	}

	if err := a.client.Get(ctx, a.Links["funding-sources"].Href, params, nil, &sources); err != nil {
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
//
// see: https://docsv2.dwolla.com/#list-mass-payments-for-an-account
func (a *Account) ListMassPayments(ctx context.Context, params *url.Values) (*MassPayments, error) {
	var payments MassPayments

	if _, ok := a.Links["self"]; !ok {
		return nil, errors.New("No self resource link")
	}

	if err := a.client.Get(ctx, fmt.Sprintf("%s/mass-payments", a.Links["self"].Href), params, nil, &payments); err != nil {
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
//
// see: https://docsv2.dwolla.com/#list-and-search-transfers-for-an-account
func (a *Account) ListTransfers(ctx context.Context, params *url.Values) (*Transfers, error) {
	var transfers Transfers

	if _, ok := a.Links["transfers"]; !ok {
		return nil, errors.New("No transfers resource link")
	}

	if err := a.client.Get(ctx, a.Links["transfers"].Href, params, nil, &transfers); err != nil {
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
