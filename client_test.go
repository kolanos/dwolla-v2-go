package dwolla

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClientEnvironment(t *testing.T) {
	productionClient := New("foobar", "barbaz", Production)

	assert.Equal(t, productionClient.APIURL(), ProductionAPIURL)
	assert.Equal(t, productionClient.AuthURL(), ProductionAuthURL)
	assert.Equal(t, productionClient.TokenURL(), ProductionTokenURL)

	sandboxClient := New("foobar", "barbaz", Sandbox)

	assert.Equal(t, sandboxClient.APIURL(), SandboxAPIURL)
	assert.Equal(t, sandboxClient.AuthURL(), SandboxAuthURL)
	assert.Equal(t, sandboxClient.TokenURL(), SandboxTokenURL)
}

func TestClientRequestToken(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "token.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)

	assert.NoError(t, c.RequestToken())
	assert.NotEmpty(t, c.Token.AccessToken)
}

func TestClientRoot(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "root.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	res, err := c.Root()

	assert.Nil(t, err)
	assert.NotNil(t, res)
}
