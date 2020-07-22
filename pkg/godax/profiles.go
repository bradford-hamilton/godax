package godax

import (
	"encoding/json"
	"net/http"
)

// Profile represents a coinbase pro profile which is equivalent to a portfolio.
/*
{
    "id": "86602c68-306a-4500-ac73-4ce56a91d83c",
    "user_id": "5844eceecf7e803e259d0365",
    "name": "default",
    "active": true,
    "is_default": true,
    "created_at": "2019-11-18T15:08:40.236309Z"
}
*/
type Profile struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Active    bool   `json:"active"`
	IsDefault bool   `json:"is_default"`
	CreatedAt string `json:"created_at"`
}

func (c *Client) listProfiles(timestamp, signature string, req *http.Request) ([]Profile, error) {
	res, err := c.do(timestamp, signature, req)
	if err != nil {
		return []Profile{}, err
	}
	defer res.Body.Close()

	var profiles []Profile
	if err := json.NewDecoder(res.Body).Decode(&profiles); err != nil {
		return []Profile{}, err
	}
	return profiles, nil
}

func (c *Client) getProfile(timestamp, signature string, req *http.Request) (Profile, error) {
	res, err := c.do(timestamp, signature, req)
	if err != nil {
		return Profile{}, err
	}
	defer res.Body.Close()

	var p Profile
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		return Profile{}, err
	}
	return p, nil
}
