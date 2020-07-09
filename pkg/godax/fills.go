package godax

import (
	"encoding/json"
	"net/http"
)

// Fill represents a filled order on coinbase
/*
{
    "trade_id": 74,
    "product_id": "BTC-USD",
    "price": "10.00",
    "size": "0.01",
    "order_id": "d50ec984-77a8-460a-b958-66f114b0de9b",
    "created_at": "2014-11-07T22:19:28.578544Z",
    "liquidity": "T",
    "fee": "0.00025",
    "settled": true,
    "side": "buy"
}
*/
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

func (c *Client) listFills(timestamp, signature string, req *http.Request) ([]Fill, error) {
	res, err := c.do(timestamp, signature, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, orderError(res)
	}

	var fills []Fill
	if err := json.NewDecoder(res.Body).Decode(&fills); err != nil {
		return nil, err
	}

	return fills, nil
}
