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

// TODO: figure file
func TestClient_ListProducts(t *testing.T) {
	tests := [...]struct {
		name    string
		fields  fields
		want    []Product
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to get a profile by ID",
			fields: defaultFields(),
			want: []Product{{
				ID:              "BAT-USDC",
				DisplayName:     "BAT/USDC",
				BaseCurrency:    "BAT",
				QuoteCurrency:   "USDC",
				BaseIncrement:   "0.00000100",
				QuoteIncrement:  "0.00000100",
				BaseMinSize:     "1.00000000",
				BaseMaxSize:     "300000.00000000",
				MinMarketFunds:  "1",
				MaxMarketFunds:  "100000",
				Status:          "online",
				StatusMessage:   "",
				CancelOnly:      false,
				LimitOnly:       false,
				PostOnly:        false,
				TradingDisabled: false,
			}, {
				ID:              "LINK-USDC",
				DisplayName:     "LINK/USDC",
				BaseCurrency:    "LINK",
				QuoteCurrency:   "USDC",
				BaseIncrement:   "1.00000000",
				QuoteIncrement:  "0.00000100",
				BaseMinSize:     "1.00000000",
				BaseMaxSize:     "800000.00000000",
				MinMarketFunds:  "10",
				MaxMarketFunds:  "100000",
				Status:          "online",
				StatusMessage:   "",
				CancelOnly:      false,
				LimitOnly:       false,
				PostOnly:        false,
				TradingDisabled: false,
			}},
			wantRaw: `[{
				"id":"BAT-USDC",
				"base_currency":"BAT",
				"quote_currency":"USDC",
				"base_min_size":"1.00000000",
				"base_max_size":"300000.00000000",
				"quote_increment":"0.00000100",
				"base_increment":"0.00000100",
				"display_name":"BAT/USDC",
				"min_market_funds":"1",
				"max_market_funds":"100000",
				"margin_enabled":false,
				"post_only":false,
				"limit_only":false,
				"cancel_only":false,
				"trading_disabled":false,
				"status":"online",
				"status_message":""
			}, {
				"id":"LINK-USDC",
				"base_currency":"LINK",
				"quote_currency":"USDC",
				"base_min_size":"1.00000000",
				"base_max_size":"800000.00000000",
				"quote_increment":"0.00000100",
				"base_increment":"1.00000000",
				"display_name":"LINK/USDC",
				"min_market_funds":"10",
				"max_market_funds":"100000",
				"margin_enabled":false,
				"post_only":false,
				"limit_only":false,
				"cancel_only":false,
				"trading_disabled":false,
				"status":"online",
				"status_message":""
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

			got, err := c.ListProducts()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListProducts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListProducts() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: figure file
func TestClient_GetProduct(t *testing.T) {
	type args struct {
		productID string
	}
	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    Product
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to get a profile by ID",
			fields: defaultFields(),
			args:   args{productID: "BTC-USD"},
			want: Product{
				ID:              "BAT-USDC",
				DisplayName:     "BAT/USDC",
				BaseCurrency:    "BAT",
				QuoteCurrency:   "USDC",
				BaseIncrement:   "0.00000100",
				QuoteIncrement:  "0.00000100",
				BaseMinSize:     "1.00000000",
				BaseMaxSize:     "300000.00000000",
				MinMarketFunds:  "1",
				MaxMarketFunds:  "100000",
				Status:          "online",
				StatusMessage:   "",
				CancelOnly:      false,
				LimitOnly:       false,
				PostOnly:        false,
				TradingDisabled: false,
			},
			wantRaw: `{
				"id":"BAT-USDC",
				"base_currency":"BAT",
				"quote_currency":"USDC",
				"base_min_size":"1.00000000",
				"base_max_size":"300000.00000000",
				"quote_increment":"0.00000100",
				"base_increment":"0.00000100",
				"display_name":"BAT/USDC",
				"min_market_funds":"1",
				"max_market_funds":"100000",
				"margin_enabled":false,
				"post_only":false,
				"limit_only":false,
				"cancel_only":false,
				"trading_disabled":false,
				"status":"online",
				"status_message":""
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

			got, err := c.GetProductByID(tt.args.productID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: figure file
func TestClient_GetProductOrderBook(t *testing.T) {
	type args struct {
		productID string
		qp        QueryParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    OrderBook
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to list trades by product",
			fields: defaultFields(),
			args:   args{qp: QueryParams{Level: "1"}},
			want: OrderBook{
				Sequence: 3,
				Bids: []OrderBookOrder{{
					Price:     "295.96",
					Size:      "4.39088265",
					NumOrders: 2,
				}},
				Asks: []OrderBookOrder{{
					Price:     "295.97",
					Size:      "25.23542881",
					NumOrders: 12,
				}},
			},
			wantRaw: `{
				"sequence": 3,
				"bids": [
					["295.96", "4.39088265", 2]
				],
				"asks": [
					["295.97", "25.23542881", 12]
				]
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

			got, err := c.GetProductOrderBook(tt.args.productID, tt.args.qp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetProductOrderBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetProductOrderBook() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: figure file
func TestClient_ListTradesByProduct(t *testing.T) {
	type args struct {
		productID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Trade
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to list trades by product",
			fields: defaultFields(),
			args:   args{productID: "BTC-USD"},
			want: []Trade{{
				Time:    "2020-07-24T22:27:41.057Z",
				TradeID: 14990863,
				Price:   "9562.51000000",
				Size:    "0.00100000",
				Side:    "sell",
			}, {
				Time:    "2020-08-18T22:25:40.027Z",
				TradeID: 34923574,
				Price:   "7645.04837465",
				Size:    "0.00100000",
				Side:    "sell",
			}},
			wantRaw: `[{
				"time": "2020-07-24T22:27:41.057Z",
				"trade_id":14990863,
				"price":"9562.51000000",
				"size":"0.00100000",
				"side":"sell"
			},{
				"time": "2020-08-18T22:25:40.027Z",
				"trade_id":34923574,
				"price":"7645.04837465",
				"size":"0.00100000",
				"side":"sell"
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

			got, err := c.ListTradesByProduct(tt.args.productID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListTradesByProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListTradesByProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: figure file
func TestClient_GetHistoricRatesForProduct(t *testing.T) {
	type args struct {
		productID string
		qp        QueryParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []HistoricRate
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to get historic rates for a product",
			fields: defaultFields(),
			args:   args{productID: "BTC-USD", qp: QueryParams{Granularity: "60"}},
			want: []HistoricRate{{
				Time:   1415398768,
				Low:    0.32,
				High:   4.2,
				Open:   0.35,
				Close:  4.2,
				Volume: 12.3,
			}, {
				Time:   1298562378,
				Low:    0.24,
				High:   1.9,
				Open:   0.44,
				Close:  8.2,
				Volume: 15.8,
			}},
			wantRaw: `[
				[1415398768, 0.32, 4.2, 0.35, 4.2, 12.3],
				[1298562378, 0.24, 1.9, 0.44, 8.2, 15.8]
			]`,
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

			got, err := c.GetHistoricRatesForProduct(tt.args.productID, tt.args.qp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetHistoricRatesForProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetHistoricRatesForProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: figure file
func TestClient_Get24HourStatsForProduct(t *testing.T) {
	type args struct {
		productID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    DayStat
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to list trades by product",
			fields: defaultFields(),
			args:   args{productID: "BTC-USD"},
			want: DayStat{
				Open:        "6745.61000000",
				High:        "7292.11000000",
				Low:         "6650.00000000",
				Volume:      "26185.51325269",
				Last:        "6813.19000000",
				Volume30Day: "1019451.11188405",
			},
			wantRaw: `{
				"open": "6745.61000000",
				"high": "7292.11000000",
				"low": "6650.00000000",
				"volume": "26185.51325269",
				"last": "6813.19000000",
				"volume_30day": "1019451.11188405"
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

			got, err := c.Get24HourStatsForProduct(tt.args.productID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Get24HourStatsForProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Get24HourStatsForProduct() = %v, want %v", got, tt.want)
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
