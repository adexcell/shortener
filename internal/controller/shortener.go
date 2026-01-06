package controller

import (
	"net/http"

	"github.com/adexcell/shortener/internal/usecase"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

const (
	postShortURL  = "/shorten"
	conversionURL = "/s/:short_url"
	analyticsURL  = "/analytics/:short_url"
)

type handler struct {
	usecase *usecase.ShortenerUsecase
}

func NewShortenHandler(u *usecase.ShortenerUsecase) Handler {
	return &handler{usecase: u}
}

func (h *handler) Register(router *ginext.Engine) {
	router.POST(postShortURL, h.PostShortURL)
	router.GET(conversionURL, h.ConversionURL)
	router.GET(analyticsURL, h.GetAnalytics)
}

type shortenRequest struct {
	URL string `json:"url" binding:"required"`
}

func (h *handler) PostShortURL(c *ginext.Context) {
	var req shortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid request"})
		return
	}

	code, err := h.usecase.Shorten(c.Request.Context(), req.URL)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to shorten url")
		c.JSON(http.StatusInternalServerError, ginext.H{"error":"db error"})
		return
	}
	
	c.JSON(http.StatusOK, ginext.H{"short_url": "/s/" + code})
}

func (h *handler) ConversionURL(c *ginext.Context) {
	code := c.Param("short_url")
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	longURL, err := h.usecase.GetOriginal(c.Request.Context(), code, ip, userAgent)
	if err != nil {
		c.JSON(http.StatusNotFound, ginext.H{"error": "not found"})
		return
	}

	zlog.Logger.Info().Str("longURL", longURL).Msg("redirect")
	c.Redirect(http.StatusMovedPermanently, longURL)
}

func (h *handler) GetAnalytics(c *ginext.Context) {
	code := c.Param("short_url")

	count, err := h.usecase.GetStats(c.Request.Context(), code)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to get count shorten url")
		c.JSON(http.StatusInternalServerError, ginext.H{"error":"db error"})
		return
	}

	c.JSON(http.StatusOK, ginext.H{"short_code": code, "clicks": count})
}
