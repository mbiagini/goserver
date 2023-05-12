package gsclient

import (
	"net/http"
)

// OauthTransport wraps an http.RoundTripper adding oauth authentication
type OauthTransport struct {

	// Base defines the implementation of http.RoundTripper wrapped by this
	// Transport. Once the oauth token is added, this base transport will be
	// called.
	Base        http.RoundTripper

	// TokenSource is the token provider that will be used by this Transport
	// to add oauth 2.0 authentication to every request.
	TokenSource *TokenSource
}

// Implements http.RoundTripper so it can be used as a Transport.
// The http.Request's context must contain a token source set. Otherwise, using
// this transport will result in a runtime error.
func (ot *OauthTransport) RoundTrip(r *http.Request) (*http.Response, error) {

	// Get saved token on token source or a new one.
	token, err := ot.TokenSource.GetToken(r.Context())
	if err != nil {
		return nil, err
	}
	r.Header.Set("Authorization", token.TokenType + " " + token.AccessToken)
	return ot.Base.RoundTrip(r)
}