package godax

// Conversion represents the return value value from a call to StableCoinConversion.
// It describes different metadata around the stablecoin conversion.
/*
{
    "id": "8942caee-f9d5-4600-a894-4811268545db",
    "amount": "10000.00",
    "from_account_id": "7849cc79-8b01-4793-9345-bc6b5f08acce",
    "to_account_id": "105c3e58-0898-4106-8283-dc5781cda07b",
    "from": "USD",
    "to": "USDC"
}
*/
type Conversion struct {
	ID            string `json:"id"`
	Amount        string `json:"amount"`
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	From          string `json:"from"`
	To            string `json:"to"`
}

// conversionReq represents the body needed in a StableCoinConversion call.
/*
{
    "from": "USD",
    "to": "USDC",
    "amount": "10000.00"
}
*/
type conversionReq struct {
	To     string `json:"to"`
	From   string `json:"from"`
	Amount string `json:"amount"`
}
