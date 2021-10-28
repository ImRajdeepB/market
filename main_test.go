package main_test

import (
	// "bytes"
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

func TestLimitBuyHandler(t *testing.T) {
	// var jsonStr = []byte(`{"type":"limit", "quantity": "10", "price": "3890"}`)
	// req, _ := http.NewRequest("POST", "/coins/ETH/buy", bytes.NewBuffer(jsonStr))
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
