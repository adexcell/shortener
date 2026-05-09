package httprouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wb-go/wbf/ginext"

	ver1 "github.com/adexcell/shortener/internal/shortener/controller/http_router/v1"
	"github.com/adexcell/shortener/pkg/logger"
	"github.com/adexcell/shortener/pkg/metrics"
	"github.com/adexcell/shortener/pkg/otel"
)

// NewRouter creates and configures the HTTP router with middleware and routes.
func ShortenRouter(r *ginext.Engine, uc ver1.Usecase, m *metrics.HTTPServer) {
	v1 := ver1.New(uc)

	r.StaticFS("/static", http.Dir("./web"))

	// Expose metrics endpoint (separate from application routes)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	api := r.Group("api/v1")
	{
		api.Use(logger.Middleware())
		api.Use(metrics.NewMiddleware(m))
		api.Use(otel.Middleware())

		api.POST("/shorten", v1.CreateShorten)
		api.GET("/s/:short_url", v1.RedirectShortLink)
		api.GET("/analytics/:short_url", v1.GetAnalytics)
	}

}
