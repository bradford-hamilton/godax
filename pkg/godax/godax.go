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

// ErrMissingOrderOrProductID TODO: does this feel weird and one offy right now?
var ErrMissingOrderOrProductID = errors.New("please provide either an orderID or productID")

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
func (c *Client) ListOrders(status, productID *string) ([]Order, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/orders"

	// TODO: come back to this with a better solution for query params.
	// TODO: To specify multiple statuses, use the status query argument multiple times: /orders?status=done&status=pending
	if status != nil {
		qp := "?status=" + *status
		if productID != nil {
			path += qp + "&product_id=" + *productID
		} else {
			path += qp
		}
	} else if productID != nil {
		path += "?status=all&product_id=" + *productID
	}

	sig, err := c.generateSig(timestamp, method, path, "")
	if err != nil {
		return nil, err
	}

	return c.listOrders(timestamp, method, path, sig)
}

// GetOrderByID gets an order by its ID. This endpoint requires either the "view" or "trade" permission.
// Orders may be queried using either the exchange assigned id or the client assigned client_oid.
// If the order is canceled the response may have status code 404 if the order had no matches. Note:
// Open orders may change state between the request and the response depending on market conditions.
func (c *Client) GetOrderByID(orderID string) (Order, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/orders/" + orderID

	sig, err := c.generateSig(timestamp, method, path, "")
	if err != nil {
		return Order{}, err
	}

	return c.getOrder(timestamp, method, path, sig)
}

// GetOrderByClientOID gets an order by its client OID. This endpoint requires either the "view" or "trade"
// permission. Orders may be queried using either the exchange assigned id or the client assigned client_oid.
// If the order is canceled the response may have status code 404 if the order had no matches. Note: Open orders
// may change state between the request and the response depending on market conditions.
func (c *Client) GetOrderByClientOID(orderClientOID string) (Order, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/orders/client:" + orderClientOID

	sig, err := c.generateSig(timestamp, method, path, "")
	if err != nil {
		return Order{}, err
	}

	return c.getOrder(timestamp, method, path, sig)
}

// ListFills gets a list of recent fills of the API key's profile. This endpoint requires either the "view"
// or "trade" permission. You can request fills for specific orders or products using the orderID and productID
// parameters. You are required to provide either a product_id or order_id.
//
// Fees are recorded in two stages. Immediately after the matching engine completes a match, the fill
// is inserted into our datastore. Once the fill is recorded, a settlement process will settle the fill and credit
// both trading counterparties. The fee field indicates the fees charged for this individual fill.
func (c *Client) ListFills(orderID, productID *string) ([]Fill, error) {
	timestamp := unixTime()
	method := http.MethodGet
	path := "/fills"

	// TODO: another query param area to clean up. Thinking a struct could
	// maybe be a better call? Look into this some.
	if orderID == nil && productID == nil {
		return nil, ErrMissingOrderOrProductID
	}
	if orderID != nil {
		path += "?order_id=" + *orderID
		if productID != nil {
			path += "&order_id=" + *productID
		}
	} else {
		path += "?product_id=" + *productID
	}

	sig, err := c.generateSig(timestamp, method, path, "")
	if err != nil {
		return []Fill{}, err
	}

	return c.listFills(timestamp, method, path, sig)
}
