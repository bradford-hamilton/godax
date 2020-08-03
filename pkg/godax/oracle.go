package godax

// Oracle represents cryptographically signed prices ready to be posted on-chain.
type Oracle struct {
	// Timestamp indicates when the latest datapoint was obtained.
	Timestamp string `json:"timestamp"`

	// Messages array contains abi-encoded values [kind, timestamp, key, value], where kind always
	// equals to 'prices', timestamp is the time when the price was obtained, key is asset ticker
	// (e.g. 'eth') and value is asset price.
	Messages []string `json:"messages"`

	// Signatures is an array of Ethereum-compatible ECDSA signatures for each message.
	Signatures []string `json:"signatures"`

	// Prices contains human-readable asset prices.
	Prices map[string]string `json:"prices"`
}
