package godax

import (
	"encoding/json"
	"net/http"
)

// ExchangeLimit represents info about your payment method transfer limits, as well as buy/sell limits per currency.
/*
{
  "limit_currency": "USD",
  "transfer_limits": {
    "ach": {
      "BAT": {
        "max": "21267.54245368",
        "remaining": "21267.54245368",
        "period_in_days": 7
      }
    },
    ...
    "buy": {
      "BAT": {
        "max": "212677.90003268",
        "remaining": "212677.90003268",
        "period_in_days": 7
      }
    },
    ...
  }
}
*/
type ExchangeLimit struct {
	LimitCurrency  string         `json:"limit_currency"`
	TransferLimits TransferLimits `json:"transfer_limits"`
}

// TransferLimits represents info about your payment method transfer limits, as well as buy/sell limits per currency.
type TransferLimits struct {
	Ach                  map[string]Limit `json:"ach"`
	AchNoBalance         map[string]Limit `json:"ach_no_balance"`
	CreditDebitCard      map[string]Limit `json:"credit_debit_card"`
	AchCurm              map[string]Limit `json:"ach_curm"`
	Secure3DBuy          map[string]Limit `json:"secure3d_buy"`
	ExchangeWithdraw     map[string]Limit `json:"exchange_withdraw"`
	ExchangeAch          map[string]Limit `json:"exchange_ach"`
	PaypalWithdrawal     map[string]Limit `json:"paypal_withdrawal"`
	InstantAchWithdrawal map[string]Limit `json:"instant_ach_withdrawal"`
	Buy                  map[string]Limit `json:"buy"`
	Sell                 map[string]Limit `json:"sell"`
}

// Limit represents the actual limit metadata per currency
/*
{
	"BAT": {
		"max": "212677.90003268",
		"remaining": "212677.90003268",
		"period_in_days": 7
	}
}
*/
type Limit struct {
	// Max represents the max you can use per currency
	Max float64 `json:"max"`

	// Remaining represents your remaining amount per currency
	Remaining float64 `json:"remaining"`

	// TODO: I'm not seeing this in the response ever, and the api actually returns
	// floats and not strings for "max" and "remaining"... Let coinbase know at some
	// point to update docs.
	// PeriodInDays int `json:"period_in_days"`
}

func (c *Client) getLimits(timestamp, signature string, req *http.Request) (ExchangeLimit, error) {
	res, err := c.do(timestamp, signature, req)
	if err != nil {
		return ExchangeLimit{}, err
	}
	defer res.Body.Close()

	var limit ExchangeLimit
	if err := json.NewDecoder(res.Body).Decode(&limit); err != nil {
		return ExchangeLimit{}, err
	}
	return limit, nil
}
