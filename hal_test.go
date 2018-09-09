package dwolla

import (
	"testing"
)

func TestLink(t *testing.T) {
	var link Link

	linkJSON := `
{
		"href": "https://api.dwolla.com/customers/123",
		"type": "application/json",
		"resource-type": "customer"
}`

	if err := Unmarshal([]byte(linkJSON), &link); err != nil {
		t.Errorf("%v", err)
	}

	if link.Href != "https://api.dwolla.com/customers/123" {
		t.Errorf("Expected https://api.dwolla.com/customers/123, got %s", link.Href)
	}
}

func TestResource(t *testing.T) {
	var resource Resource

	resourceJSON := `
	{
		"_links": {
			"self": {
				"href": "https://api.dwolla.com/customers/123",
				"type": "application/json",
				"resource-type": "customer"
			}
		}
	}`

	if err := Unmarshal([]byte(resourceJSON), &resource); err != nil {
		t.Errorf("%v", err)
	}

	if resource.Links["self"].Href != "https://api.dwolla.com/customers/123" {
		t.Errorf("Expected https://api.dwolla.com/customers/123, got %s", resource.Links["self"].Href)
	}
}
