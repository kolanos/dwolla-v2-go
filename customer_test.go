package dwolla

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomerServiceCreate(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "customer.json"))

	res, err := c.Customer.Create(ctx, &CustomerRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "janedoe@nomail.com",
		Type:      CustomerTypeUnverified,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.ID, "FC451A7A-AE30-4404-AB95-E3553FCD733F")
	assert.Equal(t, res.FirstName, "Jane")
	assert.Equal(t, res.LastName, "Doe")
	assert.Equal(t, res.Email, "janedoe@nomail.com")
	assert.Equal(t, res.Type, CustomerTypeUnverified)
	assert.Equal(t, res.Status, CustomerStatusUnverified)
	assert.Equal(t, res.Created, "2015-09-03T23:56:10.023Z")
	assert.Equal(t, res.BusinessName, "Stratapro")
	assert.Equal(t, res.CorrelationID, "777")
}

func TestCustomerServiceCreateError(t *testing.T) {
	c := newMockClient(400, filepath.Join("testdata", "validation-error.json"))

	res, err := c.Customer.Create(ctx, &CustomerRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "janedoe@nomail.com",
		Type:      CustomerTypeUnverified,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCustomerServiceCreateError_DuplicateEmail(t *testing.T) {
	c := newMockClient(400, filepath.Join("testdata", "customer-create-duplicate-error.json"))

	res, err := c.Customer.Create(ctx, &CustomerRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "janedoe@nomail.com",
		Type:      CustomerTypeUnverified,
	})

	expectedErr := ValidationError{
		Code:    "ValidationError",
		Message: "Validation error(s) present. See embedded errors list for more details.",
		Embedded: HALErrors{
			"errors": []HALError{
				{
					Code:    "Duplicate",
					Message: "A customer with the specified email already exists.",
					Path:    "/email",
					Links: Links{
						"about": Link{
							Href:         "https://api-sandbox.dwolla.com/customers/8a1187aa-05f3-4425-a62b-9ec72d1a0c6a",
							Type:         "application/vnd.dwolla.v1.hal+json",
							ResourceType: "customer",
						},
					},
				},
			},
		},
	}

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, res)
}

func TestCustomerServiceRetrieve(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "customer.json"))

	res, err := c.Customer.Retrieve(ctx, "FC451A7A-AE30-4404-AB95-E3553FCD733F")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.ID, "FC451A7A-AE30-4404-AB95-E3553FCD733F")
	assert.Equal(t, res.FirstName, "Jane")
	assert.Equal(t, res.LastName, "Doe")
	assert.Equal(t, res.Email, "janedoe@nomail.com")
	assert.Equal(t, res.Type, CustomerTypeUnverified)
	assert.Equal(t, res.Status, CustomerStatusUnverified)
	assert.Equal(t, res.Created, "2015-09-03T23:56:10.023Z")
	assert.Equal(t, res.BusinessName, "Stratapro")
	assert.Equal(t, res.CorrelationID, "777")
}

func TestCustomerServiceRetrieveError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))
	res, err := c.Customer.Retrieve(ctx, "FC451A7A-AE30-4404-AB95-E3553FCD733F")

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCustomerServiceList(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "customers.json"))
	res, err := c.Customer.List(ctx, nil)

	assert.NoError(t, err)
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

func TestCustomerServiceListError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "document-not-found.json"))
	res, err := c.Customer.List(ctx, nil)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCustomerServiceUpdate(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "customer.json"))

	res, err := c.Customer.Update(ctx, "FC451A7A-AE30-4404-AB95-E3553FCD733F", &CustomerRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "janedoe@nomail.com",
		Type:      CustomerTypeUnverified,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.ID, "FC451A7A-AE30-4404-AB95-E3553FCD733F")
	assert.Equal(t, res.FirstName, "Jane")
	assert.Equal(t, res.LastName, "Doe")
	assert.Equal(t, res.Email, "janedoe@nomail.com")
	assert.Equal(t, res.Type, CustomerTypeUnverified)
	assert.Equal(t, res.Status, CustomerStatusUnverified)
	assert.Equal(t, res.Created, "2015-09-03T23:56:10.023Z")
	assert.Equal(t, res.BusinessName, "Stratapro")
	assert.Equal(t, res.CorrelationID, "777")
}

func TestCustomerServiceUpdateError(t *testing.T) {
	c := newMockClient(400, filepath.Join("testdata", "validation-error.json"))

	res, err := c.Customer.Update(ctx, "FC451A7A-AE30-4404-AB95-E3553FCD733F", &CustomerRequest{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "janedoe@nomail.com",
		Type:      CustomerTypeUnverified,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCustomerCertifyBeneficialOwnership(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "customer.json"))

	customer := &Customer{Resource: Resource{client: c, Links: Links{"certify-beneficial-ownership": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/beneficial-ownership"}}}}
	err := customer.CertifyBeneficialOwnership(ctx)
	assert.NoError(t, err)
}

func TestCustomerCertifyBeneficialOwnershipError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	customer := &Customer{Resource: Resource{client: c}}
	err := customer.CertifyBeneficialOwnership(ctx)
	assert.Error(t, err)

	customer = &Customer{Resource: Resource{client: c, Links: Links{"certify-beneficial-ownership": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/beneficial-ownership"}}}}
	err = customer.CertifyBeneficialOwnership(ctx)
	assert.Error(t, err)
}

func TestCustomerCreateDocument(t *testing.T) {
	c := newMockClient(201, filepath.Join("testdata", "document.json"))

	customer := &Customer{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}}}

	f, _ := os.Open(filepath.Join("testdata", "document-upload-success.png"))
	res, err := customer.CreateDocument(ctx, &DocumentRequest{
		Type:     DocumentTypePassport,
		FileName: f.Name(),
		File:     f,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestCustomerCreatedTime(t *testing.T) {
	customer := &Customer{Created: "2015-09-03T23:56:10.023Z"}
	assert.NotNil(t, customer.CreatedTime())
}

func TestCustomerCreateDocumentError(t *testing.T) {
	c := newMockClient(400, filepath.Join("testdata", "validation-error.json"))

	customer := &Customer{Resource: Resource{client: c}}
	f1, _ := os.Open(filepath.Join("testdata", "document-upload-success.png"))
	defer f1.Close()

	res, err := customer.CreateDocument(ctx, &DocumentRequest{
		Type:     DocumentTypePassport,
		FileName: f1.Name(),
		File:     f1,
	})

	assert.Error(t, err)
	assert.Nil(t, res)

	customer = &Customer{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}}}
	f2, _ := os.Open(filepath.Join("testdata", "document-upload-success.png"))
	res, err = customer.CreateDocument(ctx, &DocumentRequest{
		Type:     DocumentTypePassport,
		FileName: f2.Name(),
		File:     f2,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCustomerCreateBeneficialOwner(t *testing.T) {
	c := newMockClient(201, filepath.Join("testdata", "beneficial-owner.json"))

	customer := &Customer{Resource: Resource{client: c, Links: Links{"beneficial-owners": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/beneficial-owners"}}}}
	owner, err := customer.CreateBeneficialOwner(ctx, &BeneficialOwnerRequest{
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

	assert.NoError(t, err)
	assert.NotNil(t, owner)
}

func TestCustomerCreateBeneficialOwnerError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	customer := &Customer{Resource: Resource{client: c}}
	res, err := customer.CreateBeneficialOwner(ctx, &BeneficialOwnerRequest{
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

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "No beneficial owners resource link")
	assert.Nil(t, res)

	customer = &Customer{Resource: Resource{client: c, Links: Links{"beneficial-owners": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/beneficial-owners"}}}}
	res, err = customer.CreateBeneficialOwner(ctx, &BeneficialOwnerRequest{
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

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCustomerCreateFundingSource(t *testing.T) {
	c := newMockClient(201, filepath.Join("testdata", "funding-source.json"))

	customer := &Customer{Resource: Resource{client: c, Links: Links{"funding-sources": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/funding-sources"}}}}
	source, err := customer.CreateFundingSource(ctx, &FundingSourceRequest{
		RoutingNumber:   "1234567890",
		AccountNumber:   "1234567890",
		BankAccountType: FundingSourceBankAccountTypeChecking,
		Name:            "Test Checking Account",
	})

	assert.NoError(t, err)
	assert.NotNil(t, source)
}

func TestCustomerCreateFundingSourceError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	customer := &Customer{Resource: Resource{client: c}}
	res, err := customer.CreateFundingSource(ctx, &FundingSourceRequest{
		RoutingNumber:   "1234567890",
		AccountNumber:   "1234567890",
		BankAccountType: FundingSourceBankAccountTypeChecking,
		Name:            "Test Checking Account",
	})

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "No funding sources resource link")
	assert.Nil(t, res)

	customer = &Customer{Resource: Resource{client: c, Links: Links{"funding-sources": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/funding-sources"}}}}
	res, err = customer.CreateFundingSource(ctx, &FundingSourceRequest{
		RoutingNumber:   "1234567890",
		AccountNumber:   "1234567890",
		BankAccountType: FundingSourceBankAccountTypeChecking,
		Name:            "Test Checking Account",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCustomerDeactivate(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "customer.json"))

	customer := &Customer{Resource: Resource{client: c, Links: Links{"deactivate": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}}}
	err := customer.Deactivate(ctx)

	assert.NoError(t, err)
}

func TestCustomerDeactivateError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	customer := &Customer{Resource: Resource{client: c}}
	err := customer.Deactivate(ctx)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "No deactivate resource link")

	customer.Links = Links{"deactivate": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}
	err = customer.Deactivate(ctx)

	assert.Error(t, err)
}

func TestCustomerListBeneficialOwners(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "beneficial-owners.json"))

	customer := &Customer{Resource: Resource{client: c, Links: Links{"beneficial-owners": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/beneficial-owners"}}}}
	res, err := customer.ListBeneficialOwners(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestCustomerListBeneficialOwnersError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	customer := &Customer{Resource: Resource{client: c}}
	res, err := customer.ListBeneficialOwners(ctx)

	assert.Error(t, err)
	assert.Nil(t, res)

	customer = &Customer{Resource: Resource{client: c, Links: Links{"beneficial-owners": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/beneficial-owners"}}}}
	res, err = customer.ListBeneficialOwners(ctx)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCustomerListDocuments(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "documents.json"))

	customer := &Customer{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}}}
	res, err := customer.ListDocuments(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestCustomerListDocumentsError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	customer := &Customer{Resource: Resource{client: c}}
	res, err := customer.ListDocuments(ctx)

	assert.Error(t, err)
	assert.Nil(t, res)

	customer = &Customer{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}}}
	res, err = customer.ListDocuments(ctx)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCustomerListFundingSources(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "funding-sources.json"))

	customer := &Customer{Resource: Resource{client: c, Links: Links{"funding-sources": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/funding-sources"}}}}
	res, err := customer.ListFundingSources(ctx, true)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestCustomerListFundingSourcesError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	customer := &Customer{Resource: Resource{client: c}}
	res, err := customer.ListFundingSources(ctx, true)

	assert.Error(t, err)
	assert.Nil(t, res)

	customer = &Customer{Resource: Resource{client: c, Links: Links{"funding-sources": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/funding-sources"}}}}
	res, err = customer.ListFundingSources(ctx, true)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCustomerListMassPayments(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "mass-payments.json"))

	customer := &Customer{Resource: Resource{client: c, Links: Links{"mass-payments": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/mass-payments"}}}}
	res, err := customer.ListMassPayments(ctx, nil)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestCustomerListMassPaymentsError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	customer := &Customer{Resource: Resource{client: c}}
	res, err := customer.ListMassPayments(ctx, nil)

	assert.Error(t, err)
	assert.Nil(t, res)

	customer = &Customer{Resource: Resource{client: c, Links: Links{"mass-payments": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/mass-payments"}}}}
	res, err = customer.ListMassPayments(ctx, nil)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCustomerListTransfers(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "transfers.json"))

	customer := &Customer{Resource: Resource{client: c, Links: Links{"transfers": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/transfers"}}}}
	res, err := customer.ListTransfers(ctx, nil)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestCustomerListTransfersError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	customer := &Customer{Resource: Resource{client: c}}
	res, err := customer.ListTransfers(ctx, nil)

	assert.Error(t, err)
	assert.Nil(t, res)

	customer = &Customer{Resource: Resource{client: c, Links: Links{"transfers": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/transfers"}}}}
	res, err = customer.ListTransfers(ctx, nil)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCustomerReactivate(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "customer.json"))

	customer := &Customer{Resource: Resource{client: c}}
	err := customer.Reactivate(ctx)

	assert.Error(t, err)

	customer.Links = Links{"reactivate": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}
	err = customer.Reactivate(ctx)

	assert.NoError(t, err)
}

func TestCustomerReactivateError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	customer := &Customer{Resource: Resource{client: c}}
	err := customer.Reactivate(ctx)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "No reactivate resource link")

	customer.Links = Links{"reactivate": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}
	err = customer.Reactivate(ctx)

	assert.Error(t, err)
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
	c := newMockClient(200, filepath.Join("testdata", "beneficial-ownership.json"))

	customer := &Customer{Resource: Resource{client: c, Links: Links{"beneficial-owners": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/beneficial-owners"}}}}
	res, err := customer.RetrieveBeneficialOwnership(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestCustomerRetrieveBeneficialOwnershipError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	customer := &Customer{Resource: Resource{client: c}}
	res, err := customer.RetrieveBeneficialOwnership(ctx)

	assert.Error(t, err)
	assert.Nil(t, res)

	customer = &Customer{Resource: Resource{client: c, Links: Links{"beneficial-owners": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F/beneficial-owners"}}}}
	res, err = customer.RetrieveBeneficialOwnership(ctx)

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestCustomerRetrieveIAVToken(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "iav-token.json"))

	customer := &Customer{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}}}
	res, err := customer.RetrieveIAVToken(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestCustomerRetrieveIAVTokenError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	customer := &Customer{Resource: Resource{client: c}}
	res, err := customer.RetrieveIAVToken(ctx)

	assert.Error(t, err)
	assert.Nil(t, res)

	customer = &Customer{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}}}
	res, err = customer.RetrieveIAVToken(ctx)

	assert.Error(t, err)
	assert.Nil(t, res)
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
	c := newMockClient(200, filepath.Join("testdata", "customer.json"))

	customer := &Customer{Resource: Resource{client: c}}
	err := customer.Suspend(ctx)

	assert.Error(t, err)

	customer.Links = Links{"suspend": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}
	err = customer.Suspend(ctx)

	assert.NoError(t, err)
}

func TestCustomerSuspendError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	customer := &Customer{Resource: Resource{client: c}}
	err := customer.Suspend(ctx)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "No suspend resource link")

	customer.Links = Links{"suspend": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}
	err = customer.Suspend(ctx)

	assert.Error(t, err)
}

func TestCustomerUpdate(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "customer.json"))

	customer := &Customer{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}}}
	err := customer.Update(ctx, &CustomerRequest{
		FirstName: "Foo",
		LastName:  "Bar",
	})

	assert.NoError(t, err)
}

func TestCustomerUpdateError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	customer := &Customer{Resource: Resource{client: c}}
	err := customer.Update(ctx, &CustomerRequest{
		FirstName: "Foo",
		LastName:  "Bar",
	})

	assert.Error(t, err)

	customer = &Customer{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/customers/FC451A7A-AE30-4404-AB95-E3553FCD733F"}}}}
	err = customer.Update(ctx, &CustomerRequest{
		FirstName: "Foo",
		LastName:  "Bar",
	})

	assert.Error(t, err)
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
