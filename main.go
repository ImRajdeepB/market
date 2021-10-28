package main

import (
	"github.com/rajdeepbh/market/app"
)

func main() {
	a := app.App{}
	a.Initialize()
	a.Run("80")
}
