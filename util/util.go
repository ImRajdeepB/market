package util

import ob "github.com/muzykantov/orderbook"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Asset int

const (
	BTC Asset = iota
	ETH
	CELO
	SOL
	DOT
	XTZ
)

func (s Asset) String() string {
	switch s {
	case BTC:
		return "BTC"
	case ETH:
		return "ETH"
	case CELO:
		return "CELO"
	case SOL:
		return "SOL"
	case DOT:
		return "DOT"
	case XTZ:
		return "XTZ"
	}
	return "unknown"
}

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
	// OrderBook map[Asset]*ob.OrderBook
}

func hlo() {
	var _ = Resolver{
		Assets: map[string]Pair{
			"xtz": {
				Valid:     true,
				OrderBook: ob.NewOrderBook(),
			},
		},
		// OrderBook: make(map[Asset]*ob.OrderBook),
	}
}
