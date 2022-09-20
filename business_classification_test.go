package dwolla

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBusinessClassificationServiceRetrieve(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "business-classification.json"))
	res, err := c.BusinessClassification.Retrieve(ctx, "9ed3cf58-7d6f-11e3-81a4-5404a6144203")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.ID, "9ed3a866-7d6f-11e3-a0ce-5404a6144203")
	assert.Equal(t, res.Name, "Entertainment and media")
}

func TestBusinessClassificationServiceRetrieveError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))
	res, err := c.BusinessClassification.Retrieve(ctx, "9ed3cf58-7d6f-11e3-81a4-5404a6144203")

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestBusinessClassificationServiceList(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "business-classifications.json"))
	res, err := c.BusinessClassification.List(ctx, nil)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.Total, 27)
}

func TestBusinessClassificationServiceListError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))
	res, err := c.BusinessClassification.List(ctx, nil)

	assert.Error(t, err)
	assert.Nil(t, res)
}
