package dwolla

import (
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
	c := newMockClient(201, filepath.Join("testdata", "token.json"))
	c.Token = nil

	assert.NoError(t, c.RequestToken())
	assert.NotEmpty(t, c.Token.AccessToken)
}

func TestClientRoot(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "root.json"))

	res, err := c.Root()

	assert.Nil(t, err)
	assert.NotNil(t, res)
}
