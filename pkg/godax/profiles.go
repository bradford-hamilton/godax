package godax

// Profile represents a coinbase pro profile which is equivalent to a portfolio.
type Profile struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Active    bool   `json:"active"`
	IsDefault bool   `json:"is_default"`
	CreatedAt string `json:"created_at"`
}

// TransferParams represent all the required data you must provide for a call to transfer
// cryptocurrency between accounts.
type TransferParams struct {
	// From is the profile id the API key belongs to and where the funds are sourced
	From string `json:"from"`

	// To represents the target profile id of where funds will be transferred to
	To string `json:"to"`

	// Currency - i.e. BTC or USD
	Currency string `json:"currency"`

	// Amount is the amount of currency to be transferred
	Amount string `json:"amount"`
}
