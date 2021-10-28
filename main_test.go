package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rajdeepbh/market/app"
)

var a app.App

func TestMain(m *testing.M) {
	a = app.App{}
	a.Initialize()

	code := m.Run()

	os.Exit(code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestAssetsHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/coins", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != `[{"Asset":"BTC","MarketPrice":"0"},{"Asset":"ETH","MarketPrice":"0"},{"Asset":"CELO","MarketPrice":"0"},{"Asset":"SOL","MarketPrice":"0"},{"Asset":"DOT","MarketPrice":"0"},{"Asset":"XTZ","MarketPrice":"0"}]` {
		t.Errorf("Expected an array with '0' market price. Got %s", body)
	}
}

func TestGetCoinNoLiquidity(t *testing.T) {
	req, _ := http.NewRequest("GET", "/coins/ETH", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusInternalServerError, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "orderbook: insufficient quantity to calculate price" {
		t.Errorf("Expected the 'error' key of the response to be set to 'orderbook: insufficient quantity to calculate price'. Got '%s'", m["error"])
	}
}

func TestGetDepthNoLiquidity(t *testing.T) {
	req, _ := http.NewRequest("GET", "/coins/ETH/depth", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != `{"asks":null,"bids":null}` {
		t.Errorf("Expected an asks and bids to be null. Got %s", body)
	}
}

func TestCoin(t *testing.T) {
	req1, _ := http.NewRequest("POST", "/coins/ETH/sell?type=limit&price=3970&quantity=50", nil)
	response1 := executeRequest(req1)

	checkResponseCode(t, http.StatusCreated, response1.Code)

	var m1 map[string]string
	json.Unmarshal(response1.Body.Bytes(), &m1)
	if m1["Ticker"] != "ETH" {
		t.Errorf("Expected ticker symbol to be 'ETH'. Got '%v'", m1["Ticker"])
	}
	if m1["Side"] != "sell" {
		t.Errorf("Expected side to be 'sell'. Got '%v'", m1["Side"])
	}
	if m1["OrderType"] != "limit" {
		t.Errorf("Expected order type to be 'limit'. Got '%v'", m1["OrderType"])
	}
	if m1["Quantity"] != "50" {
		t.Errorf("Expected quantity to be '50'. Got '%v'", m1["Quantity"])
	}
	if m1["Price"] != "3970" {
		t.Errorf("Expected price to be '3970'. Got '%v'", m1["Price"])
	}
	if m1["Done"] != "" {
		t.Errorf("Expected Done to be null. Got '%v'", m1["Done"])
	}
	if m1["Partial"] != "" {
		t.Errorf("Expected Partial to be null. Got '%v'", m1["Partial"])
	}
	if m1["PartialQuantityProcessed"] != "0" {
		t.Errorf("Expected PartialQuantityProcessed to be '0'. Got '%v'", m1["PartialQuantityProcessed"])
	}

	// **

	req, _ := http.NewRequest("GET", "/coins/ETH/depth", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != `{"asks":[{"price":"3970","quantity":"50"}],"bids":null}` {
		t.Errorf("Expected one ask '{\"price\":\"3970\",\"quantity\":\"50\"}' and bids to be null. Got %s", body)
	}

	// **

	req2, _ := http.NewRequest("POST", "/coins/ETH/buy?type=limit&price=3971&quantity=30", nil)
	response2 := executeRequest(req2)

	checkResponseCode(t, http.StatusCreated, response2.Code)

	var m2 map[string]string
	json.Unmarshal(response2.Body.Bytes(), &m2)
	if m2["Ticker"] != "ETH" {
		t.Errorf("Expected ticker symbol to be 'ETH'. Got '%v'", m2["Ticker"])
	}
	if m2["Side"] != "buy" {
		t.Errorf("Expected side to be 'buy'. Got '%v'", m2["Side"])
	}
	if m2["OrderType"] != "limit" {
		t.Errorf("Expected order type to be 'limit'. Got '%v'", m2["OrderType"])
	}
	if m2["Quantity"] != "30" {
		t.Errorf("Expected quantity to be '30'. Got '%v'", m2["Quantity"])
	}
	if m2["Price"] != "3971" {
		t.Errorf("Expected price to be '3971'. Got '%v'", m2["Price"])
	}
	if m2["PartialQuantityProcessed"] != "30" {
		t.Errorf("Expected PartialQuantityProcessed to be '30'. Got '%v'", m2["PartialQuantityProcessed"])
	}

	// **

	req, _ = http.NewRequest("GET", "/coins/ETH/depth", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != `{"asks":[{"price":"3970","quantity":"20"}],"bids":null}` {
		t.Errorf("Expected one ask '{\"price\":\"3970\",\"quantity\":\"20\"}' and bids to be null. Got %s", body)
	}

	// **

	req3, _ := http.NewRequest("GET", "/coins/ETH", nil)
	response3 := executeRequest(req3)

	checkResponseCode(t, http.StatusOK, response3.Code)

	var m3 map[string]map[string]map[string]string
	err := json.Unmarshal(response3.Body.Bytes(), &m3)
	if err != nil {
		t.Errorf("Expected market price to be '3970'. Got error: '%v'", err.Error())
	}
	if m3["market_data"]["current_price"]["usd"] != "3970" {
		t.Errorf("Expected market price to be '3970'. Got '%v'", m3["market_data"]["current_price"]["usd"])
	}
}

func TestGetNonExistentCoin(t *testing.T) {
	req, _ := http.NewRequest("GET", "/coins/SOMESHITCOIN", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "invalid ticker" {
		t.Errorf("Expected the 'error' key of the response to be set to 'invalid ticker'. Got '%s'", m["error"])
	}
}

func TestLimitBuyHandler(t *testing.T) {
	req, _ := http.NewRequest("POST", "/coins/ETH/buy?type=limit&quantity=10&price=3890", nil)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["Ticker"] != "ETH" {
		t.Errorf("Expected ticker symbol to be 'ETH'. Got '%v'", m["Ticker"])
	}

	if m["Side"] != "buy" {
		t.Errorf("Expected side to be 'buy'. Got '%v'", m["Side"])
	}

	if m["OrderType"] != "limit" {
		t.Errorf("Expected order type to be 'limit'. Got '%v'", m["OrderType"])
	}

	if m["Quantity"] != "10" {
		t.Errorf("Expected quantity to be '10'. Got '%v'", m["Quantity"])
	}

	if m["Price"] != "3890" {
		t.Errorf("Expected price to be '3890'. Got '%v'", m["Price"])
	}
}
