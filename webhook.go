package dwolla

// WebhookService is the webhook service interface
// see: https://docsv2.dwolla.com/#webhooks
type WebhookService interface{}

// WebhookServiceOp is an implementation of the webhook service interface
type WebhookServiceOp struct {
	client *Client
}
