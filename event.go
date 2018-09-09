package dwolla

import (
	"fmt"
	"net/url"
)

// EventService is the event service interface
//
// see: https://docsv2.dwolla.com/#events
type EventService interface {
	List(*url.Values) (*Events, error)
	Retrieve(string) (*Event, error)
}

// EventServiceOp is an implementation of the event service interface
type EventServiceOp struct {
	client *Client
}

// EventTopic is an event topic
type EventTopic string

// Event is a dwolla event
type Event struct {
	Resource
	ID         string     `json:"id"`
	Created    string     `json:"created"`
	Topic      EventTopic `json:"topic"`
	ResourceID string     `json:"resourceId"`
}

// Events is a collection of dwolla events
type Events struct {
	Collection
	client   *Client
	Embedded map[string][]Event `json:"_embedded"`
}

// List returns a collection of events
//
// see: https://docsv2.dwolla.com/#list-events
func (e *EventServiceOp) List(params *url.Values) (*Events, error) {
	var events Events

	if err := e.client.Get("events", params, nil, &events); err != nil {
		return nil, err
	}

	events.client = e.client

	for i := range events.Embedded["events"] {
		events.Embedded["events"][i].client = e.client
	}

	return &events, nil
}

// Retrieve retrieves the event matching the id
//
// see: https://docsv2.dwolla.com/#retrieve-an-event
func (e *EventServiceOp) Retrieve(id string) (*Event, error) {
	var event Event

	if err := e.client.Get(fmt.Sprintf("events/%s", id), nil, nil, &event); err != nil {
		return nil, err
	}

	event.client = e.client

	return &event, nil
}
