package godax

import (
	"fmt"
	"net/http"

	"gopkg.in/square/go-jose.v2/json"
)

// Order represents a trading account for a coinbase pro profile.
/*
	{
        "id": "d0c5340b-6d6c-49d9-b567-48c4bfca13d2",
        "price": "0.10000000",
        "size": "0.01000000",
        "product_id": "BTC-USD",
        "side": "buy",
        "stp": "dc",
        "type": "limit",
        "time_in_force": "GTC",
        "post_only": false,
        "created_at": "2016-12-08T20:02:28.53864Z",
        "fill_fees": "0.0000000000000000",
        "filled_size": "0.00000000",
        "executed_value": "0.0000000000000000",
        "status": "pending",
        "settled": false
    }
*/
type Order struct {
	// ID is the order UUID generated by coinbase.
	ID string `json:"id"`

	// CreatedAt indicates the time the order was filled.
	CreatedAt string `json:"created_at"`

	// FillFees indicates how much in fees the order used.
	FillFees string `json:"fill_fees"`

	// FilledSize represents the order's filled amount in BTC (or base currency)
	FilledSize string `json:"filled_size"`

	// ExecutedValue indicates the value of the executed order
	ExecutedValue string `json:"executed_value"`

	// Status tells you the status of the order ("pending", "filled", etc).
	Status string `json:"status"`

	// Settled indicates whether the order has been settled yet.
	Settled bool `json:"settled"`

	// OrderParams represent all the fields available to send in a PlaceOrder call.
	// They are also attached to orders returned from coinbase after creating, listing, etc.
	OrderParams
}

// OrderParams are a combination of common/shared params as well as limit and market specific params.
type OrderParams struct {
	CommonOrderParams
	LimitOrderParams
	MarketOrderParams
}

// CommonOrderParams represent the params that both limit and market orders share
type CommonOrderParams struct {
	// Side is either buy or sell
	Side string `json:"side"`

	// ProductID must be a valid product id. The list of products is available via the GetProducts method.
	ProductID string `json:"product_id"`

	// ClientOID - [optional] order UUID generated by you to identify your order
	ClientOID string `json:"client_oid,omitempty"`

	// Type - [optional] limit or market. If type is not specified, the order will default to a limit order.
	Type string `json:"type,omitempty"`

	// Price - The price must be specified in quote_increment product units. The quote increment is the
	// smallest unit of price. For example, the BTC-USD product has a quote increment of 0.01 or 1 penny.
	// Prices less than 1 penny will not be accepted, and no fractional penny prices will be accepted.
	// Not required for market orders.
	Price string `json:"price"`

	// indicates the amount of BTC (or base currency) to buy or sell. If you are placing a market order
	// this is optional, although you will need either Size OR Funds for a market order. This is required
	// for a limit order. The size must be greater than the base_min_size for the product and no larger
	// than the base_max_size. The size can be in incremented in units of base_increment.
	Size string `json:"size,omitempty"`

	// Stp - [optional] self-trade prevention flag
	Stp string `json:"stp,omitempty"`

	// Stop - [optional] either loss or entry. Requires stop_price to be defined.
	Stop string `json:"stop,omitempty"`

	// StopPrice - [optional] Only if stop is defined. Sets trigger price for stop order.
	StopPrice string `json:"stop_price,omitempty"`
}

// LimitOrderParams are embedded into OrderParams, and represent params that are only part of limit orders.
type LimitOrderParams struct {
	// TimeInForce - [optional] GTC, GTT, IOC, or FOK (default is GTC). Time in force policies provide
	// guarantees about the lifetime of an order. There are four policies: good till canceled GTC, good
	// till time GTT, immediate or cancel IOC, and fill or kill FOK.
	TimeInForce string `json:"time_in_force,omitempty"`

	// CancelAfter - [optional] - min, hour, day. Requires time_in_force to be GTT
	CancelAfter string `json:"cancel_after,omitempty"`

	// PostOnly - [optional] post only flag. Invalid when time_in_force is IOC or FOK.
	PostOnly bool `json:"post_only,omitempty"`
}

// MarketOrderParams are embedded into OrderParams, and represent params that are only part of market orders.
// NOTE: either Size (common/shared) OR Funds in this struct must be present for a market order.
type MarketOrderParams struct {
	// Funds - [optional] desired amount of quote currency to use. When specified it indicates how much of
	// the product quote currency to buy or sell. For example, a market buy for BTC-USD with funds specified
	// as 150.00 will spend 150 USD to buy BTC (including any fees). If the funds field is not specified for
	// a market buy order, size must be specified and Coinbase Pro will use available funds in your account
	// to buy bitcoin.
	Funds string `json:"funds,omitempty"`
}

// placeOrder ...
func (c *Client) placeOrder(timestamp, method, path, signature string, body []byte) (Order, error) {
	res, err := c.do(timestamp, method, path, signature, body)
	if err != nil {
		return Order{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return orderError(res)
	}

	var order Order
	if err := json.NewDecoder(res.Body).Decode(&order); err != nil {
		return Order{}, err
	}

	return order, nil
}

func orderError(res *http.Response) (Order, error) {
	var err CoinbaseErrRes
	json.NewDecoder(res.Body).Decode(&err)
	return Order{}, fmt.Errorf("status code: %d, message: %s", res.StatusCode, err.Message)
}