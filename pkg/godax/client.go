package godax

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"
)

func newClient(sandbox bool) (*Client, error) {
	c := &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
	if err := c.loadEnv(sandbox); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) setHeaders(req *http.Request, timestamp string, signature string) {
	req.Header.Set("CB-ACCESS-KEY", c.key)
	req.Header.Set("CB-ACCESS-SIGN", signature)
	req.Header.Set("CB-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("CB-ACCESS-PASSPHRASE", c.passphrase)
	req.Header.Add("User-Agent", "godax coinbase pro client")
}

func (c *Client) get(path, timestamp, signature string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseRestURL+path, nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req, timestamp, signature)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// generateSignature generates the signature for the CB-ACCESS-SIGN header.
// 1. base64 decode the client coinbase pro secret
// 2. create a sha256 HMAC using the base64 decoded secret
// 3. concatenate (timestamp + http method + coinbase pro URL path + message body), and get the bytes
//   - the timstamp used here must be the same as the one used for the CB-ACCESS-TIMESTAMP header
//   - message body can be omitted (typically for GET requests)
// 4. write the result to the hash and sum it
// 5. base64 encoded the digest
func (c *Client) generateSignature(timestamp, path, method, body string) (string, error) {
	secret, err := base64.StdEncoding.DecodeString(c.secret)
	if err != nil {
		return "", err
	}

	hash := hmac.New(sha256.New, secret)
	msg := []byte(timestamp + method + path + body)
	if _, err = hash.Write(msg); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}
