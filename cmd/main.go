package main

import "github.com/adexcell/shortener/cmd/app"

func main() {
	app := app.NewApp()
	app.Run()
}
