package util

import (
	ob "github.com/muzykantov/orderbook"
	"github.com/shopspring/decimal"
)

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

type LimitOrderResp struct {
	Ticker                   string
	Side                     string
	OrderType                string
	Quantity                 string
	Price                    string
	Done                     []*ob.Order
	Partial                  *ob.Order
	PartialQuantityProcessed decimal.Decimal
}

type MarketOrderResp struct {
	Ticker                   string
	Side                     string
	OrderType                string
	Quantity                 string
	Price                    string
	Done                     []*ob.Order
	Partial                  *ob.Order
	PartialQuantityProcessed decimal.Decimal
	QuantityLeft             decimal.Decimal
}

type Resolver struct {
	Assets       map[string]Pair
	AssetTickers []string
}
