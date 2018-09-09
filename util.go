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
	Value    string   `json:"value"`
	Currency Currency `json:"currency"`
}

// Currency represents the monetary currency
type Currency string

const (
	// USD is U.S. dollars
	USD Currency = "usd"
)

// MetaData represents key/value meta data
type MetaData map[string]interface{}

// Passport represents a passport
type Passport struct {
	Number  string `json:"number"`
	Country string `json:"country"`
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

// NewRedirectResponse returns a http redirect response
func NewRedirectResponse(url string) *http.Response {
	res := &http.Response{
		Status:     "302",
		StatusCode: 302,
		Header:     http.Header{},
	}
	res.Header.Set("Location", url)
	return res
}
