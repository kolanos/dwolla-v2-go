package dwolla

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomerRetrieve(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "customer.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	res, err := c.Customer.Retrieve("FC451A7A-AE30-4404-AB95-E3553FCD733F")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.ID, "FC451A7A-AE30-4404-AB95-E3553FCD733F")
	assert.Equal(t, res.FirstName, "Jane")
	assert.Equal(t, res.LastName, "Doe")
	assert.Equal(t, res.Email, "janedoe@nomail.com")
	assert.Equal(t, res.Type, CustomerTypeUnverified)
	assert.Equal(t, res.Status, CustomerStatusUnverified)
	assert.Equal(t, res.Created, "2015-09-03T23:56:10.023Z")
}

func TestCustomerList(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "customers.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	res, err := c.Customer.List(nil)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.Total, 1)

	cs := res.Embedded["customers"][0]

	assert.Equal(t, cs.ID, "FC451A7A-AE30-4404-AB95-E3553FCD733F")
	assert.Equal(t, cs.FirstName, "Jane")
	assert.Equal(t, cs.LastName, "Doe")
	assert.Equal(t, cs.Email, "janedoe@nomail.com")
	assert.Equal(t, cs.Type, CustomerTypeUnverified)
	assert.Equal(t, cs.Status, CustomerStatusUnverified)
	assert.Equal(t, cs.Created, "2015-09-03T23:56:10.023Z")
}
