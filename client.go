package dwolla

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Environment is a supported dwolla environment. e.g. sandbox, production, https://api-sandbox.dwolla.com
type Environment string

const (
	// Version is the version of the client
	Version string = "0.1.0"

	HeaderIdempotency = "Idempotency-Key"

	// Production is the production environment
	Production Environment = "production"
	// Sandbox is the sanbox environment
	Sandbox Environment = "sandbox"

	// ProductionAPIURL is the production api url
	ProductionAPIURL = "https://api.dwolla.com"
	// ProductionAuthURL is the production auth url
	ProductionAuthURL = "https://www.dwolla.com/oauth/v2/authenticate"
	// ProductionTokenURL is the production token url
	// Deprecated - use https://api.dwolla.com/token moving forward
	ProductionTokenURL = "https://accounts.dwolla.com/token" // #nosec G101

	// SandboxAPIURL is the sandbox api url
	SandboxAPIURL = "https://api-sandbox.dwolla.com"
	// SandboxAuthURL is the sandbox auth url
	SandboxAuthURL = "https://sandbox.dwolla.com/oauth/v2/authenticate"
	// SandboxTokenURL is the sandbox token url
	// Deprecated - use https://api-sandbox.dwolla.com moving forward
	SandboxTokenURL = "https://accounts-sandbox.dwolla.com/token" // #nosec G101
)

// Token is a dwolla auth token
type Token struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	TokenType        string `json:"token_type"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	startTime        time.Time
}

// Expired returns true if token has expired
func (t *Token) Expired() bool {
	return time.Since(t.startTime) > time.Duration(t.ExpiresIn)*time.Second
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
	BeneficialOwner        BeneficialOwnerService
	BusinessClassification BusinessClassificationService
	Customer               CustomerService
	Document               DocumentService
	Event                  EventService
	FundingSource          FundingSourceService
	KBA                    KBAService
	MassPayment            MassPaymentService
	OnDemandAuthorization  OnDemandAuthorizationService
	Transfer               TransferService
	TransferFailure        *TransferFailureServiceOp
	Webhook                WebhookService
	WebhookSubscription    WebhookSubscriptionService
}

// ClientTokenRequest is a client token request
type ClientTokenRequest struct {
	Resource
	Action string `json:"action"`
}

// ClientToken is a general use client token
type ClientToken struct {
	Token string `json:"token"`
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
	c.BeneficialOwner = &BeneficialOwnerServiceOp{c}
	c.BusinessClassification = &BusinessClassificationServiceOp{c}
	c.Customer = &CustomerServiceOp{c}
	c.Document = &DocumentServiceOp{c}
	c.Event = &EventServiceOp{c}
	c.FundingSource = &FundingSourceServiceOp{c}
	c.KBA = &KBAServiceOp{c}
	c.MassPayment = &MassPaymentServiceOp{c}
	c.OnDemandAuthorization = &OnDemandAuthorizationServiceOp{c}
	c.Transfer = &TransferServiceOp{c}
	c.TransferFailure = &TransferFailureServiceOp{client: c}
	c.Webhook = &WebhookServiceOp{c}
	c.WebhookSubscription = &WebhookSubscriptionServiceOp{c}

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
	default:
		url = string(c.Environment)
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
	default:
		url = joinURL(string(c.Environment), "/oauth/v2/authenticate")
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
	default:
		url = joinURL(string(c.Environment), "/token")
	}

	return url
}

// RequestToken requests a new auth token using client credentials
//
// See: https://developers.dwolla.com/api-reference/authorization/application-authorization
func (c *Client) RequestToken(ctx context.Context) error {
	var (
		err   error
		token Token
	)

	buf := bytes.NewBuffer([]byte("grant_type=client_credentials"))

	req, err := http.NewRequestWithContext(ctx, "POST", c.BuildAPIURL("token"), buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
	c.Token.startTime = time.Now()

	return nil
}

// EnsureToken ensures that a token exists for a request
func (c *Client) EnsureToken(ctx context.Context) error {
	if c.Token == nil {
		if err := c.RequestToken(ctx); err != nil {
			return err
		}
	}

	if c.Token.Expired() {
		if err := c.RequestToken(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Get performs a GET against the api
func (c *Client) Get(ctx context.Context, path string, params *url.Values, headers *http.Header, container interface{}) error {
	var (
		err      error
		halError HALError
	)

	if err = c.EnsureToken(ctx); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", c.BuildAPIURL(path), nil)
	if err != nil {
		return err
	}

	if headers != nil {
		req.Header = *headers
	}

	req.Header.Set("Accept", "application/vnd.dwolla.v1.hal+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token.AccessToken))

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

	if res.StatusCode > 299 {
		if err := json.Unmarshal(resBody, &halError); err != nil {
			return err
		}

		// If token is expired, we'll attempt to get a newone and reattempt
		// the request. This should probably be moved to a method to handle
		// all error scenarios.
		if halError.Code == "ExpiredAccessToken" {
			if err := c.RequestToken(ctx); err != nil {
				return err
			}

			return c.Get(ctx, path, params, headers, container)
		}

		return halError
	}

	if res.StatusCode == 202 {
		if err := json.Unmarshal(resBody, &halError); err != nil {
			return err
		}

		if halError.Code == "TryAgainLater" {
			// Get "TryAgainLater" code when Micro-deposits have not settled to destination bank.
			// A Customer can verify these amounts after micro-deposits have processed to their bank.
			return halError
		}
	}

	if container != nil {
		return json.Unmarshal(resBody, container)
	}

	return nil
}

// Post performs a POST against the api
func (c *Client) Post(ctx context.Context, path string, body interface{}, headers *http.Header, container interface{}) error {
	var (
		err             error
		halError        HALError
		validationError ValidationError
		bodyReader      io.Reader
	)

	if err = c.EnsureToken(ctx); err != nil {
		return err
	}

	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}

		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BuildAPIURL(path), bodyReader)
	if err != nil {
		return err
	}

	if headers != nil {
		req.Header = *headers
	}

	req.Header.Set("Accept", "application/vnd.dwolla.v1.hal+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token.AccessToken))
	req.Header.Set("Content-Type", "application/vnd.dwolla.v1.hal+json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// When creating a resource, Dwolla will return a 201 and a "Location"
	// header. This just cuts to the chase and retrieves the resource.
	if res.Header.Get("Location") != "" {
		return c.Get(ctx, res.Header.Get("Location"), nil, nil, container)
	}

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
			if err := c.RequestToken(ctx); err != nil {
				return err
			}

			return c.Post(ctx, path, body, headers, container)
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

// Upload performs a multipart file upload to the Dwolla API
func (c *Client) Upload(ctx context.Context, path string, documentType DocumentType, fileName string, file io.Reader, container interface{}) error {
	var (
		err             error
		halError        HALError
		validationError ValidationError
	)

	if err = c.EnsureToken(ctx); err != nil {
		return err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	err = writer.WriteField("documentType", string(documentType))
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BuildAPIURL(path), body)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.dwolla.v1.hal+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token.AccessToken))
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// When creating a resource, Dwolla will return a 201 and a "Location"
	// header. This just cuts to the chase and retrieves the resource.
	if res.Header.Get("Location") != "" {
		return c.Get(ctx, res.Header.Get("Location"), nil, nil, container)
	}

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
			if err := c.RequestToken(ctx); err != nil {
				return err
			}

			return c.Upload(ctx, path, documentType, fileName, file, container)
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

// Delete performs a DELETE against the api
func (c *Client) Delete(ctx context.Context, path string, params *url.Values, headers *http.Header) error {
	var (
		err      error
		halError HALError
	)

	if err = c.EnsureToken(ctx); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", c.BuildAPIURL(path), nil)
	if err != nil {
		return err
	}

	if headers != nil {
		req.Header = *headers
	}

	req.Header.Set("Accept", "application/vnd.dwolla.v1.hal+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token.AccessToken))

	if params != nil {
		req.URL.RawQuery = params.Encode()
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode > 299 {
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(resBody, &halError); err != nil {
			return err
		}

		// If token is expired, we'll attempt to get a newone and reattempt
		// the request. This should probably be moved to a method to handle
		// all error scenarios.
		if halError.Code == "ExpiredAccessToken" {
			if err := c.RequestToken(ctx); err != nil {
				return err
			}

			return c.Delete(ctx, path, params, headers)
		}

		return halError
	}

	return nil
}

// Root returns the dwolla root response
func (c *Client) Root(ctx context.Context) (*Resource, error) {
	if c.root != nil {
		return c.root, nil
	}

	var resource Resource

	if err := c.Get(ctx, "", nil, nil, &resource); err != nil {
		return nil, err
	}

	c.root = &resource
	return &resource, nil
}

// SandboxSimulations simulates events within the sandbox environment
//
// see: https://developers.dwolla.com/resources/testing.html#simulate-bank-transfer-processing
func (c *Client) SandboxSimulations(ctx context.Context) error {
	return c.Post(ctx, "sandbox-simulations", nil, nil, nil)
}

// CreateClientToken creates a general use client token
//
// see: https://docsv2.dwolla.com/#create-a-client-token
func (c *Client) CreateClientToken(ctx context.Context, action string, customer *Customer) (*ClientToken, error) {
	body := ClientTokenRequest{Action: action}

	if customer != nil {
		if _, ok := customer.Links["self"]; !ok {
			return nil, errors.New("No self resource link")
		}

		body.Resource = *NewResource(Links{"customer": Link{Href: customer.Links["self"].Href}}, nil)
	}

	var token ClientToken

	if err := c.Post(ctx, "client-tokens", body, nil, &token); err != nil {
		return nil, err
	}

	return &token, nil
}
