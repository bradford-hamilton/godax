package godax

// PaymentMethod represents a payment method connected the the api user's account and
// holds metadata like type, name, currency, a list of limits, etc.
/*
[
    {
        "id": "bc6d7162-d984-5ffa-963c-a493b1c1370b",
        "type": "ach_bank_account",
        "name": "Bank of America - eBan... ********7134",
        "currency": "USD",
        "primary_buy": true,
        "primary_sell": true,
        "allow_buy": true,
        "allow_sell": true,
        "allow_deposit": true,
        "allow_withdraw": true,
        "limits": {
            "buy": [
                {
                    "period_in_days": 1,
                    "total": {
                        "amount": "10000.00",
                        "currency": "USD"
                    },
                    "remaining": {
                        "amount": "10000.00",
                        "currency": "USD"
                    }
                }
            ],
            "instant_buy": [...],
            "sell": [
                {
                    "period_in_days": 1,
                    "total": {
                        "amount": "10000.00",
                        "currency": "USD"
                    },
                    "remaining": {
                        "amount": "10000.00",
                        "currency": "USD"
                    }
                }
            ],
            "deposit": [...]
        }
    }
]
*/
type PaymentMethod struct {
	ID            string   `json:"id"`
	Type          string   `json:"type"`
	Name          string   `json:"name"`
	Currency      string   `json:"currency"`
	PrimaryBuy    bool     `json:"primary_buy"`
	PrimarySell   bool     `json:"primary_sell"`
	AllowBuy      bool     `json:"allow_buy"`
	AllowSell     bool     `json:"allow_sell"`
	AllowDeposit  bool     `json:"allow_deposit"`
	AllowWithdraw bool     `json:"allow_withdraw"`
	Limits        PMLimits `json:"limits"`
}

// PMLimits represents the available payment method limits
type PMLimits struct {
	Buy        []PMLimit `json:"buy"`
	InstantBuy []PMLimit `json:"instant_buy"`
	Sell       []PMLimit `json:"sell"`
	Deposit    []PMLimit `json:"deposit"`
}

// PMLimit represents a single payment method limit
type PMLimit struct {
	PeriodInDays int      `json:"period_in_days"`
	Total        PMAmount `json:"total"`
	Remaining    PMAmount `json:"remaining"`
}

// PMAmount describes a currency and an amount. Used to describe a PMLimit.
type PMAmount struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}
