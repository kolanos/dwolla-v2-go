package dwolla

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const (
	// CustomerStatusDeactivated is when the customer has been deactivated
	CustomerStatusDeactivated CustomerStatus = "deactivated"
	// CustomerStatusDocument is when the customer needs verification document
	CustomerStatusDocument CustomerStatus = "document"
	// CustomerStatusReactivated is when a deactivated customer is reactivated
	CustomerStatusReactivated CustomerStatus = "reactivated"
	// CustomerStatusRetry is when the customer needs to retry verification
	CustomerStatusRetry CustomerStatus = "retry"
	// CustomerStatusSuspended is when the customer has been suspended
	CustomerStatusSuspended CustomerStatus = "suspended"
	// CustomerStatusUnverified is when the customer is unverified
	CustomerStatusUnverified CustomerStatus = "unverified"
	// CustomerStatusVerified is when the customer is verified
	CustomerStatusVerified CustomerStatus = "verified"
)

const (
	// CustomerTypeBusiness is when the customer is a business
	CustomerTypeBusiness CustomerType = "business"
	// CustomerTypePersonal is when the customer is an individual
	CustomerTypePersonal CustomerType = "personal"
	// CustomerTypeReceiveOnly is when the customer can only receive funds
	CustomerTypeReceiveOnly CustomerType = "receive-only"
	// CustomerTypeUnverified is when the customer is unverified
	CustomerTypeUnverified CustomerType = "unverified"
)

// CustomerService is the customerservice interface
//
// see: https://docsv2.dwolla.com/#customers
type CustomerService interface {
	Create(*CustomerRequest) (*Customer, error)
	List(*url.Values) (*Customers, error)
	Retrieve(string) (*Customer, error)
	Update(string, *CustomerRequest) (*Customer, error)
}

// CustomerServiceOp is an implementation of the customer service interface
type CustomerServiceOp struct {
	client *Client
}

// Controller is a controller of a business
type Controller struct {
	FirstName   string   `json:"firstName,omitempty"`
	LastName    string   `json:"lastName,omitempty"`
	Title       string   `json:"title,omitempty"`
	DateOfBirth string   `json:"dateOfBirth,omitempty"`
	SSN         string   `json:"ssn,omitempty"`
	Address     Address  `json:"address,omitempty"`
	Passport    Passport `json:"passport,omitempty"`
}

// CustomerStatus is the customer's status
type CustomerStatus string

// CustomerType is the customer's type
type CustomerType string

// Customer is a dwolla customer
type Customer struct {
	Resource
	ID           string         `json:"id"`
	FirstName    string         `json:"firstName"`
	LastName     string         `json:"lastName"`
	Email        string         `json:"email"`
	Type         CustomerType   `json:"type"`
	Status       CustomerStatus `json:"status"`
	Created      string         `json:"created"`
	Address1     string         `json:"address1"`
	Address2     string         `json:"address2"`
	City         string         `json:"city"`
	State        string         `json:"state"`
	PostalCode   string         `json:"postalCode"`
	BusinessName string         `json:"businessName"`
	Controller   Controller     `json:"controller"`
}

// Customers is a collection of customers
type Customers struct {
	Collection
	Embedded map[string][]Customer `json:"_embedded"`
}

// CustomerRequest is a customer create/update request
//
// We don't just use the Customer struct here because there are fields that
// are not returned by the Dwolla API. As such, we don't want fields to be
// unset during marshaling.
type CustomerRequest struct {
	FirstName              string             `json:"firstName,omitempty"`
	LastName               string             `json:"lastName,omitempty"`
	Email                  string             `json:"email,omitempty"`
	IPAddress              string             `json:"ipAddress,omitempty"`
	Type                   CustomerType       `json:"type,omitempty"`
	Status                 CustomerStatus     `json:"status,omitempty"`
	DateOfBirth            string             `json:"dateOfBirth,omitempty"`
	SSN                    string             `json:"ssn,omitempty"`
	Phone                  string             `json:"phone,omitempty"`
	Address1               string             `json:"address1,omitempty"`
	Address2               string             `json:"address2,omitempty"`
	City                   string             `json:"city,omitempty"`
	State                  string             `json:"state,omitempty"`
	PostalCode             string             `json:"postalCode,omitempty"`
	BusinessClassification string             `json:"businessClassification,omitempty"`
	BusinessType           string             `json:"businessType,omitempty"`
	BusinessName           string             `json:"businessName,omitempty"`
	DoingBusinessAs        string             `json:"doingBusinessAs,omitempty"`
	EIN                    string             `json:"ein,omitempty"`
	Website                string             `json:"website,omitempty"`
	Controller             *ControllerRequest `json:"controller,omitempty"`
}

// ControllerRequest is a controller of a business create/update request
type ControllerRequest struct {
	FirstName   string    `json:"firstName,omitempty"`
	LastName    string    `json:"lastName,omitempty"`
	Title       string    `json:"title,omitempty"`
	DateOfBirth string    `json:"dateOfBirth,omitempty"`
	SSN         string    `json:"ssn,omitempty"`
	Address     Address   `json:"address,omitempty"`
	Passport    *Passport `json:"passport,omitempty"`
}

// IAVToken is a instant account verification token
type IAVToken struct {
	Resource
	Token string `json:"token"`
}

// Create creates a dwolla customer
func (c *CustomerServiceOp) Create(body *CustomerRequest) (*Customer, error) {
	var customer Customer

	if err := c.client.Post("customers", body, nil, &customer); err != nil {
		return nil, err
	}

	customer.client = c.client

	return &customer, nil
}

// List returns a collection of customers
//
// see: https://docsv2.dwolla.com/#list-and-search-customers
func (c *CustomerServiceOp) List(params *url.Values) (*Customers, error) {
	var customers Customers

	if err := c.client.Get("customers", params, nil, &customers); err != nil {
		return nil, err
	}

	customers.client = c.client

	for i := range customers.Embedded["customers"] {
		customers.Embedded["customers"][i].client = c.client
	}

	return &customers, nil
}

// Retrieve retrieves a customer matching the id
//
// see: https://docsv2.dwolla.com/#retrieve-a-customer
func (c *CustomerServiceOp) Retrieve(id string) (*Customer, error) {
	var customer Customer

	if err := c.client.Get(fmt.Sprintf("customers/%s", id), nil, nil, &customer); err != nil {
		return nil, err
	}

	customer.client = c.client
	return &customer, nil
}

// Update updates a dwolla customer matching the id
//
// see: https://docsv2.dwolla.com/#update-a-customer
func (c *CustomerServiceOp) Update(id string, body *CustomerRequest) (*Customer, error) {
	var customer Customer

	if err := c.client.Post(fmt.Sprintf("customers/%s", id), body, nil, &customer); err != nil {
		return nil, err
	}

	customer.client = c.client

	return &customer, nil
}

// CertifyBeneficialOwnership certifies beneficial ownership
//
// see: https://docsv2.dwolla.com/#certify-beneficial-ownership
func (c *Customer) CertifyBeneficialOwnership() error {
	if _, ok := c.Links["certify-beneficial-ownership"]; !ok {
		return errors.New("No certify beneficial ownership resource link")
	}

	request := &BeneficialOwnershipRequest{Status: CertificationStatusCertified}

	return c.client.Post(c.Links["certify-beneficial-ownership"].Href, request, nil, nil)
}

// CreateDocument uploads a verification document for the customer
//
// see: https://docsv2.dwolla.com/#create-a-document
func (c *Customer) CreateDocument(body *DocumentRequest) (*Document, error) {
	var document Document

	if _, ok := c.Links["self"]; !ok {
		return nil, errors.New("No self resource link")
	}

	if err := c.client.Upload(fmt.Sprintf("%s/documents", c.Links["self"].Href), body.Type, body.FileName, body.File, &document); err != nil {
		return nil, err
	}

	document.client = c.client

	return &document, nil
}

// CreatedTime returns the created value as time.Time
func (c Customer) CreatedTime() time.Time {
	t, _ := time.Parse(time.RFC3339, c.Created)
	return t
}

// CreateBeneficialOwner creates a beneficial owner for the customer
//
// see: https://docsv2.dwolla.com/#create-a-beneficial-owner
func (c *Customer) CreateBeneficialOwner(body *BeneficialOwnerRequest) (*BeneficialOwner, error) {
	var owner BeneficialOwner

	if _, ok := c.Links["beneficial-owners"]; !ok {
		return nil, errors.New("No beneficial owners resource link")
	}

	if err := c.client.Post(c.Links["beneficial-owners"].Href, body, nil, &owner); err != nil {
		return nil, err
	}

	owner.client = c.client

	return &owner, nil
}

// CreateFundingSource creates a funding source for the customer
//
// see: https://docsv2.dwolla.com/#create-a-funding-source-for-a-customer
func (c *Customer) CreateFundingSource(body *FundingSourceRequest) (*FundingSource, error) {
	var source FundingSource

	if _, ok := c.Links["funding-sources"]; !ok {
		return nil, errors.New("No funding sources resource link")
	}

	if err := c.client.Post(c.Links["funding-sources"].Href, body, nil, &source); err != nil {
		return nil, err
	}

	source.client = c.client

	return &source, nil
}

// Deactivate deactivates a dwolla customer
func (c *Customer) Deactivate() error {
	if _, ok := c.Links["deactivate"]; !ok {
		return errors.New("No deactivate resource link")
	}

	request := &CustomerRequest{Status: CustomerStatusDeactivated}

	return c.client.Post(c.Links["deactivate"].Href, request, nil, c)
}

// ListBeneficialOwners returns the customer's beneficial owners
//
// see: https://docsv2.dwolla.com/#list-beneficial-owners
func (c *Customer) ListBeneficialOwners() (*BeneficialOwners, error) {
	var owners BeneficialOwners

	if _, ok := c.Links["beneficial-owners"]; !ok {
		return nil, errors.New("No beneficial owners resource link")
	}

	if err := c.client.Get(c.Links["beneficial-owners"].Href, nil, nil, &owners); err != nil {
		return nil, err
	}

	owners.client = c.client

	for i := range owners.Embedded["beneficial-owners"] {
		owners.Embedded["beneficial-owners"][i].client = c.client
	}

	return &owners, nil
}

// ListDocuments returns documents for customer
//
// see: https://docsv2.dwolla.com/#list-documents
func (c *Customer) ListDocuments() (*Documents, error) {
	var documents Documents

	if _, ok := c.Links["self"]; !ok {
		return nil, errors.New("No self resource link")
	}

	if err := c.client.Get(fmt.Sprintf("%s/documents", c.Links["self"].Href), nil, nil, &documents); err != nil {
		return nil, err
	}

	documents.client = c.client

	for i := range documents.Embedded["documents"] {
		documents.Embedded["documents"][i].client = c.client
	}

	return &documents, nil
}

// ListFundingSources returns the customer's funding sources
//
// see: https://docsv2.dwolla.com/#list-funding-sources-for-a-customer
func (c *Customer) ListFundingSources(removed bool) (*FundingSources, error) {
	var sources FundingSources

	if _, ok := c.Links["funding-sources"]; !ok {
		return nil, errors.New("No funding sources resource link")
	}

	params := &url.Values{}
	params.Add("removed", strconv.FormatBool(removed))

	if err := c.client.Get(c.Links["funding-sources"].Href, params, nil, &sources); err != nil {
		return nil, err
	}

	sources.client = c.client

	for i := range sources.Embedded["funding-sources"] {
		sources.Embedded["funding-sources"][i].client = c.client
	}

	return &sources, nil
}

// ListMassPayments returns the customer's mass payments
//
// see: https://docsv2.dwolla.com/#list-mass-payments-for-a-customer
func (c *Customer) ListMassPayments(params *url.Values) (*MassPayments, error) {
	var payments MassPayments

	if _, ok := c.Links["mass-payments"]; !ok {
		return nil, errors.New("No mass payments resource link")
	}

	if err := c.client.Get(c.Links["mass-payments"].Href, params, nil, &payments); err != nil {
		return nil, err
	}

	payments.client = c.client

	for i := range payments.Embedded["mass-payments"] {
		payments.Embedded["mass-payments"][i].client = c.client
	}

	return &payments, nil
}

// ListTransfers returns the customer's transfers
//
// see: https://docsv2.dwolla.com/#list-and-search-transfers-for-a-customer
func (c *Customer) ListTransfers(params *url.Values) (*Transfers, error) {
	var transfers Transfers

	if _, ok := c.Links["transfers"]; !ok {
		return nil, errors.New("No transfers resource link")
	}

	if err := c.client.Get(c.Links["transfers"].Href, params, nil, &transfers); err != nil {
		return nil, err
	}

	transfers.client = c.client

	for i := range transfers.Embedded["transfers"] {
		transfers.Embedded["transfers"][i].client = c.client
	}

	return &transfers, nil
}

// Reactivate reactivates a deactivated dwolla customer
func (c *Customer) Reactivate() error {
	if _, ok := c.Links["reactivate"]; !ok {
		return errors.New("No reactivate resource link")
	}

	request := &CustomerRequest{Status: CustomerStatusReactivated}

	return c.client.Post(c.Links["reactivate"].Href, request, nil, c)
}

// Receive returns true if customer can receive transfers
func (c *Customer) Receive() bool {
	_, ok := c.Links["receive"]
	return ok
}

// RetrieveBeneficialOwnership retrieves the customer's beneficial ownership status
func (c *Customer) RetrieveBeneficialOwnership() (*BeneficialOwnership, error) {
	var ownership BeneficialOwnership

	if _, ok := c.Links["beneficial-owners"]; !ok {
		return nil, errors.New("No beneficial owners resource link")
	}

	if err := c.client.Get(fmt.Sprintf("%s/beneficial-ownership", c.Links["self"].Href), nil, nil, &ownership); err != nil {
		return nil, err
	}

	ownership.client = c.client

	return &ownership, nil
}

// RetrieveIAVToken retrieves an instant account activation token
func (c *Customer) RetrieveIAVToken() (*IAVToken, error) {
	var token IAVToken

	if _, ok := c.Links["self"]; !ok {
		return nil, errors.New("No self resource link")
	}

	if err := c.client.Post(fmt.Sprintf("%s/iav-token", c.Links["self"].Href), nil, nil, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

// RetryVerification returns true if customer needs to retry verification
func (c *Customer) RetryVerification() bool {
	_, ok := c.Links["retry-verification"]
	return ok
}

// Send returns true if customer can send transfers
func (c *Customer) Send() bool {
	_, ok := c.Links["send"]
	return ok
}

// Suspend suspends a dwolla customer
func (c *Customer) Suspend() error {
	if _, ok := c.Links["suspend"]; !ok {
		return errors.New("No suspend resource link")
	}

	request := &CustomerRequest{Status: CustomerStatusSuspended}

	return c.client.Post(c.Links["suspend"].Href, request, nil, c)
}

// Update updates a dwolla customer
//
// see: https://docsv2.dwolla.com/#update-a-customer
func (c *Customer) Update(body *CustomerRequest) error {
	if _, ok := c.Links["self"]; !ok {
		return errors.New("No self resource link")
	}

	return c.client.Post(c.Links["self"].Href, body, nil, c)
}

// VerifyBeneficialOwners returns true if beneficial owners needed
func (c *Customer) VerifyBeneficialOwners() bool {
	_, ok := c.Links["verify-beneficial-owners"]
	return ok
}

// VerifyBusiness returns true if business needs verification document
func (c *Customer) VerifyBusiness() bool {
	_, ok := c.Links["verify-business-with-document"]
	return ok
}

// VerifyController returns true if controller needs verification document
func (c *Customer) VerifyController() bool {
	_, ok := c.Links["verify-with-document"]
	return ok
}

// VerifyControllerAndBusiness returns true if controller and business need verification document
func (c *Customer) VerifyControllerAndBusiness() bool {
	_, ok := c.Links["verify-controller-and-business-with-document"]
	return ok
}
