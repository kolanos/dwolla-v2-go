package dwolla

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventServiceRetrieve(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "event.json"))
	res, err := c.Event.Retrieve(ctx, "9ed3cf58-7d6f-11e3-81a4-5404a6144203")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.ID, "81f6e13c-557c-4449-9331-da5c65e61095")
}

func TestEventServiceRetrieveError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))
	res, err := c.Event.Retrieve(ctx, "9ed3cf58-7d6f-11e3-81a4-5404a6144203")

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestEventServiceList(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "events.json"))
	res, err := c.Event.List(ctx, nil)

	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.Equal(t, res.Total, 3)
}

func TestEventServiceListError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))
	res, err := c.Event.List(ctx, nil)

	assert.Error(t, err)
	assert.Nil(t, res)
}
