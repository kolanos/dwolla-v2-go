package dwolla

import (
	"errors"
	"fmt"
)

// WebhookSubscriptionService is the webhook subscription service interface
// see: https://docsv2.dwolla.com/#webhook-subscriptions
type WebhookSubscriptionService interface {
	Create(*WebhookSubscriptionRequest) (*WebhookSubscription, error)
	Retrieve(string) (*WebhookSubscription, error)
	List() (*WebhookSubscriptions, error)
	Remove(string) error
}

// WebhookSubscriptionServiceOp is an implementation of the webhook
// subscription service interface
type WebhookSubscriptionServiceOp struct {
	client *Client
}

// WebhookSubscription is a webhook subscription
type WebhookSubscription struct {
	Resource
	ID      string `json:"id"`
	URL     string `json:"url"`
	Created string `json:"created"`
}

// WebhookSubscriptions is a collection of webhook subscriptions
type WebhookSubscriptions struct {
	Collection
	Embedded map[string][]WebhookSubscription `json:"_embedded"`
	Total    int                              `json:"total"`
}

// WebhookSubscriptionRequest is a webhook subscription request
type WebhookSubscriptionRequest struct {
	URL    string `json:"url,omitempty"`
	Secret string `json:"secret,omitempty"`
	Paused bool   `json:"paused,omitempty"`
}

// Create creates a webhook subscription
func (w *WebhookSubscriptionServiceOp) Create(body *WebhookSubscriptionRequest) (*WebhookSubscription, error) {
	var subscription WebhookSubscription

	if err := w.client.Post("webhook-subscriptions", body, nil, &subscription); err != nil {
		return nil, err
	}

	subscription.client = w.client

	return &subscription, nil
}

// Retrieve retrieves the webhook subscription matching id
func (w *WebhookSubscriptionServiceOp) Retrieve(id string) (*WebhookSubscription, error) {
	var subscription WebhookSubscription

	if err := w.client.Get(fmt.Sprintf("webhook-subscriptions/%s", id), nil, nil, &subscription); err != nil {
		return nil, err
	}

	subscription.client = w.client

	return &subscription, nil
}

// List returns a list of webhook subscriptions
func (w *WebhookSubscriptionServiceOp) List() (*WebhookSubscriptions, error) {
	var subscriptions WebhookSubscriptions

	if err := w.client.Get("webhook-subscriptions", nil, nil, &subscriptions); err != nil {
		return nil, err
	}

	subscriptions.client = w.client
	for i := range subscriptions.Embedded["webhook-subscriptions"] {
		subscriptions.Embedded["webhook-subscriptions"][i].client = w.client
	}

	return &subscriptions, nil
}

// Remove removes the webhook subscription matching the id
func (w *WebhookSubscriptionServiceOp) Remove(id string) error {
	return w.client.Delete(fmt.Sprintf("webhook-subscriptions/%s", id), nil, nil)
}

// Pause pauses the webhook subscription
func (w *WebhookSubscription) Pause() error {
	if _, ok := w.Links["self"]; !ok {
		return errors.New("No self resource link")
	}

	body := &WebhookSubscriptionRequest{Paused: true}

	return w.client.Post(w.Links["self"].Href, body, nil, w)
}

// Remove removes the webhook subscription
func (w *WebhookSubscription) Remove() error {
	if _, ok := w.Links["self"]; !ok {
		return errors.New("No self resource link")
	}

	return w.client.Delete(w.Links["self"].Href, nil, nil)
}

// Unpause unpauses the webhook subscription
func (w *WebhookSubscription) Unpause() error {
	if _, ok := w.Links["self"]; !ok {
		return errors.New("No self resource link")
	}

	body := &WebhookSubscriptionRequest{Paused: false}

	return w.client.Post(w.Links["self"].Href, body, nil, w)
}

// RetrieveWebhooks returns webhooks for this webhook subscription
func (w *Webhook) RetrieveWebhooks() (*Webhooks, error) {
	var webhooks Webhooks

	if _, ok := w.Links["webhooks"]; !ok {
		return nil, errors.New("No webhooks resource link")
	}

	if err := w.client.Get(w.Links["webhooks"].Href, nil, nil, &webhooks); err != nil {
		return nil, err
	}

	webhooks.client = w.client

	for i := range webhooks.Embedded["webhooks"] {
		webhooks.Embedded["webhooks"][i].client = w.client
	}

	return &webhooks, nil
}
