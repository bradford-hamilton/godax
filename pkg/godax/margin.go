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

// ExitPlan represents a liquidation strategy that can be performed to get your equity
// percentage back to an acceptable level
/*
{
    "id": "239f4dc6-72b6-11ea-b311-168e5016c449",
    "userId": "5cf6e115aaf44503db300f1e",
    "profileId": "8058d771-2d88-4f0f-ab6e-299c153d4308",
    "accountsList": [
        {
            "id": "434e1152-8eb5-4bfa-89a1-92bb1dcaf0c3",
            "currency": "BTC",
            "amount": "0.00221897"
        },
        {
            "id": "6d326768-71f2-4068-99dc-7075c78f6402",
            "currency": "USD",
            "amount": "-1.9004458409934425"
        },
        {
            "id": "120c8fcf-94da-4b45-9c43-18f114880f7a",
            "currency": "USDC",
            "amount": "1.003328032382"
        }
    ],
    "equityPercentage": "0.8744507743595747",
    "totalAssetsUsd": "15.137057447382",
    "totalLiabilitiesUsd": "1.9004458409934425",
    "strategiesList": [{
        "type": "",
        "amount": "",
        "product": "",
        "strategy": "",
        "accountId": "",
        "orderId": ""
    }],
    "createdAt": "2020-03-30 18:41:59.547863064 +0000 UTC m=+260120.906569441"
}
*/
type ExitPlan struct {
	ID           string `json:"id"`
	UserID       string `json:"userId"`
	ProfileID    string `json:"profileId"`
	AccountsList []struct {
		ID       string `json:"id"`
		Currency string `json:"currency"`
		Amount   string `json:"amount"`
	} `json:"accountsList"`
	EquityPercentage    string         `json:"equityPercentage"`
	TotalAssetsUSD      string         `json:"totalAssetsUsd"`
	TotalLiabilitiesUSD string         `json:"totalLiabilitiesUsd"`
	StrategiesList      []ExitStrategy `json:"strategiesList"`
	CreatedAt           string         `json:"createdAt"`
}

// ExitStrategy is an exit plan strategy
/*
{
    "type": "",
    "amount": "",
    "product": "",
    "strategy": "",
    "accountId": "",
    "orderId": ""
}
*/
type ExitStrategy struct {
	Type      string `json:"type"`
	Amount    string `json:"amount"`
	Product   string `json:"product"`
	Strategy  string `json:"strategy"`
	AccountID string `json:"accountId"`
	OrderID   string `json:"orderId"`
}
