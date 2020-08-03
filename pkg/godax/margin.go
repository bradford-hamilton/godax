package godax

// MarginProfile represents a snapshot of a profile's margin capabilities.
type MarginProfile struct {
	ProfileID             string                `json:"profile_id"`
	MarginInitialEquity   string                `json:"margin_initial_equity"`
	MarginWarningEquity   string                `json:"margin_warning_equity"`
	MarginCallEquity      string                `json:"margin_call_equity"`
	EquityPercentage      float64               `json:"equity_percentage"`
	SellingPower          float64               `json:"selling_power"`
	BuyingPower           float64               `json:"buying_power"`
	BorrowPower           float64               `json:"borrow_power"`
	InterestRate          string                `json:"interest_rate"`
	InterestPaid          string                `json:"interest_paid"`
	CollateralCurrencies  []string              `json:"collateral_currencies"`
	CollateralHoldValue   string                `json:"collateral_hold_value"`
	LastLiquidationAt     string                `json:"last_liquidation_at"`
	AvailableBorrowLimits AvailableBorrowLimits `json:"available_borrow_limits"`
	BorrowLimit           string                `json:"borrow_limit"`
	TopUpAmounts          TopUpAmounts          `json:"top_up_amounts"`
}

// AvailableBorrowLimits describes your marginable and non-marginale limits.
type AvailableBorrowLimits struct {
	MarginableLimit    float64 `json:"marginable_limit"`
	NonMarginableLimit float64 `json:"nonmarginable_limit"`
}

// TopUpAmounts shows borrowable and non-borrowable usd amounts.
type TopUpAmounts struct {
	BorrowableUsd    string `json:"borrowable_usd"`
	NonBorrowableUsd string `json:"non_borrowable_usd"`
}

// BuyingPower represents a buying power and selling power for a particular product.
// Used for marshalling response from GetBuyingPower
type BuyingPower struct {
	BuyingPower            float64 `json:"buying_power"`
	SellingPower           float64 `json:"selling_power"`
	BuyingPowerExplanation string  `json:"buying_power_explanation"`
}

// CurrencyWithdrawalPower represents the withdrawal power for a specific currency
type CurrencyWithdrawalPower struct {
	ProfileID       string `json:"profile_id"`
	WithdrawalPower string `json:"withdrawal_power"`
}

// AllWithdrawalPower represents the max amount of each currency that you can withdraw
// from your margin profile. Used for calls to GetAllWithdrawalPower.
type AllWithdrawalPower struct {
	ProfileID                  string                      `json:"profile_id"`
	MarginableWithdrawalPowers []MarginableWithdrawalPower `json:"marginable_withdrawal_powers"`
}

// MarginableWithdrawalPower describes a currency and the mount of withdrawal power you have.
type MarginableWithdrawalPower struct {
	Currency        string `json:"currency"`
	WithdrawalPower string `json:"withdrawal_power"`
}

// ExitPlan represents a liquidation strategy that can be performed to get your equity
// percentage back to an acceptable level
type ExitPlan struct {
	ID                  string            `json:"id"`
	UserID              string            `json:"userId"`
	ProfileID           string            `json:"profileId"`
	AccountsList        []ExitPlanAccount `json:"accountsList"`
	EquityPercentage    string            `json:"equityPercentage"`
	TotalAssetsUSD      string            `json:"totalAssetsUsd"`
	TotalLiabilitiesUSD string            `json:"totalLiabilitiesUsd"`
	StrategiesList      []ExitStrategy    `json:"strategiesList"`
	CreatedAt           string            `json:"createdAt"`
}

// ExitPlanAccount represents an account within an exit plan.
type ExitPlanAccount struct {
	ID       string `json:"id"`
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
}

// ExitStrategy is a strategy that can be performed to get your equity
// percentage back to an acceptable level
type ExitStrategy struct {
	Type      string `json:"type"`
	Amount    string `json:"amount"`
	Product   string `json:"product"`
	Strategy  string `json:"strategy"`
	AccountID string `json:"accountId"`
	OrderID   string `json:"orderId"`
}

// LiquidationEvent represents a liquididation event.
type LiquidationEvent struct {
	EventID   string     `json:"event_id"`
	EventTime string     `json:"event_time"`
	Orders    []LiqOrder `json:"orders"`
}

// LiqOrder is the order metadata attached to a liquidation event.
type LiqOrder struct {
	ID            string `json:"id"`
	Size          string `json:"size"`
	ProductID     string `json:"product_id"`
	ProfileID     string `json:"profile_id"`
	Side          string `json:"side"`
	Type          string `json:"type"`
	PostOnly      bool   `json:"post_only"`
	CreatedAt     string `json:"created_at"`
	DoneAt        string `json:"done_at"`
	DoneReason    string `json:"done_reason"`
	FillFees      string `json:"fill_fees"`
	FilledSize    string `json:"filled_size"`
	ExecutedValue string `json:"executed_value"`
	Status        string `json:"status"`
	Settled       bool   `json:"settled"`
}

// RefreshAmount represents amount in USD of loans that will be renewed in
// the next day and then the day after.
type RefreshAmount struct {
	OneDayRenewalAmount string `json:"oneDayRenewalAmount"`
	TwoDayRenewalAmount string `json:"twoDayRenewalAmount"`
}

// MarginStatus represents the current status of an account's margin.
type MarginStatus struct {
	Tier     int  `json:"tier"`
	Enabled  bool `json:"enabled"`
	Eligible bool `json:"eligible"`
}
