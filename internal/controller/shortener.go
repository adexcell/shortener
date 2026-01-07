// Package controller содержит API http_v1.
package controller

import (
	"errors"
	"net/http"

	"github.com/adexcell/shortener/internal/domain"
	"github.com/adexcell/shortener/pkg/logger"
	"github.com/adexcell/shortener/pkg/router"
)

const (
	postShortURL  = "/shorten"
	conversionURL = "/s/:short_url"
	analyticsURL  = "/analytics/:short_url"
)

type handler struct {
	usecase domain.ShortenerUsecase
	log     logger.Log
}

func NewShortenHandler(u domain.ShortenerUsecase, l logger.Log) router.Handler {
	return &handler{usecase: u, log: l}
}

func (h *handler) Register(router *router.Router) {
	router.POST(postShortURL, h.PostShortURL)
	router.GET(conversionURL, h.ConversionURL)
	router.GET(analyticsURL, h.GetAnalytics)
}

type shortenRequest struct {
	// URL - полная ссылка
	URL string `json:"url" binding:"required"`
	// кастомное имя сокращенной ссылки
	Alias string `json:"alias"`
}

func (h *handler) PostShortURL(c *router.Context) {
	var req shortenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, router.H{"error": "invalid request"})
		return
	}

	dto, err := shortenerToControllerDTO(req.Alias, req.URL)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, router.H{"error": err.Error()})
		return
	}

	code, err := h.usecase.Shorten(c.Request.Context(), dto.ShortCode, dto.LongURL)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			c.JSON(http.StatusConflict, router.H{"error": domain.ErrAlreadyExists})
			return
		}
		h.log.Error().Err(err).Msg("failed to shorten url")
		c.JSON(http.StatusInternalServerError, router.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, router.H{"short_url": code})
}

func (h *handler) ConversionURL(c *router.Context) {
	code := c.Param("short_url")
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	dto, err := statsToControllerDTO(code, ip, userAgent)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, router.H{"error": err.Error()})
		return
	}

	longURL, err := h.usecase.GetOriginal(
		c.Request.Context(),
		dto.ShortCode,
		dto.IP,
		dto.UserAgent,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, router.H{"error": "not found"})
		return
	}

	h.log.Info().Str("longURL", longURL).Msg("redirect")
	c.Redirect(http.StatusFound, longURL)
}

func (h *handler) GetAnalytics(c *router.Context) {
	code := c.Param("short_url")

	var stats domain.Stats
	stats, err := h.usecase.GetStats(c.Request.Context(), code)
	if err != nil {
		h.log.Error().Err(err).Msg("failed to get count shorten url")
		c.JSON(http.StatusInternalServerError, router.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, statsToResponse(stats))
}
