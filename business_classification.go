package dwolla

import (
	"context"
	"fmt"
	"net/url"
)

// BusinessClassificationService is the business classification interface
type BusinessClassificationService interface {
	Retrieve(context.Context, string) (*BusinessClassification, error)
	List(context.Context, *url.Values) (*BusinessClassifications, error)
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

// BusinessClassifications is a business classification
type BusinessClassifications struct {
	Collection
	Embedded map[string][]BusinessClassification `json:"_embedded"`
}

// IndustryClassification is a industry subclassification
type IndustryClassification struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Retrieve retrieves a business classification matching the id
//
// see: https://docsv2.dwolla.com/#retrieve-a-business-classification
func (b *BusinessClassificationServiceOp) Retrieve(ctx context.Context, id string) (*BusinessClassification, error) {
	var classification BusinessClassification

	if err := b.client.Get(ctx, fmt.Sprintf("business-classifications/%s", id), nil, nil, &classification); err != nil {
		return nil, err
	}

	classification.client = b.client

	return &classification, nil
}

// List returns a collection of business classifications
//
// see: https://docsv2.dwolla.com/#list-business-classifications
func (b *BusinessClassificationServiceOp) List(ctx context.Context, params *url.Values) (*BusinessClassifications, error) {
	var classifications BusinessClassifications

	if err := b.client.Get(ctx, "business-classifications", params, nil, &classifications); err != nil {
		return nil, err
	}

	classifications.client = b.client

	for i := range classifications.Embedded["business-classifications"] {
		classifications.Embedded["business-classifications"][i].client = b.client
	}

	return &classifications, nil
}
