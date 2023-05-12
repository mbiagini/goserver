package gsclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"goserver/utils/gsmiddleware"
	"io"
	"net/http"
	"net/url"
	"time"
)

// DefaultTimeout specifies the default time limit for requests made by a Client.
// The timeout includes connection time, any redirects, and reading the response body.
var defaultTimeout = 30 * time.Second

// clients contains all configured clients.
var clients map[string]*Client

// ClientConfig contains all the information needed to construct a new Client.
type ClientConfig struct {

	// Key used in the internal clients map to store and retrieve the client.
	Key            string             `json:"key"`

	// Basepath to be used by the client in every of its http requests.
	Basepath       string             `json:"basepath"`

	// Transport specifies any middleware used for the outgoing http requests.
	// If not defined, the DefaultTransport function will set this field.
	Transport      *http.RoundTripper

	// Timeout specifies the duration in seconds to wait for an http response
	// from the client.
	// If not set, defaultTimeout will be used.
	Timeout        *int               `json:"timeout"`

	// DefaultHeaders define headers to be sent in every request.
	DefaultHeaders *http.Header       `json:"default_headers"`

	// DefaultParams define query parameters to be sent in every request.
	DefaultParams  *url.Values        `json:"default_params"`

	// SkipSSL may be used to skip ssl verify when calling an https endpoint
	// whose certificate we can't validate.
	// This should be used cautiously.
	SkipSSL        bool               `json:"skip_ssl"`

	// TokenSourceKey is used to configure a token source in the client. This
	// token source should already be configured when adding this client.
	// 
	// If the client does not required Oauth 2.0 authentication, this field
	// must be omited.
	TokenSourceKey *string            `json:"token_source_key"`
}

// Client contains all the information needed to interact with an external API via
// HTTP protocol.
type Client struct {

	// Key used in the internal clients map to store and retrieve the client.
	Key            string

	// Basepath to be used by the client in every of its http requests.
	Basepath       string

	// HttpClient is the actual http.Client used to make http calls.
	HttpClient     http.Client

	// DefaultHeaders define headers to be sent in every request.
	DefaultHeaders http.Header

	// DefaultParams define query parameters to be sent in every request.
	DefaultParams  url.Values

	// TokenSource is used to set an ouath 2.0 token provider, which will be used
	// to authenticate all http requests made to this client. Once the token source
	// is configured, the token handling is automatically.
	TokenSource    *TokenSource
}

// Initialization. Only initializes the clients map.
func init() {
	clients = make(map[string]*Client)
}

// GetClient retrieves a previously configured Client from the internal map. The
// only way this will return false is by a lack of client configuration.
func GetClient(key string) (c *Client, ok bool) {
	c, ok = clients[key]
	return c, ok
}

// NewClient receives a ClientConfig struct containing all the information needed to
// create a new Client. If no transport is provided, the DefaultTransport is used. 
// Likewise, if no timeout is given, defaultTimeout is used.
// Once the Client is created, it is added to the internal map with the provided
// key.
//
// If the TokenSourceKey is present and doesn't reference an already configured
// token source, an error is returned.
//
// Note: if there's already an existing Client in the internal map with the same
// key, this is replaced with the new one.
func NewClient(cc ClientConfig) (*Client, error) {

	// Get token source by its key, if any.
	var tokenSource *TokenSource
	if cc.TokenSourceKey != nil {
		ts, ok := GetTokenSource(*cc.TokenSourceKey)
		if !ok {
			return nil, fmt.Errorf("client configuration failed: couldn't find token source with key %s", *cc.TokenSourceKey)
		}
		tokenSource = ts
	}

	transport := DefaultTransport(cc.SkipSSL, tokenSource)
	if cc.Transport != nil {
		transport = *cc.Transport
	}
	timeout := defaultTimeout
	if cc.Timeout != nil {
		timeout = time.Duration(*cc.Timeout) * time.Second
	}
	headers := make(map[string][]string)
	if cc.DefaultHeaders != nil {
		headers = *cc.DefaultHeaders
	}
	params := make(map[string][]string)
	if cc.DefaultParams != nil {
		params = *cc.DefaultParams
	}
	client := &Client{
		Key: cc.Key,
		Basepath: cc.Basepath,
		HttpClient: http.Client{
			Transport: transport,
			Timeout: timeout,
		},
		DefaultHeaders: headers,
		DefaultParams: params,
		TokenSource: tokenSource,
	}

	// Save (or overwrite) new client to internal map.
	clients[cc.Key] = client

	return client, nil
}

// ConfigureClients wraps NewClient function allowing the configuration of
// multiple clients at once.
//
// Function intended to be called from an application's configuration file.
func ConfigureClients(ccs []ClientConfig) error {
	for _, cc := range ccs {
		_, err := NewClient(cc)
		if err != nil {
			return err
		}
	}
	return nil
}

// DefaultTransport return the default http.RoundTripper implementation of this
// package. This includes a logging, a tracing, and an oauth transports.
//
// SSL verification can be skipped by sending "skipSSL" true.
// To enable oauth transport, a token source must be provided.
func DefaultTransport(skipSSL bool, ts *TokenSource) http.RoundTripper {
	base := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSSL},
	}
	logging := &gsmiddleware.HttpLogTransport{
		Base: *base,
	}
	tracing := &gsmiddleware.TracingTransport{
		Base: logging,
	}
	if ts == nil {
		return tracing
	}
	return &OauthTransport{
		Base: tracing,
		TokenSource: ts,
	}
}

// SetRequestDefaults adds the default headers and query parameters
// configured for the given Client, if any, to an http.Request struct.
func (c *Client) SetRequestDefaults(r http.Request) {
	
	// Add default headers, if any.
	for k, values := range c.DefaultHeaders {
		for _, v := range values {
			r.Header.Set(k, v)
		}
	}

	// Get the existing query parameters from the URL.
	queryParams, _ := url.ParseQuery(r.URL.RawQuery)

	// Add default query parameters, if any.
	for k, values := range c.DefaultParams {
		for _, v := range values {
			queryParams.Add(k, v)
		}
	}

	// Assign encoded query string to http request.
	r.URL.RawQuery = queryParams.Encode()
}

// NewJSONRequest wraps http.NewRequestWithContext, adding JSON marshalling
// instead of receiving an io.Reader directly.
// In addition, receives only the operation path to be invoked, appending it
// to the client's basepath.
func (c *Client) NewJSONRequest(ctx context.Context, method string, path string, v any) (*http.Request, error) {

	var reader io.Reader
	if v != nil {
		byteSlice, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(byteSlice)
	}

	url := c.Basepath + path
	return http.NewRequestWithContext(ctx, method, url, reader)
}

// NewRequest wraps http.NewRequestWithContext receiving only the operation 
// path to be invoked, appending it to the client's basepath.
func (c *Client) NewRequest(ctx context.Context, method string, path string) (*http.Request, error) {
	return c.NewJSONRequest(ctx, method, path, nil)
}