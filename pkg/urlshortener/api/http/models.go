package http

type ShortenRequest struct {
	Url string `json:"url" binding:"required"`
}

type ShortenResponse struct {
	Message string `json:"message"`
}

type UnshortenResponse struct {
	Message string `json:"message"`
}
