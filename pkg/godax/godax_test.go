package godax

import (
	"net/http"
	"reflect"
	"testing"
)

const (
	baseRestURL = "https://test.pro.coinbase"
	baseWsURL   = "wss://test.ws-feed.pro.coinbase.com"
	key         = "super_secret_key_123_abc"
	secret      = "MTIzYWJjU3VwZXJTZWNyZXRTZWNyZXQ="
	passphrase  = "1q2w3e4r"
)

func TestClient_ListAccounts(t *testing.T) {
	type fields struct {
		baseRestURL string
		baseWsURL   string
		key         string
		secret      string
		passphrase  string
		httpClient  *http.Client
	}
	genFields := func() fields {
		return fields{
			baseRestURL: baseRestURL,
			baseWsURL:   baseWsURL,
			key:         key,
			secret:      secret,
			passphrase:  passphrase,
		}
	}
	tests := []struct {
		name    string
		fields  fields
		mock    HTTPClient
		want    []ListAccount
		wantRaw string
		wantErr bool
	}{
		{
			name:    "when a successful call is made to ListAccounts with no results",
			fields:  genFields(),
			want:    []ListAccount{},
			wantRaw: `[]`,
		},
		{
			name:   "when a successful call is made to ListAccounts with one account",
			fields: genFields(),
			want: []ListAccount{{
				ID:        "f1f2404a-7de7-4cf6-81f9-5cb0256c8cea",
				Currency:  "BTC",
				Balance:   "10.01",
				Available: "15.449977",
				Hold:      "wat",
			}},
			wantRaw: `[{
                "id": "f1f2404a-7de7-4cf6-81f9-5cb0256c8cea",
                "currency": "BTC",
                "balance": "10.01",
                "available": "15.449977",
                "hold": "wat"
            }]`,
		},
		{
			name:   "when a successful call is made to ListAccounts with many accounts",
			fields: genFields(),
			want: []ListAccount{{
				ID:        "766b7a10-06bb-4b1d-a4b3-679d025352ad",
				Currency:  "BTC",
				Balance:   "0.00000000000",
				Available: "123.456789",
				Hold:      "0.101",
			}, {
				ID:        "543c3da9-a71d-4a4c-b6e7-edc132ff704e",
				Currency:  "ETH",
				Balance:   "1.300006",
				Available: "9000.67685938262624",
				Hold:      "0.101",
			}, {
				ID:        "dcbe61c2-a1bd-444c-b41a-3c6b2363afd6",
				Currency:  "BAT",
				Balance:   "9999.677773333",
				Available: "9000.67685938262624",
				Hold:      "0.101",
			}},
			wantRaw: `[{
                "id": "766b7a10-06bb-4b1d-a4b3-679d025352ad",
                "currency": "BTC",
                "balance": "0.00000000000",
                "available": "123.456789",
                "hold": "0.101"
            },{
                "id": "543c3da9-a71d-4a4c-b6e7-edc132ff704e",
                "currency": "ETH",
                "balance": "1.300006",
                "available": "9000.67685938262624",
                "hold": "0.101"
            },{
                "id": "dcbe61c2-a1bd-444c-b41a-3c6b2363afd6",
                "currency": "BAT",
                "balance": "9999.677773333",
                "available": "9000.67685938262624",
                "hold": "0.101"
            }]`,
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

			got, err := c.ListAccounts()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListAccounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetAccount(t *testing.T) {
	type fields struct {
		baseRestURL string
		baseWsURL   string
		key         string
		secret      string
		passphrase  string
		httpClient  *http.Client
	}
	genFields := func() fields {
		return fields{
			baseRestURL: baseRestURL,
			baseWsURL:   baseWsURL,
			key:         key,
			secret:      secret,
			passphrase:  passphrase,
		}
	}
	type args struct {
		accountID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Account
		wantRaw string
		wantErr bool
	}{
		{
			name:    "when a successful call is made to GetAccount and no account is found",
			fields:  genFields(),
			args:    args{accountID: "1q2w3e4r"},
			want:    Account{},
			wantRaw: `{}`,
		},
		{
			name:   "when a successful call is made to GetAccount and an account is found",
			fields: genFields(),
			args:   args{accountID: "a1b2c3d4"},
			want: Account{
				ID:        "a1b2c3d4",
				Balance:   "1.100",
				Holds:     "0.100",
				Available: "101.56",
				Currency:  "USD",
			},
			wantRaw: `{
                "id": "a1b2c3d4",
                "balance": "1.100",
                "holds": "0.100",
                "available": "101.56",
                "currency": "USD"
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

			got, err := c.GetAccount(tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func validateHeaders(t *testing.T, client *Client) {
	compareHeader(t, client, "CB-ACCESS-KEY", key)
	compareHeader(t, client, "CB-ACCESS-PASSPHRASE", passphrase)
	compareHeader(t, client, "User-Agent", userAgent)
	validateHeaderPresent(t, client, "CB-ACCESS-SIGN")
	validateHeaderPresent(t, client, "CB-ACCESS-TIMESTAMP")
}

func compareHeader(t *testing.T, c *Client, wantHeader string, wantContent string) {
	if c.httpClient.(*MockClient).Requests[0].Header.Get(wantHeader) != wantContent {
		t.Errorf(
			"%s header should be %s, was '%s'\n",
			wantHeader,
			wantContent,
			c.httpClient.(*MockClient).Requests[0].Header.Get(wantHeader),
		)
	}
}

func validateHeaderPresent(t *testing.T, c *Client, wantHeader string) {
	if c.httpClient.(*MockClient).Requests[0].Header.Get(wantHeader) == "" {
		t.Errorf("%s header should not be empty\n", wantHeader)
	}
}
