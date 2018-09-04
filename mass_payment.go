package dwolla

import (
	"fmt"
	"net/url"
)

// MassPaymentService is the mass payment service interface
// see: https://docsv2.dwolla.com/#mass-payments
type MassPaymentService interface {
	Create(*MassPayment) (*MassPayment, error)
	Retrieve(string) (*MassPayment, error)
	Update(string, MassPaymentStatus) (*MassPayment, error)
}

// MassPaymentServiceOp is an implementation of the mass payment interface
type MassPaymentServiceOp struct {
	client *Client
}

// MassPaymentStatus is a mass payment status
type MassPaymentStatus string

const (
	// MassPaymentStatusDeferred is when a mass payment is deferred
	MassPaymentStatusDeferred MassPaymentStatus = "deferred"
	// MassPaymentStatusPending is when the mass payment is pending
	MassPaymentStatusPending MassPaymentStatus = "pending"
	// MassPaymentStatusProcessing is when the mass payment is processing
	MassPaymentStatusProcessing MassPaymentStatus = "processing"
	// MassPaymentStatusComplete is when the mass payment is complete
	MassPaymentStatusComplete MassPaymentStatus = "complete"
	// MassPaymentStatusCancelled is when the mass payment is cancelled
	MassPaymentStatusCancelled MassPaymentStatus = "cancelled"
)

// MassPayment is a dwolla mass payment
type MassPayment struct {
	Resource
	ID            string                       `json:"id,omitempty"`
	Status        MassPaymentStatus            `json:"status,omitempty"`
	Items         []MassPaymentItem            `json:"items,omitempty"`
	Embedded      map[string][]MassPaymentItem `json:"_embedded,omitempty"`
	Created       string                       `json:"created,omitempty"`
	MetaData      MetaData                     `json:"metadata,omitempty"`
	Total         Amount                       `json:"total,omitempty"`
	TotalFees     Amount                       `json:"totalFees,omitempty"`
	CorrelationID string                       `json:"correlationId,omitempty"`
}

// MassPayments is a collection of mass payments
type MassPayments struct {
	Collection
	Embedded map[string][]MassPayment `json:"_embedded"`
}

// MassPaymentItemStatus is a mass payment item status
type MassPaymentItemStatus string

const (
	// MassPaymentItemStatusPending is when a mass payment item is pending
	MassPaymentItemStatusPending MassPaymentItemStatus = "pending"
	// MassPaymentItemStatusSuccess is when amass payment item is successful
	MassPaymentItemStatusSuccess MassPaymentItemStatus = "success"
	// MassPaymentItemStatusFailed is when a mass payment item failed
	MassPaymentItemStatusFailed MassPaymentItemStatus = "failed"
)

// MassPaymentItem is a dwolla mass payment item
type MassPaymentItem struct {
	Resource
	ID            string                `json:"id,omitempty"`
	Status        MassPaymentItemStatus `json:"status,omitempty"`
	Amount        Amount                `json:"amount,omitempty"`
	MetaData      MetaData              `json:"metadata,omitempty"`
	CorrelationID string                `json:"correlationId,omitempty"`
	Embedded      HALErrors             `json:"_embedded,omitempty"`
}

// MassPaymentItems is a collection of mass payment items
type MassPaymentItems struct {
	Collection
	Embedded map[string][]MassPaymentItem `json:"_embedded"`
	Total    int                          `json:"total"`
}

// Create initiates a mass payment
// see: https://docsv2.dwolla.com/#initiate-a-mass-payment
func (m *MassPaymentServiceOp) Create(body *MassPayment) (*MassPayment, error) {
	var payment MassPayment

	if err := m.client.Post("mass-payments", body, nil, &payment); err != nil {
		return nil, err
	}

	payment.client = m.client

	return &payment, nil
}

// Retrieve retrieves the mass payment matching the id
// see: https://docsv2.dwolla.com/#retrieve-a-mass-payment
func (m *MassPaymentServiceOp) Retrieve(id string) (*MassPayment, error) {
	var payment MassPayment

	if err := m.client.Get(fmt.Sprintf("mass-payments/%s", id), nil, nil, &payment); err != nil {
		return nil, err
	}

	payment.client = m.client

	return &payment, nil
}

// Update updates a mass payment's status
// see: https://docsv2.dwolla.com/#update-a-mass-payment
func (m *MassPaymentServiceOp) Update(id string, status MassPaymentStatus) (*MassPayment, error) {
	var payment MassPayment

	body := &MassPayment{Status: status}

	if err := m.client.Post(fmt.Sprintf("mass-payments/%s", id), body, nil, &payment); err != nil {
		return nil, err
	}

	payment.client = m.client

	return &payment, nil
}

// ListItems returns a collection of items for the mass payment
// see: https://docsv2.dwolla.com/#list-items-for-a-mass-payment
func (m *MassPayment) ListItems(params *url.Values) (*MassPaymentItems, error) {
	var items MassPaymentItems

	if _, ok := m.Links["items"]; !ok {
		return nil, fmt.Errorf("No items resource link")
	}

	if err := m.client.Get(m.Links["items"].Href, params, nil, &items); err != nil {
		return nil, err
	}

	items.client = m.client

	for i := range items.Embedded["items"] {
		items.Embedded["items"][i].client = m.client
	}

	return &items, nil
}

// RetrieveItem returns a mass payment item matching id
// see: https://docsv2.dwolla.com/#retrieve-a-mass-payment-item
func (m *MassPayment) RetrieveItem(id string) (*MassPaymentItem, error) {
	var item MassPaymentItem

	if err := m.client.Get(fmt.Sprintf("mass-payment-items/%s", id), nil, nil, &item); err != nil {
		return nil, err
	}

	item.client = m.client

	return &item, nil
}

// RetrieveSource retrieves the mass payment funding source
func (m *MassPayment) RetrieveSource() (*FundingSource, error) {
	if _, ok := m.Links["source"]; !ok {
		return nil, fmt.Errorf("No source resource link")
	}

	return m.client.FundingSource.Retrieve(m.Links["source"].Href)
}

// RetrieveDestination retrieves the destination for the item
func (m *MassPaymentItem) RetrieveDestination() (*Customer, error) {
	if _, ok := m.Links["destination"]; !ok {
		return nil, fmt.Errorf("No destination resource link")
	}

	return m.client.Customer.Retrieve(m.Links["destination"].Href)
}

// RetrieveMassPayment retrieves the mass payment for the item
func (m *MassPaymentItem) RetrieveMassPayment() (*MassPayment, error) {
	if _, ok := m.Links["mass-payment"]; !ok {
		return nil, fmt.Errorf("No mass payment resource link")
	}

	return m.client.MassPayment.Retrieve(m.Links["mass-payment"].Href)
}

// RetrieveTransfer retrieves the transfer for the item
func (m *MassPaymentItem) RetrieveTransfer() (*Transfer, error) {
	if _, ok := m.Links["transfer"]; !ok {
		return nil, fmt.Errorf("No transfer resource link")
	}

	return m.client.Transfer.Retrieve(m.Links["transfer"].Href)
}
