package dwolla

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBeneficialOwnerServiceRetrieve(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "beneficial-owner.json"))
	res, err := c.BeneficialOwner.Retrieve(ctx, "00cb67f2-768c-4ee3-ac81-73bc4faf9c2b")

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.ID, "00cb67f2-768c-4ee3-ac81-73bc4faf9c2b")
}

func TestBeneficialOwnerServiceRetrieveError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))
	res, err := c.BeneficialOwner.Retrieve(ctx, "00cb67f2-768c-4ee3-ac81-73bc4faf9c2b")

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestBeneficialOwnerServiceRemove(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "beneficial-owner.json"))
	err := c.BeneficialOwner.Remove(ctx, "00cb67f2-768c-4ee3-ac81-73bc4faf9c2b")

	assert.NoError(t, err)
}

func TestBeneficialOwnerServiceUpdate(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "beneficial-owner.json"))
	res, err := c.BeneficialOwner.Update(ctx, "00cb67f2-768c-4ee3-ac81-73bc4faf9c2b", &BeneficialOwnerRequest{
		FirstName: "John",
		LastName:  "Doe",
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.ID, "00cb67f2-768c-4ee3-ac81-73bc4faf9c2b")
}

func TestBeneficialOwnerServiceUpdateError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))
	res, err := c.BeneficialOwner.Update(ctx, "00cb67f2-768c-4ee3-ac81-73bc4faf9c2b", &BeneficialOwnerRequest{
		FirstName: "John",
		LastName:  "Doe",
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestBeneficialOwnerCreateDocument(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "document.json"))

	owner := &BeneficialOwner{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/beneficial-owners/07d59716-ef22-4fe6-98e8-f3190233dfb8"}}}}

	f, _ := os.Open(filepath.Join("testdata", "document-upload-success.png"))
	res, err := owner.CreateDocument(ctx, &DocumentRequest{
		Type:     DocumentTypePassport,
		FileName: f.Name(),
		File:     f,
	})

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestBeneficialOwnerCreateDocumentError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	owner := &BeneficialOwner{Resource: Resource{client: c}}
	f1, _ := os.Open(filepath.Join("testdata", "document-upload-success.png"))
	defer f1.Close()
	res, err := owner.CreateDocument(ctx, &DocumentRequest{
		Type:     DocumentTypePassport,
		FileName: f1.Name(),
		File:     f1,
	})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "No self resource link")
	assert.Nil(t, res)

	owner = &BeneficialOwner{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/beneficial-owners/07d59716-ef22-4fe6-98e8-f3190233dfb8"}}}}
	f2, _ := os.Open(filepath.Join("testdata", "document-upload-success.png"))
	res, err = owner.CreateDocument(ctx, &DocumentRequest{
		Type:     DocumentTypePassport,
		FileName: f2.Name(),
		File:     f2,
	})

	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestBeneficialOwnerListDocuments(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "documents.json"))

	owner := &BeneficialOwner{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/beneficial-owners/07d59716-ef22-4fe6-98e8-f3190233dfb8"}}}}
	res, err := owner.ListDocuments(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestBeneficialOwnerListDocumentsError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	owner := &BeneficialOwner{Resource: Resource{client: c}}
	res, err := owner.ListDocuments(ctx)
	assert.Error(t, err)
	assert.Nil(t, res)

	owner = &BeneficialOwner{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/beneficial-owners/07d59716-ef22-4fe6-98e8-f3190233dfb8"}}}}
	res, err = owner.ListDocuments(ctx)
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestBeneficialOwnerRemove(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "documents.json"))

	owner := &BeneficialOwner{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/beneficial-owners/07d59716-ef22-4fe6-98e8-f3190233dfb8"}}}}
	err := owner.Remove(ctx)

	assert.NoError(t, err)
}

func TestBeneficialOwnerRemoveError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	owner := &BeneficialOwner{Resource: Resource{client: c}}
	err := owner.Remove(ctx)
	assert.Error(t, err)

	owner = &BeneficialOwner{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/beneficial-owners/07d59716-ef22-4fe6-98e8-f3190233dfb8"}}}}
	err = owner.Remove(ctx)
	assert.Error(t, err)
}

func TestBeneficialOwnerUpdate(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "beneficial-owner.json"))

	owner := &BeneficialOwner{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/beneficial-owners/07d59716-ef22-4fe6-98e8-f3190233dfb8"}}}}
	err := owner.Update(ctx, &BeneficialOwnerRequest{
		FirstName: "John",
		LastName:  "Doe",
	})

	assert.NoError(t, err)
}

func TestBeneficialOwnerUpdateError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	owner := &BeneficialOwner{Resource: Resource{client: c}}
	err := owner.Update(ctx, &BeneficialOwnerRequest{
		FirstName: "John",
		LastName:  "Doe",
	})
	assert.Error(t, err)

	owner = &BeneficialOwner{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/beneficial-owners/07d59716-ef22-4fe6-98e8-f3190233dfb8"}}}}
	err = owner.Update(ctx, &BeneficialOwnerRequest{
		FirstName: "John",
		LastName:  "Doe",
	})
	assert.Error(t, err)
}

func TestBeneficialOwnershipCertify(t *testing.T) {
	c := newMockClient(200, filepath.Join("testdata", "beneficial-ownership.json"))

	ownership := &BeneficialOwnership{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/beneficial-owners/07d59716-ef22-4fe6-98e8-f3190233dfb8"}}}}
	err := ownership.Certify(ctx)

	assert.NoError(t, err)
}

func TestBeneficialOwnershipCertifyError(t *testing.T) {
	c := newMockClient(404, filepath.Join("testdata", "resource-not-found.json"))

	ownership := &BeneficialOwnership{Resource: Resource{client: c}}
	err := ownership.Certify(ctx)
	assert.Error(t, err)

	ownership = &BeneficialOwnership{Resource: Resource{client: c, Links: Links{"self": Link{Href: "https://api-sandbox.dwolla.com/beneficial-owners/07d59716-ef22-4fe6-98e8-f3190233dfb8"}}}}
	err = ownership.Certify(ctx)
	assert.Error(t, err)
}
