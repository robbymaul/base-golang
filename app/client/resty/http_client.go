package restyclient

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

type RestyClient struct {
	Client *resty.Client
}

func NewRestyClient() *RestyClient {
	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(2).
		SetTransport(&http.Transport{
			IdleConnTimeout:     60 * time.Second,
			TLSHandshakeTimeout: 30 * time.Second,
		})

	return &RestyClient{Client: client}
}
