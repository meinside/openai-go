package openai

import (
	"net"
	"net/http"
	"time"
)

const (
	// timeout seconds
	DialTimeout           = 5 * 60 * time.Second
	KeepAlive             = 60 * time.Second
	IdleConnTimeout       = 60 * time.Second
	TLSHandshakeTimeout   = 10 * time.Second
	ResponseHeaderTimeout = DialTimeout
	ExpectContinueTimeout = 1 * time.Second
)

// Client struct which holds its API key, Organization ID, and HTTP client.
type Client struct {
	APIKey         string `json:"api_key"`
	OrganizationID string `json:"organization_id"`

	httpClient *http.Client

	beta *string

	Verbose bool
}

// NewClient returns a new API client
func NewClient(apiKey, organizationID string) *Client {
	return &Client{
		APIKey:         apiKey,
		OrganizationID: organizationID,

		// for reusing http client
		httpClient: &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   DialTimeout,
					KeepAlive: KeepAlive,
				}).DialContext,
				IdleConnTimeout:       IdleConnTimeout,
				TLSHandshakeTimeout:   TLSHandshakeTimeout,
				ResponseHeaderTimeout: ResponseHeaderTimeout,
				ExpectContinueTimeout: ExpectContinueTimeout,
			},
		},
	}
}

// SetBetaHeader sets the beta HTTP header for beta features.
func (c *Client) SetBetaHeader(beta string) *Client {
	c.beta = &beta
	return c
}
