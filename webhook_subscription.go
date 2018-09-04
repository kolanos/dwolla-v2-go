package dwolla

// WebhookSubscriptionService is the webhook subscription service interface
// see: https://docsv2.dwolla.com/#webhook-subscriptions
type WebhookSubscriptionService interface{}

// WebhookSubscriptionServiceOp is an implementation of the webhook
// subscription service interface
type WebhookSubscriptionServiceOp struct {
	client *Client
}
