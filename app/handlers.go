package app

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	ob "github.com/muzykantov/orderbook"
	"github.com/rajdeepbh/market/util"
	"github.com/shopspring/decimal"
)

func (a *App) AssetsHandler(w http.ResponseWriter, r *http.Request) {
	assets := []*util.AssetMarketPrice{}
	for _, asset := range a.Resolver.AssetTickers {
		buy_market_price, _ := a.Resolver.Assets[asset].OrderBook.CalculateMarketPrice(ob.Buy, decimal.NewFromFloat(1))
		bmp := buy_market_price.String()

		assets = append(assets, &util.AssetMarketPrice{
			Asset:       asset,
			MarketPrice: bmp,
		})
	}

	respondWithJSON(w, http.StatusOK, assets)
}

func (a *App) DepthHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	asset := vars["asset"]
	if !a.Resolver.Assets[asset].Valid {
		respondWithError(w, http.StatusNotFound, "invalid ticker")
		return
	}

	asks, bids := a.Resolver.Assets[asset].OrderBook.Depth()

	mmp := make(map[string][]*ob.PriceLevel)
	mmp["asks"] = asks
	mmp["bids"] = bids

	respondWithJSON(w, http.StatusOK, mmp)
}

func (a *App) CoinHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	asset := vars["asset"]
	if !a.Resolver.Assets[asset].Valid {
		respondWithError(w, http.StatusNotFound, "invalid ticker")
		return
	}

	buy_market_price, err := a.Resolver.Assets[asset].OrderBook.CalculateMarketPrice(ob.Buy, decimal.NewFromFloat(1))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	v := map[string]map[string]map[string]string{
		"market_data": {
			"current_price": {
				"usd": buy_market_price.String(),
			},
		},
	}

	respondWithJSON(w, http.StatusOK, v)
}

func (a *App) BuyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	asset := vars["asset"]
	if !a.Resolver.Assets[asset].Valid {
		respondWithError(w, http.StatusNotFound, "invalid ticker")
		return
	}
	order_type := r.URL.Query()["type"]
	price := r.URL.Query()["price"]
	quantity := r.URL.Query()["quantity"]

	ProcessOrder(w, r, asset, ob.Buy, order_type, quantity, price, &a.Resolver)
}

func (a *App) SellHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	asset := vars["asset"]
	if !a.Resolver.Assets[asset].Valid {
		respondWithError(w, http.StatusNotFound, "invalid ticker")
		return
	}
	order_type := r.URL.Query()["type"]
	price := r.URL.Query()["price"]
	quantity := r.URL.Query()["quantity"]

	ProcessOrder(w, r, asset, ob.Sell, order_type, quantity, price, &a.Resolver)
}

func ProcessOrder(w http.ResponseWriter, r *http.Request, asset string, side ob.Side, order_type []string, quantity []string, price []string, resolver *util.Resolver) {
	if len(order_type) != 1 {
		respondWithError(w, http.StatusBadRequest, "unknown order type")
		return
	}
	_quantity, err := strconv.ParseFloat(quantity[0], 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if order_type[0] == "limit" {
		if len(price) == 1 && len(quantity) == 1 {
			_price, err := strconv.ParseFloat(price[0], 64)
			if err != nil {
				respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
			done, partial, partialQuantityProcessed, err := resolver.Assets[asset].OrderBook.ProcessLimitOrder(side, uuid.NewString(), decimal.NewFromFloat(_quantity), decimal.NewFromFloat(_price))
			if err != nil {
				respondWithError(w, http.StatusBadRequest, err.Error())
				return
			}

			respondWithJSON(w, http.StatusCreated, util.LimitOrderResp{
				Ticker:                   asset,
				Side:                     side.String(),
				OrderType:                order_type[0],
				Quantity:                 quantity[0],
				Price:                    price[0],
				Done:                     done,
				Partial:                  partial,
				PartialQuantityProcessed: partialQuantityProcessed,
			})
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	} else if order_type[0] == "market" {
		if len(quantity) == 1 {
			done, partial, partialQuantityProcessed, quantityLeft, err := resolver.Assets[asset].OrderBook.ProcessMarketOrder(side, decimal.NewFromFloat(_quantity))
			if err != nil {
				respondWithError(w, http.StatusNotFound, err.Error())
				return
			}

			respondWithJSON(w, http.StatusCreated, util.MarketOrderResp{
				Ticker:                   asset,
				Side:                     side.String(),
				OrderType:                order_type[0],
				Quantity:                 quantity[0],
				Price:                    price[0],
				Done:                     done,
				Partial:                  partial,
				PartialQuantityProcessed: partialQuantityProcessed,
				QuantityLeft:             quantityLeft,
			})
		} else {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	} else {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
}
