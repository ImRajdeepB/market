package util

import ob "github.com/muzykantov/orderbook"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type AssetMarketPrice struct {
	Asset       string
	MarketPrice string
}

type Pair struct {
	Valid     bool
	OrderBook *ob.OrderBook
}

type Resolver struct {
	Assets       map[string]Pair
	AssetTickers []string
}
