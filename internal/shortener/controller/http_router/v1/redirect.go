package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"

	"github.com/adexcell/shortener/internal/shortener/controller/http_router/response"
	"github.com/adexcell/shortener/internal/shortener/domain"
	"github.com/adexcell/shortener/internal/shortener/dto"
)

// RedirectShortLink handles the GET request to retrieve analytics aboute shorten URL usage.
func (h *Handler) RedirectShortLink(c *ginext.Context) {
	input := dto.RedirectInput{
		ShortenCode: c.Param("short_url"),
		IP:          c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
		ClickedAt:   time.Now().UTC(),
	}

	location, err := h.usecase.GetOriginalURL(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrShortNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		
		response.InternalServerError(c, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, location)
}
