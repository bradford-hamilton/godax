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
