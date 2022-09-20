package dwolla

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentServiceRetrieve(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "document.json"))
	res, err := c.Document.Retrieve(ctx, "9ed3cf58-7d6f-11e3-81a4-5404a6144203")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.ID, "56502f7a-fa59-4a2f-8579-0f8bc9d7b9cc")
}

func TestDocumentServiceRetrieveError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))
	res, err := c.Document.Retrieve(ctx, "9ed3cf58-7d6f-11e3-81a4-5404a6144203")

	assert.Error(t, err)
	assert.Nil(t, res)
}
