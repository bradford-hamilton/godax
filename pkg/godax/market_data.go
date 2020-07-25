package godax

import (
	"encoding/json"
	"net/http"
)

// Product represents a coinbase pro product, for example "BTC-USD". Only a maximum of
// one of trading_disabled, cancel_only, post_only, limit_only can be true at once. If
// none are true, the product is trading normally. Product ID will not change once assigned
// to a product but all other fields ares subject to change.
/*
{
    "id": "BTC-USD",
    "display_name": "BTC/USD",
    "base_currency": "BTC",
    "quote_currency": "USD",
    "base_increment": "0.00000001",
    "quote_increment": "0.01000000",
    "base_min_size": "0.00100000",
    "base_max_size": "280.00000000",
    "min_market_funds": "5",
    "max_market_funds": "1000000",
    "status": "online",
    "status_message": "",
    "cancel_only": false,
    "limit_only": false,
    "post_only": false,
    "trading_disabled": false
}
*/
type Product struct {
	// ID represents the product ID, ie. BTC-USD
	ID string `json:"id"`

	// DisplayName represents human friendly name, example: "BTC/USD"
	DisplayName string `json:"display_name"`

	// BaseCurrency is as titled, example: "BTC"
	BaseCurrency string `json:"base_currency"`

	// BaseCurrency is as titled, example: "USD"
	QuoteCurrency string `json:"quote_currency"`

	// BaseIncrement specifies the minimum increment for the base_currency
	BaseIncrement string `json:"base_increment"`

	// QuoteIncrement specifies the min order price as well as the price increment.
	// The order price must be a multiple of this increment (i.e. if the increment is
	// 0.01, order prices of 0.001 or 0.021 would be rejected).
	QuoteIncrement string `json:"quote_increment"`

	// BaseMinSize describes the minimum order size
	BaseMinSize string `json:"base_min_size"`

	// BaseMaxSize describes the maximum order size
	BaseMaxSize string `json:"base_max_size"`

	// MinMarketFunds describes the minimum funds allowed in a market order
	MinMarketFunds string `json:"min_market_funds"`

	// MaxMarketFunds describes the maximum funds allowed in a market order
	MaxMarketFunds string `json:"max_market_funds"`

	// Status is the product's current status, example: "online"
	Status string `json:"status"`

	// StatusMessage provides any extra information regarding the status if available
	StatusMessage string `json:"status_message"`

	// CancelOnly indicates indicates whether this product only accepts cancel requests for orders
	CancelOnly bool `json:"cancel_only"`

	// LimitOnly indicates whether this product only accepts limit orders. When LimitOnly is true, matching can occur if a limit order crosses the book.
	LimitOnly bool `json:"limit_only"`

	// PostOnly indicates whether only maker orders can be placed. No orders will be matched when post_only mode is active
	PostOnly bool `json:"post_only"`

	// TradingDisabled indicates whether trading is currently restricted on this product.
	// This includes whether both new orders and order cancelations are restricted
	TradingDisabled bool `json:"trading_disabled"`
}

// OrderBook represents a list of orders for a product. TODO: maybe notify coinbase. The docs say sequence is a string as
// it appears below, however it is an int.
/*
{
    "sequence": "3",
    "bids": [
        [ price, size, num-orders ],
        [ "295.96", "4.39088265", 2 ],
        ...
    ],
    "asks": [
        [ price, size, num-orders ],
        [ "295.97", "25.23542881", 12 ],
        ...
    ]
}
*/
type OrderBook struct {
	Sequence int              `json:"sequence"`
	Bids     []OrderBookOrder `json:"bids"`
	Asks     []OrderBookOrder `json:"asks"`
}

// OrderBookOrder represents the price, size, and number of orders for a product.
type OrderBookOrder struct {
	Price     string `json:"price"`
	Size      string `json:"size"`
	NumOrders int    `json:"num_orders"`
}

// Trade represents a trade that has happened for a product
type Trade struct {
	Time    string `json:"time"`
	TradeID int    `json:"trade_id"`
	Price   string `json:"price"`
	Size    string `json:"size"`
	Side    string `json:"side"`
}

// TODO: it seems when you ask for level 3, the shape of the bids and asks no longer apply :(
// The NumOrders field comes back as a string UUID, and so I am not sure what that is. May reach out
// to coinbase on this one as well.

// UnmarshalJSON is a custom unmarshaller for an OrderBook. Unfortunately the coinbase pro
// API returns different types in the bids & asks JSON arrays, so we handle that here.
// This approach should provide us with all the standard JSON errors if something goes wrong.
func (o *OrderBookOrder) UnmarshalJSON(b []byte) error {
	var msg []json.RawMessage
	if err := json.Unmarshal(b, &msg); err != nil {
		return err
	}
	if err := json.Unmarshal(msg[0], &o.Price); err != nil {
		return err
	}
	if err := json.Unmarshal(msg[1], &o.Size); err != nil {
		return err
	}
	if err := json.Unmarshal(msg[2], &o.NumOrders); err != nil {
		return err
	}
	return nil
}

func (c *Client) listProducts(timestamp, signature string, req *http.Request) ([]Product, error) {
	res, err := c.do(timestamp, signature, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var p []Product
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		return nil, err
	}
	return p, nil
}

func (c *Client) getProductByID(timestamp, signature string, req *http.Request) (Product, error) {
	res, err := c.do(timestamp, signature, req)
	if err != nil {
		return Product{}, err
	}
	defer res.Body.Close()

	var p Product
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		return Product{}, err
	}
	return p, nil
}

func (c *Client) getProductOrderBook(timestamp, signature string, req *http.Request) (OrderBook, error) {
	res, err := c.do(timestamp, signature, req)
	if err != nil {
		return OrderBook{}, err
	}
	defer res.Body.Close()

	var ob OrderBook
	if err := json.NewDecoder(res.Body).Decode(&ob); err != nil {
		return OrderBook{}, err
	}
	return ob, nil
}

func (c *Client) listTradesByProduct(timestamp, signature string, req *http.Request) ([]Trade, error) {
	res, err := c.do(timestamp, signature, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var trades []Trade
	if err := json.NewDecoder(res.Body).Decode(&trades); err != nil {
		return nil, err
	}
	return trades, nil
}
