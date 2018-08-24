package dwolla

// Topic is an event topic
type Topic string

// Event is a dwolla event
type Event struct {
	Resource
	client     *Client
	ID         string `json:"id"`
	Created    string `json:"created"`
	Topic      Topic  `json:"topic"`
	ResourceID string `json:"resourceId"`
}

// Events is a collection of dwolla events
type Events struct {
	Collection
	client   *Client
	Embedded map[string][]Event `json:"_embedded"`
}
