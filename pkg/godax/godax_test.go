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

type fields struct {
	baseRestURL string
	baseWsURL   string
	key         string
	secret      string
	passphrase  string
	httpClient  *http.Client
}

func defaultFields() fields {
	return fields{
		baseRestURL: baseRestURL,
		baseWsURL:   baseWsURL,
		key:         key,
		secret:      secret,
		passphrase:  passphrase,
	}
}

func validateHeaders(t *testing.T, client *Client) {
	compareHeader(t, client, "CB-ACCESS-KEY", key)
	compareHeader(t, client, "CB-ACCESS-PASSPHRASE", passphrase)
	compareHeader(t, client, "User-Agent", userAgent)
	compareHeader(t, client, "Content-Type", "application/json")
	compareHeader(t, client, "Accept", "application/json")
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

// TODO: figure file
func TestClient_ListFills(t *testing.T) {
	orderID := "d50ec984-77a8-460a-b958-66f114b0de9b"
	productID := "BTC-USD"

	type args struct{ qp QueryParams }

	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    []Fill
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to list fills it returns a slice of fills",
			fields: defaultFields(),
			args:   args{qp: QueryParams{OrderID: orderID, ProductID: productID}},
			want: []Fill{{
				TradeID:   74,
				ProductID: "BTC-USD",
				Price:     "10.00",
				Size:      "0.01",
				OrderID:   "d50ec984-77a8-460a-b958-66f114b0de9b",
				CreatedAt: "2014-11-07T22:19:28.578544Z",
				Liquidity: "T",
				Fee:       "0.00025",
				Settled:   true,
				Side:      "buy",
			}},
			wantRaw: `[{
				"trade_id": 74,
				"product_id": "BTC-USD",
				"price": "10.00",
				"size": "0.01",
				"order_id": "d50ec984-77a8-460a-b958-66f114b0de9b",
				"created_at": "2014-11-07T22:19:28.578544Z",
				"liquidity": "T",
				"fee": "0.00025",
				"settled": true,
				"side": "buy"
			}]`,
		},
		{
			name:    "when neither an orderID or productID is provided by the caller, it should return an error",
			fields:  defaultFields(),
			args:    args{qp: QueryParams{}},
			want:    []Fill{},
			wantErr: true,
			wantRaw: `[]`,
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

			qp := QueryParams{OrderID: tt.args.qp[OrderID], ProductID: tt.args.qp[ProductID]}
			got, err := c.ListFills(qp)
			if err != nil {
				if tt.wantErr && err == ErrMissingOrderOrProductID {
					return
				}
				t.Errorf("Client.ListFills() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListFills() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: accounts_test.go ???
func TestClient_ListCoinbaseAccounts(t *testing.T) {
	tests := [...]struct {
		name    string
		fields  fields
		want    []CoinbaseAccount
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to list coinbase (non-pro) accounts",
			fields: defaultFields(),
			want: []CoinbaseAccount{{
				ID:       "2ae3354e-f1c3-5771-8a37-6228e9d239db",
				Name:     "USD Wallet",
				Balance:  "0.00",
				Currency: "USD",
				Type:     "fiat",
				Primary:  false,
				Active:   true,
				WireDepositInformation: WireDepositInfo{
					AccountNumber: "0199003122",
					RoutingNumber: "026013356",
					BankName:      "Metropolitan Commercial Bank",
					BankAddress:   "99 Park Ave 4th Fl New York, NY 10016",
					BankCountry: BankCountry{
						Code: "US",
						Name: "United States",
					},
					AccountName:    "Coinbase, Inc",
					AccountAddress: "548 Market Street, #23008, San Francisco, CA 94104",
					Reference:      "BAOCAEUX",
				},
			}, {
				ID:       "2a11354e-f133-5771-8a37-622be9b239db",
				Name:     "EUR Wallet",
				Balance:  "0.00",
				Currency: "EUR",
				Type:     "fiat",
				Primary:  false,
				Active:   true,
				SepaDepositInformation: SepaDepositInfo{
					IBAN:            "EE957700771001355096",
					Swift:           "LHVBEE22",
					BankName:        "AS LHV Pank",
					BankAddress:     "Tartu mnt 2, 10145 Tallinn, Estonia",
					BankCountryName: "Estonia",
					AccountName:     "Coinbase UK, Ltd.",
					AccountAddress:  "9th Floor, 107 Cheapside, London, EC2V 6DN, United Kingdom",
					Reference:       "CBAEUXOVFXOXYX",
				},
			}},
			wantErr: false,
			wantRaw: `[{
				"id": "2ae3354e-f1c3-5771-8a37-6228e9d239db",
				"name": "USD Wallet",
				"balance": "0.00",
				"currency": "USD",
				"type": "fiat",
				"primary": false,
				"active": true,
				"wire_deposit_information": {
					"account_number": "0199003122",
					"routing_number": "026013356",
					"bank_name": "Metropolitan Commercial Bank",
					"bank_address": "99 Park Ave 4th Fl New York, NY 10016",
					"bank_country": {
						"code": "US",
						"name": "United States"
					},
					"account_name": "Coinbase, Inc",
					"account_address": "548 Market Street, #23008, San Francisco, CA 94104",
					"reference": "BAOCAEUX"
				}
			}, {
				"id": "2a11354e-f133-5771-8a37-622be9b239db",
				"name": "EUR Wallet",
				"balance": "0.00",
				"currency": "EUR",
				"type": "fiat",
				"primary": false,
				"active": true,
				"sepa_deposit_information": {
					"iban": "EE957700771001355096",
					"swift": "LHVBEE22",
					"bank_name": "AS LHV Pank",
					"bank_address": "Tartu mnt 2, 10145 Tallinn, Estonia",
					"bank_country_name": "Estonia",
					"account_name": "Coinbase UK, Ltd.",
					"account_address": "9th Floor, 107 Cheapside, London, EC2V 6DN, United Kingdom",
					"reference": "CBAEUXOVFXOXYX"
				}
			}]`,
		},
	}
	for _, tt := range tests {
		mockClient := MockResponse(tt.wantRaw)

		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				baseRestURL: tt.fields.baseRestURL,
				baseWsURL:   tt.fields.baseWsURL,
				key:         tt.fields.key,
				secret:      tt.fields.secret,
				passphrase:  tt.fields.passphrase,
				httpClient:  mockClient,
			}

			got, err := c.ListCoinbaseAccounts()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListCoinbaseAccounts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListCoinbaseAccounts() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: figure file
func TestClient_GetTrailingVolume(t *testing.T) {
	tests := [...]struct {
		name    string
		fields  fields
		want    []UserAccount
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to get a users trailing volumes",
			fields: defaultFields(),
			want: []UserAccount{{
				ProductID:      "BTC-USD",
				ExchangeVolume: "11800.00000000",
				Volume:         "100.00000000",
				RecordedAt:     "1973-11-29T00:05:01.123456Z",
			}},
			wantRaw: `[{
				"product_id": "BTC-USD",
				"exchange_volume": "11800.00000000",
				"volume": "100.00000000",
				"recorded_at": "1973-11-29T00:05:01.123456Z"
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

			got, err := c.GetTrailingVolume()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetTrailingVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetTrailingVolume() = %v, want %v", got, tt.want)
			}
		})
	}
}

// FIXME: time_test.go ???
func TestClient_GetServerTime(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    ServerTime
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to get the coinbase pro server time",
			fields: defaultFields(),
			want: ServerTime{
				ISO:   "2015-01-07T23:47:25.201Z",
				Epoch: 1420674445.201,
			},
			wantRaw: `{
				"iso": "2015-01-07T23:47:25.201Z",
				"epoch": 1420674445.201
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

			got, err := c.GetServerTime()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetServerTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetServerTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
