package godax

// MarginProfile ...
/*
{
    "profile_id": "8058d771-2d88-4f0f-ab6e-299c153d4308",
    "margin_initial_equity": "0.33",
    "margin_warning_equity": "0.2",
    "margin_call_equity": "0.15",
    "equity_percentage": 0.8725562096924747,
    "selling_power": 0.00221896,
    "buying_power": 23.51,
    "borrow_power": 23.51,
    "interest_rate": "0",
    "interest_paid": "0.3205913399694425",
    "collateral_currencies": [
        "BTC",
        "USD",
        "USDC"
    ],
    "collateral_hold_value": "1.0050000000000000",
    "last_liquidation_at": "2019-11-21T14:58:49.879Z",
    "available_borrow_limits": {
        "marginable_limit": 23.51,
        "nonmarginable_limit": 7.75
    },
    "borrow_limit": "5000",
    "top_up_amounts": {
        "borrowable_usd": "0",
        "non_borrowable_usd": "0"
    }
}
*/
type MarginProfile struct {
	ProfileID             string   `json:"profile_id"`
	MarginInitialEquity   string   `json:"margin_initial_equity"`
	MarginWarningEquity   string   `json:"margin_warning_equity"`
	MarginCallEquity      string   `json:"margin_call_equity"`
	EquityPercentage      float64  `json:"equity_percentage"`
	SellingPower          float64  `json:"selling_power"`
	BuyingPower           float64  `json:"buying_power"`
	BorrowPower           float64  `json:"borrow_power"`
	InterestRate          string   `json:"interest_rate"`
	InterestPaid          string   `json:"interest_paid"`
	CollateralCurrencies  []string `json:"collateral_currencies"`
	CollateralHoldValue   string   `json:"collateral_hold_value"`
	LastLiquidationAt     string   `json:"last_liquidation_at"`
	AvailableBorrowLimits struct {
		MarginableLimit    float64 `json:"marginable_limit"`
		NonMarginableLimit float64 `json:"nonmarginable_limit"`
	} `json:"available_borrow_limits"`
	BorrowLimit  string `json:"borrow_limit"`
	TopUpAmounts struct {
		BorrowableUsd    float64 `json:"borrowable_usd"`
		NonBorrowableUsd float64 `json:"non_borrowable_usd"`
	} `json:"top_up_amounts"`
}

// BuyingPower represents a buying power and selling power for a particular product. Used for marshalling response
// from GetBuyingPower
/*
{
    "buying_power": 23.53,
    "selling_power": 0.00221896,
    "buying_power_explanation": "This is the line of credit available to you on the BTC-USD market, given how much collateral assets you currently have in your portfolio."
}
*/
type BuyingPower struct {
	BuyingPower            float64 `json:"buying_power"`
	SellingPower           float64 `json:"selling_power"`
	BuyingPowerExplanation string  `json:"buying_power_explanation"`
}

// CurrencyWithdrawalPower represents the withdrawal power for a specific currency
/*
{
    "profile_id": "8058d771-2d88-4f0f-ab6e-299c153d4308",
    "withdrawal_power": "7.77569088416849750000"
}
*/
type CurrencyWithdrawalPower struct {
	ProfileID       string `json:"profile_id"`
	WithdrawalPower string `json:"withdrawal_power"`
}

// AllWithdrawalPower represents the max amount of each currency that you can withdraw from your margin profile. Used
// for calls to GetAllWithdrawalPower.
/*
{
    "profile_id": "8058d771-2d88-4f0f-ab6e-299c153d4308",
    "marginable_withdrawal_powers": [
        {
            "currency": "ETH",
            "withdrawal_power": "0.0000000000000000"
        },
        {
            "currency": "BTC",
            "withdrawal_power": "0.00184821818021342913"
        },
        {
            "currency": "USD",
            "withdrawal_power": "7.77601796034649750000"
        },
        {
            "currency": "USDC",
            "withdrawal_power": "1.00332803238200000000"
        }
    ]
}
*/
type AllWithdrawalPower struct {
	ProfileID                  string `json:"profile_id"`
	MarginableWithdrawalPowers []struct {
		Currency        string `json:"currency"`
		WithdrawalPower string `json:"withdrawal_power"`
	} `json:"marginable_withdrawal_powers"`
}
