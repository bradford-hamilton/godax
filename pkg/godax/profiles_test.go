package godax

import (
	"reflect"
	"testing"
)

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
