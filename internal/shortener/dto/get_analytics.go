package dto

// GetAnalyticsInput represents the input data to retrieve an analytics about specific shorten URL.
type GetAnalyticsInput struct {
	ShortCode string `json:"shorten" binding:"lte=20,gte=6"`
}

// GetAnalyticsOutput represents the detailed analytics about of a usage of shorten URL.
type GetAnalyticsOutput struct {
	TotalClicks int            `json:"total_clicks"`
	ByDate      map[string]int `json:"by_date"`
	ByBrowser   map[string]int `json:"by_browser"`
}
