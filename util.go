package dwolla

import (
	"net/http"
)

// Address represents a street address
type Address struct {
	Address1            string `json:"address1"`
	Address2            string `json:"address2,omitempty"`
	Address3            string `json:"address3,omitempty"`
	City                string `json:"city"`
	StateProvinceRegion string `json:"stateProvinceRegion"`
	PostalCode          string `json:"postalCode,omitempty"`
	Country             string `json:"country"`
}

// Amount is a monetary value
type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

// MockHTTPClient mocks an http client
type MockHTTPClient struct {
	err error
	res *http.Response
}

// Do mocks an http request/response
func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.res, m.err
}

// SetResponse sets the mocked response
func (m *MockHTTPClient) SetResponse(res *http.Response, err error) {
	m.res, m.err = res, err
}
