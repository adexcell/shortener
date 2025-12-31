package main

import (
	"net/http"

	"github.com/adexcell/shortener/internal/shortener"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	zlog.Init()

	zlog.Logger.Info().Msg("create router")
	httprouter := ginext.New("debug")

	zlog.Logger.Info().Msg("register shorten handler")
	shortenHandler := shortener.NewShortenHandler()
	shortenHandler.Register(httprouter)

	httprouter.GET("/", Hello)

	httprouter.Run()
}

func Hello(c *ginext.Context) {
	c.JSON(http.StatusOK, ginext.H{
			"msg": "hello_world",
		})
}
