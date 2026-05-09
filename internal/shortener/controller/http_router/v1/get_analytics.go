package v1

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"

	"github.com/adexcell/shortener/internal/shortener/controller/http_router/response"
	"github.com/adexcell/shortener/internal/shortener/dto"
)

// GetAnalytics handles the GET request to retrieve analytics about shorten URL usage.
func (h *Handler) GetAnalytics(c *ginext.Context) {
	input := dto.GetAnalyticsInput{ShortCode: c.Param("short_url")}

	output, err := h.usecase.GetAnalytics(c.Request.Context(), input)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, output)
}
