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
}
