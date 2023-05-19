package server

type ShortenRequest struct {
	URL        string `json:"url" validate:"required,url"`
	Identifier string `json:"identifier,omitempty" validate:"omitempty,alphanum"`
}

type ShortenResponce struct {
	ShortURL string `json:"short_url,omitempty"`
	Message  string `json:"message,omitempty"`
}
