package dto

// CreateShortenInput represents the input data required to create a new shorten URL.
type CreateShortenInput struct {
	ShortenCode string `json:"shorten_code,omitempty" binding:"lte=20,gte=6"`
	OriginalURL string `json:"original_url"           binding:"required,url"`
}

// CreateShortenOutput represents the result of a successfull shorten URL creation.
type CreateShortenOutput struct {
	ShortenCode string `json:"shorten_code"`
}
