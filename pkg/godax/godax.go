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
// This endpoint requires either the "view" or "trade" permission.
// This endpoint has a custom rate limit by profile ID: 25 requests per second, up to 50 requests per second in bursts
func (c *Client) ListAccounts() ([]ListAccount, error) {
	return c.listAccounts()
}

// GetAccount retrieves information for a single account.
// Use this endpoint when you know the account_id.
// API key must belong to the same profile as the account.
// This endpoint requires either the "view" or "trade" permission.
func (c *Client) GetAccount(accountID string) (Account, error) {
	return c.getAccount(accountID)
}
