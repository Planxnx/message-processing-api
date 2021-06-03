package coingecko

import "time"

const (
	CurrencyUSD = "usd"
	CurrencyTHB = "thb"
	CurrencyBTC = "btc"
	CurrencyBNB = "bnb"
)

const (
	DefaultTimeout = 10 * time.Second
	DefaultExpired = 5 * time.Minute
)

const (
	CacheKeyCoinsList  = "coins-list"
	CacheKeyCoinPrices = "coin-prices"
)
