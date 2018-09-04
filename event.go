package dwolla

// EventService is the event service interface
// see: https://docsv2.dwolla.com/#events
type EventService interface{}

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
