package openai

import (
	"net"
	"net/http"
	"time"
)

const (
	// timeout seconds
	DialTimeoutSeconds           = 60
	KeepAliveSeconds             = 120
	IdleConnTimeoutSeconds       = 90
	TLSHandshakeTimeoutSeconds   = DialTimeoutSeconds
	ResponseHeaderTimeoutSeconds = DialTimeoutSeconds
	ExpectContinueTimeoutSeconds = 1
)

// Client struct which holds its API key, Organization ID, and HTTP client.
type Client struct {
	APIKey       string `json:"api_key"`
	Organization string `json:"organization"`

	httpClient *http.Client

	Verbose bool
}

// NewClient returns a new API client
func NewClient(apiKey, organization string) *Client {
	return &Client{
		APIKey:       apiKey,
		Organization: organization,

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
