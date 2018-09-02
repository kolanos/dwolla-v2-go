package dwolla

import (
	"encoding/json"
	"fmt"
)

// Error is a hal error
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Path    string `json:"path"`
}

// Error implements the error interface
func (e Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Errors is an array of embedded hal errors
type Errors map[string][]Error

// ValidationError is a dwolla validation error
type ValidationError struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Path     string `json:"path"`
	Embedded Errors `json:"_embedded"`
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("[%s] %s", v.Code, v.Message)
}

// Link is a hal resource link
type Link struct {
	Href         string `json:"href"`
	ResourceType string `json:"resource-type"`
	Type         string `json:"type"`
}

// Links is a group of resource links
type Links map[string]Link

// Resource is a hal resource
type Resource struct {
	Links  Links `json:"_links,omitempty"`
	client *Client
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
	if err := json.Unmarshal(data, container); err != nil {
		return err
	}

	return nil
}
