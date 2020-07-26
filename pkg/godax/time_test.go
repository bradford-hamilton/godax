package godax

import (
	"reflect"
	"testing"
)

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
