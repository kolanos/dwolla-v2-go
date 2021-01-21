package dwolla

import (
	"encoding/json"
	"fmt"
)

// HALError is a hal error
type HALError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Path    string `json:"path"`
}

// HALError implements the error interface
func (e HALError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// HALErrors is an array of embedded hal errors
type HALErrors map[string][]HALError

// ValidationError is a dwolla validation error
type ValidationError struct {
	Code     string    `json:"code"`
	Message  string    `json:"message"`
	Path     string    `json:"path"`
	Embedded HALErrors `json:"_embedded"`
}

// Error implements the error interface
func (v ValidationError) Error() string {
	return fmt.Sprintf("[%s] %s", v.Code, v.Embedded["_embedded"][0].Message)
}

// Link is a hal resource link
type Link struct {
	Href         string `json:"href"`
	ResourceType string `json:"resource-type,omitempty"`
	Type         string `json:"type,omitempty"`
}

// Links is a group of resource links
type Links map[string]Link

// Resource is a hal resource
type Resource struct {
	Links  Links `json:"_links,omitempty"`
	client *Client
}

// NewResource is constructor for Resource
func NewResource(links Links, client *Client) *Resource {
	return &Resource{Links: links, client: client}
}

// Embedded is a hal embedded resource
type Embedded map[string][]Resource

// Collection is a collection of hal resources
type Collection struct {
	Links    Links    `json:"_links"`
	Embedded Embedded `json:"_embedded"`
	Total    int      `json:"total"`
	client   *Client
}

// Unmarshal unmarhsals a hal object into a struct
func Unmarshal(data []byte, container interface{}) error {
	return json.Unmarshal(data, container)
}
