package dwolla

import (
	"errors"
	"fmt"
)

const (
	// BeneficialOwnerStatusDocument is when the beneficial owner needs verification document
	BeneficialOwnerStatusDocument BeneficialOwnerStatus = "document"

	// BeneficialOwnerStatusIncomplete is when the beneficial owner is incomplete
	BeneficialOwnerStatusIncomplete BeneficialOwnerStatus = "incomplete"

	// BeneficialOwnerStatusVerified is when the beneficial owner is verified
	BeneficialOwnerStatusVerified BeneficialOwnerStatus = "verified"
)

const (
	// CertificationStatusCertified is when the ownership status is certified
	CertificationStatusCertified CertificationStatus = "certified"

	// CertificationStatusRecertify is when the ownership status needs
	// to be recertified
	CertificationStatusRecertify CertificationStatus = "recertify"

	// CertificationStatusUncertified is when the ownership status is uncertified
	CertificationStatusUncertified CertificationStatus = "uncertified"
)

// BeneficialOwnerService is the beneficial owner service interface
//
// see: https://docsv2.dwolla.com/#beneficial-owners
type BeneficialOwnerService interface {
	Remove(string) error
	Retrieve(string) (*BeneficialOwner, error)
	Update(string, *BeneficialOwnerRequest) (*BeneficialOwner, error)
}

// BeneficialOwnerServiceOp is an implementation of the beneficial owner
// service
type BeneficialOwnerServiceOp struct {
	client *Client
}

// BeneficialOwnerStatus is the status of the beneficial owner
type BeneficialOwnerStatus string

// BeneficialOwner is a beneficial owner
type BeneficialOwner struct {
	Resource
	ID                 string                `json:"id"`
	FirstName          string                `json:"firstName"`
	LastName           string                `json:"lastName"`
	Address            Address               `json:"address"`
	Passport           Passport              `json:"passport"`
	VerificationStatus BeneficialOwnerStatus `json:"verificationStatus"`
}

// BeneficialOwners is a collection of beneficial owners
type BeneficialOwners struct {
	Collection
	Embedded map[string][]BeneficialOwner `json:"_embedded"`
}

// BeneficialOwnerRequest is a beneficial owner request
type BeneficialOwnerRequest struct {
	FirstName   string    `json:"firstName,omitempty"`
	LastName    string    `json:"lastName,omitempty"`
	DateOfBirth string    `json:"dateOfBirth,omitempty"`
	SSN         string    `json:"ssn,omitempty"`
	Address     Address   `json:"address,omitempty"`
	Passport    *Passport `json:"passport,omitempty"`
}

// CertificationStatus is the beneficial ownership certification status
type CertificationStatus string

// BeneficialOwnership is the beneficial ownership status
type BeneficialOwnership struct {
	Resource
	Status CertificationStatus `json:"status"`
}

// BeneficialOwnershipRequest is a beneficial ownership request
type BeneficialOwnershipRequest struct {
	Status CertificationStatus `json:"status,omitempty"`
}

// Remove removes a beneficial owner matching the id
//
// see: https://docsv2.dwolla.com/#remove-a-beneficial-owner
func (b *BeneficialOwnerServiceOp) Remove(id string) error {
	return b.client.Delete(fmt.Sprintf("beneficial-owners/%s", id), nil, nil)
}

// Retrieve retrieves a beneficial owner matching the id
//
// see: https://docsv2.dwolla.com/#retrieve-a-beneficial-owner
func (b *BeneficialOwnerServiceOp) Retrieve(id string) (*BeneficialOwner, error) {
	var owner BeneficialOwner

	if err := b.client.Get(fmt.Sprintf("beneficial-owners/%s", id), nil, nil, &owner); err != nil {
		return nil, err
	}

	owner.client = b.client

	return &owner, nil
}

// Update updates a beneficial owner matching the id
//
// see: https://docsv2.dwolla.com/#update-a-beneficial-owner
func (b *BeneficialOwnerServiceOp) Update(id string, body *BeneficialOwnerRequest) (*BeneficialOwner, error) {
	var owner BeneficialOwner

	if err := b.client.Post(fmt.Sprintf("beneficial-owners/%s", id), body, nil, &owner); err != nil {
		return nil, err
	}

	owner.client = b.client

	return &owner, nil
}

// CreateDocument uploads a document for the beneficial owner
//
// see: https://docsv2.dwolla.com/#create-a-document-for-a-beneficial-owner
func (b *BeneficialOwner) CreateDocument(body *DocumentRequest) (*Document, error) {
	var document Document

	if _, ok := b.Links["self"]; !ok {
		return nil, errors.New("No self resource link")
	}

	if err := b.client.Upload(fmt.Sprintf("%s/documents", b.Links["self"].Href), body.Type, body.FileName, body.File, &document); err != nil {
		return nil, err
	}

	document.client = b.client

	return &document, nil
}

// ListDocuments returns documents for beneficial owner
//
// see: https://docsv2.dwolla.com/#list-documents-for-beneficial-owners
func (b *BeneficialOwner) ListDocuments() (*Documents, error) {
	var documents Documents

	if _, ok := b.Links["self"]; !ok {
		return nil, errors.New("No self resource link")
	}

	if err := b.client.Get(fmt.Sprintf("%s/documents", b.Links["self"].Href), nil, nil, &documents); err != nil {
		return nil, err
	}

	documents.client = b.client

	for i := range documents.Embedded["documents"] {
		documents.Embedded["documents"][i].client = b.client
	}

	return &documents, nil
}

// Remove removes the beneficial owner
//
// see: https://docsv2.dwolla.com/#remove-a-beneficial-owner
func (b *BeneficialOwner) Remove() error {
	if _, ok := b.Links["self"]; !ok {
		return errors.New("No self resource link")
	}

	return b.client.Delete(b.Links["self"].Href, nil, nil)
}

// Update updates the dwolla beneficial owner
//
// see: https://docsv2.dwolla.com/#update-a-beneficial-owner
func (b *BeneficialOwner) Update(body *BeneficialOwnerRequest) error {
	if _, ok := b.Links["self"]; !ok {
		return errors.New("No self resource link")
	}

	return b.client.Post(b.Links["self"].Href, body, nil, b)
}

// Certify certifies beneficial ownership
//
// see: https://docsv2.dwolla.com/#certify-beneficial-ownership
func (b *BeneficialOwnership) Certify() error {
	if _, ok := b.Links["self"]; !ok {
		return errors.New("No self resource link")
	}

	request := &BeneficialOwnershipRequest{Status: CertificationStatusCertified}

	return b.client.Post(b.Links["self"].Href, request, nil, b)
}
