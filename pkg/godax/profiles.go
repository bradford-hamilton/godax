package godax

// Profile represents a coinbase pro profile which is equivalent to a portfolio.
/*
{
    "id": "86602c68-306a-4500-ac73-4ce56a91d83c",
    "user_id": "5844eceecf7e803e259d0365",
    "name": "default",
    "active": true,
    "is_default": true,
    "created_at": "2019-11-18T15:08:40.236309Z"
}
*/
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
/*
{
    "from": "86602c68-306a-4500-ac73-4ce56a91d83c",
    "to": "e87429d3-f0a7-4f28-8dff-8dd93d383de1",
    "currency": "BTC",
    "amount": "1000.00"
}
*/
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
