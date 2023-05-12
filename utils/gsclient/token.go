package gsclient

import (
	"time"
)

// defaultExpiryDelta determines how earlier a token should be considered
// expired than its actual expiration time. It is used to avoid late expirations
// due to client-server time mismatches.
const defaultExpiryDelta = 60 * time.Second

// Token represents the credentials used to authorize the requests to access
// protected resources on the Oauth 2.0 provider's backend.
type Token struct {

	// TokenType is the type of the token (should be "Bearer").
	TokenType    string

	// AccessToken is the token that authorizes and authenticates the requests.
	AccessToken  string

	// CreationTime is the time in which the token was stored in memory.
	CreationTime time.Time

	// Expiry is the optional expiration time of the access token.
	//
	// If zero, the same access token will be reused forever.
	Expiry       time.Time

	// expiryDelta is used internally to calculate when a token is considered
	// expired, by substracting from Expiry.
	expiryDelta  time.Duration
}

// TokenCDO is the raw struct returned by all RFC compliant providers.
type TokenCDO struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// expired reports whether the token is expired. This is determined by comparing the
// current time to t.Expiry (now > t.Expiry - t.ExpiryDelta).
// An Expiry of Zero (January 1, Year 1, 00:00:00) indicates that the token should
// never expire.
func (t *Token) expired() bool {
	if t.Expiry.IsZero() {
		return false
	}
	return t.Expiry.Round(0).Add(-t.expiryDelta).Before(time.Now())
}

// Valid reports whether t is non-nil, has an AccessToken, and this is not expired.
func (t *Token) Valid() bool {
	return t != nil && t.AccessToken != "" && !t.expired()
}