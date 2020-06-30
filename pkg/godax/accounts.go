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

// AccountHold ...
/*
	{
        "id": "82dcd140-c3c7-4507-8de4-2c529cd1a28f",
        "account_id": "e0b3f39a-183d-453e-b754-0c13e5bab0b3",
        "created_at": "2014-11-06T10:34:47.123456Z",
        "updated_at": "2014-11-06T10:40:47.123456Z",
        "amount": "4.23",
        "type": "order",
        "ref": "0a205de4-dd35-4370-a285-fe8fc375a273",
    }
*/
type AccountHold struct {
	// ID - the hold ID
	ID string `json:"id"`

	// Account ID - the account ID associated with the coinbase pro profile
	AccountID string `json:"account_id"`

	// CreatedAt - when this hold happened
	CreatedAt string `json:"created_at"`

	// Updated - the last time this hold was updated
	UpdatedAt string `json:"updated_at"`

	// Amount - the amount in the hold
	Amount string `json:"amount"`

	// Type - type of the hold will indicate why the hold exists. The hold type is order
	// for holds related to open orders and transfer for holds related to a withdraw.
	Type string `json:"type"`

	// Ref - The ref field contains the id of the order or transfer which created the hold.
	Ref string `json:"ref"`
}

// listAccounts gets a list of trading accounts from the profile associated with the API key.
func (c *Client) listAccounts(method, path, timestamp, signature string) ([]ListAccount, error) {
	res, err := c.get(path, timestamp, signature)
	if err != nil {
		return []ListAccount{}, err
	}
	defer res.Body.Close()

	var accounts []ListAccount
	if err := json.NewDecoder(res.Body).Decode(&accounts); err != nil {
		return []ListAccount{}, err
	}

	return accounts, nil
}

// getAccount retrieves information for a single account.
func (c *Client) getAccount(accountID, method, path, timestamp, signature string) (Account, error) {
	res, err := c.get(path, timestamp, signature)
	if err != nil {
		return Account{}, err
	}
	defer res.Body.Close()

	var account Account
	if err := json.NewDecoder(res.Body).Decode(&account); err != nil {
		return Account{}, err
	}

	return account, nil
}

// getAccountHistory lists account activity of the API key's profile
func (c *Client) getAccountHistory(accountID, method, path, timestamp, signature string) ([]AccountActivity, error) {
	res, err := c.get(path, timestamp, signature)
	if err != nil {
		return []AccountActivity{}, err
	}
	defer res.Body.Close()

	var activities []AccountActivity
	if err := json.NewDecoder(res.Body).Decode(&activities); err != nil {
		return []AccountActivity{}, err
	}

	return activities, nil
}

// getAccountHolds lists holds of an account that belong to the same profile as the API key
func (c *Client) getAccountHolds(accountID, method, path, timestamp, signature string) ([]AccountHold, error) {
	res, err := c.get(path, timestamp, signature)
	if err != nil {
		return []AccountHold{}, err
	}
	defer res.Body.Close()

	var holds []AccountHold
	if err := json.NewDecoder(res.Body).Decode(&holds); err != nil {
		return []AccountHold{}, err
	}

	return holds, nil
}
