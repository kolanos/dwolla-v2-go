package dwolla

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnDemandAuthorizationServiceCreate(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "on-demand-authorization.json"))
	res, err := c.OnDemandAuthorization.Create(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestOnDemandAuthorizationServiceCreateError(t *testing.T) {
	c := newMockClient(400, filepath.Join("testdata", "validation-error.json"))
	res, err := c.OnDemandAuthorization.Create(ctx)

	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestOnDemandAuthorizationServiceRetrieve(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "on-demand-authorization.json"))
	res, err := c.OnDemandAuthorization.Retrieve(ctx, "30e7c028-0bdf-e511-80de-0aa34a9b2388")

	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, res.BodyText, "I agree that future payments to Company ABC inc. will be processed by the Dwolla payment system from the selected account above. In order to cancel this authorization, I will change my payment settings within my Company ABC inc. account.")
	assert.Equal(t, res.ButtonText, "Agree & Continue")
}

func TestOnDemandAuthorizationServiceRetrieveError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))
	res, err := c.OnDemandAuthorization.Retrieve(ctx, "30e7c028-0bdf-e511-80de-0aa34a9b2388")

	assert.Error(t, err)
	assert.Nil(t, res)
}
