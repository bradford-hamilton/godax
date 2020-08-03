package godax

// Conversion represents the return value value from a call to StableCoinConversion.
// It describes different metadata around the stablecoin conversion.
type Conversion struct {
	ID            string `json:"id"`
	Amount        string `json:"amount"`
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	From          string `json:"from"`
	To            string `json:"to"`
}
