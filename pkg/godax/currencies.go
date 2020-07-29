package godax

// Currency represents an available currency on coinbase pro.
/*
{
    "id": "BTC",
    "name": "Bitcoin",
    "min_size": "0.00000001"
}
*/
type Currency struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	MinSize string `json:"min_size"`
}
