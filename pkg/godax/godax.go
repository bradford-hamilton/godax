// Package godax implements an http client for interacting with the coinbase pro API.
// Rest API
// PUBLIC ENDPOINTS
// 		We throttle public endpoints by IP: 3 requests per second, up to 6 requests per second in bursts.
// 		Some endpoints may have custom rate limits.
// PRIVATE ENDPOINTS
// 		We throttle private endpoints by profile ID: 5 requests per second, up to 10 requests per second
// 		in bursts. Some endpoints may have custom rate limits.
package godax

import (
	"encoding/json"
	"errors"
	"io/ioutil"
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

// Available query params
const (
	OrderID   Param = "order_id"
	ProductID Param = "product_id"
	Status    Param = "status"

	// Level is used when calling GetProductOrderBook
	Level Param = "level"

	// These params are used when calling GetHistoricRatesForProduct
	Start       Param = "start"
	End         Param = "end"
	Granularity Param = "granularity"
)

var noBody = []byte{}

// QueryParams represent the available query params for any given coinbase pro call.
type QueryParams map[Param]string

// ErrMissingOrderOrProductID TODO: does this feel weird and one offy right now?
var (
	ErrMissingOrderOrProductID = errors.New("please provide either an order_id or product_id in your query params")
	ErrMissingConversionParams = errors.New("please provide all of the following params: to, from, amount")
	ErrCoinbaseProAPIChange    = errors.New("there appears to have been a coinbase pro API change. Please open a new issue on godax, thanks")
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
	method := http.MethodGet
	path := "/accounts"

	var accounts []ListAccount
	if err := c.exec(unixTime(), method, path, noBody, nil, &accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}

// GetAccount retrieves information for a single account. Use this endpoint when you know the
// account_id. API key must belong to the same profile as the account. This endpoint requires
// either the "view" or "trade" permission.
func (c *Client) GetAccount(accountID string) (Account, error) {
	method := http.MethodGet
	path := "/accounts/" + accountID

	var account Account
	if err := c.exec(unixTime(), method, path, noBody, nil, &account); err != nil {
		return Account{}, err
	}

	return account, nil
}

// GetAccountHistory lists account activity of the API key's profile. Account activity either increases
// or decreases your account balance. If an entry is the result of a trade (match, fee), the details
// field on an AccountActivity will contain additional information about the trade. Items are paginated
// and sorted latest first. This endpoint requires either the "view" or "trade" permission.
// TODO: paginate
func (c *Client) GetAccountHistory(accountID string) ([]AccountActivity, error) {
	method := http.MethodGet
	path := "/accounts/" + accountID + "/ledger"

	var activities []AccountActivity
	if err := c.exec(unixTime(), method, path, noBody, nil, &activities); err != nil {
		return nil, err
	}

	return activities, nil
}

// GetAccountHolds lists holds of an account that belong to the same profile as the API key.
// Holds are placed on an account for any active orders or pending withdraw requests. As an
// order is filled, the hold amount is updated. If an order is canceled, any remaining hold
// is removed. For a withdraw, once it is completed, the hold is removed. This endpoint
// requires either the "view" or "trade" permission.
// TODO: paginate
func (c *Client) GetAccountHolds(accountID string) ([]AccountHold, error) {
	method := http.MethodGet
	path := "/accounts/" + accountID + "/holds"

	var holds []AccountHold
	if err := c.exec(unixTime(), method, path, noBody, nil, &holds); err != nil {
		return nil, err
	}

	return holds, nil
}

// PlaceOrder allows you to place two types of orders: limit and market. Orders can only be
// placed if your account has sufficient funds. Once an order is placed, your account funds
// will be put on hold for the duration of the order. How much and which funds are put on
// hold depends on the order type and parameters specified. This endpoint requires the
// "trade" permission.
func (c *Client) PlaceOrder(order OrderParams) (Order, error) {
	method := http.MethodPost
	path := "/orders"

	body, err := json.Marshal(order)
	if err != nil {
		return Order{}, err
	}

	var o Order
	if err := c.exec(unixTime(), method, path, body, nil, &o); err != nil {
		return Order{}, err
	}

	return o, nil
}

// TODO: do the below two methods need to rerturn anything but an error? Probably not.

// CancelOrderByID cancels a previously placed order. Order must belong to the profile that
// the API key belongs to. If the order had no matches during its lifetime its record may be
// purged. This means the order details will not be available with GetOrderByID or GetOrderByClientOID.
// The product ID of the order is not required so if you don't have it you can pass nil here.
// The request will be more performant if you include it. This endpoint requires the "trade" permission.
func (c *Client) CancelOrderByID(orderID string, qp QueryParams) (canceledOrderID string, err error) {
	method := http.MethodDelete
	path := "/orders/" + orderID

	req, sig, err := c.createAndSignReq(unixTime(), method, path, noBody, &qp)
	if err != nil {
		return "", err
	}

	res, err := c.do(unixTime(), sig, req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	return orderID, nil
}

// CancelOrderByClientOID cancels a previously placed order. Order must belong to the profile that
// the API key belongs to. If the order had no matches during its lifetime its record may be
// purged. This means the order details will not be available with GetOrderByID or GetOrderByClientOID.
// The product ID of the order is not required so if you don't have it you can pass nil here.
// The request will be more performant if you include it. This endpoint requires the "trade" permission.
func (c *Client) CancelOrderByClientOID(clientOID string, qp QueryParams) (canceledOrderID string, err error) {
	method := http.MethodDelete
	path := "/orders/client:" + clientOID

	req, sig, err := c.createAndSignReq(unixTime(), method, path, noBody, &qp)
	if err != nil {
		return "", err
	}

	res, err := c.do(unixTime(), sig, req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	id, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(id), nil
}

// CancelAllOrders cancel all open orders from the profile that the API key belongs to. The response is
// a list of ids of the canceled orders. This endpoint requires the "trade" permission. The productID
// param is opitonal and a pointer, so you can pass nil. If you do provide the productID here, you will
// only cancel orders open for that specific product.
func (c *Client) CancelAllOrders(qp QueryParams) (canceledOrderIDs []string, err error) {
	method := http.MethodDelete
	path := "/orders"

	if err := c.exec(unixTime(), method, path, noBody, &qp, &canceledOrderIDs); err != nil {
		return nil, err
	}

	return canceledOrderIDs, nil
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
	method := http.MethodGet
	path := "/orders"

	var orders []Order
	if err := c.exec(unixTime(), method, path, noBody, &qp, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrderByID gets an order by its ID. This endpoint requires either the "view" or "trade" permission.
// Orders may be queried using either the exchange assigned id or the client assigned client_oid.
// If the order is canceled the response may have status code 404 if the order had no matches. Note:
// Open orders may change state between the request and the response depending on market conditions.
func (c *Client) GetOrderByID(orderID string) (Order, error) {
	method := http.MethodGet
	path := "/orders/" + orderID

	var order Order
	if err := c.exec(unixTime(), method, path, noBody, nil, &order); err != nil {
		return Order{}, err
	}

	return order, nil
}

// GetOrderByClientOID gets an order by its client OID. This endpoint requires either the "view" or "trade"
// permission. Orders may be queried using either the exchange assigned id or the client assigned client_oid.
// If the order is canceled the response may have status code 404 if the order had no matches. Note: Open orders
// may change state between the request and the response depending on market conditions.
func (c *Client) GetOrderByClientOID(orderClientOID string) (Order, error) {
	method := http.MethodGet
	path := "/orders/client:" + orderClientOID

	var order Order
	if err := c.exec(unixTime(), method, path, noBody, nil, &order); err != nil {
		return Order{}, err
	}

	return order, nil
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
	method := http.MethodGet
	path := "/fills"

	var fills []Fill
	if err := c.exec(unixTime(), method, path, noBody, &qp, &fills); err != nil {
		return nil, err
	}

	return fills, nil
}

// GetCurrentExchangeLimits will return information on your payment method transfer limits, as well as buy/sell limits per currency.
func (c *Client) GetCurrentExchangeLimits() (ExchangeLimit, error) {
	method := http.MethodGet
	path := "/users/self/exchange-limits"

	var limit ExchangeLimit
	if err := c.exec(unixTime(), method, path, noBody, nil, &limit); err != nil {
		return ExchangeLimit{}, err
	}

	return limit, nil
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
	method := http.MethodPost
	path := "/conversions"

	body, err := json.Marshal(conversionReq{From: from, To: to, Amount: amount})
	if err != nil {
		return Conversion{}, err
	}

	var conv Conversion
	if err := c.exec(unixTime(), method, path, body, nil, &conv); err != nil {
		return Conversion{}, err
	}

	return conv, nil
}

// ListPaymentMethods gets a list of your payment methods.
func (c *Client) ListPaymentMethods() ([]PaymentMethod, error) {
	method := http.MethodGet
	path := "/payment-methods"

	var pm []PaymentMethod
	if err := c.exec(unixTime(), method, path, noBody, nil, &pm); err != nil {
		return nil, err
	}

	return pm, nil
}

// ListCoinbaseAccounts lists your user's coinbase (non-pro) accounts.
// This endpoint requires either the "view" or "transfer" permission.
func (c *Client) ListCoinbaseAccounts() ([]CoinbaseAccount, error) {
	method := http.MethodGet
	path := "/coinbase-accounts"

	var accts []CoinbaseAccount
	if err := c.exec(unixTime(), method, path, noBody, nil, &accts); err != nil {
		return nil, err
	}

	return accts, nil
}

// GetCurrentFees returns your current maker & taker fee rates, as well as your 30-day trailing volume.
// Quoted rates are subject to change.
func (c *Client) GetCurrentFees() (Fees, error) {
	method := http.MethodGet
	path := "/fees"

	var fees Fees
	if err := c.exec(unixTime(), method, path, noBody, nil, &fees); err != nil {
		return Fees{}, err
	}

	return fees, nil
}

// GetTrailingVolume returns your 30-day trailing volume for all products of the API key's profile.
// This is a cached value that's calculated every day at midnight UTC.
// This endpoint requires either the "view" or "trade" permission.
func (c *Client) GetTrailingVolume() ([]UserAccount, error) {
	method := http.MethodGet
	path := "/users/self/trailing-volume"

	var userActs []UserAccount
	if err := c.exec(unixTime(), method, path, noBody, nil, &userActs); err != nil {
		return nil, err
	}

	return userActs, nil
}

// ListProfiles lists the api key user's profiles which are equivilant to portfolios.
// This endpoint requires the "view" permission and is accessible by any profile's API key.
func (c *Client) ListProfiles() ([]Profile, error) {
	method := http.MethodGet
	path := "/profiles"

	var profiles []Profile
	if err := c.exec(unixTime(), method, path, noBody, nil, &profiles); err != nil {
		return nil, err
	}

	return profiles, nil
}

// GetProfile gets a single profile by profile id. This endpoint requires the "view" permission
// and is accessible by any profile's API key.
func (c *Client) GetProfile(profileID string) (Profile, error) {
	method := http.MethodGet
	path := "/profiles/" + profileID

	var prof Profile
	if err := c.exec(unixTime(), method, path, noBody, nil, &prof); err != nil {
		return Profile{}, err
	}

	return prof, nil
}

// ProfileTransfer transfers funds from API key's profile to another user owned profile.
// This endpoint requires the "transfer" permission.
func (c *Client) ProfileTransfer(transfer TransferParams) error {
	method := http.MethodPost
	path := "/profiles/transfer"

	body, err := json.Marshal(transfer)
	if err != nil {
		return err
	}
	req, sig, err := c.createAndSignReq(unixTime(), method, path, body, nil)
	if err != nil {
		return err
	}
	if _, err := c.do(unixTime(), sig, req); err != nil {
		return err
	}

	return nil
}

// ListProducts gets a list of available currency pairs for trading. The Market Data API is an
// unauthenticated set of endpoints for retrieving market data. These endpoints provide snapshots
// of market data.
func (c *Client) ListProducts() ([]Product, error) {
	method := http.MethodGet
	path := "/products"

	var products []Product
	if err := c.exec(unixTime(), method, path, noBody, nil, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// GetProductByID gets market data for a specific currency pair.
func (c *Client) GetProductByID(productID string) (Product, error) {
	method := http.MethodGet
	path := "/products/" + productID

	var product Product
	if err := c.exec(unixTime(), method, path, noBody, nil, &product); err != nil {
		return Product{}, err
	}

	return product, nil
}

// GetProductOrderBook gets a list of open orders for a product. The amount of detail shown can
// be customized with the level parameter. By default, only the inside (i.e. best) bid and ask
// are returned. This is equivalent to a book depth of 1 level. If you would like to see a
// larger order book, specify the level query parameter. If a level is not aggregated, then
// all of the orders at each price will be returned. Aggregated levels return only one size
// for each active price (as if there was only a single order for that size at the level). You
// can provide the "level" query param here which defaults to "1" if you do not specify.
// Level 	Description
// 1	 	Only the best bid and ask
// 2	 	Top 50 bids and asks (aggregated)
// 3	 	Full order book (non aggregated)
// This request is NOT paginated. The entire book is returned in one response. Level 1 and
// Level 2 are recommended for polling. For the most up-to-date data, consider using the
// websocket stream. Level 3 is only recommended for users wishing to maintain a full real-time
// order book using the websocket stream. Abuse of Level 3 via polling will cause your access
// to be limited or blocked.
func (c *Client) GetProductOrderBook(productID string, qp QueryParams) (OrderBook, error) {
	method := http.MethodGet
	path := "/products/" + productID + "/book"

	var ob OrderBook
	if err := c.exec(unixTime(), method, path, noBody, &qp, &ob); err != nil {
		return OrderBook{}, err
	}

	return ob, nil
}

// GetProductTicker returns snapshot information about the last trade (tick), best bid/ask and 24h volume.
// Polling is discouraged in favor of connecting via the websocket stream and listening for match messages.
func (c *Client) GetProductTicker(productID string) (Ticker, error) {
	method := http.MethodGet
	path := "/products/" + productID + "/ticker"

	var ticker Ticker
	if err := c.exec(unixTime(), method, path, noBody, nil, &ticker); err != nil {
		return Ticker{}, err
	}

	return ticker, nil
}

// ListTradesByProduct lists the latest trades for a product. The trade side indicates the
// maker order side. The maker order is the order that was open on the order book. Buy side
// indicates a down-tick because the maker was a buy order and their order was removed.
// Conversely, sell side indicates an up-tick.
func (c *Client) ListTradesByProduct(productID string) ([]Trade, error) {
	method := http.MethodGet
	path := "/products/" + productID + "/trades"

	var trades []Trade
	if err := c.exec(unixTime(), method, path, noBody, nil, &trades); err != nil {
		return nil, err
	}

	return trades, nil
}

// GetHistoricRatesForProduct gets historic rates for a product. Rates are returned in grouped
// buckets based on requested granularity. Historical rate data may be incomplete. No data is
// published for intervals where there are no ticks.
// PARAMETERS
// Param		Description
// start		Start time in ISO 8601
// end			End time in ISO 8601
// granularity	Desired timeslice in seconds
// If either one of the start or end fields are not provided then both fields will be ignored. If a custom
// time range is not declared then one ending now is selected. The granularity field must be one of the
// following values: {60, 300, 900, 3600, 21600, 86400}. Otherwise, your request will be rejected. These
// values correspond to timeslices representing one minute, five minutes, fifteen minutes, one hour, six
// hours, and one day, respectively. If data points are readily available, your response may contain as many
// as 300 candles and some of those candles may precede your declared start value. The maximum number of data
// points for a single request is 300 candles. If your selection of start/end time and granularity will result
// in more than 300 data points, your request will be rejected. If you wish to retrieve fine granularity data
// over a larger time range, you will need to make multiple requests with new start/end ranges.
func (c *Client) GetHistoricRatesForProduct(productID string, qp QueryParams) ([]HistoricRate, error) {
	method := http.MethodGet
	path := "/products/" + productID + "/candles"

	var rates []HistoricRate
	if err := c.exec(unixTime(), method, path, noBody, &qp, &rates); err != nil {
		return nil, err
	}

	return rates, nil
}

// Get24HourStatsForProduct gets 24 hr stats for the product. Volume is in base currency units. Open, high,
// low are in quote currency units.
func (c *Client) Get24HourStatsForProduct(productID string) (DayStat, error) {
	method := http.MethodGet
	path := "/products/" + productID + "/stats"

	var dayStat DayStat
	if err := c.exec(unixTime(), method, path, noBody, nil, &dayStat); err != nil {
		return DayStat{}, err
	}

	return dayStat, nil
}

// ListCurrencies lists known currencies. Currency codes will conform to the ISO 4217 standard where possible.
// Currencies which have or had no representation in ISO 4217 may use a custom code.
// Code		Description
// BTC		Bitcoin
// ETH		Ether
// LTC		Litecoin
func (c *Client) ListCurrencies() ([]Currency, error) {
	method := http.MethodGet
	path := "/currencies"

	var currencies []Currency
	if err := c.exec(unixTime(), method, path, noBody, nil, &currencies); err != nil {
		return nil, err
	}

	return currencies, nil
}

// GetServerTime fetches the current coinbase pro server time. This endpoint does not require authentication.
func (c *Client) GetServerTime() (ServerTime, error) {
	method := http.MethodGet
	path := "/time"

	var srvTime ServerTime
	if err := c.exec(unixTime(), method, path, noBody, nil, &srvTime); err != nil {
		return ServerTime{}, err
	}

	return srvTime, nil
}

// CreateReport creates reports that provide batches of historic information about your
// profile in various human and machine readable forms. This endpoint requires either
// the "view" or "trade" permission. The report will be generated when resources are
// available. Report status can be queried via the /reports/:report_id endpoint. The
// file_url field will be available once the report has successfully been created and
// is available for download. Expired reports: reports are only available for download
// for a few days after being created. Once a report expires, the report is no longer
// available for download and is deleted.
func (c *Client) CreateReport(report ReportParams) (ReportStatus, error) {
	method := http.MethodPost
	path := "/reports"

	body, err := json.Marshal(report)
	if err != nil {
		return ReportStatus{}, err
	}

	var r ReportStatus
	if err := c.exec(unixTime(), method, path, body, nil, &r); err != nil {
		return ReportStatus{}, err
	}

	return r, nil
}

// GetReportStatus Once a report request has been accepted for processing, the status
// is available by polling the report resource endpoint. The final report will be
// uploaded and available at file_url once the status indicates ready.
// STATUS		DESCRIPTION
// pending		The report request has been accepted and is awaiting processing
// creating		The report is being created
// ready		The report is ready for download from file_url
// This endpoint requires either the "view" or "trade" permission.
func (c *Client) GetReportStatus(reportID string) (ReportStatus, error) {
	method := http.MethodGet
	path := "/reports/" + reportID

	var r ReportStatus
	if err := c.exec(unixTime(), method, path, noBody, nil, &r); err != nil {
		return ReportStatus{}, err
	}
	return r, nil
}
