package coingecko

import (
	"fmt"
	"net/http"
	"strings"

	cache "github.com/patrickmn/go-cache"
	coingecko "github.com/superoo7/go-gecko/v3"
	cointypes "github.com/superoo7/go-gecko/v3/types"
)

type CoinGecko struct {
	client *coingecko.Client
	cache  *cache.Cache
}

func New() (*CoinGecko, error) {
	httpClient := &http.Client{
		Timeout: DefaultTimeout,
	}
	client := coingecko.NewClient(httpClient)

	_, err := client.Ping()
	if err != nil {
		return nil, err
	}

	return &CoinGecko{
		client: client,
		cache:  cache.New(DefaultExpired, 2*DefaultExpired),
	}, nil
}

func (cg *CoinGecko) ValidateSupportedCoin(symbol string) bool {
	coinsList, err := cg.GetCoinsList()
	if err != nil {
		return false
	}

	for _, coin := range coinsList {
		if coin.Symbol == strings.ToLower(symbol) {
			return true
		}
	}

	return false
}

func (cg *CoinGecko) GetCoinsList() ([]cointypes.CoinsListItem, error) {

	cached, found := cg.cache.Get(CacheKeyCoinsList)
	if found {
		if coinsList, ok := cached.([]cointypes.CoinsListItem); ok {
			fmt.Println("use coin list cached!")
			return coinsList, nil
		}
	}

	coinsList, err := cg.client.CoinsList()
	if err != nil {
		return nil, err
	}

	go cg.cache.Add(CacheKeyCoinsList, []cointypes.CoinsListItem(*coinsList), DefaultExpired)

	return []cointypes.CoinsListItem(*coinsList), nil
}

func (cg *CoinGecko) GetCoinPrices(symbol string) (map[string]map[string]float32, error) {

	cached, found := cg.cache.Get(CacheKeyCoinPrices + "-" + symbol)
	if found {
		if v, ok := cached.(map[string]map[string]float32); ok {
			fmt.Println("use prices cached!")
			return v, nil
		}
	}

	coinsList, _ := cg.GetCoinsList()

	targetIDs := make([]string, 0)
	for _, coin := range coinsList {
		if coin.Symbol == strings.ToLower(symbol) {
			targetIDs = append(targetIDs, coin.ID)
		}
	}

	prices, err := cg.client.SimplePrice(targetIDs, []string{CurrencyTHB, CurrencyUSD})
	if err != nil {
		return nil, err
	}

	go cg.cache.Add(CacheKeyCoinPrices + "-" + symbol, *prices, DefaultExpired)

	return *prices, nil
}
