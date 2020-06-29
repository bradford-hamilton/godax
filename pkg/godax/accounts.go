package godax

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// Account represents a trading account for a coinbase pro profile.
/*
	{
		"id": "71452118-efc7-4cc4-8780-a5e22d4baa53",
		"currency": "BTC",
		"balance": "0.0000000000000000",
		"available": "0.0000000000000000",
		"hold": "0.0000000000000000",
		"profile_id": "75da88c5-05bf-4f54-bc85-5c775bd68254"
	}
*/
type Account struct {
	// ID - the account ID associated with the coinbase pro profile
	ID string `json:"id"`
	// Currency - the currency of the account
	Currency string `json:"currency"`
	// Balance - the total funds in the account
	Balance string `json:"balance"`
	// Available - funds available to withdraw or trade
	Available string `json:"available"`
	// Hold - funds on hold (not available for use)
	Hold string `json:"hold"`
}

// listAccounts gets a list of trading accounts from the profile associated with the API key.
func (c *Client) listAccounts() ([]Account, error) {
	path := "/accounts"

	req, err := http.NewRequest(http.MethodGet, c.baseRestURL+path, nil)
	if err != nil {
		return nil, err
	}

	ts := strconv.FormatInt(time.Now().Unix(), 10)
	sig, err := c.generateSignature(ts, path, http.MethodGet, "")
	if err != nil {
		return nil, err
	}

	c.setHeaders(req, ts, sig)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var accounts []Account
	json.NewDecoder(res.Body).Decode(&accounts)

	return accounts, nil
}
