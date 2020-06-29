package godax

import "net/http"

// Client is the main export of godax. All its fields are unexported.
// This file contains all the exported methods available to use.
type Client struct {
	baseRestURL string
	baseWsURL   string
	key         string
	secret      string
	passphrase  string
	httpClient  *http.Client
}

// NewClient returns a godax Client that is hooked up to the live REST and web socket APIs.
func NewClient() (*Client, error) {
	return newClient(false)
}

// NewSandboxClient returns a godax Client that is hooked up to the sandbox REST and web socket APIs.
func NewSandboxClient() (*Client, error) {
	return newClient(true)
}

// ListAccounts gets a list of trading accounts from the profile associated with the API key.
func (c *Client) ListAccounts() ([]Account, error) {
	return c.listAccounts()
}
