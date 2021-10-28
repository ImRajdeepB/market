package app

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	ob "github.com/muzykantov/orderbook"
	"github.com/rajdeepbh/market/util"
)

type App struct {
	Router   *mux.Router
	Resolver util.Resolver
}

func (a *App) Initialize() {
	a.Resolver = util.Resolver{
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

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) Run(port string) {
	log.Fatal(http.ListenAndServe(":"+port, a.Router))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "*")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/coins", func(w http.ResponseWriter, r *http.Request) {
		a.AssetsHandler(w, r)
	}).Methods("GET", http.MethodOptions)
	a.Router.HandleFunc("/coins/{asset}", func(w http.ResponseWriter, r *http.Request) {
		a.CoinHandler(w, r)
	}).Methods("GET", http.MethodOptions)
	a.Router.HandleFunc("/coins/{asset}/depth", func(w http.ResponseWriter, r *http.Request) {
		a.DepthHandler(w, r)
	}).Methods("GET", http.MethodOptions)
	a.Router.HandleFunc("/coins/{asset}/buy", func(w http.ResponseWriter, r *http.Request) {
		a.BuyHandler(w, r)
	}).Methods("POST", http.MethodOptions)
	a.Router.HandleFunc("/coins/{asset}/sell", func(w http.ResponseWriter, r *http.Request) {
		a.SellHandler(w, r)
	}).Methods("POST", http.MethodOptions)
}
