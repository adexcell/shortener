package shortener

import (
	"net/http"

	"github.com/adexcell/shortener/internal/controllers"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

const (
	postShortURL  = "/shorten"
	conversionURL = "/s/:short_url"
	analyticsURL  = "/analytics/:short_url"
)

type handler struct {
}

func NewShortenHandler() controllers.Handler {
	return &handler{}
}

func (h *handler) Register(router *ginext.Engine) {
	router.POST(postShortURL, h.PostShortURL)
	router.GET(conversionURL, h.ConversionURL)
	router.GET(analyticsURL, h.GetAnalytics)
}

func (h *handler) PostShortURL(c *ginext.Context) {
	zlog.Logger.Info().Msg("post notify")
	c.JSON(http.StatusOK, ginext.H{
		"msg": "post notify",
	})
}

func (h *handler) ConversionURL(c *ginext.Context) {
	var URL string
	c.ShouldBindJSON(&URL)

	zlog.Logger.Info().Msg("get notify")
	c.JSON(http.StatusOK, ginext.H{
		"msg": "get notify",
		"URL":  URL,
	})
}

func (h *handler) GetAnalytics(c *ginext.Context) {
	var id int
	c.ShouldBindJSON(&id)

	zlog.Logger.Info().Msg("delete notify")
	c.JSON(http.StatusOK, ginext.H{
		"msg": "delete notify",
		"id":  id,
	})
}
