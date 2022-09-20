package dwolla

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountServiceRetrieve(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "account.json"))
	c.root = &Resource{Links: Links{"account": Link{Href: "foobar"}}}
	res, err := c.Account.Retrieve(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.ID, "ca32853c-48fa-40be-ae75-77b37504581b")
	assert.Equal(t, res.Name, "Jane Doe")
}

func TestAccountServiceRetrieveError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))
	res, err := c.Account.Retrieve(ctx)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "[NotFound] Resource Not Found")
	assert.Nil(t, res)

	c = newMockClient(200, filepath.Join("testdata", "account.json"))
	res, err = c.Account.Retrieve(ctx)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "No account resource link")
	assert.Nil(t, res)

	c = newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))
	c.root = &Resource{Links: Links{"account": Link{Href: "foobar"}}}
	res, err = c.Account.Retrieve(ctx)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "[NotFound] Resource Not Found")
	assert.Nil(t, res)
}

func TestAccountCreateFundingSource(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "funding-source.json"))
	a := &Account{Resource: Resource{client: c}}
	fs := &FundingSourceRequest{
		RoutingNumber:   "222222226",
		AccountNumber:   "0123456789",
		BankAccountType: FundingSourceBankAccountTypeChecking,
		Name:            "My Checking Account",
	}

	res, err := a.CreateFundingSource(ctx, fs)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestAccountCreateFundingSourceError(t *testing.T) {
	c := newMockClient(400, filepath.Join("testdata", "validation-error.json"))
	a := &Account{Resource: Resource{client: c}}
	res, err := a.CreateFundingSource(ctx, &FundingSourceRequest{
		RoutingNumber:   "222222226",
		AccountNumber:   "0123456789",
		BankAccountType: FundingSourceBankAccountTypeChecking,
		Name:            "My Checking Account",
	})

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "[ValidationError] There was a validation error.")
	assert.Nil(t, res)
}

func TestAccountListFundingSources(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "funding-sources.json"))
	a := &Account{Resource: Resource{client: c, Links: Links{"funding-sources": Link{Href: "foobar"}}}}
	res, err := a.ListFundingSources(ctx, false)

	assert.NoError(t, err)
	assert.NotNil(t, res)

	fs := res.Embedded["funding-sources"][0]

	assert.Equal(t, fs.ID, "04173e17-6398-4d36-a167-9d98c4b1f1c3")
	assert.Equal(t, fs.Status, FundingSourceStatusVerified)
	assert.Equal(t, fs.Type, FundingSourceTypeBank)
	assert.Equal(t, fs.BankAccountType, FundingSourceBankAccountTypeChecking)
	assert.Equal(t, fs.Name, "My Account - Checking")
	assert.Equal(t, fs.Created, "2017-09-25T20:03:41.000Z")
	assert.Equal(t, fs.Removed, false)
	assert.Equal(t, fs.Channels, []string{"ach"})
	assert.Equal(t, fs.BankName, "First Midwestern Bank")
}

func TestAccountListFundingSourcesError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))
	a := &Account{Resource: Resource{client: c, Links: Links{"funding-sources": Link{Href: "foobar"}}}}
	res, err := a.ListFundingSources(ctx, false)

	assert.Equal(t, err.Error(), "[NotFound] Resource Not Found")
	assert.Error(t, err)
	assert.Nil(t, res)

	c = newMockClient(200, filepath.Join("testdata", "funding-sources.json"))
	a = &Account{Resource: Resource{client: c}}
	res, err = a.ListFundingSources(ctx, false)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "No funding sources resource link")
	assert.Nil(t, res)
}

func TestAccountListMassPayments(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "mass-payments.json"))
	a := &Account{Resource: Resource{client: c, Links: Links{"self": Link{Href: "foobar"}}}}
	res, err := a.ListMassPayments(ctx, nil)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.Total, 1)

	mp := res.Embedded["mass-payments"][0]

	assert.Equal(t, mp.ID, "b4b5a699-5278-4727-9f81-a50800ea9abc")
	assert.Equal(t, mp.Status, MassPaymentStatusComplete)
	assert.Equal(t, mp.Created, "2015-09-03T14:14:10.000Z")
	assert.Equal(t, mp.MetaData, MetaData{"UserJobId": "some ID"})
	assert.Equal(t, mp.CorrelationID, "8a2cdc8d-629d-4a24-98ac-40b735229fe2")
}

func TestAccountListMassPaymentsError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))
	a := &Account{Resource: Resource{client: c, Links: Links{"self": Link{Href: "foobar"}}}}
	res, err := a.ListMassPayments(ctx, nil)

	assert.Error(t, err)
	assert.Nil(t, res)

	c = newMockClient(200, filepath.Join("testdata", "mass-payments.json"))
	a = &Account{Resource: Resource{client: c}}
	res, err = a.ListMassPayments(ctx, nil)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "No self resource link")
	assert.Nil(t, res)
}

func TestAccountListTransfers(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "transfers.json"))
	a := &Account{Resource: Resource{client: c, Links: Links{"transfers": Link{Href: "foobar"}}}}
	res, err := a.ListTransfers(ctx, nil)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestAccountListTransfersError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))
	a := &Account{Resource: Resource{client: c, Links: Links{"transfers": Link{Href: "foobar"}}}}
	res, err := a.ListTransfers(ctx, nil)

	assert.Error(t, err)
	assert.Nil(t, res)

	c = newMockClient(200, filepath.Join("testdata", "transfers.json"))
	a = &Account{Resource: Resource{client: c}}
	res, err = a.ListTransfers(ctx, nil)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "No transfers resource link")
	assert.Nil(t, res)
}
