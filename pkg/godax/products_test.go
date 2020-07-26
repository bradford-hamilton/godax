package godax

import (
	"reflect"
	"testing"
)

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
