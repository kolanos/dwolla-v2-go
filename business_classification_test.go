package dwolla

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBusinessClassificationRetrieve(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "business-classification.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	res, err := c.BusinessClassification.Retrieve("9ed3cf58-7d6f-11e3-81a4-5404a6144203")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.ID, "9ed3a866-7d6f-11e3-a0ce-5404a6144203")
	assert.Equal(t, res.Name, "Entertainment and media")
}

func TestBusinessClassificationList(t *testing.T) {
	f, _ := os.Open(filepath.Join("testdata", "business-classifications.json"))
	mr := &http.Response{Body: f, StatusCode: 200}
	mc := &MockHTTPClient{err: nil, res: mr}

	c := NewWithHTTPClient("foobar", "barbaz", Sandbox, mc)
	c.Token = &Token{}

	res, err := c.BusinessClassification.List(nil)

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.Total, 27)
}
