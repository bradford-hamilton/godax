package godax

// Fill represents a filled order on coinbase.
type Fill struct {
	// TradeID is the identifier (int) for the trade that created the fill.
	TradeID int `json:"trade_id"`

	// ProductID - the product ID related to the activity.
	ProductID string `json:"product_id"`

	// Price is specified in quote_increment product units.
	Price string `json:"price"`

	// Size indicates the amount of BTC (or base currency) to buy or sell.
	Size string `json:"size"`

	// OrderID is the identifier (int) for the order that created the fill.
	OrderID string `json:"order_id"`

	// CreatedAt - when this fill was created.
	CreatedAt string `json:"created_at"`

	// Liquidity field indicates if the fill was the result of a liquidity provider or liquidity taker.
	// M indicates Maker and T indicates Taker
	Liquidity string `json:"liquidity"`

	// Fee indicates fees incurred as a result of the fill.
	Fee string `json:"fee"`

	// Settled indicates whether the fill has was settled yet.
	Settled bool `json:"settled"`

	// Side is either buy or sell
	Side string `json:"side"`
}
