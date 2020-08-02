package godax

// TODO: uncomment these once all margin methods are ready.

// func TestClient_GetMarginProfile(t *testing.T) {
// 	type args struct {
// 		qp QueryParams
// 	}
// 	tests := [...]struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    MarginProfile
// 		wantRaw string
// 		wantErr bool
// 	}{
// 		{
// 			name:   "when a successful call is made to get my margin profile",
// 			fields: defaultFields(),
// 			args:   args{QueryParams{ProductIDParam: "BTC-USD"}},
// 			want: MarginProfile{
// 				ProfileID:            "8058d771-2d88-4f0f-ab6e-299c153d4308",
// 				MarginInitialEquity:  "0.33",
// 				MarginWarningEquity:  "0.2",
// 				MarginCallEquity:     "0.15",
// 				EquityPercentage:     0.8725562096924747,
// 				SellingPower:         0.00221896,
// 				BuyingPower:          23.51,
// 				BorrowPower:          23.51,
// 				InterestRate:         "0",
// 				InterestPaid:         "0.3205913399694425",
// 				CollateralCurrencies: []string{"BTC", "USD", "USDC"},
// 				CollateralHoldValue:  "1.0050000000000000",
// 				LastLiquidationAt:    "2019-11-21T14:58:49.879Z",
// 				AvailableBorrowLimits: struct {
// 					MarginableLimit    float64 "json:\"marginable_limit\""
// 					NonMarginableLimit float64 "json:\"nonmarginable_limit\""
// 				}{
// 					MarginableLimit:    23.51,
// 					NonMarginableLimit: 7.75,
// 				},
// 				BorrowLimit: "5000",
// 				TopUpAmounts: struct {
// 					BorrowableUsd    string "json:\"borrowable_usd\""
// 					NonBorrowableUsd string "json:\"non_borrowable_usd\""
// 				}{
// 					BorrowableUsd:    "0",
// 					NonBorrowableUsd: "0",
// 				},
// 			},
// 			wantRaw: `{
//                 "profile_id": "8058d771-2d88-4f0f-ab6e-299c153d4308",
//                 "margin_initial_equity": "0.33",
//                 "margin_warning_equity": "0.2",
//                 "margin_call_equity": "0.15",
//                 "equity_percentage": 0.8725562096924747,
//                 "selling_power": 0.00221896,
//                 "buying_power": 23.51,
//                 "borrow_power": 23.51,
//                 "interest_rate": "0",
//                 "interest_paid": "0.3205913399694425",
//                 "collateral_currencies": [
//                     "BTC",
//                     "USD",
//                     "USDC"
//                 ],
//                 "collateral_hold_value": "1.0050000000000000",
//                 "last_liquidation_at": "2019-11-21T14:58:49.879Z",
//                 "available_borrow_limits": {
//                     "marginable_limit": 23.51,
//                     "nonmarginable_limit": 7.75
//                 },
//                 "borrow_limit": "5000",
//                 "top_up_amounts": {
//                     "borrowable_usd": "0",
//                     "non_borrowable_usd": "0"
//                 }
//             }`,
// 		},
// 	}
// 	for _, tt := range tests {
// 		mockClient := MockResponse(tt.wantRaw)

// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &Client{
// 				baseRestURL: tt.fields.baseRestURL,
// 				baseWsURL:   tt.fields.baseWsURL,
// 				key:         tt.fields.key,
// 				secret:      tt.fields.secret,
// 				passphrase:  tt.fields.passphrase,
// 				httpClient:  mockClient,
// 			}

// 			got, err := c.GetMarginProfile(tt.args.qp)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Client.GetMarginProfile() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if len(c.httpClient.(*MockClient).Requests) != 1 {
// 				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
// 			}

// 			validateHeaders(t, c)

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Client.GetMarginProfile() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestClient_GetBuyingPower(t *testing.T) {
// 	type args struct {
// 		qp QueryParams
// 	}
// 	tests := [...]struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    BuyingPower
// 		wantRaw string
// 		wantErr bool
// 	}{
// 		{
// 			name:   "when a successful call is made to get buying power for a product",
// 			fields: defaultFields(),
// 			args:   args{QueryParams{ProductIDParam: "BTC-USD"}},
// 			want: BuyingPower{
// 				BuyingPower:            23.53,
// 				SellingPower:           0.00221896,
// 				BuyingPowerExplanation: "This is the line of credit available to you on the BTC-USD market, given how much collateral assets you currently have in your portfolio.",
// 			},
// 			wantRaw: `{
//                 "buying_power": 23.53,
//                 "selling_power": 0.00221896,
//                 "buying_power_explanation": "This is the line of credit available to you on the BTC-USD market, given how much collateral assets you currently have in your portfolio."
//             }`,
// 		},
// 	}
// 	for _, tt := range tests {
// 		mockClient := MockResponse(tt.wantRaw)

// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &Client{
// 				baseRestURL: tt.fields.baseRestURL,
// 				baseWsURL:   tt.fields.baseWsURL,
// 				key:         tt.fields.key,
// 				secret:      tt.fields.secret,
// 				passphrase:  tt.fields.passphrase,
// 				httpClient:  mockClient,
// 			}

// 			got, err := c.GetBuyingPower(tt.args.qp)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Client.GetBuyingPower() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if len(c.httpClient.(*MockClient).Requests) != 1 {
// 				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
// 			}

// 			validateHeaders(t, c)

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Client.GetBuyingPower() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestClient_GetWithdrawalPowerForCurrency(t *testing.T) {
// 	type args struct {
// 		qp QueryParams
// 	}
// 	tests := [...]struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    []CurrencyWithdrawalPower
// 		wantRaw string
// 		wantErr bool
// 	}{
// 		{
// 			name:   "when a successful call is made to get buying power for a product",
// 			fields: defaultFields(),
// 			args:   args{QueryParams{CurrencyParam: "BTC"}},
// 			want: []CurrencyWithdrawalPower{{
// 				ProfileID:       "8058d771-2d88-4f0f-ab6e-299c153d4308",
// 				WithdrawalPower: "7.77569088416849750000",
// 			}},
// 			wantRaw: `[{
//                 "profile_id": "8058d771-2d88-4f0f-ab6e-299c153d4308",
//                 "withdrawal_power": "7.77569088416849750000"
//             }]`,
// 		},
// 	}
// 	for _, tt := range tests {
// 		mockClient := MockResponse(tt.wantRaw)

// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &Client{
// 				baseRestURL: tt.fields.baseRestURL,
// 				baseWsURL:   tt.fields.baseWsURL,
// 				key:         tt.fields.key,
// 				secret:      tt.fields.secret,
// 				passphrase:  tt.fields.passphrase,
// 				httpClient:  mockClient,
// 			}

// 			got, err := c.GetWithdrawalPowerForCurrency(tt.args.qp)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Client.GetWithdrawalPowerForCurrency() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if len(c.httpClient.(*MockClient).Requests) != 1 {
// 				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
// 			}

// 			validateHeaders(t, c)

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Client.GetWithdrawalPowerForCurrency() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestClient_GetAllWithdrawalPower(t *testing.T) {
// 	tests := [...]struct {
// 		name    string
// 		fields  fields
// 		want    []AllWithdrawalPower
// 		wantRaw string
// 		wantErr bool
// 	}{
// 		{
// 			name:   "when a successful call is made to get all withdrawal power for an account",
// 			fields: defaultFields(),
// 			want: []AllWithdrawalPower{{
// 				ProfileID: "8058d771-2d88-4f0f-ab6e-299c153d4308",
// 				MarginableWithdrawalPowers: []struct {
// 					Currency        string `json:"currency"`
// 					WithdrawalPower string `json:"withdrawal_power"`
// 				}{
// 					{
// 						Currency:        "ETH",
// 						WithdrawalPower: "0.0000000000000000",
// 					},
// 					{
// 						Currency:        "BTC",
// 						WithdrawalPower: "0.00184821818021342913",
// 					},
// 					{
// 						Currency:        "USD",
// 						WithdrawalPower: "7.77601796034649750000",
// 					},
// 					{
// 						Currency:        "USDC",
// 						WithdrawalPower: "1.00332803238200000000",
// 					},
// 				},
// 			}},
// 			wantRaw: `[{
//                 "profile_id": "8058d771-2d88-4f0f-ab6e-299c153d4308",
//                 "marginable_withdrawal_powers": [{
//                     "currency": "ETH",
//                     "withdrawal_power": "0.0000000000000000"
//                 }, {
//                     "currency": "BTC",
//                     "withdrawal_power": "0.00184821818021342913"
//                 }, {
//                     "currency": "USD",
//                     "withdrawal_power": "7.77601796034649750000"
//                 }, {
//                     "currency": "USDC",
//                     "withdrawal_power": "1.00332803238200000000"
//                 }]
//             }]`,
// 		},
// 	}
// 	for _, tt := range tests {
// 		mockClient := MockResponse(tt.wantRaw)

// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &Client{
// 				baseRestURL: tt.fields.baseRestURL,
// 				baseWsURL:   tt.fields.baseWsURL,
// 				key:         tt.fields.key,
// 				secret:      tt.fields.secret,
// 				passphrase:  tt.fields.passphrase,
// 				httpClient:  mockClient,
// 			}

// 			got, err := c.GetAllWithdrawalPower()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Client.GetAllWithdrawalPower() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if len(c.httpClient.(*MockClient).Requests) != 1 {
// 				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
// 			}

// 			validateHeaders(t, c)

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Client.GetAllWithdrawalPower() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestClient_GetMarginExitPlan(t *testing.T) {
// 	tests := [...]struct {
// 		name    string
// 		fields  fields
// 		want    ExitPlan
// 		wantRaw string
// 		wantErr bool
// 	}{
// 		{
// 			name:   "when a successful call is made to get an exit plan",
// 			fields: defaultFields(),
// 			want: ExitPlan{
// 				ID:        "239f4dc6-72b6-11ea-b311-168e5016c449",
// 				UserID:    "5cf6e115aaf44503db300f1e",
// 				ProfileID: "8058d771-2d88-4f0f-ab6e-299c153d4308",
// 				AccountsList: []struct {
// 					ID       string "json:\"id\""
// 					Currency string "json:\"currency\""
// 					Amount   string "json:\"amount\""
// 				}{
// 					{
// 						ID:       "434e1152-8eb5-4bfa-89a1-92bb1dcaf0c3",
// 						Currency: "BTC",
// 						Amount:   "0.00221897",
// 					},
// 					{
// 						ID:       "6d326768-71f2-4068-99dc-7075c78f6402",
// 						Currency: "USD",
// 						Amount:   "-1.9004458409934425",
// 					},
// 					{
// 						ID:       "120c8fcf-94da-4b45-9c43-18f114880f7a",
// 						Currency: "USDC",
// 						Amount:   "1.003328032382",
// 					},
// 				},
// 				EquityPercentage:    "0.8744507743595747",
// 				TotalAssetsUSD:      "15.137057447382",
// 				TotalLiabilitiesUSD: "1.9004458409934425",
// 				StrategiesList: []ExitStrategy{{
// 					Type:      "not sure",
// 					Amount:    "not sure",
// 					Product:   "not sure",
// 					Strategy:  "not sure",
// 					AccountID: "not sure",
// 					OrderID:   "not sure",
// 				}},
// 				CreatedAt: "2020-03-30 18:41:59.547863064 +0000 UTC m=+260120.906569441",
// 			},
// 			wantRaw: `{
//                 "id": "239f4dc6-72b6-11ea-b311-168e5016c449",
//                 "userId": "5cf6e115aaf44503db300f1e",
//                 "profileId": "8058d771-2d88-4f0f-ab6e-299c153d4308",
//                 "accountsList": [{
//                     "id": "434e1152-8eb5-4bfa-89a1-92bb1dcaf0c3",
//                     "currency": "BTC",
//                     "amount": "0.00221897"
//                 }, {
//                     "id": "6d326768-71f2-4068-99dc-7075c78f6402",
//                     "currency": "USD",
//                     "amount": "-1.9004458409934425"
//                 }, {
//                     "id": "120c8fcf-94da-4b45-9c43-18f114880f7a",
//                     "currency": "USDC",
//                     "amount": "1.003328032382"
//                 }],
//                 "equityPercentage": "0.8744507743595747",
//                 "totalAssetsUsd": "15.137057447382",
//                 "totalLiabilitiesUsd": "1.9004458409934425",
//                 "strategiesList": [{
//                     "type": "not sure",
//                     "amount": "not sure",
//                     "product": "not sure",
//                     "strategy": "not sure",
//                     "accountId": "not sure",
//                     "orderId": "not sure"
//                 }],
//                 "createdAt": "2020-03-30 18:41:59.547863064 +0000 UTC m=+260120.906569441"
//             }`,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockClient := MockResponse(tt.wantRaw)

// 			c := &Client{
// 				baseRestURL: tt.fields.baseRestURL,
// 				baseWsURL:   tt.fields.baseWsURL,
// 				key:         tt.fields.key,
// 				secret:      tt.fields.secret,
// 				passphrase:  tt.fields.passphrase,
// 				httpClient:  mockClient,
// 			}

// 			got, err := c.GetMarginExitPlan()
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Client.GetMarginExitPlan() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if len(c.httpClient.(*MockClient).Requests) != 1 {
// 				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
// 			}

// 			validateHeaders(t, c)

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Client.GetMarginExitPlan() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestClient_ListLiquidationHistory(t *testing.T) {
// 	type args struct {
// 		qp QueryParams
// 	}
// 	tests := [...]struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    []LiquidationHistory
// 		wantRaw string
// 		wantErr bool
// 	}{
// 		{
// 			name:   "when a successful call is made to list your liquidation history",
// 			fields: defaultFields(),
// 			args:   args{QueryParams{CurrencyParam: "BTC"}},
// 			want: []LiquidationHistory{{
// 				EventID:   "6d0edaf1-0c6f-11ea-a88c-0a04debd8c33",
// 				EventTime: "2019-11-21T14:58:49.879Z",
// 				Orders: []struct {
// 					ID            string "json:\"id\""
// 					Size          string "json:\"size\""
// 					ProductID     string "json:\"product_id\""
// 					ProfileID     string "json:\"profile_id\""
// 					Side          string "json:\"side\""
// 					Type          string "json:\"type\""
// 					PostOnly      bool   "json:\"post_only\""
// 					CreatedAt     string "json:\"created_at\""
// 					DoneAt        string "json:\"done_at\""
// 					DoneReason    string "json:\"done_reason\""
// 					FillFees      string "json:\"fill_fees\""
// 					FilledSize    string "json:\"filled_size\""
// 					ExecutedValue string "json:\"executed_value\""
// 					Status        string "json:\"status\""
// 					Settled       bool   "json:\"settled\""
// 				}{{
// 					ID:            "6c8d0d4e-0c6f-11ea-947d-0a04debd8c33",
// 					Size:          "0.02973507",
// 					ProductID:     "BTC-USD",
// 					ProfileID:     "8058d771-2d88-4f0f-ab6e-299c153d4308",
// 					Side:          "sell",
// 					Type:          "market",
// 					PostOnly:      false,
// 					CreatedAt:     "2019-11-21 14:58:49.582305+00",
// 					DoneAt:        "2019-11-21 14:58:49.596+00",
// 					DoneReason:    "filled",
// 					FillFees:      "1.1529981537990000",
// 					FilledSize:    "0.02973507",
// 					ExecutedValue: "230.5996307598000000",
// 					Status:        "done",
// 					Settled:       true,
// 				}},
// 			}},
// 			wantRaw: `[{
//                 "event_id": "6d0edaf1-0c6f-11ea-a88c-0a04debd8c33",
//                 "event_time": "2019-11-21T14:58:49.879Z",
//                 "orders": [{
//                     "id": "6c8d0d4e-0c6f-11ea-947d-0a04debd8c33",
//                     "size": "0.02973507",
//                     "product_id": "BTC-USD",
//                     "profile_id": "8058d771-2d88-4f0f-ab6e-299c153d4308",
//                     "side": "sell",
//                     "type": "market",
//                     "post_only": false,
//                     "created_at": "2019-11-21 14:58:49.582305+00",
//                     "done_at": "2019-11-21 14:58:49.596+00",
//                     "done_reason": "filled",
//                     "fill_fees": "1.1529981537990000",
//                     "filled_size": "0.02973507",
//                     "executed_value": "230.5996307598000000",
//                     "status": "done",
//                     "settled": true
//                 }]
//             }]`,
// 		},
// 	}
// 	for _, tt := range tests {
// 		mockClient := MockResponse(tt.wantRaw)

// 		t.Run(tt.name, func(t *testing.T) {
// 			c := &Client{
// 				baseRestURL: tt.fields.baseRestURL,
// 				baseWsURL:   tt.fields.baseWsURL,
// 				key:         tt.fields.key,
// 				secret:      tt.fields.secret,
// 				passphrase:  tt.fields.passphrase,
// 				httpClient:  mockClient,
// 			}

// 			got, err := c.ListLiquidationHistory(tt.args.qp)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Client.ListLiquidationHistory() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if len(c.httpClient.(*MockClient).Requests) != 1 {
// 				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
// 			}

// 			validateHeaders(t, c)

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Client.ListLiquidationHistory() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestClient_GetPositionRefreshAmounts(t *testing.T) {
// 	type args struct {
// 		qp QueryParams
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    RefreshAmount
// 		wantRaw string
// 		wantErr bool
// 	}{
// 		{
// 			name:   "when a successful call is made to get position refresh amounts",
// 			fields: defaultFields(),
// 			args:   args{},
// 			want: RefreshAmount{
// 				OneDayRenewalAmount: "0",
// 				TwoDayRenewalAmount: "417.93",
// 			},
// 			wantRaw: `{
//                 "oneDayRenewalAmount": "0",
//                 "twoDayRenewalAmount": "417.93"
//             }`,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockClient := MockResponse(tt.wantRaw)

// 			c := &Client{
// 				baseRestURL: tt.fields.baseRestURL,
// 				baseWsURL:   tt.fields.baseWsURL,
// 				key:         tt.fields.key,
// 				secret:      tt.fields.secret,
// 				passphrase:  tt.fields.passphrase,
// 				httpClient:  mockClient,
// 			}

// 			got, err := c.GetPositionRefreshAmounts(tt.args.qp)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Client.GetPositionRefreshAmounts() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if len(c.httpClient.(*MockClient).Requests) != 1 {
// 				t.Errorf("should have made one request, but made: %d", len(c.httpClient.(*MockClient).Requests))
// 			}

// 			validateHeaders(t, c)

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Client.GetPositionRefreshAmounts() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
