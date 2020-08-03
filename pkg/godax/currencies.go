package godax

// Currency represents an available currency on coinbase pro.
type Currency struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	MinSize string `json:"min_size"`
}
