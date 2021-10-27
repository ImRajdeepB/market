package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	ob "github.com/muzykantov/orderbook"
	"github.com/shopspring/decimal"
)

func index(w http.ResponseWriter, req *http.Request) {
	// fmt.Fprintf(w, "hello\n")
	// fmt.Fprintf(w, "ETH:3900, BTC: 60000\n")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	// resp := make(map[string]map[string]map[string]int)
	// resp["market_data"]["current_price"]["usd"] = 60000

	currencies_map := make(map[string]int)
	currencies_map["usd"] = 60000
	mdata := make(map[string]map[string]int)
	mdata["current_price"] = currencies_map
	adata := make(map[string]map[string]map[string]int)
	adata["market_data"] = mdata
	jsonResp, err := json.Marshal(adata)
	// jsonResp, err := json.Marshal(Root{
	// 	market_data: MarketData{
	// 		current_price: currencies_map,
	// 	},
	// })
	if err != nil {
		fmt.Printf("Error happened in JSON marshal. Err: %s\n", err)
	}
	w.Write(jsonResp)
}

func index2(w http.ResponseWriter, req *http.Request) {
	// fmt.Fprintf(w, "hello\n")
	// fmt.Fprintf(w, "ETH:3900, BTC: 60000\n")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	// resp := make(map[string]map[string]map[string]int)
	// resp["market_data"]["current_price"]["usd"] = 60000

	currencies_map := make(map[string]int)
	currencies_map["usd"] = 4000
	mdata := make(map[string]map[string]int)
	mdata["current_price"] = currencies_map
	adata := make(map[string]map[string]map[string]int)
	adata["market_data"] = mdata
	jsonResp, err := json.Marshal(adata)
	// jsonResp, err := json.Marshal(Root{
	// 	market_data: MarketData{
	// 		current_price: currencies_map,
	// 	},
	// })
	if err != nil {
		fmt.Printf("Error happened in JSON marshal. Err: %s\n", err)
	}
	w.Write(jsonResp)
}

func main() {

	http.HandleFunc("/coins/bitcoin", index)
	http.HandleFunc("/coins/ethereum", index2)
	// http.HandleFunc("/headers", headers)
	http.ListenAndServe(":80", nil)

	orderBook := ob.NewOrderBook()
	// fmt.Println(orderBook)
	// ob.NewOrder("s", ob.Sell, 4.0, 5.0, time.Now())
	done, _, _, err := orderBook.ProcessLimitOrder(ob.Sell, "uinqueID", decimal.New(55, 0), decimal.New(100, 0))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(done))
	// fmt.Println(orderBook)
	orderBook.ProcessLimitOrder(ob.Buy, "uinqueID1", decimal.New(7, 0), decimal.New(98, 0))
	// fmt.Println(orderBook)
	orderBook.ProcessLimitOrder(ob.Buy, "uinqueID2", decimal.New(3, 0), decimal.New(120, 0))
	fmt.Println(orderBook)
	asks, bids := orderBook.Depth()
	fmt.Println("asks:")
	for _, ask := range asks {
		fmt.Println(ask.Price, ask.Quantity)
	}
	fmt.Println("bids:")
	for _, bid := range bids {
		fmt.Println(bid.Price, bid.Quantity)
	}
	fmt.Println(orderBook.CalculateMarketPrice(ob.Buy, decimal.New(1, 0)))

}
