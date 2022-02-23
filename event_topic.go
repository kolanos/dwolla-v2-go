package dwolla

// EventTopic is an event topic
type EventTopic string

const (
	EventTopic_BankTransferCreated                EventTopic = "bank_transfer_created"
	EventTopic_BankTransferCreationFailed         EventTopic = "bank_transfer_creation_failed"
	EventTopic_BankTransferCancelled              EventTopic = "bank_transfer_cancelled"
	EventTopic_BankTransferFailed                 EventTopic = "bank_transfer_failed"
	EventTopic_BankTransfer_Completed             EventTopic = "bank_transfer_completed"
	EventTopic_TransferCreated                    EventTopic = "transfer_created"
	EventTopic_Cancelled                          EventTopic = "transfer_cancelled"
	EventTopic_TransferFailed                     EventTopic = "transfer_failed"
	EventTopic_Transfer_Completed                 EventTopic = "transfer_completed"
	EventTopic_CustomerBankTransferCreated        EventTopic = "customer_bank_transfer_created"
	EventTopic_CustomerBankTransferCreationFailed EventTopic = "customer_bank_transfer_creation_failed"
	EventTopic_CustomerBankTransferCancelled      EventTopic = "customer_bank_transfer_cancelled"
	EventTopic_CustomerBankTransferFailed         EventTopic = "customer_bank_transfer_failed"
	EventTopic_CustomerBankTransferCompleted      EventTopic = "customer_bank_transfer_completed"
	EventTopic_CustomerTransferCreated            EventTopic = "customer_transfer_created"
	EventTopic_CustomerTransferCancelled          EventTopic = "customer_transfer_cancelled"
	EventTopic_CustomerTransferFailed             EventTopic = "customer_transfer_failed"
	EventTopic_CustomerTransferCompleted          EventTopic = "customer_transfer_completed"
	EventTopic_FundingSourceAdded                 EventTopic = "funding_source_added"
	EventTopic_FundingSourceRemoved               EventTopic = "funding_source_removed"
	EventTopic_FundingSourceVerified              EventTopic = "funding_source_verified"
	EventTopic_CustomerCreated                    EventTopic = "customer_created"
	EventTopic_CustomerSuspended                  EventTopic = "customer_suspended"
	EventTopic_CustomerActivated                  EventTopic = "customer_activated"
	EventTopic_CustomerDeactivated                EventTopic = "customer_deactivated"
	EventTopic_CustomerVerified                   EventTopic = "customer_verified"
	EventTopic_CustomerFundingSourceAdded         EventTopic = "customer_funding_source_added"
	EventTopic_CustomerFundingSourceRemoved       EventTopic = "customer_funding_source_removed"
	EventTopic_CustomerFundingSourceVerified      EventTopic = "customer_funding_source_verified"
	EventTopic_CustomerFundingSourceUpdated       EventTopic = "customer_funding_source_updated"
)
