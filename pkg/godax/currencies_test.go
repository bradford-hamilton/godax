package godax

import (
	"reflect"
	"testing"
)

func TestClient_ListCurrencies(t *testing.T) {
	tests := [...]struct {
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
