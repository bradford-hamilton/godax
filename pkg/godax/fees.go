package godax

import (
	"encoding/json"
	"net/http"
)

// Fees represent the current maker & taker fees, and usd volume.
/*
{
    "maker_fee_rate": "0.0015",
    "taker_fee_rate": "0.0025",
    "usd_volume": "25000.00"
}
*/
type Fees struct {
	MakerFeeRate string `json:"maker_fee_rate"`
	TakerFeeRate string `json:"taker_fee_rate"`
	USDVolume    string `json:"usd_volume"`
}

func (c *Client) getCurrentFees(timestamp, signature string, req *http.Request) (Fees, error) {
	res, err := c.do(timestamp, signature, req)
	if err != nil {
		return Fees{}, err
	}
	defer res.Body.Close()

	var fees Fees
	if err := json.NewDecoder(res.Body).Decode(&fees); err != nil {
		return Fees{}, err
	}
	return fees, nil
}
