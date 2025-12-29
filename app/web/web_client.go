package web

type ClientRequest struct {
	ApiKey    string `json:"apiKey"`
	SecretKey string `json:"secretKey"`
}
