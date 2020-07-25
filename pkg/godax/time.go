package godax

import (
	"encoding/json"
	"net/http"
)

// ServerTime represents a coinbase pro server time by providing ISO and epoch times.
/*
{
    "iso": "2015-01-07T23:47:25.201Z",
    "epoch": 1420674445.201
}
*/
type ServerTime struct {
	ISO   string  `json:"iso"`
	Epoch float64 `json:"epoch"`
}

func (c *Client) getServerTime(timestamp, signature string, req *http.Request) (ServerTime, error) {
	res, err := c.do(timestamp, signature, req)
	if err != nil {
		return ServerTime{}, err
	}
	defer res.Body.Close()

	var srvTime ServerTime
	if err := json.NewDecoder(res.Body).Decode(&srvTime); err != nil {
		return ServerTime{}, err
	}
	return srvTime, nil
}
