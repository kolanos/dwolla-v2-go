package dwolla

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var ErrNoID = errors.New("unable to extract ID")

// ACHDetails contains data sent to the bank account
type ACHDetails struct {
	Destination Addenda `json:"destination,omitempty"`
	Source      Addenda `json:"source,omitempty"`
}

// TransferACHDetails is an alias of ACHDetails for backwards compatibility
type TransferACHDetails = ACHDetails

// Addenda is a transfer addenda
type Addenda struct {
	Addenda TransferAddendaValues `json:"addenda,omitempty"`
}

// TransferAddenda is an alias of Addenda for backwards compatibility
type TransferAddenda = Addenda

// AddendaValues is the addenda values
type AddendaValues struct {
	Values []string `json:"values,omitempty"`
}

// TransferAddendaValues is an alias of AddendaValues for backwards compatibility
type TransferAddendaValues = AddendaValues

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

// Clearing is a transfer clearing schedule
type Clearing struct {
	Destination string `json:"destination,omitempty"`
	Source      string `json:"source,omitempty"`
}

// TransferClearing is an alias of Clearing for backwards compatiblity
type TransferClearing = Clearing

// MetaData represents key/value meta data
type MetaData map[string]interface{}

// Passport represents a passport
type Passport struct {
	Number  string `json:"number"`
	Country string `json:"country"`
}

// IDFromHREF takes an HREF link and returns the ID at the end of the HREF.
// This is useful for processing webhooks where you have an HREF, but need
// to make calls using this SDK, which expects bare IDs.
//
// If the input HREF is malformed, or this function is unable to extract the ID,
// ErrNoID will be returned.
func IDFromHREF(href string) (string, error) {
	lastIDX := strings.LastIndex(href, "/")
	if lastIDX < 0 {
		return "", ErrNoID
	}

	return href[lastIDX:], nil
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
	f, _ := os.Open(filepath.Clean(file))
	mr := &http.Response{Body: f, StatusCode: status}
	mc := &mockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{ExpiresIn: 3600, startTime: time.Now()}
	return c
}

func joinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}
