package godax

import (
	"encoding/json"
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

// AccountActivity represents an increase or decrease in your account balance.
/*
   {
       "id": "100",
       "created_at": "2014-11-07T08:19:27.028459Z",
       "amount": "0.001",
       "balance": "239.669",
       "type": "fee",
       "details": {
           "order_id": "d50ec984-77a8-460a-b958-66f114b0de9b",
           "trade_id": "74",
           "product_id": "BTC-USD"
       }
   }
*/
type AccountActivity struct {
	// ID - the account ID associated with the coinbase pro profile
	ID string `json:"id"`
	// CreatedAt - when did this activity happen
	CreatedAt string `json:"created_at"`
	// Amount - the amount used in this activity
	Amount string `json:"amount"`
	// Balance - the total funds available
	Balance string `json:"balance"`
	// Type can be one of the following:
	// "transfer"   - Funds moved to/from Coinbase to Coinbase Pro
	// "match"      - Funds moved as a result of a trade
	// "fee"        - Fee as a result of a trade
	// "rebate"     - Fee rebate as per our fee schedule
	// "conversion"	- Funds converted between fiat currency and a stablecoin/
	Type string `json:"type"`
	// Details - If an entry is the result of a trade (match, fee),
	// the details field will contain additional information about the trade.
	Details ActivityDetail `json:"details"`
}

// ActivityDetail describes important activity metadata (order, trade, and product IDs)
type ActivityDetail struct {
	// OrderID - the order ID related to the activity
	OrderID string `json:"order_id"`
	// TradeID - the trade ID related to the activity
	TradeID string `json:"trade_id"`
	// ProductID - the product ID related to the activity
	ProductID string `json:"product_id"`
}

// listAccounts gets a list of trading accounts from the profile associated with the API key.
func (c *Client) listAccounts(method, path, timestamp, signature string) ([]ListAccount, error) {
	res, err := c.get(path, timestamp, signature)
	if err != nil {
		return []ListAccount{}, err
	}
	defer res.Body.Close()

	var acts []ListAccount
	if err := json.NewDecoder(res.Body).Decode(&acts); err != nil {
		return []ListAccount{}, err
	}

	return acts, nil
}

// getAccount retrieves information for a single account.
func (c *Client) getAccount(accountID, method, path, timestamp, signature string) (Account, error) {
	res, err := c.get(path, timestamp, signature)
	if err != nil {
		return Account{}, err
	}
	defer res.Body.Close()

	var act Account
	if err := json.NewDecoder(res.Body).Decode(&act); err != nil {
		return Account{}, err
	}

	return act, nil
}

// getAccountHistory lists account activity of the API key's profile
func (c *Client) getAccountHistory(accountID, method, path, timestamp, signature string) ([]AccountActivity, error) {
	res, err := c.get(path, timestamp, signature)
	if err != nil {
		return []AccountActivity{}, err
	}
	defer res.Body.Close()

	var aa []AccountActivity
	if err := json.NewDecoder(res.Body).Decode(&aa); err != nil {
		return []AccountActivity{}, err
	}

	return aa, nil
}
