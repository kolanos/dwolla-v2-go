package dwolla

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
	ID        string                `json:"id"`
	FirstName string                `json:"firstName"`
	LastName  string                `json:"lastName"`
	Address   Address               `json:"address"`
	Status    BeneficialOwnerStatus `json:"verificationStatus"`
}
