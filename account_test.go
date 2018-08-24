package dwolla

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountGet(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "account.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.root = &Resource{Links: Links{"account": Link{HREF: "foobar"}}}
	c.Token = &Token{}

	res, err := c.Account.Get()

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.ID, "ca32853c-48fa-40be-ae75-77b37504581b")
	assert.Equal(t, res.Name, "Jane Doe")
}

func TestAccountListFundingSources(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "funding-sources.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	a := &Account{client: c, Resource: Resource{Links: Links{"funding-sources": Link{HREF: "foobar"}}}}
	res, err := a.ListFundingSources(nil)

	assert.Nil(t, err)
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

func TestAccountListMassPayments(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "mass-payments.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	a := &Account{client: c, Resource: Resource{Links: Links{"mass-payments": Link{HREF: "foobar"}}}}
	res, err := a.ListMassPayments(nil)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.Total, 1)

	mp := res.Embedded["mass-payments"][0]

	assert.Equal(t, mp.ID, "b4b5a699-5278-4727-9f81-a50800ea9abc")
	assert.Equal(t, mp.Status, "complete")
	assert.Equal(t, mp.Created, "2015-09-03T14:14:10.000Z")
	assert.Equal(t, mp.MetaData, map[string]string(map[string]string{"UserJobId": "some ID"}))
	assert.Equal(t, mp.CorrelationID, "8a2cdc8d-629d-4a24-98ac-40b735229fe2")
}

func TestAccountListTransfers(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "transfers.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	a := &Account{client: c, Resource: Resource{Links: Links{"transfers": Link{HREF: "foobar"}}}}
	res, err := a.ListTransfers(nil)

	assert.Nil(t, err)
	assert.NotNil(t, res)
}
