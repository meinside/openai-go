package openai

import (
	"net"
	"net/http"
	"time"
)

const (
	// timeout seconds
	DialTimeoutSeconds           = 180
	KeepAliveSeconds             = 60
	IdleConnTimeoutSeconds       = 60
	TLSHandshakeTimeoutSeconds   = 10
	ResponseHeaderTimeoutSeconds = DialTimeoutSeconds
	ExpectContinueTimeoutSeconds = 1
)

// Client struct which holds its API key, Organization ID, and HTTP client.
type Client struct {
	APIKey         string `json:"api_key"`
	OrganizationID string `json:"organization_id"`

	httpClient *http.Client

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
					Timeout:   DialTimeoutSeconds * time.Second,
					KeepAlive: KeepAliveSeconds * time.Second,
				}).DialContext,
				IdleConnTimeout:       IdleConnTimeoutSeconds * time.Second,
				TLSHandshakeTimeout:   TLSHandshakeTimeoutSeconds * time.Second,
				ResponseHeaderTimeout: ResponseHeaderTimeoutSeconds * time.Second,
				ExpectContinueTimeout: ExpectContinueTimeoutSeconds * time.Second,
			},
		},
	}
}
