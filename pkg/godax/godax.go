package godax

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

// Client is the main export of godax. All its fields are unexported.
// This file contains all the exported methods available to use.
type Client struct {
	baseRestURL string
	baseWsURL   string
	key         string
	secret      string
	passphrase  string
	httpClient  HTTPClient
}

// Param is a type alias for a string. The hope is that godax can offer all of the available
// query params as constants, but also if you need to - you can create your own Param to use.
type Param string

// Available query params for all coinbase pro calls
const (
	OrderID   Param = "order_id"
	ProductID Param = "product_id"
	Status    Param = "status"
)

var noBody = []byte{}

// QueryParams represent the available query params for any given coinbase pro call.
type QueryParams map[Param]string

// ErrMissingOrderOrProductID TODO: does this feel weird and one offy right now?
var (
	ErrMissingOrderOrProductID = errors.New("please provide either an order_id or product_id in your query params")
	ErrMissingConversionParams = errors.New("please provide all of the following params: to, from, amount")
)

// NewClient returns a godax Client that is hooked up to the live REST and web socket APIs.
func NewClient() (*Client, error) {
	return newClient(false)
}

// NewSandboxClient returns a godax Client that is hooked up to the sandbox REST and web socket APIs.
func NewSandboxClient() (*Client, error) {
	return newClient(true)
}

func unixTime() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// ListAccounts gets a list of trading accounts from the profile associated with the API key.
// This endpoint requires either the "view" or "trade" permission. This endpoint has a custom
// rate limit by profile ID: 25 requests per second, up to 50 requests per second in bursts
func (c *Client) ListAccounts() ([]ListAccount, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/accounts"

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, nil)
	if err != nil {
		return nil, err
	}

	return c.listAccounts(timestamp, sig, req)
}

// GetAccount retrieves information for a single account. Use this endpoint when you know the
// account_id. API key must belong to the same profile as the account. This endpoint requires
// either the "view" or "trade" permission.
func (c *Client) GetAccount(accountID string) (Account, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/accounts/" + accountID

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, nil)
	if err != nil {
		return Account{}, err
	}

	return c.getAccount(timestamp, sig, req)
}

// GetAccountHistory lists account activity of the API key's profile. Account activity either increases
// or decreases your account balance. If an entry is the result of a trade (match, fee), the details
// field on an AccountActivity will contain additional information about the trade. Items are paginated
// and sorted latest first. This endpoint requires either the "view" or "trade" permission.
// TODO: paginate
func (c *Client) GetAccountHistory(accountID string) ([]AccountActivity, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/accounts/" + accountID + "/ledger"

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, nil)
	if err != nil {
		return nil, err
	}

	return c.getAccountHistory(timestamp, sig, req)
}

// GetAccountHolds lists holds of an account that belong to the same profile as the API key.
// Holds are placed on an account for any active orders or pending withdraw requests. As an
// order is filled, the hold amount is updated. If an order is canceled, any remaining hold
// is removed. For a withdraw, once it is completed, the hold is removed. This endpoint
// requires either the "view" or "trade" permission.
// TODO: paginate
func (c *Client) GetAccountHolds(accountID string) ([]AccountHold, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/accounts/" + accountID + "/holds"

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, nil)
	if err != nil {
		return nil, err
	}

	return c.getAccountHolds(timestamp, sig, req)
}

// PlaceOrder allows you to place two types of orders: limit and market. Orders can only be
// placed if your account has sufficient funds. Once an order is placed, your account funds
// will be put on hold for the duration of the order. How much and which funds are put on
// hold depends on the order type and parameters specified. This endpoint requires the
// "trade" permission.
func (c *Client) PlaceOrder(order OrderParams) (Order, error) {
	timestamp := unixTime()
	method := http.MethodPost
	path := "/orders"

	body, err := json.Marshal(order)
	if err != nil {
		return Order{}, err
	}

	req, sig, err := c.createAndSignRequest(timestamp, method, path, body, nil)
	if err != nil {
		return Order{}, err
	}

	return c.placeOrder(timestamp, sig, req)
}

// CancelOrderByID cancels a previously placed order. Order must belong to the profile that
// the API key belongs to. If the order had no matches during its lifetime its record may be
// purged. This means the order details will not be available with GetOrderByID or GetOrderByClientOID.
// The product ID of the order is not required so if you don't have it you can pass nil here.
// The request will be more performant if you include it. This endpoint requires the "trade" permission.
func (c *Client) CancelOrderByID(orderID string, qp QueryParams) (canceledOrderID string, err error) {
	timestamp := unixTime()
	method := http.MethodDelete
	path := "/orders/" + orderID

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, &qp)
	if err != nil {
		return "", err
	}

	return c.cancelOrder(timestamp, sig, req)
}

// CancelOrderByClientOID cancels a previously placed order. Order must belong to the profile that
// the API key belongs to. If the order had no matches during its lifetime its record may be
// purged. This means the order details will not be available with GetOrderByID or GetOrderByClientOID.
// The product ID of the order is not required so if you don't have it you can pass nil here.
// The request will be more performant if you include it. This endpoint requires the "trade" permission.
func (c *Client) CancelOrderByClientOID(clientOID string, qp QueryParams) (canceledOrderID string, err error) {
	timestamp := unixTime()
	method := http.MethodDelete
	path := "/orders/client:" + clientOID

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, &qp)
	if err != nil {
		return "", err
	}

	return c.cancelOrder(timestamp, sig, req)
}

// CancelAllOrders cancel all open orders from the profile that the API key belongs to. The response is
// a list of ids of the canceled orders. This endpoint requires the "trade" permission. The productID
// param is opitonal and a pointer, so you can pass nil. If you do provide the productID here, you will
// only cancel orders open for that specific product.
func (c *Client) CancelAllOrders(qp QueryParams) (canceledOrderIDs []string, err error) {
	timestamp := unixTime()
	method := http.MethodDelete
	path := "/orders"

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, &qp)
	if err != nil {
		return nil, err
	}

	return c.cancelAllOrders(timestamp, sig, req)
}

// ListOrders lists your current open orders from the profile that the API key belongs to. Only open or un-settled
// orders are returned. As soon as an order is no longer open and settled, it will no longer appear in the default
// request. This endpoint requires either the "view" or "trade" permission.
// Valid status args to filter return orders: [open, pending, active].
//
// Orders which are no longer resting on the order book, will be marked with the done status. There is a small window
// between an order being done and settled. An order is settled when all of the fills have settled and the remaining
// holds (if any) have been removed.
//
// For high-volume trading it is strongly recommended that you maintain your own list
// of open orders and use one of the streaming market data feeds to keep it updated. You should poll the open orders
// endpoint once when you start trading to obtain the current state of any open orders. executed_value is the cumulative
// match size * price and is only present for orders placed after 2016-05-20. Open orders may change state between the
// request and the response depending on market conditions.
func (c *Client) ListOrders(qp QueryParams) ([]Order, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/orders"

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, &qp)
	if err != nil {
		return nil, err
	}

	return c.listOrders(timestamp, sig, req)
}

// GetOrderByID gets an order by its ID. This endpoint requires either the "view" or "trade" permission.
// Orders may be queried using either the exchange assigned id or the client assigned client_oid.
// If the order is canceled the response may have status code 404 if the order had no matches. Note:
// Open orders may change state between the request and the response depending on market conditions.
func (c *Client) GetOrderByID(orderID string) (Order, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/orders/" + orderID

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, nil)
	if err != nil {
		return Order{}, err
	}

	return c.getOrder(timestamp, sig, req)
}

// GetOrderByClientOID gets an order by its client OID. This endpoint requires either the "view" or "trade"
// permission. Orders may be queried using either the exchange assigned id or the client assigned client_oid.
// If the order is canceled the response may have status code 404 if the order had no matches. Note: Open orders
// may change state between the request and the response depending on market conditions.
func (c *Client) GetOrderByClientOID(orderClientOID string) (Order, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/orders/client:" + orderClientOID

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, nil)
	if err != nil {
		return Order{}, err
	}

	return c.getOrder(timestamp, sig, req)
}

// ListFills gets a list of recent fills of the API key's profile. This endpoint requires either the "view"
// or "trade" permission. You can request fills for specific orders or products using the orderID and productID
// parameters. You are required to provide either a product_id or order_id.
//
// Fees are recorded in two stages. Immediately after the matching engine completes a match, the fill
// is inserted into our datastore. Once the fill is recorded, a settlement process will settle the fill and credit
// both trading counterparties. The fee field indicates the fees charged for this individual fill.
func (c *Client) ListFills(qp QueryParams) ([]Fill, error) {
	if qp[ProductID] == "" && qp[OrderID] == "" {
		return nil, ErrMissingOrderOrProductID
	}
	timestamp := unixTime()
	method := http.MethodGet
	path := "/fills"

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, &qp)
	if err != nil {
		return nil, err
	}

	return c.listFills(timestamp, sig, req)
}

// GetCurrentExchangeLimits will return information on your payment method transfer limits, as well as buy/sell limits per currency.
func (c *Client) GetCurrentExchangeLimits() (ExchangeLimit, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/users/self/exchange-limits"

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, nil)
	if err != nil {
		return ExchangeLimit{}, err
	}

	return c.getLimits(timestamp, sig, req)
}

// StableCoinConversion creates a stablecoin conversion. One example is converting $10,000.00 USD to 10,000.00 USDC.
// A successful conversion will be assigned a conversion id which comes back on the Conversion as ID. The corresponding
// ledger entries for a conversion will reference this conversion id. Params:
// from:	A valid currency id
// to:		A valid currency id
// amount:	Amount of from to convert to to
func (c *Client) StableCoinConversion(from string, to string, amount string) (Conversion, error) {
	if from == "" || to == "" || amount == "" {
		return Conversion{}, ErrMissingConversionParams
	}
	conv := conversionReq{From: from, To: to, Amount: amount}
	timestamp := unixTime()
	method := http.MethodPost
	path := "/conversions"

	body, err := json.Marshal(conv)
	if err != nil {
		return Conversion{}, err
	}

	req, sig, err := c.createAndSignRequest(timestamp, method, path, body, nil)
	if err != nil {
		return Conversion{}, err
	}

	return c.stableCoinConversion(timestamp, sig, req)
}

// ListPaymentMethods gets a list of your payment methods.
func (c *Client) ListPaymentMethods() ([]PaymentMethod, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/payment-methods"

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, nil)
	if err != nil {
		return nil, err
	}

	return c.listPaymentMethods(timestamp, sig, req)
}

// ListCoinbaseAccounts lists your user's coinbase (non-pro) accounts.
// This endpoint requires either the "view" or "transfer" permission.
func (c *Client) ListCoinbaseAccounts() ([]CoinbaseAccount, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/coinbase-accounts"

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, nil)
	if err != nil {
		return nil, err
	}

	return c.listCoinbaseAccounts(timestamp, sig, req)
}

// GetCurrentFees returns your current maker & taker fee rates, as well as your 30-day trailing volume.
// Quoted rates are subject to change.
func (c *Client) GetCurrentFees() (Fees, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/fees"

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, nil)
	if err != nil {
		return Fees{}, err
	}

	return c.getCurrentFees(timestamp, sig, req)
}

// GetTrailingVolume returns your 30-day trailing volume for all products of the API key's profile.
// This is a cached value that's calculated every day at midnight UTC.
// This endpoint requires either the "view" or "trade" permission.
func (c *Client) GetTrailingVolume() ([]UserAccount, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/users/self/trailing-volume"

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, nil)
	if err != nil {
		return nil, err
	}

	return c.getTrailingVolume(timestamp, sig, req)
}

// ListProfiles lists the api key user's profiles which are equivilant to portfolios.
// This endpoint requires the "view" permission and is accessible by any profile's API key.
func (c *Client) ListProfiles() ([]Profile, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/profiles"

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, nil)
	if err != nil {
		return nil, err
	}

	return c.listProfiles(timestamp, sig, req)
}

// GetProfile gets a single profile by profile id. This endpoint requires the "view" permission
// and is accessible by any profile's API key.
func (c *Client) GetProfile(profileID string) (Profile, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/profiles/" + profileID

	req, sig, err := c.createAndSignRequest(timestamp, method, path, noBody, nil)
	if err != nil {
		return Profile{}, err
	}

	return c.getProfile(timestamp, sig, req)
}

// ProfileTransfer transfers funds from API key's profile to another user owned profile.
// This endpoint requires the "transfer" permission.
func (c *Client) ProfileTransfer(transfer TransferParams) error {
	timestamp := unixTime()
	method := http.MethodPost
	path := "/profiles/transfer"

	body, err := json.Marshal(transfer)
	if err != nil {
		return err
	}

	req, sig, err := c.createAndSignRequest(timestamp, method, path, body, nil)
	if err != nil {
		return err
	}

	return c.profileTransfer(timestamp, sig, req)
}
