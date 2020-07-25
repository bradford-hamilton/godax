package godax

import (
	"encoding/json"
	"net/http"
)

// Currency represents an available currency on coinbase pro.
/*
{
    "id": "BTC",
    "name": "Bitcoin",
    "min_size": "0.00000001"
}
*/
type Currency struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	MinSize string `json:"min_size"`
}

func (c *Client) listCurrencies(timestamp, signature string, req *http.Request) ([]Currency, error) {
	res, err := c.do(timestamp, signature, req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var currencies []Currency
	if err := json.NewDecoder(res.Body).Decode(&currencies); err != nil {
		return nil, err
	}
	return currencies, nil
}
