package dwolla

import (
	"fmt"
	"net/http"
	"os"
	"time"
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

// String returns a string representation of the amount
func (a Amount) String() string {
	return fmt.Sprintf("%s %s", a.Value, a.Currency)
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

type mockHTTPClient struct {
	err error
	res *http.Response
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.res, m.err
}

func (m *mockHTTPClient) SetResponse(res *http.Response, err error) {
	m.res, m.err = res, err
}

func newRedirectResponse(url string) *http.Response {
	res := &http.Response{
		Status:     "302",
		StatusCode: 302,
		Header:     http.Header{},
	}
	res.Header.Set("Location", url)
	return res
}

func newMockClient(status int, file string) *Client {
	f, _ := os.Open(file)
	mr := &http.Response{Body: f, StatusCode: status}
	mc := &mockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{ExpiresIn: 3600, startTime: time.Now()}
	return c
}
