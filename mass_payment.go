package dwolla

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

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

const (
	// MassPaymentItemStatusPending is when a mass payment item is pending
	MassPaymentItemStatusPending MassPaymentItemStatus = "pending"
	// MassPaymentItemStatusSuccess is when amass payment item is successful
	MassPaymentItemStatusSuccess MassPaymentItemStatus = "success"
	// MassPaymentItemStatusFailed is when a mass payment item failed
	MassPaymentItemStatusFailed MassPaymentItemStatus = "failed"
)

// MassPaymentService is the mass payment service interface
//
// see: https://docsv2.dwolla.com/#mass-payments
type MassPaymentService interface {
	Create(context.Context, *MassPayment) (*MassPayment, error)
	Retrieve(context.Context, string) (*MassPayment, error)
	Update(context.Context, string, MassPaymentStatus) (*MassPayment, error)
}

// MassPaymentServiceOp is an implementation of the mass payment interface
type MassPaymentServiceOp struct {
	client *Client
}

// MassPaymentStatus is a mass payment status
type MassPaymentStatus string

// MassPayment is a dwolla mass payment
type MassPayment struct {
	Resource
	ID            string                       `json:"id,omitempty"`
	Status        MassPaymentStatus            `json:"status,omitempty"`
	ACHDetails    *ACHDetails                  `json:"achDetails,omitempty"`
	Clearing      *Clearing                    `json:"clearing,omitempty"`
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
//
// see: https://docsv2.dwolla.com/#initiate-a-mass-payment
func (m *MassPaymentServiceOp) Create(ctx context.Context, body *MassPayment) (*MassPayment, error) {
	var payment MassPayment

	if err := m.client.Post(ctx, "mass-payments", body, nil, &payment); err != nil {
		return nil, err
	}

	payment.client = m.client

	return &payment, nil
}

// Retrieve retrieves the mass payment matching the id
//
// see: https://docsv2.dwolla.com/#retrieve-a-mass-payment
func (m *MassPaymentServiceOp) Retrieve(ctx context.Context, id string) (*MassPayment, error) {
	var payment MassPayment

	if err := m.client.Get(ctx, fmt.Sprintf("mass-payments/%s", id), nil, nil, &payment); err != nil {
		return nil, err
	}

	payment.client = m.client

	return &payment, nil
}

// Update updates a mass payment's status
//
// see: https://docsv2.dwolla.com/#update-a-mass-payment
func (m *MassPaymentServiceOp) Update(ctx context.Context, id string, status MassPaymentStatus) (*MassPayment, error) {
	var payment MassPayment

	body := &MassPayment{Status: status}

	if err := m.client.Post(ctx, fmt.Sprintf("mass-payments/%s", id), body, nil, &payment); err != nil {
		return nil, err
	}

	payment.client = m.client

	return &payment, nil
}

// ListItems returns a collection of items for the mass payment
//
// see: https://docsv2.dwolla.com/#list-items-for-a-mass-payment
func (m *MassPayment) ListItems(ctx context.Context, params *url.Values) (*MassPaymentItems, error) {
	var items MassPaymentItems

	if _, ok := m.Links["items"]; !ok {
		return nil, errors.New("No items resource link")
	}

	if err := m.client.Get(ctx, m.Links["items"].Href, params, nil, &items); err != nil {
		return nil, err
	}

	items.client = m.client

	for i := range items.Embedded["items"] {
		items.Embedded["items"][i].client = m.client
	}

	return &items, nil
}

// RetrieveItem returns a mass payment item matching id
//
// see: https://docsv2.dwolla.com/#retrieve-a-mass-payment-item
func (m *MassPayment) RetrieveItem(ctx context.Context, id string) (*MassPaymentItem, error) {
	var item MassPaymentItem

	if err := m.client.Get(ctx, fmt.Sprintf("mass-payment-items/%s", id), nil, nil, &item); err != nil {
		return nil, err
	}

	item.client = m.client

	return &item, nil
}

// RetrieveSource retrieves the mass payment funding source
func (m *MassPayment) RetrieveSource(ctx context.Context) (*FundingSource, error) {
	if _, ok := m.Links["source"]; !ok {
		return nil, errors.New("No source resource link")
	}

	return m.client.FundingSource.Retrieve(ctx, m.Links["source"].Href)
}

// RetrieveDestination retrieves the destination for the item
func (m *MassPaymentItem) RetrieveDestination(ctx context.Context) (*Customer, error) {
	if _, ok := m.Links["destination"]; !ok {
		return nil, errors.New("No destination resource link")
	}

	return m.client.Customer.Retrieve(ctx, m.Links["destination"].Href)
}

// RetrieveMassPayment retrieves the mass payment for the item
func (m *MassPaymentItem) RetrieveMassPayment(ctx context.Context) (*MassPayment, error) {
	if _, ok := m.Links["mass-payment"]; !ok {
		return nil, errors.New("No mass payment resource link")
	}

	return m.client.MassPayment.Retrieve(ctx, m.Links["mass-payment"].Href)
}

// RetrieveTransfer retrieves the transfer for the item
func (m *MassPaymentItem) RetrieveTransfer(ctx context.Context) (*Transfer, error) {
	if _, ok := m.Links["transfer"]; !ok {
		return nil, errors.New("No transfer resource link")
	}

	return m.client.Transfer.Retrieve(ctx, m.Links["transfer"].Href)
}
