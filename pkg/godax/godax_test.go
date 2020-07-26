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

// FIXME: orders_test.go
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

// FIXME: orders_test.go
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

// FIXME: orders_test.go
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

// FIXME: orders_test.go
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

// FIXME: orders_test.go
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

// FIXME: orders_test.go
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

// FIXME: orders_test.go
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

// FIXME: conversions_test.go
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

// FIXME: payment_methods_test.go
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

// FIXME: fees_test.go
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

// FIXME: profiles_test.go
func TestClient_ListProfiles(t *testing.T) {
	tests := [...]struct {
		name    string
		fields  fields
		want    []Profile
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to list profiles",
			fields: defaultFields(),
			want: []Profile{{
				ID:        "86602c68-306a-4500-ac73-4ce56a91d83c",
				UserID:    "5844eceecf7e803e259d0365",
				Name:      "default",
				Active:    true,
				IsDefault: true,
				CreatedAt: "2019-11-18T15:08:40.236309Z",
			}},
			wantRaw: `[{
				"id": "86602c68-306a-4500-ac73-4ce56a91d83c",
				"user_id": "5844eceecf7e803e259d0365",
				"name": "default",
				"active": true,
				"is_default": true,
				"created_at": "2019-11-18T15:08:40.236309Z"
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

			got, err := c.ListProfiles()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListProfiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListProfiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

// FIXME: profiles_test.go
func TestClient_GetProfile(t *testing.T) {
	type args struct {
		profileID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Profile
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to get a profile by ID",
			fields: defaultFields(),
			args:   args{profileID: "123abc"},
			want: Profile{
				ID:        "86602c68-306a-4500-ac73-4ce56a91d83c",
				UserID:    "5844eceecf7e803e259d0365",
				Name:      "default",
				Active:    true,
				IsDefault: true,
				CreatedAt: "2019-11-18T15:08:40.236309Z",
			},
			wantRaw: `{
				"id": "86602c68-306a-4500-ac73-4ce56a91d83c",
				"user_id": "5844eceecf7e803e259d0365",
				"name": "default",
				"active": true,
				"is_default": true,
				"created_at": "2019-11-18T15:08:40.236309Z"
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

			got, err := c.GetProfile(tt.args.profileID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: profiles_test.go
func TestClient_ProfileTransfer(t *testing.T) {
	type args struct {
		transfer TransferParams
	}
	tests := [...]struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		wantRaw string
	}{
		{
			name:   "when a successful market order is made to transfer currency across profiles",
			fields: defaultFields(),
			args: args{transfer: TransferParams{
				From:     "86602c68-306a-4500-ac73-4ce56a91d83c",
				To:       "e87429d3-f0a7-4f28-8dff-8dd93d383de1",
				Currency: "BTC-USD",
				Amount:   "0.05",
			}},
			wantErr: false,
			wantRaw: "",
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

			err := c.ProfileTransfer(tt.args.transfer)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ProfileTransfer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)
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

// FIXME: currencies_test.go
func TestClient_ListCurrencies(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    []Currency
		wantRaw string
		wantErr bool
	}{
		{
			name:   "when a successful call is made to list currencies",
			fields: defaultFields(),
			want: []Currency{{
				ID:      "BTC",
				Name:    "Bitcoin",
				MinSize: "0.00000001",
			}, {
				ID:      "USD",
				Name:    "United States Dollar",
				MinSize: "0.01000000",
			}},
			wantRaw: `[{
				"id": "BTC",
				"name": "Bitcoin",
				"min_size": "0.00000001"
			}, {
				"id": "USD",
				"name": "United States Dollar",
				"min_size": "0.01000000"
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

			got, err := c.ListCurrencies()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.ListCurrencies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(c.httpClient.(*MockClient).Requests) != 1 {
				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
			}

			validateHeaders(t, c)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.ListCurrencies() = %v, want %v", got, tt.want)
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
