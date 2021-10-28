package main_test

import (
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

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}
