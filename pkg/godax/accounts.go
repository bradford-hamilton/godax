package godax

import (
	"encoding/json"
	"net/http"
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

// CoinbaseAccount represents a Coinbase (non-pro) account.
/*
[
    {
        "id": "fc3a8a57-7142-542d-8436-95a3d82e1622",
        "name": "ETH Wallet",
        "balance": "0.00000000",
        "currency": "ETH",
        "type": "wallet",
        "primary": false,
        "active": true
    },
    {
        "id": "2ae3354e-f1c3-5771-8a37-6228e9d239db",
        "name": "USD Wallet",
        "balance": "0.00",
        "currency": "USD",
        "type": "fiat",
        "primary": false,
        "active": true,
        "wire_deposit_information": {
            "account_number": "0199003122",
            "routing_number": "026013356",
            "bank_name": "Metropolitan Commercial Bank",
            "bank_address": "99 Park Ave 4th Fl New York, NY 10016",
            "bank_country": {
                "code": "US",
                "name": "United States"
            },
            "account_name": "Coinbase, Inc",
            "account_address": "548 Market Street, #23008, San Francisco, CA 94104",
            "reference": "BAOCAEUX"
        }
    },
    {
        "id": "1bfad868-5223-5d3c-8a22-b5ed371e55cb",
        "name": "BTC Wallet",
        "balance": "0.00000000",
        "currency": "BTC",
        "type": "wallet",
        "primary": true,
        "active": true
    },
    {
        "id": "2a11354e-f133-5771-8a37-622be9b239db",
        "name": "EUR Wallet",
        "balance": "0.00",
        "currency": "EUR",
        "type": "fiat",
        "primary": false,
        "active": true,
        "sepa_deposit_information": {
            "iban": "EE957700771001355096",
            "swift": "LHVBEE22",
            "bank_name": "AS LHV Pank",
            "bank_address": "Tartu mnt 2, 10145 Tallinn, Estonia",
            "bank_country_name": "Estonia",
            "account_name": "Coinbase UK, Ltd.",
            "account_address": "9th Floor, 107 Cheapside, London, EC2V 6DN, United Kingdom",
            "reference": "CBAEUXOVFXOXYX"
        }
    },
    ...
]
*/
type CoinbaseAccount struct {
	ID                     string          `json:"id"`
	Name                   string          `json:"name"`
	Balance                string          `json:"balance"`
	Currency               string          `json:"currency"`
	Type                   string          `json:"type"`
	Primary                bool            `json:"primary"`
	Active                 bool            `json:"active"`
	WireDepositInformation WireDepositInfo `json:"wire_deposit_information"`
	SepaDepositInformation SepaDepositInfo `json:"sepa_deposit_information"`
	UKDepositInformation   UKDepositInfo   `json:"uk_deposit_information"`

	// These came back on response, but were not in docs, possibly notify coinbase to update
	HoldBalance         string `json:"hold_balance"`
	HoldCurrency        string `json:"hold_currency"`
	AvailableOnConsumer bool   `json:"available_on_consumer"`
}

// WireDepositInfo describes all the metadata needed for a wire deposit to a bank.
type WireDepositInfo struct {
	AccountNumber  string      `json:"account_number"`
	RoutingNumber  string      `json:"routing_number"`
	BankName       string      `json:"bank_name"`
	BankAddress    string      `json:"bank_address"`
	BankCountry    BankCountry `json:"bank_country"`
	AccountName    string      `json:"account_name"`
	AccountAddress string      `json:"account_address"`
	Reference      string      `json:"reference"`
}

// SepaDepositInfo describes all the metadata needed for a Sepa desposit.
type SepaDepositInfo struct {
	IBAN            string `json:"iban"`
	Swift           string `json:"swift"`
	BankName        string `json:"bank_name"`
	BankAddress     string `json:"bank_address"`
	BankCountryName string `json:"bank_country_name"`
	AccountName     string `json:"account_name"`
	AccountAddress  string `json:"account_address"`
	Reference       string `json:"reference"`
}

// UKDepositInfo describes all the metadata needed for a UK desposit.
type UKDepositInfo struct {
	SortCode      string `json:"sort_code"`
	AccountName   string `json:"account_name"`
	AccountNumber string `json:"account_number"`
	BankName      string `json:"bank_name"`
	Reference     string `json:"reference"`
}

// BankCountry represents a bank country object in a WireDepositInfo
type BankCountry struct {
	// Code defines the country code
	Code string `json:"code"`

	// Name defines the country name
	Name string `json:"name"`
}

// UserAccount is used for fetching trailing volumes for your user
/*
[
    {
        "product_id": "BTC-USD",
        "exchange_volume": "11800.00000000",
        "volume": "100.00000000",
        "recorded_at": "1973-11-29T00:05:01.123456Z"
    },
    {
        "product_id": "LTC-USD",
        "exchange_volume": "51010.04100000",
        "volume": "2010.04100000",
        "recorded_at": "1973-11-29T00:05:02.123456Z"
    }
]
*/
type UserAccount struct {
	ProductID      string `json:"product_id"`
	ExchangeVolume string `json:"exchange_volume"`
	Volume         string `json:"volume"`
	RecordedAt     string `json:"recorded_at"`
}

// listAccounts gets a list of trading accounts from the profile associated with the API key.
func (c *Client) listAccounts(timestamp, signature string, req *http.Request) ([]ListAccount, error) {
	res, err := c.do(timestamp, signature, req)
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
func (c *Client) getAccount(timestamp, signature string, req *http.Request) (Account, error) {
	res, err := c.do(timestamp, signature, req)
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
func (c *Client) getAccountHistory(timestamp, signature string, req *http.Request) ([]AccountActivity, error) {
	res, err := c.do(timestamp, signature, req)
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
func (c *Client) getAccountHolds(timestamp, signature string, req *http.Request) ([]AccountHold, error) {
	res, err := c.do(timestamp, signature, req)
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

// listCoinbaseAccounts gets a list of trading accounts from the profile associated with the API key.
func (c *Client) listCoinbaseAccounts(timestamp, signature string, req *http.Request) ([]CoinbaseAccount, error) {
	res, err := c.do(timestamp, signature, req)
	if err != nil {
		return []CoinbaseAccount{}, err
	}
	defer res.Body.Close()

	var cbAccounts []CoinbaseAccount
	if err := json.NewDecoder(res.Body).Decode(&cbAccounts); err != nil {
		return []CoinbaseAccount{}, err
	}
	return cbAccounts, nil
}

// getTrailingVolume returns your 30-day trailing volume for all products of the API key's profile.
func (c *Client) getTrailingVolume(timestamp, signature string, req *http.Request) ([]UserAccount, error) {
	res, err := c.do(timestamp, signature, req)
	if err != nil {
		return []UserAccount{}, err
	}
	defer res.Body.Close()

	var userActs []UserAccount
	if err := json.NewDecoder(res.Body).Decode(&userActs); err != nil {
		return []UserAccount{}, err
	}
	return userActs, nil
}
