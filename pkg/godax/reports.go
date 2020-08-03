package godax

// ReportParams describe the body needed in a call to generate a report.
type ReportParams struct {
	// Type	is either "fills" or "account"
	Type string `json:"type"`

	// StartDate is the starting date for the report (inclusive)
	StartDate string `json:"start_date"`

	// EndDate is the ending date for the report (inclusive)
	EndDate string `json:"end_date"`

	// ProductID is the ID of the product to generate a fills report for.
	// E.g. BTC-USD. Required if type is fills
	ProductID string `json:"product_id,omitempty"`

	// AccountID is the ID of the account to generate an account report for. Required if type is account
	AccountID string `json:"account_id,omitempty"`

	// Format must either be "pdf" or "csv" (defualt is "pdf")
	Format string `json:"format,omitempty"`

	// Email is the email address to send the report to (optional)
	Email string `json:"email,omitempty"`
}

// ReportStatus represents the returned response after creating a report.
type ReportStatus struct {
	ID          string       `json:"id"`
	Type        string       `json:"type"`
	Status      string       `json:"status"`
	CreatedAt   string       `json:"created_at"`
	CompletedAt string       `json:"completed_at"`
	ExpiresAt   string       `json:"expires_at"`
	FileURL     string       `json:"file_url"`
	Params      ReportParams `json:"params"`
}
