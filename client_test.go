package dwolla

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ctx = context.TODO()

func TestClientEnvironment(t *testing.T) {
	productionClient := New("foobar", "barbaz", Production)
	assert.Equal(t, ProductionAPIURL, productionClient.APIURL())
	assert.Equal(t, ProductionAuthURL, productionClient.AuthURL())
	assert.Equal(t, ProductionTokenURL, productionClient.TokenURL())

	sandboxClient := New("foobar", "barbaz", Sandbox)
	assert.Equal(t, SandboxAPIURL, sandboxClient.APIURL())
	assert.Equal(t, SandboxAuthURL, sandboxClient.AuthURL())
	assert.Equal(t, SandboxTokenURL, sandboxClient.TokenURL())

	testClient := New("foobar", "barbaz", "https://api-test.dwolla.com")
	assert.Equal(t, "https://api-test.dwolla.com", testClient.APIURL())
	assert.Equal(t, "https://api-test.dwolla.com/oauth/v2/authenticate", testClient.AuthURL())
	assert.Equal(t, "https://api-test.dwolla.com/token", testClient.TokenURL())
}

func TestClientRequestToken(t *testing.T) {
	c := newMockClient(201, filepath.Join("testdata", "token.json"))
	c.Token = nil

	assert.NoError(t, c.RequestToken(ctx))
	assert.NotEmpty(t, c.Token.AccessToken)
}

func TestClientRequestTokenLive(t *testing.T) {
	if os.Getenv("DWOLLA_API_KEY") == "" || os.Getenv("DWOLLA_API_SECRET") == "" {
		t.Skip("Test requires DWOLLA_API_KEY and DWOLLA_API_SECRET environment variables")
	}

	c := New(os.Getenv("DWOLLA_API_KEY"), os.Getenv("DWOLLA_API_SECRET"), Sandbox)

	assert.NoError(t, c.RequestToken(ctx))
	assert.NotEmpty(t, c.Token.AccessToken)
	assert.NotEmpty(t, c.Token.ExpiresIn)
}

func TestClientRoot(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "root.json"))

	res, err := c.Root(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestClientRootLive(t *testing.T) {
	if os.Getenv("DWOLLA_API_KEY") == "" || os.Getenv("DWOLLA_API_SECRET") == "" {
		t.Skip("Test requires DWOLLA_API_KEY and DWOLLA_API_SECRET environment variables")
	}

	c := New(os.Getenv("DWOLLA_API_KEY"), os.Getenv("DWOLLA_API_SECRET"), Sandbox)

	res, err := c.Root(ctx)

	assert.NoError(t, err)
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

func TestCreateCustomerAndFundingSource(t *testing.T) {
	t.Skip()

	if os.Getenv("DWOLLA_API_KEY") == "" || os.Getenv("DWOLLA_API_SECRET") == "" {
		t.Skip("Test requires DWOLLA_API_KEY and DWOLLA_API_SECRET environment variables")
	}

	c := New(os.Getenv("DWOLLA_API_KEY"), os.Getenv("DWOLLA_API_SECRET"), Sandbox)
	ts := time.Now().Unix()

	// Create customer
	customer, err := c.Customer.Create(ctx, &CustomerRequest{
		FirstName:    "Unit",
		LastName:     "Test",
		Email:        fmt.Sprintf("unittest_%d@example.com", ts),
		BusinessName: fmt.Sprintf("Unit Test %d", ts),
	})
	require.NoError(t, err)
	require.NotNil(t, customer)

	// Create funding source
	fundingSourceCreated, err := customer.CreateFundingSource(ctx, &FundingSourceRequest{
		RoutingNumber:   "222222226",
		AccountNumber:   "296797",
		BankAccountType: FundingSourceBankAccountTypeChecking,
		Name:            "Dingo Central Checking",
	})
	require.NoError(t, err)
	require.NotNil(t, fundingSourceCreated)

	// Initiate Micro Deposits
	microDeposit, err := fundingSourceCreated.InitiateMicroDeposits(ctx, "")
	require.NoError(t, err)
	require.NotNil(t, microDeposit)
	assert.Equal(t, MicroDepositStatusPending, microDeposit.Status)

	// Verify Micro Deposits
	fundingSource, err := c.FundingSource.Retrieve(ctx, fundingSourceCreated.ID)
	require.NoError(t, err)
	require.NotNil(t, fundingSource)
	err = fundingSource.VerifyMicroDeposits(ctx, &MicroDepositRequest{
		Amount1: Amount{Value: "0.03", Currency: "USD"},
		Amount2: Amount{Value: "0.09", Currency: "USD"},
	})
	require.NoError(t, err)

	// Check funding source's status
	fundingSource, err = c.FundingSource.Retrieve(ctx, fundingSourceCreated.ID)
	require.NoError(t, err)
	require.NotNil(t, fundingSource)
	assert.Equal(t, FundingSourceStatusVerified, fundingSource.Status)
}
