package godax

// PaymentMethod represents a payment method connected the the api user's account and
// holds metadata like type, name, currency, a list of limits, etc.
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
