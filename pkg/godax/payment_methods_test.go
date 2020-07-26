package godax

import (
	"reflect"
	"testing"
)

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
