package dto

// ShortenTask represents a task to update URL in a cache or store.
type ShortenTask struct {
	Shorten     string `json:"shorten_url"`
	OriginalURL string `json:"OriginalURL"`
}
