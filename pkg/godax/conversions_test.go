package godax

import (
	"reflect"
	"testing"
)

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
