package dwolla

import (
	"io"
)

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
