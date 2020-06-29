package godax

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// ListAccount represents a trading account for a coinbase pro profile.
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
type ListAccount struct {
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

// Account describes information for a single account
/*
	{
		"id": "a1b2c3d4",
		"balance": "1.100",
		"holds": "0.100",
		"available": "1.00",
		"currency": "USD"
	}
*/
type Account struct {
	// ID - the account ID associated with the coinbase pro profile
	ID string `json:"id"`
	// Balance - the total funds in the account
	Balance string `json:"balance"`
	// Holds - funds on hold (not available for use)
	Holds string `json:"holds"`
	// Available - funds available to withdraw or trade
	Available string `json:"available"`
	// Currency - the currency of the account
	Currency string `json:"currency"`
}

// listAccounts gets a list of trading accounts from the profile associated with the API key.
func (c *Client) listAccounts() ([]ListAccount, error) {
	path := "/accounts"

	req, err := http.NewRequest(http.MethodGet, c.baseRestURL+path, nil)
	if err != nil {
		return []ListAccount{}, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sig, err := c.generateSignature(timestamp, path, http.MethodGet, "")
	if err != nil {
		return []ListAccount{}, err
	}

	c.setHeaders(req, timestamp, sig)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return []ListAccount{}, err
	}
	defer res.Body.Close()

	var acts []ListAccount
	json.NewDecoder(res.Body).Decode(&acts)

	return acts, nil
}

// getAccount retrieves information for a single account.
func (c *Client) getAccount(accountID string) (Account, error) {
	path := "/accounts/" + accountID

	req, err := http.NewRequest(http.MethodGet, c.baseRestURL+path, nil)
	if err != nil {
		return Account{}, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	sig, err := c.generateSignature(timestamp, path, http.MethodGet, "")
	if err != nil {
		return Account{}, err
	}

	c.setHeaders(req, timestamp, sig)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return Account{}, err
	}
	defer res.Body.Close()

	var act Account
	json.NewDecoder(res.Body).Decode(&act)

	return act, nil
}
