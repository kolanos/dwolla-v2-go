package dwolla

import (
	"fmt"
	"io"
)

// DocumentService is the document service interface
//
// see: https://docsv2.dwolla.com/#documents
type DocumentService interface {
	Retrieve(string) (*Document, error)
}

// DocumentServiceOp is an implementation of the document service
type DocumentServiceOp struct {
	client *Client
}

// Document is a dwolla verification document
type Document struct {
	Resource
	ID      string `json:"id"`
	Status  string `json:"status"`
	Type    string `json:"type"`
	Created string `json:"created"`
}

// Documents is a collection of dwolla documents
type Documents struct {
	Collection
	Embedded map[string][]Document `json:"_embedded"`
}

// DocumentRequest is a verification document request
type DocumentRequest struct {
	DocumentType string
	FileName     string
	File         io.Reader
}

// Retrieve retrieves a document matching the id
func (d *DocumentServiceOp) Retrieve(id string) (*Document, error) {
	var document Document

	if err := d.client.Get(fmt.Sprintf("documents/%s", id), nil, nil, &document); err != nil {
		return nil, err
	}

	document.client = d.client
	return &document, nil
}
