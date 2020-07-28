package godax

import (
	"encoding/json"
	"net/http"
)

// ReportParams describe the body needed in a call to generate a report.
/*
{
    "type": "fills",
    "start_date": "2014-11-01T00:00:00.000Z",
    "end_date": "2014-11-30T23:59:59.000Z",
    ...
}
*/
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
/*
{
    "id": "0428b97b-bec1-429e-a94c-59232926778d",
    "type": "fills",
    "status": "pending",
    "created_at": "2015-01-06T10:34:47.000Z",
    "completed_at": undefined,
    "expires_at": "2015-01-13T10:35:47.000Z",
    "file_url": undefined,
    "params": {
        "start_date": "2014-11-01T00:00:00.000Z",
        "end_date": "2014-11-30T23:59:59.000Z"
    }
*/
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

func (c *Client) createReport(timestamp, signature string, req *http.Request) (ReportStatus, error) {
	res, err := c.do(timestamp, signature, req)
	if err != nil {
		return ReportStatus{}, err
	}
	defer res.Body.Close()

	var report ReportStatus
	if err := json.NewDecoder(res.Body).Decode(&report); err != nil {
		return ReportStatus{}, err
	}
	return report, nil
}

func (c *Client) getReportStatus(timestamp, signature string, req *http.Request) (ReportStatus, error) {
	res, err := c.do(timestamp, signature, req)
	if err != nil {
		return ReportStatus{}, err
	}
	defer res.Body.Close()

	var rs ReportStatus
	if err := json.NewDecoder(res.Body).Decode(&rs); err != nil {
		return ReportStatus{}, err
	}
	return rs, nil
}
