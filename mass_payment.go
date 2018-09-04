package dwolla

import (
	"fmt"
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
