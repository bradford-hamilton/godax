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

func (c *Client) getProduct(timestamp, signature string, req *http.Request) (Product, error) {
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
