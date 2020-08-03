package godax

// ListAccount represents a trading account for a coinbase pro profile
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

// ActivityDetail describes important activity metadata (order, trade, and product IDs).
type ActivityDetail struct {
	// OrderID - the order ID related to the activity
	OrderID string `json:"order_id"`

	// TradeID - the trade ID related to the activity
	TradeID string `json:"trade_id"`

	// ProductID - the product ID related to the activity
	ProductID string `json:"product_id"`
}

// AccountHold describes a hold on your coinbase pro account.
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

// BankCountry represents a bank country object in a WireDepositInfo.
type BankCountry struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// UserAccount is used for fetching trailing volumes for your user
type UserAccount struct {
	ProductID      string `json:"product_id"`
	ExchangeVolume string `json:"exchange_volume"`
	Volume         string `json:"volume"`
	RecordedAt     string `json:"recorded_at"`
}
