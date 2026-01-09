package main

import (
	"log"

	"github.com/adexcell/shortener/cmd/app"
)

// @title           Shortener API
// @version         1.0
// @description     URL Shortener Service with Analytics.
// @host            localhost:8080
// @BasePath        /

func main() {
	app, err := app.NewApp()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	
	if err := app.Run(); err != nil {
		log.Fatalf("error: %v", err)
	}
}
