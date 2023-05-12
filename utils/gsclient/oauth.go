package gsclient

import (
	"bytes"
	"context"
	"fmt"
	"goserver/utils/gsvalidation"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// defaultGrantType specifies the default grant type value to use when requesting
// a new Oauth 2.0 token.
const defaultGrantType = "client_credentials"

// sources is the internal map of token sources. It's initialized empty and
// populated by adding new sources by configuration.
var sources map[string]*TokenSource

// TokenSourceConfig contains the configuration properties needed to add a new 
// token source to the internal map for later use.
type TokenSourceConfig struct {

	// Key is the entry used to set and access the token source in the source map.
	Key          string         `json:"key"`

	// ExpiryDelta is used to configure the token source's expiry delta. If not
	// specified, defaultExpiryDelta is used.
	ExpiryDelta  *time.Duration `json:"expiry_delta"`

	// ClientID is the application's ID.
	ClientID     string         `json:"client_id"`

	// ClientSecret is the application's secret.
	ClientSecret string         `json:"client_secret"`

	// GrantType configures the grant type of the token source.
	GrantType    *string        `json:"grant_type"`

	// Scopes specifies optional requested permissions.
	Scopes       []string       `json:"scopes"`

	// ClientKey is the key that allows to access the client used to request a new
	// token. This client should already be configured when adding this token
	// configuration.
	ClientKey    string         `json:"client_key"`
}

type TokenSource struct {

	// ExpiryDelta is used to calculate when a token is considered expired, by
	// subtracting it from the expiration date (Expiry).
	ExpiryDelta  time.Duration

	// ClientID is the application's ID.
	ClientID     string

	// ClientSecret is the application's secret.
	ClientSecret string

	// GrantType is the grant type value to use when requesting a new token.
	// If not set, defaultGrantType will be used.
	GrantType    string

	// Scopes specifies optional requested permissions.
	Scopes       []string

	// Client is the client used to request a new token.
	Client       Client

	// Token contains the Oauth 2.0 token information with the actual
	// authorization code added in those requests that require authentication.
	Token        *Token

	// mu is a Mutex to synchronize the access to the stored token. This is to
	// prevent that multiple concurrent requests try to generate a new token each.
	mu           sync.Mutex
}

// Initialization. Only initializes the sources map.
func init() {
	sources = make(map[string]*TokenSource)
}

// GetTokenSource retrieves a previously configured TokenSource from the internal map.
// The only way this will return false is by a lack of token source configuration.
func GetTokenSource(key string) (ts *TokenSource, ok bool) {
	ts, ok = sources[key]
	return ts, ok
}

// NewTokenSource receives a TokenSourceConfig struct containing all the information
// needed to create a new TokenSource.
//
// If the given ClientKey doesn't reference an already configured Client, an error 
// is returned.
//
// Note: if there's already an existing TokenSource in the internal map with the same
// key, this is replaced with the new one.
func NewTokenSource(tc TokenSourceConfig) (*TokenSource, error) {
	expiryDelta := defaultExpiryDelta
	if tc.ExpiryDelta != nil {
		expiryDelta = *tc.ExpiryDelta
	}
	grantType := defaultGrantType
	if tc.GrantType != nil {
		grantType = *tc.GrantType
	}
	client, ok := GetClient(tc.ClientKey)
	if !ok {
		return nil, fmt.Errorf("token source configuration failed: couldn't find client with key %s", tc.ClientKey)
	}
	tokenSource := &TokenSource{
		ExpiryDelta: expiryDelta,
		ClientID: tc.ClientID,
		ClientSecret: tc.ClientSecret,
		GrantType: grantType,
		Scopes: tc.Scopes,
		Client: *client,
	}
	sources[tc.Key] = tokenSource
	return tokenSource, nil
}

// ConfigureTokenSources wraps NewTokenSource function allowing the configuration of
// multiple token sources at once.
//
// Function intended to be called from an application's configuration file.
func ConfigureTokenSources(tcs []TokenSourceConfig) error {
	for _, tc := range tcs {
		_, err := NewTokenSource(tc)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetToken retrieves a valid token for the given token source. If there's a saved
// one in memory and it's valid, returns it. Otherwise, calls RenewToken to get
// a new one from the token source.
//
// The access to the stored token is synchronized
//
// Receives a context to pass to RenewToken, if needed.
func (ts *TokenSource) GetToken(ctx context.Context) (*Token, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if ts.Token.Valid() {
		return ts.Token, nil
	}
	return ts.RenewToken(ctx)
}

// RenewToken retrieves a new token from the token source and saves it to memory.
// Uses the Client stored in the token source for the retrieval and the given
// context to create the http.Request (to preserve tracing_id, etc.).
func (ts *TokenSource) RenewToken(ctx context.Context) (*Token, error) {
	
	endpoint := ts.Client.Basepath

	data := url.Values{}
	data.Set("client_id", ts.ClientID)
	data.Set("client_secret", ts.ClientSecret)
	data.Set("grant_type", ts.GrantType)
	data.Set("scope", strings.Join(ts.Scopes, ","))

	reader := bytes.NewBufferString(data.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, reader)
	if err != nil {
		return nil, fmt.Errorf("error creating request to call %s: %s", endpoint, err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := ts.Client.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	var tokenCDO TokenCDO 

	respType, err := gsvalidation.DecodeJSONResponseBody(resp, &tokenCDO, nil)
	switch respType {
	case gsvalidation.OK_RESPONSE:
		creationTime := time.Now()
		token := &Token{
			TokenType: tokenCDO.TokenType,
			AccessToken: tokenCDO.AccessToken,
			CreationTime: creationTime,
			Expiry: creationTime.Add(time.Duration(tokenCDO.ExpiresIn) * time.Second),
		}
		ts.Token = token
		return token, nil
	default:
		return nil, err
	}

}