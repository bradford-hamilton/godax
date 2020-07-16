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

func TestClient_ListAccounts(t *testing.T) {
	tests := [...]struct {
		name    string
		fields  fields
		mock    HTTPClient
		want    []ListAccount
		wantRaw string
		wantErr bool
	}{
		{
			name:    "when a successful call is made to ListAccounts with no results",
			fields:  defaultFields(),
			want:    []ListAccount{},
			wantRaw: `[]`,
		},
		{
			name:   "when a successful call is made to ListAccounts with one account",
			fields: defaultFields(),
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
			fields: defaultFields(),
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
	type args struct {
		accountID string
	}
	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    Account
		wantRaw string
		wantErr bool
	}{
		{
			name:    "when a successful call is made to GetAccount and no account is found",
			fields:  defaultFields(),
			args:    args{accountID: "1q2w3e4r"},
			want:    Account{},
			wantRaw: `{}`,
		},
		{
			name:   "when a successful call is made to GetAccount and an account is found",
			fields: defaultFields(),
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

func TestClient_GetAccountHistory(t *testing.T) {
	type args struct {
		accountID string
	}
	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    []AccountActivity
		wantRaw string
		wantErr bool
	}{
		{
			name:    "when a successful call is made to GetAccountHistory and no history is found",
			fields:  defaultFields(),
			args:    args{accountID: "1q2w3e4r"},
			want:    []AccountActivity{},
			wantRaw: `[]`,
		},
		{
			name:   "when a successful call is made to GetAccountHistory and one history is found",
			fields: defaultFields(),
			args:   args{accountID: "a1b2c3d4"},
			want: []AccountActivity{{
				ID:        "100",
				CreatedAt: "2014-11-07T08:19:27.028459Z",
				Amount:    "0.001",
				Balance:   "239.669",
				Type:      "fee",
				Details: ActivityDetail{
					OrderID:   "d50ec984-77a8-460a-b958-66f114b0de9b",
					TradeID:   "74",
					ProductID: "BTC-USD",
				},
			}},
			wantRaw: `[{
                "id": "100",
                "created_at": "2014-11-07T08:19:27.028459Z",
                "amount": "0.001",
                "balance": "239.669",
                "type": "fee",
                "details": {
                    "order_id": "d50ec984-77a8-460a-b958-66f114b0de9b",
                    "trade_id": "74",
                    "product_id": "BTC-USD"
                }
            }]`,
		},
		{
			name:   "when a successful call is made to GetAccountHistory and multiple histories are found",
			fields: defaultFields(),
			args:   args{accountID: "a1b2c3d4"},
			want: []AccountActivity{{
				ID:        "100",
				CreatedAt: "2014-11-07T08:19:27.028459Z",
				Amount:    "0.001",
				Balance:   "239.669",
				Type:      "fee",
				Details: ActivityDetail{
					OrderID:   "d50ec984-77a8-460a-b958-66f114b0de9b",
					TradeID:   "74",
					ProductID: "BTC-USD",
				},
			}, {
				ID:        "80",
				CreatedAt: "2015-12-04T08:19:27.028459Z",
				Amount:    "0.011",
				Balance:   "4059.212345",
				Type:      "fee",
				Details: ActivityDetail{
					OrderID:   "8b9258f8-811b-429b-810d-71fede464b29",
					TradeID:   "99",
					ProductID: "BTC-ETH",
				},
			}},
			wantRaw: `[{
                "id": "100",
                "created_at": "2014-11-07T08:19:27.028459Z",
                "amount": "0.001",
                "balance": "239.669",
                "type": "fee",
                "details": {
                    "order_id": "d50ec984-77a8-460a-b958-66f114b0de9b",
                    "trade_id": "74",
                    "product_id": "BTC-USD"
                }
            },{
                "id": "80",
                "created_at": "2015-12-04T08:19:27.028459Z",
                "amount": "0.011",
                "balance": "4059.212345",
                "type": "fee",
                "details": {
                    "order_id": "8b9258f8-811b-429b-810d-71fede464b29",
                    "trade_id": "99",
                    "product_id": "BTC-ETH"
                }
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

			got, err := c.GetAccountHistory(tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetAccountHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetAccountHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetAccountHolds(t *testing.T) {
	type args struct {
		accountID string
	}
	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    []AccountHold
		wantRaw string
		wantErr bool
	}{
		{
			name:    "when a successful call is made to GetAccountHolds and no holds are found",
			fields:  defaultFields(),
			args:    args{accountID: "1q2w3e4r"},
			want:    []AccountHold{},
			wantRaw: `[]`,
		},
		{
			name:   "when a successful call is made to GetAccountHolds and one hold is found",
			fields: defaultFields(),
			args:   args{accountID: "1q2w3e4r"},
			want: []AccountHold{{
				ID:        "82dcd140-c3c7-4507-8de4-2c529cd1a28f",
				AccountID: "e0b3f39a-183d-453e-b754-0c13e5bab0b3",
				CreatedAt: "2014-11-06T10:34:47.123456Z",
				UpdatedAt: "2014-11-06T10:40:47.123456Z",
				Amount:    "4.23",
				Type:      "order",
				Ref:       "0a205de4-dd35-4370-a285-fe8fc375a273",
			}},
			wantRaw: `[{
                "id": "82dcd140-c3c7-4507-8de4-2c529cd1a28f",
                "account_id": "e0b3f39a-183d-453e-b754-0c13e5bab0b3",
                "created_at": "2014-11-06T10:34:47.123456Z",
                "updated_at": "2014-11-06T10:40:47.123456Z",
                "amount": "4.23",
                "type": "order",
                "ref": "0a205de4-dd35-4370-a285-fe8fc375a273"
            }]`,
		},
		{
			name:   "when a successful call is made to GetAccountHolds and multiple holds are found",
			fields: defaultFields(),
			args:   args{accountID: "1q2w3e4r"},
			want: []AccountHold{{
				ID:        "82dcd140-c3c7-4507-8de4-2c529cd1a28f",
				AccountID: "e0b3f39a-183d-453e-b754-0c13e5bab0b3",
				CreatedAt: "2014-11-06T10:34:47.123456Z",
				UpdatedAt: "2014-11-06T10:40:47.123456Z",
				Amount:    "4.23",
				Type:      "order",
				Ref:       "0a205de4-dd35-4370-a285-fe8fc375a273",
			}, {
				ID:        "3d58f10b-3d9a-4d38-bb51-c8800f5ad4ca",
				AccountID: "b6f8fee0-f47f-481a-98ee-08d397681edb",
				CreatedAt: "2015-10-06T10:34:47.123456Z",
				UpdatedAt: "2015-10-06T10:40:47.123456Z",
				Amount:    "4.23",
				Type:      "order",
				Ref:       "0a205de4-dd35-4370-a285-fe8fc375a273",
			}},
			wantRaw: `[{
                "id": "82dcd140-c3c7-4507-8de4-2c529cd1a28f",
                "account_id": "e0b3f39a-183d-453e-b754-0c13e5bab0b3",
                "created_at": "2014-11-06T10:34:47.123456Z",
                "updated_at": "2014-11-06T10:40:47.123456Z",
                "amount": "4.23",
                "type": "order",
                "ref": "0a205de4-dd35-4370-a285-fe8fc375a273"
            },
            {
                "id": "3d58f10b-3d9a-4d38-bb51-c8800f5ad4ca",
                "account_id": "b6f8fee0-f47f-481a-98ee-08d397681edb",
                "created_at": "2015-10-06T10:34:47.123456Z",
                "updated_at": "2015-10-06T10:40:47.123456Z",
                "amount": "4.23",
                "type": "order",
                "ref": "0a205de4-dd35-4370-a285-fe8fc375a273"
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

			got, err := c.GetAccountHolds(tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetAccountHolds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetAccountHolds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_PlaceOrder(t *testing.T) {
	type args struct {
		order OrderParams
	}
	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    Order
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful market order is made with PlaceOrder",
			fields: defaultFields(),
			args: args{order: OrderParams{
				CommonOrderParams: CommonOrderParams{
					Side:      "buy",
					ProductID: "ETH-BTC",
					ClientOID: "0f5d8030-8908-4ef5-b3e8-7131bbe35588",
					Type:      "market",
					Size:      "0.15",
				},
			}},
			want: Order{
				ID:            "918e9893-ecb1-4a51-8b7b-3f29f9baff9f",
				CreatedAt:     "2020-07-02T19:03:55.864229Z",
				FillFees:      "0",
				FilledSize:    "0",
				ExecutedValue: "0",
				Status:        "pending",
				Settled:       false,
				OrderParams: OrderParams{
					CommonOrderParams: CommonOrderParams{
						Side:      "buy",
						ProductID: "ETH-BTC",
						Type:      "market",
						Size:      "0.15",
						Stp:       "cn",
					},
					LimitOrderParams: LimitOrderParams{
						PostOnly: false,
					},
					MarketOrderParams: MarketOrderParams{
						Funds: "49.20762366",
					},
				},
			},
			wantRaw: `{
				"id": "918e9893-ecb1-4a51-8b7b-3f29f9baff9f",
				"size": "0.15",
				"product_id": "ETH-BTC",
				"side": "buy",
				"stp": "cn",
				"funds": "49.20762366",
				"type": "market",
				"post_only": false,
				"created_at": "2020-07-02T19:03:55.864229Z",
				"fill_fees": "0",
				"filled_size": "0",
				"executed_value": "0",
				"status": "pending",
				"settled": false
			}`,
		},
		{
			name:   "when a successful limit order is made with PlaceOrder",
			fields: defaultFields(),
			args: args{order: OrderParams{
				CommonOrderParams: CommonOrderParams{
					Side:      "buy",
					ProductID: "ETH-BTC",
					Type:      "limit",
					Size:      "0.15",
					Price:     "1",
				},
			}},
			want: Order{
				ID:            "5033b47c-bb58-4089-b062-b830fc93c207",
				CreatedAt:     "2020-07-02T20:57:22.077332Z",
				FillFees:      "0",
				FilledSize:    "0",
				ExecutedValue: "0",
				Status:        "pending",
				Settled:       false,
				OrderParams: OrderParams{
					CommonOrderParams: CommonOrderParams{
						Side:      "buy",
						ProductID: "ETH-BTC",
						Type:      "limit",
						Price:     "1",
						Size:      "0.15",
						Stp:       "cn",
					},
					LimitOrderParams: LimitOrderParams{
						PostOnly:    false,
						TimeInForce: "GTC",
					},
				},
			},
			wantRaw: `{
				"id": "5033b47c-bb58-4089-b062-b830fc93c207",
				"price": "1",
				"size": "0.15",
				"product_id": "ETH-BTC",
				"side": "buy",
				"stp": "cn",
				"type": "limit",
				"time_in_force": "GTC",
				"post_only": false,
				"created_at": "2020-07-02T20:57:22.077332Z",
				"fill_fees": "0",
				"filled_size": "0",
				"executed_value": "0",
				"status": "pending",
				"settled":false
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

			got, err := c.PlaceOrder(tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.PlaceOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.PlaceOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_CancelOrderByID(t *testing.T) {
	productID := "BTC-USD"

	type args struct {
		orderID string
		qp      QueryParams
	}

	tests := [...]struct {
		name                string
		fields              fields
		args                args
		wantCanceledOrderID string
		wantRaw             string
		wantErr             bool
	}{
		{
			name:                "when a successful cancel order has been made with no product ID",
			fields:              defaultFields(),
			args:                args{orderID: "c6dfb02e-7f65-4e02-8fa3-866d46ed15b3", qp: QueryParams{}},
			wantCanceledOrderID: "c6dfb02e-7f65-4e02-8fa3-866d46ed15b3",
			wantRaw:             "c6dfb02e-7f65-4e02-8fa3-866d46ed15b3",
		},
		{
			name:                "when a successful cancel order has been made with a product ID",
			fields:              defaultFields(),
			args:                args{orderID: "4f92c553-7c71-4b3a-8878-f415e6a2f0d8", qp: QueryParams{ProductID: productID}},
			wantCanceledOrderID: "4f92c553-7c71-4b3a-8878-f415e6a2f0d8",
			wantRaw:             "4f92c553-7c71-4b3a-8878-f415e6a2f0d8",
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

			gotCanceledOrderID, err := c.CancelOrderByID(tt.args.orderID, tt.args.qp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CancelOrderByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if gotCanceledOrderID != tt.wantCanceledOrderID {
				t.Errorf("Client.CancelOrderByID() = %v, want %v", gotCanceledOrderID, tt.wantCanceledOrderID)
			}
		})
	}
}

func TestClient_CancelOrderByClientOID(t *testing.T) {
	productID := "BTC-USD"

	type args struct {
		clientOID string
		qp        QueryParams
	}

	tests := [...]struct {
		name                string
		fields              fields
		args                args
		wantCanceledOrderID string
		wantRaw             string
		wantErr             bool
	}{
		{
			name:                "when a successful cancel order has been made with no product ID",
			fields:              defaultFields(),
			args:                args{clientOID: "408290f2-f13e-465d-a2ff-98a29d130bd4", qp: QueryParams{}},
			wantCanceledOrderID: "408290f2-f13e-465d-a2ff-98a29d130bd4",
			wantRaw:             "408290f2-f13e-465d-a2ff-98a29d130bd4",
		},
		{
			name:                "when a successful cancel order has been made with a product ID",
			fields:              defaultFields(),
			args:                args{clientOID: "52e06257-dc1f-4e82-b115-c81f5f07a9d8", qp: QueryParams{ProductID: productID}},
			wantCanceledOrderID: "52e06257-dc1f-4e82-b115-c81f5f07a9d8",
			wantRaw:             "52e06257-dc1f-4e82-b115-c81f5f07a9d8",
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

			gotCanceledOrderID, err := c.CancelOrderByClientOID(tt.args.clientOID, tt.args.qp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CancelOrderByClientOID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if gotCanceledOrderID != tt.wantCanceledOrderID {
				t.Errorf("Client.CancelOrderByClientOID() = %v, want %v", gotCanceledOrderID, tt.wantCanceledOrderID)
			}
		})
	}
}

func TestClient_CancelAllOrders(t *testing.T) {
	productID := "BTC-USD"

	type args struct{ qp QueryParams }

	tests := [...]struct {
		name                 string
		fields               fields
		args                 args
		wantCanceledOrderIDs []string
		wantRaw              string
		wantErr              bool
	}{
		{
			name:   "when a successful cancel all orders call has been made with no product ID",
			fields: defaultFields(),
			args:   args{QueryParams{}},
			wantCanceledOrderIDs: []string{
				"920dfecf-2dde-491d-9dd1-ca9f335a0663",
				"0189696d-7b3e-4e9d-aa27-5df9f620466f",
				"91315780-ced8-43a7-856d-52d5daf9d574",
			},
			wantRaw: `[
				"920dfecf-2dde-491d-9dd1-ca9f335a0663",
				"0189696d-7b3e-4e9d-aa27-5df9f620466f",
				"91315780-ced8-43a7-856d-52d5daf9d574"
			]`,
		},
		{
			name:   "when a successful cancel all orders call has been made with a product ID",
			fields: defaultFields(),
			args:   args{QueryParams{ProductID: productID}},
			wantCanceledOrderIDs: []string{
				"db5b5cb9-3a86-4d44-b62d-c2c2c39d1446",
				"7fc4a8bf-e22b-46f2-a881-5bedd2bc1571",
				"85495e53-1e07-4ae0-b6e6-205bbf1b0552",
			},
			wantRaw: `[
				"db5b5cb9-3a86-4d44-b62d-c2c2c39d1446",
				"7fc4a8bf-e22b-46f2-a881-5bedd2bc1571",
				"85495e53-1e07-4ae0-b6e6-205bbf1b0552"
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

			gotCanceledOrderIDs, err := c.CancelAllOrders(tt.args.qp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CancelAllOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(gotCanceledOrderIDs, tt.wantCanceledOrderIDs) {
				t.Errorf("Client.CancelAllOrders() = %v, want %v", gotCanceledOrderIDs, tt.wantCanceledOrderIDs)
			}
		})
	}
}

func TestClient_ListOrders(t *testing.T) {
	productID := "BTC-USD"
	status := "pending"

	type args struct{ qp QueryParams }

	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    []Order
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful list orders call has been made with no product ID",
			fields: defaultFields(),
			args:   args{QueryParams{Status: "", ProductID: ""}},
			want: []Order{{
				ID:            "b5bebfc6-4ee1-463c-8a57-7aa98eeb3f7e",
				CreatedAt:     "2020-07-02T19:03:55.864229Z",
				FillFees:      "0",
				FilledSize:    "0",
				ExecutedValue: "0",
				Status:        "pending",
				Settled:       false,
				OrderParams: OrderParams{
					CommonOrderParams: CommonOrderParams{
						Side:      "buy",
						ProductID: "ETH-BTC",
						Type:      "market",
						Size:      "0.15",
						Stp:       "cn",
					},
					LimitOrderParams: LimitOrderParams{
						PostOnly: false,
					},
					MarketOrderParams: MarketOrderParams{
						Funds: "49.20762366",
					},
				},
			}},
			wantRaw: `[{
				"id": "b5bebfc6-4ee1-463c-8a57-7aa98eeb3f7e",
				"size": "0.15",
				"product_id": "ETH-BTC",
				"side": "buy",
				"stp": "cn",
				"funds": "49.20762366",
				"type": "market",
				"post_only": false,
				"created_at": "2020-07-02T19:03:55.864229Z",
				"fill_fees": "0",
				"filled_size": "0",
				"executed_value": "0",
				"status": "pending",
				"settled": false
			}]`,
		},
		{
			name:   "when a successful list orders call has been made with a product ID",
			fields: defaultFields(),
			args:   args{QueryParams{ProductID: productID}},
			want: []Order{{
				ID:            "5b598eca-27cd-4ffc-bc17-8d5eaae67541",
				CreatedAt:     "2020-07-02T19:03:55.864229Z",
				FillFees:      "0",
				FilledSize:    "0",
				ExecutedValue: "0",
				Status:        "pending",
				Settled:       false,
				OrderParams: OrderParams{
					CommonOrderParams: CommonOrderParams{
						Side:      "buy",
						ProductID: "BTC-LTC",
						Type:      "market",
						Size:      "0.15",
						Stp:       "cn",
					},
					LimitOrderParams: LimitOrderParams{
						PostOnly: false,
					},
					MarketOrderParams: MarketOrderParams{
						Funds: "49.20762366",
					},
				},
			}, {
				ID:            "fe82e31d-105b-4a1d-847c-d43e280d5966",
				CreatedAt:     "2020-07-02T19:03:55.864229Z",
				FillFees:      "0",
				FilledSize:    "0",
				ExecutedValue: "0",
				Status:        "pending",
				Settled:       false,
				OrderParams: OrderParams{
					CommonOrderParams: CommonOrderParams{
						Side:      "buy",
						ProductID: "BTC-LTC",
						Type:      "market",
						Size:      "0.15",
						Stp:       "cn",
					},
					LimitOrderParams: LimitOrderParams{
						PostOnly: false,
					},
					MarketOrderParams: MarketOrderParams{
						Funds: "49.20762366",
					},
				},
			}},
			wantRaw: `[{
				"id": "5b598eca-27cd-4ffc-bc17-8d5eaae67541",
				"size": "0.15",
				"product_id": "BTC-LTC",
				"side": "buy",
				"stp": "cn",
				"funds": "49.20762366",
				"type": "market",
				"post_only": false,
				"created_at": "2020-07-02T19:03:55.864229Z",
				"fill_fees": "0",
				"filled_size": "0",
				"executed_value": "0",
				"status": "pending",
				"settled": false
			},{
				"id": "fe82e31d-105b-4a1d-847c-d43e280d5966",
				"size": "0.15",
				"product_id": "BTC-LTC",
				"side": "buy",
				"stp": "cn",
				"funds": "49.20762366",
				"type": "market",
				"post_only": false,
				"created_at": "2020-07-02T19:03:55.864229Z",
				"fill_fees": "0",
				"filled_size": "0",
				"executed_value": "0",
				"status": "pending",
				"settled": false
			}]`,
		},
		{
			name:   "when a successful list orders call has been made with a status filter",
			fields: defaultFields(),
			args:   args{QueryParams{Status: status}},
			want: []Order{{
				ID:            "2a795d2e-f77a-4ac6-8998-767fe36fbaed",
				CreatedAt:     "2020-07-02T19:03:55.864229Z",
				FillFees:      "0",
				FilledSize:    "0",
				ExecutedValue: "0",
				Status:        "pending",
				Settled:       false,
				OrderParams: OrderParams{
					CommonOrderParams: CommonOrderParams{
						Side:      "buy",
						ProductID: "BTC-XTZ",
						Type:      "market",
						Size:      "0.15",
						Stp:       "cn",
					},
					LimitOrderParams: LimitOrderParams{
						PostOnly: false,
					},
					MarketOrderParams: MarketOrderParams{
						Funds: "49.20762366",
					},
				},
			}, {
				ID:            "5e8333a9-c86a-4877-b0ac-9b35769f74af",
				CreatedAt:     "2020-07-02T19:03:55.864229Z",
				FillFees:      "0",
				FilledSize:    "0",
				ExecutedValue: "0",
				Status:        "pending",
				Settled:       false,
				OrderParams: OrderParams{
					CommonOrderParams: CommonOrderParams{
						Side:      "buy",
						ProductID: "BTC-XTZ",
						Type:      "market",
						Size:      "0.15",
						Stp:       "cn",
					},
					LimitOrderParams: LimitOrderParams{
						PostOnly: false,
					},
					MarketOrderParams: MarketOrderParams{
						Funds: "49.20762366",
					},
				},
			}},
			wantRaw: `[{
				"id": "2a795d2e-f77a-4ac6-8998-767fe36fbaed",
				"size": "0.15",
				"product_id": "BTC-XTZ",
				"side": "buy",
				"stp": "cn",
				"funds": "49.20762366",
				"type": "market",
				"post_only": false,
				"created_at": "2020-07-02T19:03:55.864229Z",
				"fill_fees": "0",
				"filled_size": "0",
				"executed_value": "0",
				"status": "pending",
				"settled": false
			},{
				"id": "5e8333a9-c86a-4877-b0ac-9b35769f74af",
				"size": "0.15",
				"product_id": "BTC-XTZ",
				"side": "buy",
				"stp": "cn",
				"funds": "49.20762366",
				"type": "market",
				"post_only": false,
				"created_at": "2020-07-02T19:03:55.864229Z",
				"fill_fees": "0",
				"filled_size": "0",
				"executed_value": "0",
				"status": "pending",
				"settled": false
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

			got, err := c.ListOrders(tt.args.qp)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListOrders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListOrders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetOrderByID(t *testing.T) {
	type args struct{ orderID string }

	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    Order
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to GetOrderByID an Order is returned",
			fields: defaultFields(),
			args:   args{orderID: "68e6a28f-ae28-4788-8d4f-5ab4e5e5ae08"},
			want: Order{
				ID:            "68e6a28f-ae28-4788-8d4f-5ab4e5e5ae08",
				CreatedAt:     "2016-12-08T20:09:05.508883Z",
				FillFees:      "0.0249376391550000",
				FilledSize:    "0.01291771",
				ExecutedValue: "9.9750556620000000",
				Status:        "done",
				Settled:       true,
				OrderParams: OrderParams{
					CommonOrderParams: CommonOrderParams{
						Side:      "buy",
						ProductID: "BTC-USD",
						Type:      "market",
						Size:      "1.00000000",
						Stp:       "dc",
					},
					LimitOrderParams: LimitOrderParams{
						PostOnly: false,
					},
					MarketOrderParams: MarketOrderParams{
						Funds: "9.9750623400000000",
					},
				},
			},
			wantRaw: `{
				"id": "68e6a28f-ae28-4788-8d4f-5ab4e5e5ae08",
				"size": "1.00000000",
				"product_id": "BTC-USD",
				"side": "buy",
				"stp": "dc",
				"funds": "9.9750623400000000",
				"specified_funds": "10.0000000000000000",
				"type": "market",
				"post_only": false,
				"created_at": "2016-12-08T20:09:05.508883Z",
				"done_at": "2016-12-08T20:09:05.527Z",
				"done_reason": "filled",
				"fill_fees": "0.0249376391550000",
				"filled_size": "0.01291771",
				"executed_value": "9.9750556620000000",
				"status": "done",
				"settled": true
			}`,
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

			got, err := c.GetOrderByID(tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetOrderByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetOrderByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetOrderByClientOID(t *testing.T) {
	type args struct{ orderClientOID string }

	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    Order
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to GetOrderByID an Order is returned",
			fields: defaultFields(),
			args:   args{orderClientOID: "3aaa5f6c-9ca2-4e5e-bca5-cebff877a170"},
			want: Order{
				ID:            "3aaa5f6c-9ca2-4e5e-bca5-cebff877a170",
				CreatedAt:     "2016-12-08T20:09:05.508883Z",
				FillFees:      "0.0249376391550000",
				FilledSize:    "0.01291771",
				ExecutedValue: "9.9750556620000000",
				Status:        "done",
				Settled:       true,
				OrderParams: OrderParams{
					CommonOrderParams: CommonOrderParams{
						Side:      "buy",
						ProductID: "BTC-USD",
						Type:      "market",
						Size:      "1.00000000",
						Stp:       "dc",
					},
					LimitOrderParams: LimitOrderParams{
						PostOnly: false,
					},
					MarketOrderParams: MarketOrderParams{
						Funds: "9.9750623400000000",
					},
				},
			},
			wantRaw: `{
				"id": "3aaa5f6c-9ca2-4e5e-bca5-cebff877a170",
				"size": "1.00000000",
				"product_id": "BTC-USD",
				"side": "buy",
				"stp": "dc",
				"funds": "9.9750623400000000",
				"specified_funds": "10.0000000000000000",
				"type": "market",
				"post_only": false,
				"created_at": "2016-12-08T20:09:05.508883Z",
				"done_at": "2016-12-08T20:09:05.527Z",
				"done_reason": "filled",
				"fill_fees": "0.0249376391550000",
				"filled_size": "0.01291771",
				"executed_value": "9.9750556620000000",
				"status": "done",
				"settled": true
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

			got, err := c.GetOrderByClientOID(tt.args.orderClientOID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetOrderByClientOID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetOrderByClientOID() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestClient_StableCoinConversion(t *testing.T) {
	type args struct {
		from   string
		to     string
		amount string
	}
	tests := [...]struct {
		name    string
		fields  fields
		args    args
		want    Conversion
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to convert to/from a stablecoin",
			fields: defaultFields(),
			args:   args{to: "USD", from: "USDC", amount: "10"},
			want: Conversion{
				ID:            "8942caee-f9d5-4600-a894-4811268545db",
				Amount:        "10000.00",
				FromAccountID: "7849cc79-8b01-4793-9345-bc6b5f08acce",
				ToAccountID:   "105c3e58-0898-4106-8283-dc5781cda07b",
				From:          "USD",
				To:            "USDC",
			},
			wantErr: false,
			wantRaw: `{
				"id": "8942caee-f9d5-4600-a894-4811268545db",
				"amount": "10000.00",
				"from_account_id": "7849cc79-8b01-4793-9345-bc6b5f08acce",
				"to_account_id": "105c3e58-0898-4106-8283-dc5781cda07b",
				"from": "USD",
				"to": "USDC"
			}`,
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
			got, err := c.StableCoinConversion(tt.args.from, tt.args.to, tt.args.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.StableCoinConversion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.StableCoinConversion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ListPaymentMethods(t *testing.T) {
	tests := [...]struct {
		name    string
		fields  fields
		want    []PaymentMethod
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to list payment methods",
			fields: defaultFields(),
			want: []PaymentMethod{{
				ID:            "bc6d7162-d984-5ffa-963c-a493b1c1370b",
				Type:          "ach_bank_account",
				Name:          "Bank of America - eBan... ********7134",
				Currency:      "USD",
				PrimaryBuy:    true,
				PrimarySell:   true,
				AllowBuy:      true,
				AllowSell:     true,
				AllowDeposit:  true,
				AllowWithdraw: true,
				Limits: PMLimits{
					Buy: []PMLimit{{
						PeriodInDays: 1,
						Total: PMAmount{
							Amount:   "10000.00",
							Currency: "USD",
						},
						Remaining: PMAmount{
							Amount:   "10000.00",
							Currency: "USD",
						},
					}},
					Sell: []PMLimit{{
						PeriodInDays: 1,
						Total: PMAmount{
							Amount:   "10000.00",
							Currency: "USD",
						},
						Remaining: PMAmount{
							Amount:   "10000.00",
							Currency: "USD",
						},
					}},
				},
			}},
			wantErr: false,
			wantRaw: `[{
				"id": "bc6d7162-d984-5ffa-963c-a493b1c1370b",
				"type": "ach_bank_account",
				"name": "Bank of America - eBan... ********7134",
				"currency": "USD",
				"primary_buy": true,
				"primary_sell": true,
				"allow_buy": true,
				"allow_sell": true,
				"allow_deposit": true,
				"allow_withdraw": true,
				"limits": {
					"buy": [{
						"period_in_days": 1,
						"total": {
							"amount": "10000.00",
							"currency": "USD"
						},
						"remaining": {
							"amount": "10000.00",
							"currency": "USD"
						}
					}],
					"sell": [{
						"period_in_days": 1,
						"total": {
							"amount": "10000.00",
							"currency": "USD"
						},
						"remaining": {
							"amount": "10000.00",
							"currency": "USD"
						}
					}]
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

			got, err := c.ListPaymentMethods()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListPaymentMethods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListPaymentMethods() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func TestClient_GetCurrentFees(t *testing.T) {
	tests := [...]struct {
		name    string
		fields  fields
		want    Fees
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to get current fees",
			fields: defaultFields(),
			want: Fees{
				MakerFeeRate: "0.0015",
				TakerFeeRate: "0.0025",
				USDVolume:    "25000.00",
			},
			wantRaw: `{
				"maker_fee_rate": "0.0015",
				"taker_fee_rate": "0.0025",
				"usd_volume": "25000.00"
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

			got, err := c.GetCurrentFees()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetCurrentFees() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetCurrentFees() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetTrailingVolume(t *testing.T) {
	tests := []struct {
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
