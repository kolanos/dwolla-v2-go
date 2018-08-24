package dwolla

import (
	"fmt"
	"net/url"
)

// BusinessClassificationService is the business classification interface
type BusinessClassificationService interface {
	Get(string) (*BusinessClassification, error)
	List(*url.Values) (*BusinessClassifications, error)
}

// BusinessClassificationServiceOp is an implementation of the business
// classification interface
type BusinessClassificationServiceOp struct {
	client *Client
}

// BusinessClassification is a business industry type
type BusinessClassification struct {
	Resource
	ID       string                              `json:"id"`
	Name     string                              `json:"name"`
	Embedded map[string][]IndustryClassification `json:"_embedded"`
}

type BusinessClassifications struct {
	Collection
	Embedded map[string][]BusinessClassification `json:"_embedded"`
}

// IndustryClassification is a industry subclassification
type IndustryClassification struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Get returns a business classification matching the id
func (b *BusinessClassificationServiceOp) Get(id string) (*BusinessClassification, error) {
	var classification BusinessClassification

	if err := b.client.Get(fmt.Sprintf("business-classifications/%s", id), nil, nil, &classification); err != nil {
		return nil, err
	}

	classification.client = b.client
	return &classification, nil
}

// List returns a collection of business classifications
func (b *BusinessClassificationServiceOp) List(params *url.Values) (*BusinessClassifications, error) {
	var classifications BusinessClassifications

	if err := b.client.Get("business-classifications", params, nil, &classifications); err != nil {
		return nil, err
	}

	classifications.client = b.client

	for i := range classifications.Embedded["business-classifications"] {
		classifications.Embedded["business-classifications"][i].client = b.client
	}

	return &classifications, nil
}
