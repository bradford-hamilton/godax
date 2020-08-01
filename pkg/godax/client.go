package godax

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

const userAgent = "godax coinbase pro client"

// HTTPClient is a simple http interface that both live and test code can implement
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// CoinbaseErrRes represents the shape that comes back when a status code is non-200
type CoinbaseErrRes struct {
	// Message contains the error information from coinbase
	Message string `json:"message"`
}

func newClient(sandbox bool) (*Client, error) {
	c := &Client{
		httpClient: retryablehttp.NewClient().StandardClient(),
	}
	if err := c.loadEnv(sandbox); err != nil {
		return nil, err
	}
	return c, nil
}

// exec note: decodes into the value pointed at by v. In other words, `v` cannot be nil or a non-pointer.
func (c *Client) exec(timestamp, method, path string, body []byte, qp *QueryParams, v interface{}) error {
	req, sig, err := c.createAndSignReq(timestamp, method, path, body, qp)
	if err != nil {
		return err
	}
	res, err := c.do(timestamp, sig, req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}
	return nil
}

func (c *Client) createAndSignReq(timestamp, method, path string, body []byte, qp *QueryParams) (*http.Request, string, error) {
	req, err := http.NewRequest(method, c.baseRestURL+path, bytes.NewBuffer(body))
	if err != nil {
		return nil, "", err
	}
	if qp != nil {
		c.setQueryParams(req, *qp)
	}
	sig, err := c.generateSig(timestamp, method, req.URL.RequestURI(), string(body))
	if err != nil {
		return nil, "", err
	}
	return req, sig, nil
}

// generateSig generates the signature for the CB-ACCESS-SIGN header.
// 1. base64 decode the client coinbase pro secret
// 2. create a sha256 HMAC using the base64 decoded secret
// 3. concatenate (timestamp + http method + coinbase pro URL path + message body), and get the bytes
//   - the timstamp used here must be the same as the one used for the CB-ACCESS-TIMESTAMP header
//   - message body can be omitted (typically for GET requests)
// 4. write the result to the hash and sum it
// 5. base64 encoded the digest
func (c *Client) generateSig(timestamp, method, path, body string) (string, error) {
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

func (c *Client) do(timestamp string, signature string, req *http.Request) (*http.Response, error) {
	c.setHeaders(req, timestamp, signature)
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, coinbaseError(res)
	}
	return res, nil
}

func (c *Client) setHeaders(req *http.Request, timestamp string, signature string) {
	req.Header.Set("CB-ACCESS-KEY", c.key)
	req.Header.Set("CB-ACCESS-SIGN", signature)
	req.Header.Set("CB-ACCESS-TIMESTAMP", timestamp)
	req.Header.Set("CB-ACCESS-PASSPHRASE", c.passphrase)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
}

func (c *Client) setQueryParams(req *http.Request, qp QueryParams) {
	q := req.URL.Query()
	for k, v := range qp {
		if v != "" {
			q.Add(string(k), v)
		}
	}
	req.URL.RawQuery = q.Encode()
}

func coinbaseError(res *http.Response) error {
	var err CoinbaseErrRes
	if err := json.NewDecoder(res.Body).Decode(&err); err != nil {
		return err
	}
	return fmt.Errorf("status code: %d, message: %s", res.StatusCode, err.Message)
}
