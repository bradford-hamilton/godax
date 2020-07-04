package godax

import (
	"encoding/json"
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

// NewClient returns a godax Client that is hooked up to the live REST and web socket APIs.
func NewClient() (*Client, error) {
	return newClient(false)
}

// NewSandboxClient returns a godax Client that is hooked up to the sandbox REST and web socket APIs.
func NewSandboxClient() (*Client, error) {
	return newClient(true)
}

// ListAccounts gets a list of trading accounts from the profile associated with the API key.
// This endpoint requires either the "view" or "trade" permission. This endpoint has a custom
// rate limit by profile ID: 25 requests per second, up to 50 requests per second in bursts
func (c *Client) ListAccounts() ([]ListAccount, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/accounts"

	sig, err := c.generateSig(timestamp, method, path, "")
	if err != nil {
		return []ListAccount{}, err
	}

	return c.listAccounts(timestamp, method, path, sig)
}

// GetAccount retrieves information for a single account. Use this endpoint when you know the
// account_id. API key must belong to the same profile as the account. This endpoint requires
// either the "view" or "trade" permission.
func (c *Client) GetAccount(accountID string) (Account, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/accounts/" + accountID

	sig, err := c.generateSig(timestamp, method, path, "")
	if err != nil {
		return Account{}, err
	}

	return c.getAccount(accountID, timestamp, method, path, sig)
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

	sig, err := c.generateSig(timestamp, method, path, "")
	if err != nil {
		return []AccountActivity{}, err
	}

	return c.getAccountHistory(accountID, timestamp, method, path, sig)
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

	sig, err := c.generateSig(timestamp, method, path, "")
	if err != nil {
		return []AccountHold{}, err
	}

	return c.getAccountHolds(accountID, timestamp, method, path, sig)
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

	sig, err := c.generateSig(timestamp, method, path, string(body))
	if err != nil {
		return Order{}, err
	}

	return c.placeOrder(timestamp, method, path, sig, body)
}

func unixTime() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// CancelOrderByID cancels a previously placed order. Order must belong to the profile that
// the API key belongs to. If the order had no matches during its lifetime its record may be
// purged. This means the order details will not be available with GetOrderByID or GetOrderByClientOID.
// The product ID of the order is not required so if you don't have it you can pass nil here.
// The request will be more performant if you include it. This endpoint requires the "trade" permission.
func (c *Client) CancelOrderByID(orderID string, productID *string) (canceledOrderID string, err error) {
	timestamp := unixTime()
	method := http.MethodDelete
	path := "/orders/" + orderID
	if productID != nil {
		path += "?product_id=" + *productID
	}

	sig, err := c.generateSig(timestamp, method, path, "")
	if err != nil {
		return "", err
	}

	return c.cancelOrder(timestamp, method, path, sig)
}

// CancelOrderByClientOID cancels a previously placed order. Order must belong to the profile that
// the API key belongs to. If the order had no matches during its lifetime its record may be
// purged. This means the order details will not be available with GetOrderByID or GetOrderByClientOID.
// The product ID of the order is not required so if you don't have it you can pass nil here.
// The request will be more performant if you include it. This endpoint requires the "trade" permission.
func (c *Client) CancelOrderByClientOID(clientOID string, productID *string) (canceledOrderID string, err error) {
	timestamp := unixTime()
	method := http.MethodDelete
	path := "/orders/client:" + clientOID
	if productID != nil {
		path += "?product_id=" + *productID
	}

	sig, err := c.generateSig(timestamp, method, path, "")
	if err != nil {
		return "", err
	}

	return c.cancelOrder(timestamp, method, path, sig)
}

// CancelAllOrders cancel all open orders from the profile that the API key belongs to. The response is
// a list of ids of the canceled orders. This endpoint requires the "trade" permission. The productID
// param is opitonal and a pointer, so you can pass nil. If you do provide the productID here, you will
// only cancel orders open for that specific product.
func (c *Client) CancelAllOrders(productID *string) (canceledOrderIDs []string, err error) {
	timestamp := unixTime()
	method := http.MethodDelete
	path := "/orders"
	if productID != nil {
		path += "?product_id=" + *productID
	}

	sig, err := c.generateSig(timestamp, method, path, "")
	if err != nil {
		return nil, err
	}

	return c.cancelAllOrders(timestamp, method, path, sig)
}
