package dwolla

// BeneficialOwnerService is the beneficial owner service interface
// see: https://docsv2.dwolla.com/#beneficial-owners
type BeneficialOwnerService interface {
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

const (
	// BeneficialOwnerStatusDocument is when the beneficial owner needs verification document
	BeneficialOwnerStatusDocument BeneficialOwnerStatus = "document"
	// BeneficialOwnerStatusIncomplete is when the beneficial owner is incomplete
	BeneficialOwnerStatusIncomplete BeneficialOwnerStatus = "incomplete"
	// BeneficialOwnerStatusVerified is when the beneficial owner is verified
	BeneficialOwnerStatusVerified BeneficialOwnerStatus = "verified"
)

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

// BeneficialOwnership is the beneficial ownership status
type BeneficialOwnership struct {
	Resource
	Status string `json:"status"`
}

// BeneficialOwnerRequest is a beneficial owner request
type BeneficialOwnerRequest struct {
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	DateOfBirth string   `json:"dateOfBirth"`
	SSN         string   `json:"ssn,omitempty"`
	Address     Address  `json:"address"`
	Passport    Passport `json:"passport,omitempty"`
}
