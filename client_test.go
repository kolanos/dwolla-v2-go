package dwolla

import (
	"os"
	"path/filepath"
	"testing"
	"time"

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
	c := newMockClient(201, filepath.Join("testdata", "token.json"))
	c.Token = nil

	assert.NoError(t, c.RequestToken())
	assert.NotEmpty(t, c.Token.AccessToken)
}

func TestClientRequestTokenLive(t *testing.T) {
	if os.Getenv("DWOLLA_API_KEY") == "" || os.Getenv("DWOLLA_API_SECRET") == "" {
		t.Skip("Test requires DWOLLA_API_KEY and DWOLLA_API_SECRET environment variables")
	}

	c := New(os.Getenv("DWOLLA_API_KEY"), os.Getenv("DWOLLA_API_SECRET"), Sandbox)

	assert.NoError(t, c.RequestToken())
	assert.NotEmpty(t, c.Token.AccessToken)
}

func TestClientRoot(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "root.json"))

	res, err := c.Root()

	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestClientRootLive(t *testing.T) {
	if os.Getenv("DWOLLA_API_KEY") == "" || os.Getenv("DWOLLA_API_SECRET") == "" {
		t.Skip("Test requires DWOLLA_API_KEY and DWOLLA_API_SECRET environment variables")
	}

	c := New(os.Getenv("DWOLLA_API_KEY"), os.Getenv("DWOLLA_API_SECRET"), Sandbox)

	res, err := c.Root()

	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestTokenExpired(t *testing.T) {
	token := &Token{}
	assert.True(t, token.Expired())

	token = &Token{ExpiresIn: 0, startTime: time.Now()}
	assert.True(t, token.Expired())

	token = &Token{ExpiresIn: 1, startTime: time.Now()}
	assert.False(t, token.Expired())

	token = &Token{ExpiresIn: 3600, startTime: time.Now()}
	assert.False(t, token.Expired())
}
