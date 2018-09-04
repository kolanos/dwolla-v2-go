package dwolla

import (
	"fmt"
)

// WebhookService is the webhook service interface
// see: https://docsv2.dwolla.com/#webhooks
type WebhookService interface {
	Retrieve(string) (*Webhook, error)
}

// WebhookServiceOp is an implementation of the webhook service interface
type WebhookServiceOp struct {
	client *Client
}

// Webhook is a dwolla webhook
type Webhook struct {
	Resource
	ID             string           `json:"id"`
	Topic          EventTopic       `json:"topic"`
	AccountID      string           `json:"accountId"`
	EventID        string           `json:"eventId"`
	SubscriptionID string           `json:"subscriptionId"`
	Attempts       []WebhookAttempt `json:"attempts"`
}

// Webhooks is a collection of webhooks
type Webhooks struct {
	Collection
	Embedded map[string][]Webhook `json:"_embedded"`
	Total    int                  `json:"total"`
}

// WebhookAttempt is a webhook attempt
type WebhookAttempt struct {
	ID       string          `json:"id"`
	Request  WebhookRequest  `json:"request"`
	Response WebhookResponse `json:"response"`
}

// WebhookHeader is a webhook request/response header
type WebhookHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// WebhookRequest is a webhook request
type WebhookRequest struct {
	Timestamp string          `json:"timestamp"`
	UrL       string          `json:"url"`
	Headers   []WebhookHeader `json:"headers"`
	Body      string          `json:"body"`
}

// WebhookResponse is a webhook response
type WebhookResponse struct {
	Timestamp  string          `json:"timestamp"`
	Headers    []WebhookHeader `json:"headers"`
	StatusCode int             `json:"statusCode"`
	Body       string          `json:"body"`
}

// WebhookRetry is a webhook retry
type WebhookRetry struct {
	Resource
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
}

// WebhookRetries is a collection of webhook retries
type WebhookRetries struct {
	Collection
	Embedded map[string][]WebhookRetry `json:"_embedded"`
	Total    int                       `json:"total"`
}

// Retrieve retrieves the webhook with matching id
// see: https://docsv2.dwolla.com/#retrieve-a-webhook
func (w *WebhookServiceOp) Retrieve(id string) (*Webhook, error) {
	var webhook Webhook

	if err := w.client.Get(fmt.Sprintf("webhooks/%s", id), nil, nil, &webhook); err != nil {
		return nil, err
	}

	webhook.client = w.client

	return &webhook, nil
}

// ListRetries returns a collection of retries for this webhook
// see: https://docsv2.dwolla.com/#list-retries-for-a-webhook
func (w *Webhook) ListRetries() (*WebhookRetries, error) {
	var retries WebhookRetries

	if _, ok := w.Links["retry"]; !ok {
		return nil, fmt.Errorf("No retry resource link")
	}

	if err := w.client.Get(w.Links["retry"].Href, nil, nil, &retries); err != nil {
		return nil, err
	}

	retries.client = w.client

	for i := range retries.Embedded["retries"] {
		retries.Embedded["retries"][i].client = w.client
	}

	return &retries, nil
}

// Retry retries the webhook
// see: https://docsv2.dwolla.com/#retry-a-webhook
func (w *Webhook) Retry() (*WebhookRetry, error) {
	var retry WebhookRetry

	if _, ok := w.Links["retry"]; !ok {
		return nil, fmt.Errorf("No retry resource link")
	}

	if err := w.client.Post(w.Links["retry"].Href, nil, nil, &retry); err != nil {
		return nil, err
	}

	retry.client = w.client

	return &retry, nil
}
