package godax

import (
	"reflect"
	"testing"
)

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
