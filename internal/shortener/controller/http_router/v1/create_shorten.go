package v1

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/wb-go/wbf/ginext"

	"github.com/adexcell/shortener/internal/shortener/controller/http_router/response"
	"github.com/adexcell/shortener/internal/shortener/domain"
	"github.com/adexcell/shortener/internal/shortener/domain/validation"
	"github.com/adexcell/shortener/internal/shortener/dto"
)

// CreateShorten handles the POST request to create new shorten URL.
func (h *Handler) CreateShorten(c *ginext.Context) {
	input := dto.CreateShortenInput{}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			errMap := validation.ExtractErrors(verrs)
			response.ValidationsError(c, errMap)

			return
		}

		response.BadRequest(c, err.Error())
		return
	}

	output, err := h.usecase.CreateShorten(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrShortCodeAlreadyExists) {
			response.Conflict(c, err.Error())
			return
		}
		response.InternalServerError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, output)
}
