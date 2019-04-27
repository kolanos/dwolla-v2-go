package dwolla

import (
	"context"
	"fmt"
	"io"
)

const (
	// DocumentScanIDTypeNotSupported is when the scanned I.D. type is not
	// supported
	DocumentScanIDTypeNotSupported DocumentFailureReason = "ScanIdTypeNotSupported"
	// DocumentNameMismatch is when the scanned document name does not match
	DocumentScanNameMismatch DocumentFailureReason = "ScanNameMismatch"
	// DocumentScanNotReadable is when the scanned document is not readable
	DocumentScanNotReadable DocumentFailureReason = "ScanNotReadable"
	// DocumentScanNotUploaded is when the scanned document failed to upload
	DocumentScanNotUploaded DocumentFailureReason = "ScanNotUploaded"
	// DocumentScanFailedOther is when the scanned document was rejected for
	// another reason
	DocumentScanFailedOther DocumentFailureReason = "ScanFailedOther"
	// DocumentFailedOther is when the document was rejected for another
	// reason
	DocumentFailedOther DocumentFailureReason = "FailedOther"
)

const (
	// DocumentStatusPending is when the document is pending review
	DocumentStatusPending DocumentStatus = "pending"
	// DocumentStatusReviewed is when the document has been reviewed
	DocumentStatusReviewed DocumentStatus = "reviewed"
)

const (
	// DocumentTypePassport is a passport
	DocumentTypePassport DocumentType = "passport"
	// DocumentTypeLicense is a state-issued driver's license
	DocumentTypeLicense DocumentType = "license"
	// DocumentTypeIDCard is a U.S. government issued photo I.D. card
	DocumentTypeIDCard DocumentType = "idCard"
	// DocumentTypeOther is an EIN Letter / IRS-issued SS4 Confirmation Letter
	DocumentTypeOther DocumentType = "other"
)

// DocumentService is the document service interface
//
// see: https://docsv2.dwolla.com/#documents
type DocumentService interface {
	Retrieve(context.Context, string) (*Document, error)
}

// DocumentServiceOp is an implementation of the document service
type DocumentServiceOp struct {
	client *Client
}

// Document is a dwolla verification document
type Document struct {
	Resource
	ID            string                `json:"id"`
	Status        DocumentStatus        `json:"status"`
	Type          DocumentType          `json:"type"`
	Created       string                `json:"created"`
	FailureReason DocumentFailureReason `json:"failureReason"`
}

// DocumentFailureReason is the reason document verification failed
type DocumentFailureReason string

// Documents is a collection of dwolla documents
type Documents struct {
	Collection
	Embedded map[string][]Document `json:"_embedded"`
}

// DocumentStatus is the status of the document
type DocumentStatus string

// DocumentType is the type of document
type DocumentType string

// DocumentRequest is a verification document request
type DocumentRequest struct {
	Type     DocumentType
	FileName string
	File     io.Reader
}

// Retrieve retrieves a document matching the id
func (d *DocumentServiceOp) Retrieve(ctx context.Context, id string) (*Document, error) {
	var document Document

	if err := d.client.Get(ctx, fmt.Sprintf("documents/%s", id), nil, nil, &document); err != nil {
		return nil, err
	}

	document.client = d.client

	return &document, nil
}
