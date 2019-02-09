package dwolla

import (
	"fmt"
)

// OnDemandAuthorizationService is the on-demand authorization interface
//
// see: https://docsv2.dwolla.com/#create-an-on-demand-transfer-authorization
type OnDemandAuthorizationService interface {
	Create() (*OnDemandAuthorization, error)
	Retrieve(string) (*OnDemandAuthorization, error)
}

// OnDemandAuthorizationServiceOp is an implementation of the on-demand authorization interface
type OnDemandAuthorizationServiceOp struct {
	client *Client
}

// OnDemandAuthorization is a dwolla on-demand transfer authorization
type OnDemandAuthorization struct {
	Resource
	BodyText   string `json:"bodyText"`
	ButtonText string `json:"buttonText"`
}

// Create creates an on-demand transfer authorization
func (o *OnDemandAuthorizationServiceOp) Create() (*OnDemandAuthorization, error) {
	var authorization OnDemandAuthorization

	if err := o.client.Post("on-demand-authorizations", nil, nil, &authorization); err != nil {
		return nil, err
	}

	authorization.client = o.client

	return &authorization, nil
}

// Retrieve returns a on-demand authorization matching the id
func (o *OnDemandAuthorizationServiceOp) Retrieve(id string) (*OnDemandAuthorization, error) {
	var authorization OnDemandAuthorization

	if err := o.client.Get(fmt.Sprintf("on-demand-authorizations/%s", id), nil, nil, &authorization); err != nil {
		return nil, err
	}

	authorization.client = o.client

	return &authorization, nil
}
