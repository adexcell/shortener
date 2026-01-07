package main

import "github.com/adexcell/shortener/cmd/app"

// @title           Shortener API
// @version         1.0
// @description     URL Shortener Service with Analytics.
// @host            localhost:8080
// @BasePath        /

func main() {
	app := app.NewApp()
	app.Run()
}
