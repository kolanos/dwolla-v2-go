package dwolla

import (
	"context"
	"fmt"
	"time"
)

type TransferFailureCode string

const (
	TransferFailureCodeInsufficientFunds                     TransferFailureCode = "R01"
	TransferFailureCodeBankAccountClosed                     TransferFailureCode = "R02"
	TransferFailureCodeNoAccount                             TransferFailureCode = "R03"
	TransferFailureCodeInvalidBankAccountNumberStructure     TransferFailureCode = "R04"
	TransferFailureCodeUnauthorizedDebitToConsumerAccount    TransferFailureCode = "R05"
	TransferFailureCodeReturnedPerODFIRequest                TransferFailureCode = "R06"
	TransferFailureCodeAuthorizationRevokedByCustomer        TransferFailureCode = "R07"
	TransferFailureCodePaymentStopped                        TransferFailureCode = "R08"
	TransferFailureCodeUncollectedFunds                      TransferFailureCode = "R09"
	TransferFailureCodeCustomerAdvisesNotAuthorized          TransferFailureCode = "R10"
	TransferFailureCodeBankAccountFrozen                     TransferFailureCode = "R16"
	TransferFailureCodeNonTransactionAccount                 TransferFailureCode = "R20"
	TransferFailureCodeCreditEntryRefusedByReceiver          TransferFailureCode = "R23"
	TransferFailureCodeCorporateCustomerAdvisesNotAuthorized TransferFailureCode = "R29"
)

// TransferFailureService is the transfer service interface
//
// see: https://docsv2.dwolla.com/#transfers
type TransferFailureService interface {
	Retrieve(ctx context.Context, transferID string) (*TransferFailure, error)
}

// TransferFailureServiceOp is an implementation of the transfer service interface
type TransferFailureServiceOp struct {
	client *Client
}

// TransferFailure is a dwolla transfer
type TransferFailure struct {
	Resource
	Code        TransferFailureCode `json:"code"`
	Description string              `json:"description"`
	Explanation string              `json:"explanation"`
	Created     string              `json:"created"`
}

func (t TransferFailureServiceOp) Retrieve(ctx context.Context, transferID string) (*TransferFailure, error) {
	var transferFailure TransferFailure

	if err := t.client.Get(ctx, fmt.Sprintf("transfers/%s/failure", transferID), nil, nil, &transferFailure); err != nil {
		return nil, err
	}

	transferFailure.client = t.client

	return &transferFailure, nil
}

// CreatedTime returns the created value as time.Time
func (tf *TransferFailure) CreatedTime() time.Time {
	t, _ := time.Parse(time.RFC3339, tf.Created)
	return t
}
