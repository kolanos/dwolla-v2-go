package dwolla

// WebhookSubscriptionService is the webhook subscription service interface
// see: https://docsv2.dwolla.com/#webhook-subscriptions
type WebhookSubscriptionService interface {
	Retrieve(string) (*WebhookSubscription, error)
}

// WebhookSubscriptionServiceOp is an implementation of the webhook
// subscription service interface
type WebhookSubscriptionServiceOp struct {
	client *Client
}

// WebhookSubscription is a webhook subscription
type WebhookSubscription struct{}

// Retrieve retrieves the webhook subscription matching id
func (w *WebhookSubscriptionServiceOp) Retrieve(id string) (*WebhookSubscription, error) {
	var subscription WebhookSubscription
	return &subscription, nil
}
