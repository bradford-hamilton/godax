package godax

import (
	"fmt"
	"os"
)

const (
	// Sandbox coinbase pro RESTful and web socket endpoints
	sandboxREST = "https://api-public.sandbox.pro.coinbase.com"
	sandboxWS   = "wss://ws-feed-public.sandbox.pro.coinbase.com"

	// Live coinbase pro RESTful and web socket endpoints
	liveREST = "https://api.pro.coinbase.com"
	liveWS   = "wss://ws-feed.pro.coinbase.com"

	// Every godax Client needs access to the following ENV variables
	coinbaseProKey        = "COINBASE_PRO_KEY"
	coinbaseProSecret     = "COINBASE_PRO_SECRET"
	coinbaseProPassphrase = "COINBASE_PRO_PASSPHRASE"
)

func (c *Client) loadEnv(sandbox bool) error {
	k := os.Getenv(coinbaseProKey)
	s := os.Getenv(coinbaseProSecret)
	p := os.Getenv(coinbaseProPassphrase)

	if k == "" {
		return errMissingEnv(coinbaseProKey)
	}
	if s == "" {
		return errMissingEnv(coinbaseProSecret)
	}
	if p == "" {
		return errMissingEnv(coinbaseProPassphrase)
	}

	if sandbox {
		c.baseRestURL = sandboxREST
		c.baseWsURL = sandboxWS
	} else {
		c.baseRestURL = liveREST
		c.baseWsURL = liveWS
	}

	c.key = k
	c.secret = s
	c.passphrase = p

	return nil
}

func errMissingEnv(envVar string) error {
	return fmt.Errorf("please provide a %s", envVar)
}
