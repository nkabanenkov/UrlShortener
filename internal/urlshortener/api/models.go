package api

type ShortenRequest struct {
	Url string `json:"url" binding:"required"`
}

type ShortenAnswer struct {
	Message string `json:"message"`
}

type UnshortenAnswer struct {
	Message string `json:"message"`
}
