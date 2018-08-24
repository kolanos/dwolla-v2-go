package dwolla

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Environment is a supported dwolla environment
type Environment string

const (
	// Version is the version of the client
	Version string = "0.1.0"

	// Production is the production environment
	Production Environment = "production"
	// Sandbox is the sanbox environment
	Sandbox Environment = "sandbox"

	// ProductionAPIURL is the production api url
	ProductionAPIURL = "https://api.dwolla.com"
	// ProductionAuthURL is the production auth url
	ProductionAuthURL = "https://www.dwolla.com/oauth/v2/authenticate"
	// ProductionTokenURL is the production token url
	ProductionTokenURL = "https://www.dwolla.com/oauth/v2/token"

	// SandboxAPIURL is the sandbox api url
	SandboxAPIURL = "https://api-sandbox.dwolla.com"
	// SandboxAuthURL is the sandbox auth url
	SandboxAuthURL = "https://sandbox.dwolla.com/oauth/v2/authenticate"
	// SandboxTokenURL is the sandbox token url
	SandboxTokenURL = "https://sandbox.dwolla.com/oauth/v2/token"
)

// Token is a dwolla auth token
type Token struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	TokenType        string `json:"token_type"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// HTTPClient is the http client interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is the dwolla client
type Client struct {
	Key         string
	Secret      string
	Environment Environment
	HTTPClient  HTTPClient
	Token       *Token

	root                   *Resource
	Account                AccountService
	BusinessClassification BusinessClassificationService
	Customer               CustomerService
	FundingSource          FundingSourceService
}

// New initializes a new dwolla client
func New(key, secret string, environment Environment) *Client {
	return NewWithHTTPClient(key, secret, environment, nil)
}

// NewWithHTTPClient initializes the client with specified http client
func NewWithHTTPClient(key, secret string, environment Environment, httpClient HTTPClient) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{
		Key:         key,
		Secret:      secret,
		Environment: environment,
		HTTPClient:  httpClient,
	}

	c.Account = &AccountServiceOp{c}
	c.BusinessClassification = &BusinessClassificationServiceOp{c}
	c.Customer = &CustomerServiceOp{c}
	c.FundingSource = &FundingSourceServiceOp{c}

	return c
}

// APIURL returns the api url for the environment
func (c Client) APIURL() string {
	var url string

	switch c.Environment {
	case Production:
		url = ProductionAPIURL
	case Sandbox:
		url = SandboxAPIURL
	}

	return url
}

// BuildAPIURL builds an api url with a given path
func (c Client) BuildAPIURL(path string) string {
	apiURL := c.APIURL()

	if strings.HasPrefix(path, apiURL) {
		return path
	}

	if !strings.HasPrefix(path, "/") {
		path = fmt.Sprintf("/%s", path)
	}

	return fmt.Sprintf("%s%s", apiURL, path)
}

// AuthURL returns the auth url for the environment
func (c Client) AuthURL() string {
	var url string

	switch c.Environment {
	case Production:
		url = ProductionAuthURL
	case Sandbox:
		url = SandboxAuthURL
	}

	return url
}

// TokenURL returns the token url for the environment
func (c Client) TokenURL() string {
	var url string

	switch c.Environment {
	case Production:
		url = ProductionTokenURL
	case Sandbox:
		url = SandboxTokenURL
	}

	return url
}

// RequestToken requests a new auth token using client credentials
func (c *Client) RequestToken() error {
	var (
		err   error
		token Token
	)

	buf := bytes.NewBuffer([]byte("grant_type=client_credentials"))

	req, err := http.NewRequest("POST", c.TokenURL(), buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", fmt.Sprintf("dwolla-v2-go/%s", Version))

	req.SetBasicAuth(c.Key, c.Secret)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resBody, &token)
	if err != nil {
		return err
	}

	if token.Error != "" {
		return fmt.Errorf("[%s] %s", token.Error, token.ErrorDescription)
	}

	c.Token = &token

	return nil
}

// Get performs a GET against the api
func (c *Client) Get(path string, params *url.Values, headers *http.Header, container interface{}) error {
	var (
		err      error
		halError Error
	)

	if c.Token == nil {
		if err := c.RequestToken(); err != nil {
			return err
		}
	}
	req, err := http.NewRequest("GET", c.BuildAPIURL(path), nil)
	if err != nil {
		return err
	}

	if headers != nil {
		req.Header = *headers
	}

	req.Header.Set("Accept", "application/vnd.dwolla.v1.hal+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token.AccessToken))
	req.Header.Set("User-Agent", fmt.Sprintf("dwolla-v2-go/%s", Version))

	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	//fmt.Println(string(resBody))

	if res.StatusCode > 299 {
		if err := json.Unmarshal(resBody, &halError); err != nil {
			return err
		}

		// If token is expired, we'll attempt to get a newone and reattempt
		// the request. This should probably be moved to a method to handle
		// all error scenarios.
		if halError.Code == "ExpiredAccessToken" {
			if err := c.RequestToken(); err == nil {
				return c.Get(path, params, headers, container)
			}
		}

		return halError
	}

	if container != nil {
		return json.Unmarshal(resBody, container)
	}

	return nil
}

// Post performs a POST against the api
func (c *Client) Post(path string, body interface{}, headers *http.Header, container interface{}) error {
	var (
		err             error
		halError        Error
		validationError ValidationError
	)

	if c.Token == nil {
		if err := c.RequestToken(); err != nil {
			return err
		}
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.BuildAPIURL(path), bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}

	if headers != nil {
		req.Header = *headers
	}

	req.Header.Set("Accept", "application/vnd.dwolla.v1.hal+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token.AccessToken))
	req.Header.Set("Content-Type", "application/vnd.dwolla.v1.hal+json")
	req.Header.Set("User-Agent", "dwolla-v2-go")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode > 299 {
		if err := json.Unmarshal(resBody, &halError); err != nil {
			return err
		}

		// If token is expired, we'll attempt to get a newone and reattempt
		// the request. This should probably be moved to a method to handle
		// all error scenarios.
		if halError.Code == "ExpiredAccessToken" {
			if err := c.RequestToken(); err == nil {
				return c.Post(path, body, headers, container)
			}
		}

		if halError.Code == "ValidationError" {
			if err := json.Unmarshal(resBody, &validationError); err != nil {
				return err
			}

			return validationError
		}

		return halError
	}

	if container != nil {
		return json.Unmarshal(resBody, container)
	}

	return nil
}

// Root returns the dwolla root response
func (c *Client) Root() (*Resource, error) {
	if c.root != nil {
		return c.root, nil
	}

	var resource Resource

	if err := c.Get("", nil, nil, &resource); err != nil {
		return nil, err
	}

	c.root = &resource

	return &resource, nil
}
