package godax

// ServerTime represents a coinbase pro server time by providing ISO and epoch times.
type ServerTime struct {
	ISO   string  `json:"iso"`
	Epoch float64 `json:"epoch"`
}
