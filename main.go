package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	ob "github.com/muzykantov/orderbook"
	"github.com/rajdeepbh/market/util"
	"github.com/shopspring/decimal"
)

func DepthHandler(w http.ResponseWriter, r *http.Request, resolver *util.Resolver) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	vars := mux.Vars(r)
	asset := vars["asset"]
	if !resolver.Assets[asset].Valid {
		fmt.Fprintf(w, "{\"error\": \"invalid ticker\"}\n")
		return
	}

	asks, bids := resolver.Assets[asset].OrderBook.Depth()

	mmp := make(map[string][]*ob.PriceLevel)
	mmp["asks"] = asks
	mmp["bids"] = bids
	jsonResp, err := json.Marshal(mmp)
	if err != nil {
		fmt.Fprintf(w, "{\"error\": \"%s\"}\n", err)
		return
	}
	w.Write(jsonResp)
}

func AssetsHandler(w http.ResponseWriter, r *http.Request, resolver *util.Resolver) {
	// w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	assets := []*util.AssetMarketPrice{}
	for _, asset := range resolver.AssetTickers {
		buy_market_price, err := resolver.Assets[asset].OrderBook.CalculateMarketPrice(ob.Buy, decimal.NewFromFloat(1))
		bmp := buy_market_price.String()
		if err != nil {
			// bmp = "0"
		}
		assets = append(assets, &util.AssetMarketPrice{
			Asset:       asset,
			MarketPrice: bmp,
		})
	}
	jsonResp, err := json.Marshal(assets)
	if err != nil {
		fmt.Fprintf(w, "{\"error\": \"%s\"}\n", err)
	}
	w.Write(jsonResp)
}

func CoinHandler(w http.ResponseWriter, r *http.Request, resolver *util.Resolver) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	vars := mux.Vars(r)
	asset := vars["asset"]
	println(asset)
	if !resolver.Assets[asset].Valid {
		fmt.Fprintf(w, "{\"error\": \"invalid ticker\"}\n")
		return
	}

	buy_market_price, err := resolver.Assets[asset].OrderBook.CalculateMarketPrice(ob.Buy, decimal.NewFromFloat(1))
	if err != nil {
		// fmt.Fprintf(w, "{\"error\": \"%s\"}\n", err)
		// return
	}
	fmt.Println("buy_market_price", buy_market_price)
	v := map[string]map[string]map[string]string{
		"market_data": {
			"current_price": {
				"usd": buy_market_price.String(),
			},
		},
	}
	jsonResp, err := json.Marshal(v)
	if err != nil {
		fmt.Fprintf(w, "{\"error\": \"%s\"}\n", err)
		return
	}
	w.Write(jsonResp)
}

func ProcessOrder(w http.ResponseWriter, r *http.Request, asset string, side ob.Side, order_type []string, quantity []string, price []string, resolver *util.Resolver) {
	fmt.Println(asset, side, order_type, quantity, price)
	if len(order_type) != 1 {
		fmt.Fprintf(w, "{\"error\": \"unknown order type\"}\n")
		return
	}
	_quantity, err := strconv.ParseFloat(quantity[0], 64)
	if err != nil {
		fmt.Fprintf(w, "{\"error\": \"%s\"}\n", err)
		return
	}
	if order_type[0] == "limit" {
		if len(price) == 1 && len(quantity) == 1 {
			_price, err := strconv.ParseFloat(price[0], 64)
			if err != nil {
				fmt.Fprintf(w, "{\"error\": \"%s\"}\n", err)
				return
			}
			done, partial, partialQuantityProcessed, err := resolver.Assets[asset].OrderBook.ProcessLimitOrder(side, uuid.NewString(), decimal.NewFromFloat(_quantity), decimal.NewFromFloat(_price))
			if err != nil {
				fmt.Fprintf(w, "{\"error\": \"%s\"}\n", err)
				return
			}
			fmt.Println("d", done, "p", partial, "p", partialQuantityProcessed)
		} else {
			fmt.Fprintf(w, "{\"error\": \"%s\"}\n", err)
			return
		}
	} else if order_type[0] == "market" {
		if len(quantity) == 1 {
			done, partial, partialQuantityProcessed, quantityLeft, err := resolver.Assets[asset].OrderBook.ProcessMarketOrder(side, decimal.NewFromFloat(_quantity))
			if err != nil {
				fmt.Fprintf(w, "{\"error\": \"%s\"}\n", err)
				return
			}
			fmt.Println("d", done, "p", partial, "p", partialQuantityProcessed, "q", quantityLeft)
		} else {
			fmt.Fprintf(w, "{\"error\": \"%s\"}\n", err)
			return
		}
	} else {
		fmt.Fprintf(w, "{\"error\": \"unknown order type\"}\n")
		return
	}

}

func BuyHandler(w http.ResponseWriter, r *http.Request, resolver *util.Resolver) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	vars := mux.Vars(r)
	asset := vars["asset"]
	if !resolver.Assets[asset].Valid {
		fmt.Fprintf(w, "{\"error\": \"invalid ticker\"}\n")
		return
	}
	order_type := r.URL.Query()["type"]
	price := r.URL.Query()["price"]
	quantity := r.URL.Query()["quantity"]

	ProcessOrder(w, r, asset, ob.Buy, order_type, quantity, price, resolver)
}

func SellHandler(w http.ResponseWriter, r *http.Request, resolver *util.Resolver) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	vars := mux.Vars(r)
	asset := vars["asset"]
	if !resolver.Assets[asset].Valid {
		fmt.Fprintf(w, "{\"error\": \"invalid ticker\"}\n")
		return
	}
	order_type := r.URL.Query()["type"]
	price := r.URL.Query()["price"]
	quantity := r.URL.Query()["quantity"]

	ProcessOrder(w, r, asset, ob.Sell, order_type, quantity, price, resolver)
}

func main() {
	resolver := util.Resolver{
		Assets: map[string]util.Pair{
			"BTC":  {Valid: true, OrderBook: ob.NewOrderBook()},
			"ETH":  {Valid: true, OrderBook: ob.NewOrderBook()},
			"CELO": {Valid: true, OrderBook: ob.NewOrderBook()},
			"SOL":  {Valid: true, OrderBook: ob.NewOrderBook()},
			"DOT":  {Valid: true, OrderBook: ob.NewOrderBook()},
			"XTZ":  {Valid: true, OrderBook: ob.NewOrderBook()},
		},
		AssetTickers: []string{"BTC", "ETH", "CELO", "SOL", "DOT", "XTZ"},
	}
	router := mux.NewRouter()

	router.HandleFunc("/coins", func(w http.ResponseWriter, r *http.Request) {
		AssetsHandler(w, r, &resolver)
	}).Methods("GET", http.MethodOptions)
	router.HandleFunc("/coins/{asset}", func(w http.ResponseWriter, r *http.Request) {
		CoinHandler(w, r, &resolver)
	}).Methods("GET", http.MethodOptions)
	router.HandleFunc("/coins/{asset}/depth", func(w http.ResponseWriter, r *http.Request) {
		DepthHandler(w, r, &resolver)
	}).Methods("GET", http.MethodOptions)
	router.HandleFunc("/coins/{asset}/buy", func(w http.ResponseWriter, r *http.Request) {
		BuyHandler(w, r, &resolver)
	}).Methods("POST", http.MethodOptions)
	router.HandleFunc("/coins/{asset}/sell", func(w http.ResponseWriter, r *http.Request) {
		SellHandler(w, r, &resolver)
	}).Methods("POST", http.MethodOptions)

	fmt.Println("starting server")
	fmt.Println("error starting server", http.ListenAndServe(":80", router))
}
