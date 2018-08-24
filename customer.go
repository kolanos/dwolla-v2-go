package dwolla

import (
	"fmt"
	"net/url"
	"strconv"
)

// CustomerService is the customerservice interface
// see: https://docsv2.dwolla.com/#customers
type CustomerService interface {
	Get(string) (*Customer, error)
	List(*url.Values) (*Customers, error)
}

// CustomerServiceOp is an implementation of the customer service interface
type CustomerServiceOp struct {
	client *Client
}

// Controller is a controller of a business
type Controller struct {
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	Title       string   `json:"title"`
	DateOfBirth string   `json:"dateOfBirth,omitempty"`
	SSN         string   `json:"ssn,omitempty"`
	Address     Address  `json:"address"`
	Passport    Passport `json:"passport,omitempty"`
}

// Passport is a controller's passport
type Passport struct {
	Number  string `json:"number"`
	Country string `json:"country"`
}

// CustomerStatus is the customer's status
type CustomerStatus string

const (
	// CustomerStatusDeactivated is when the customer has been deactivated
	CustomerStatusDeactivated CustomerStatus = "deactivated"
	// CustomerStatusDocument is when the customer needs verification document
	CustomerStatusDocument CustomerStatus = "document"
	// CustomerStatusRetry is when the customer needs to retry verification
	CustomerStatusRetry CustomerStatus = "retry"
	// CustomerStatusSuspended is when the customer has been suspended
	CustomerStatusSuspended CustomerStatus = "suspended"
	// CustomerStatusUnverified is when the customer is unverified
	CustomerStatusUnverified CustomerStatus = "unverified"
	// CustomerStatusVerified is when the customer is verified
	CustomerStatusVerified CustomerStatus = "verified"
)

// CustomerType is the customer's type
type CustomerType string

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

// CustomerCreate is a customer create request
type CustomerCreate struct {
	FirstName              string       `json:"firstName"`
	LastName               string       `json:"lastName"`
	Email                  string       `json:"email"`
	IPAddress              string       `json:"ipAddress,omitempty"`
	Type                   CustomerType `json:"type"`
	DateOfBirth            string       `json:"dateOfBirth,omitempty"`
	SSN                    string       `json:"ssn,omitempty"`
	Phone                  string       `json:"phone,omitempty"`
	Address1               string       `json:"address1,omitempty"`
	Address2               string       `json:"address2,omitempty"`
	City                   string       `json:"city,omitempty"`
	State                  string       `json:"state,omitempty"`
	PostalCode             string       `json:"postalCode,omitempty"`
	BusinessClassification string       `json:businessClassification,omitempty"`
	BusinessType           string       `json:"businessType,omitempty"`
	BusinessName           string       `json:"businessName,omitempty"`
	DoingBusinessAs        string       `json:"doingBusinessAs,omitempty"`
	EIN                    string       `json:"ein,omitempty"`
	Website                string       `json:"website,omitempty"`
}

// Get returns a customer matching the id
func (c *CustomerServiceOp) Get(id string) (*Customer, error) {
	var customer Customer

	if err := c.client.Get(fmt.Sprintf("customers/%s", id), nil, nil, &customer); err != nil {
		return nil, err
	}

	customer.client = c.client
	return &customer, nil
}

// List returns a collection of customers
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

// ListFundingSources returns the customer's funding sources
func (c *Customer) ListFundingSources(removed bool) (*FundingSources, error) {
	var sources FundingSources

	params := &url.Values{}
	params.Add("removed", strconv.FormatBool(removed))

	if err := c.client.Get(c.Links["funding-sources"].HREF, params, nil, &sources); err != nil {
		return nil, err
	}

	sources.client = c.client

	for i := range sources.Embedded["funding-sources"] {
		sources.Embedded["funding-sources"][i].client = c.client
	}

	return &sources, nil
}

// ListTransfers returns the customer's transfers
func (c *Customer) ListTransfers() (*Transfers, error) {
	var transfers Transfers

	if err := c.client.Get(c.Links["transfers"].HREF, nil, nil, &transfers); err != nil {
		return nil, err
	}

	transfers.client = c.client

	for i := range transfers.Embedded["transfers"] {
		transfers.Embedded["transfers"][i].client = c.client
	}

	return &transfers, nil
}
