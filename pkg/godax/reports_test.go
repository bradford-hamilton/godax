package godax

import (
	"reflect"
	"testing"
)

func TestClient_CreateReport(t *testing.T) {
	type args struct {
		report ReportParams
	}
	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    ReportStatus
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to create a report",
			fields: defaultFields(),
			args: args{report: ReportParams{
				Type:      "fills",
				StartDate: "2014-11-01T00:00:00.000Z",
				EndDate:   "2014-11-30T23:59:59.000Z",
			}},
			want: ReportStatus{
				ID:          "0428b97b-bec1-429e-a94c-59232926778d",
				Type:        "fills",
				Status:      "pending",
				CreatedAt:   "2015-01-06T10:34:47.000Z",
				CompletedAt: "",
				ExpiresAt:   "2015-01-13T10:35:47.000Z",
				FileURL:     "https://gdax-reports-sandbox.s3.amazonaws.com/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx/account.pdf?AWSAccessKeyId=abc123&Expires=1596334414&Signature=neatsignature456",
				Params: ReportParams{
					StartDate: "2014-11-01T00:00:00.000Z",
					EndDate:   "2014-11-30T23:59:59.000Z",
				},
			},
			wantRaw: `{
                "id": "0428b97b-bec1-429e-a94c-59232926778d",
                "type": "fills",
                "status": "pending",
                "created_at": "2015-01-06T10:34:47.000Z",
                "completed_at": "",
                "expires_at": "2015-01-13T10:35:47.000Z",
                "file_url": "https://gdax-reports-sandbox.s3.amazonaws.com/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx/account.pdf?AWSAccessKeyId=abc123&Expires=1596334414&Signature=neatsignature456",
                "params": {
                    "start_date": "2014-11-01T00:00:00.000Z",
                    "end_date": "2014-11-30T23:59:59.000Z"
                }
            }`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := MockResponse(tt.wantRaw)

			c := &Client{
				baseRestURL: tt.fields.baseRestURL,
				baseWsURL:   tt.fields.baseWsURL,
				key:         tt.fields.key,
				secret:      tt.fields.secret,
				passphrase:  tt.fields.passphrase,
				httpClient:  mockClient,
			}

			got, err := c.CreateReport(tt.args.report)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.CreateReport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetReportStatus(t *testing.T) {
	type args struct {
		reportID string
	}
	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    ReportStatus
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful market order is made to create a report",
			fields: defaultFields(),
			args:   args{reportID: "0428b97b-bec1-429e-a94c-59232926778d"},
			want: ReportStatus{
				ID:          "0428b97b-bec1-429e-a94c-59232926778d",
				Type:        "account",
				Status:      "ready",
				CreatedAt:   "2016-08-26T11:56:12.000Z",
				CompletedAt: "",
				ExpiresAt:   "2016-09-14T11:22:11.000Z",
				FileURL:     "https://gdax-reports-sandbox.s3.amazonaws.com/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx/account.pdf?AWSAccessKeyId=abc123&Expires=1596334414&Signature=neatsignature456",
				Params: ReportParams{
					StartDate: "2014-11-01T00:00:00.000Z",
					EndDate:   "2014-11-30T23:59:59.000Z",
					AccountID: "f1f2404a-7de7-4cf6-81f9-5cb0256c8cea",
					Format:    "pdf",
					Email:     "brad.lamson@gmail.com",
				},
			},
			wantRaw: `{
                "id": "0428b97b-bec1-429e-a94c-59232926778d",
                "type": "account",
                "status": "ready",
                "created_at": "2016-08-26T11:56:12.000Z",
                "completed_at": "",
                "expires_at": "2016-09-14T11:22:11.000Z",
                "file_url": "https://gdax-reports-sandbox.s3.amazonaws.com/xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx/account.pdf?AWSAccessKeyId=abc123&Expires=1596334414&Signature=neatsignature456",
                "params": {
                    "start_date": "2014-11-01T00:00:00.000Z",
					"end_date": "2014-11-30T23:59:59.000Z",
					"format": "pdf",
					"email": "brad.lamson@gmail.com",
					"account_id": "f1f2404a-7de7-4cf6-81f9-5cb0256c8cea"
                }
            }`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := MockResponse(tt.wantRaw)

			c := &Client{
				baseRestURL: tt.fields.baseRestURL,
				baseWsURL:   tt.fields.baseWsURL,
				key:         tt.fields.key,
				secret:      tt.fields.secret,
				passphrase:  tt.fields.passphrase,
				httpClient:  mockClient,
			}

			got, err := c.GetReportStatus(tt.args.reportID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetReportStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetReportStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
