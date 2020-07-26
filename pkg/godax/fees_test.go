package godax

import (
	"reflect"
	"testing"
)

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
