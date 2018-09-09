package dwolla

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomerServiceCreate(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "customer.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	res, err := c.Customer.Create(&CustomerRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "janedoe@nomail.com",
		Type:      CustomerTypeUnverified,
	})

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

func TestCustomerServiceRetrieve(t *testing.T) {
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

func TestCustomerServiceList(t *testing.T) {
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

func TestCustomerServiceUpdate(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "customer.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	res, err := c.Customer.Update("FC451A7A-AE30-4404-AB95-E3553FCD733F", &CustomerRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "janedoe@nomail.com",
		Type:      CustomerTypeUnverified,
	})

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

func TestCustomerCertifyBeneficialOwnership(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "customer.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	customer := &Customer{}
	customer.client = c
	err := customer.CertifyBeneficialOwnership()

	assert.Error(t, err)

	customer.Links = Links{"certify-beneficial-ownership": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/beneficial-ownership"}}
	err = customer.CertifyBeneficialOwnership()

	assert.Nil(t, err)
}

func TestCustomerCreateDocument(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "document.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	customer := &Customer{}
	customer.client = c
	customer.Links = Links{"self": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}

	f2, _ := os.Open(filepath.Join("testdata", "document-upload-success.png"))

	doc, err := customer.CreateDocument(&DocumentRequest{
		DocumentType: "passport",
		FileName:     f2.Name(),
		File:         f2,
	})

	assert.Nil(t, err)
	assert.NotNil(t, doc)
}

func TestCustomerCreateBeneficialOwner(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "beneficial-owner.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	customer := &Customer{}
	customer.client = c
	customer.Links = Links{"beneficial-owners": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/beneficial-owners"}}

	owner, err := customer.CreateBeneficialOwner(&BeneficialOwnerRequest{
		FirstName:   "beneficial",
		LastName:    "owner",
		DateOfBirth: "1980-01-01",
		SSN:         "555-55-5555",
		Address: Address{
			Address1:            "123 Main St.",
			Address2:            "Apt 123",
			City:                "Des Moines",
			StateProvinceRegion: "IA",
			Country:             "US",
			PostalCode:          "50309",
		},
	})

	assert.Nil(t, err)
	assert.NotNil(t, owner)
}

func TestCustomerCreateFundingSource(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "funding-source.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	customer := &Customer{}
	customer.client = c
	customer.Links = Links{"funding-sources": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/funding-sources"}}

	source, err := customer.CreateFundingSource(&FundingSourceRequest{
		RoutingNumber:   "1234567890",
		AccountNumber:   "1234567890",
		BankAccountType: FundingSourceBankAccountTypeChecking,
		Name:            "Test Checking Account",
	})

	assert.Nil(t, err)
	assert.NotNil(t, source)
}

func TestCustomerDeactivate(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "customer.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	customer := &Customer{}
	customer.client = c

	err := customer.Deactivate()

	assert.Error(t, err)

	customer.Links = Links{"deactivate": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}
	err = customer.Deactivate()

	assert.Nil(t, err)
}

func TestCustomerListBeneficialOwners(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "beneficial-owners.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	customer := &Customer{}
	customer.client = c

	owners, err := customer.ListBeneficialOwners()

	assert.Error(t, err)
	assert.Nil(t, owners)

	customer.Links = Links{"beneficial-owners": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/beneficial-owners"}}
	owners, err = customer.ListBeneficialOwners()

	assert.Nil(t, err)
	assert.NotNil(t, owners)
}

func TestCustomerListDocuments(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "documents.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	customer := &Customer{}
	customer.client = c

	documents, err := customer.ListDocuments()

	assert.Error(t, err)
	assert.Nil(t, documents)

	customer.Links = Links{"self": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}
	documents, err = customer.ListDocuments()

	assert.Nil(t, err)
	assert.NotNil(t, documents)
}

func TestCustomerListFundingSources(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "funding-sources.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	customer := &Customer{}
	customer.client = c

	sources, err := customer.ListFundingSources(true)

	assert.Error(t, err)
	assert.Nil(t, sources)

	customer.Links = Links{"funding-sources": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/funding-sources"}}
	sources, err = customer.ListFundingSources(true)

	assert.Nil(t, err)
	assert.NotNil(t, sources)
}

func TestCustomerListMassPayments(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "mass-payments.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	customer := &Customer{}
	customer.client = c

	payments, err := customer.ListMassPayments(nil)

	assert.Error(t, err)
	assert.Nil(t, payments)

	customer.Links = Links{"mass-payments": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/mass-payments"}}
	payments, err = customer.ListMassPayments(nil)

	assert.Nil(t, err)
	assert.NotNil(t, payments)
}

func TestCustomerListTransfers(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "transfers.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	customer := &Customer{}
	customer.client = c

	transfers, err := customer.ListTransfers(nil)

	assert.Error(t, err)
	assert.Nil(t, transfers)

	customer.Links = Links{"transfers": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/transfers"}}
	transfers, err = customer.ListTransfers(nil)

	assert.Nil(t, err)
	assert.NotNil(t, transfers)
}

func TestCustomerReactivate(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "customer.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	customer := &Customer{}
	customer.client = c

	err := customer.Reactivate()

	assert.Error(t, err)

	customer.Links = Links{"reactivate": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}
	err = customer.Reactivate()

	assert.Nil(t, err)
}

func TestCustomerReceive(t *testing.T) {
	customer := &Customer{}

	res := customer.Receive()

	assert.False(t, res)

	customer.Links = Links{"receive": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/transfers"}}
	res = customer.Receive()

	assert.True(t, res)
}

func TestCustomerRetrieveBeneficialOwnership(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "beneficial-ownership.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	customer := &Customer{}
	customer.client = c

	ownership, err := customer.RetrieveBeneficialOwnership()

	assert.Error(t, err)
	assert.Nil(t, ownership)

	customer.Links = Links{"beneficial-owners": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/beneficial-owners"}}
	ownership, err = customer.RetrieveBeneficialOwnership()

	assert.Nil(t, err)
	assert.NotNil(t, ownership)
}

func TestCustomerRetrieveIAVToken(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "iav-token.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	customer := &Customer{}
	customer.client = c

	token, err := customer.RetrieveIAVToken()

	assert.Error(t, err)
	assert.Nil(t, token)

	customer.Links = Links{"self": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}
	token, err = customer.RetrieveIAVToken()

	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestCustomerRetryVerification(t *testing.T) {
	customer := &Customer{}

	res := customer.RetryVerification()

	assert.False(t, res)

	customer.Links = Links{"retry-verification": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}
	res = customer.RetryVerification()

	assert.True(t, res)
}

func TestCustomerSend(t *testing.T) {
	customer := &Customer{}

	res := customer.Send()

	assert.False(t, res)

	customer.Links = Links{"send": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/transfers"}}
	res = customer.Send()

	assert.True(t, res)
}

func TestCustomerSuspend(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "customer.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	customer := &Customer{}
	customer.client = c

	err := customer.Suspend()

	assert.Error(t, err)

	customer.Links = Links{"suspend": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}
	err = customer.Suspend()

	assert.Nil(t, err)
}

func TestCustomerUpdate(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "customer.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	customer := &Customer{}
	customer.client = c

	err := customer.Update(&CustomerRequest{
		FirstName: "Foo",
		LastName:  "Bar",
	})

	assert.Error(t, err)

	customer.Links = Links{"self": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}
	err = customer.Update(&CustomerRequest{
		FirstName: "Foo",
		LastName:  "Bar",
	})

	assert.Nil(t, err)
}

func TestCustomerVerifyBeneficialOwners(t *testing.T) {
	customer := &Customer{}

	res := customer.VerifyBeneficialOwners()

	assert.False(t, res)

	customer.Links = Links{"verify-beneficial-owners": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/beneficial-owners"}}
	res = customer.VerifyBeneficialOwners()

	assert.True(t, res)
}

func TestCustomerVerifyBusiness(t *testing.T) {
	customer := &Customer{}

	res := customer.VerifyBusiness()

	assert.False(t, res)

	customer.Links = Links{"verify-business-with-document": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/documents"}}
	res = customer.VerifyBusiness()

	assert.True(t, res)
}

func TestCustomerVerifyController(t *testing.T) {
	customer := &Customer{}

	res := customer.VerifyController()

	assert.False(t, res)

	customer.Links = Links{"verify-with-document": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/documents"}}
	res = customer.VerifyController()

	assert.True(t, res)
}

func TestCustomerVerifyControllerAndBusiness(t *testing.T) {
	customer := &Customer{}

	res := customer.VerifyControllerAndBusiness()

	assert.False(t, res)

	customer.Links = Links{"verify-controller-and-business-with-document": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/documents"}}
	res = customer.VerifyControllerAndBusiness()

	assert.True(t, res)
}
