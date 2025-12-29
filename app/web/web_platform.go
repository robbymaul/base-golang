package web

import (
	"time"
)

type CreatePlatformRequest struct {
	//Code          string                   `json:"code"`
	Name          string                   `json:"name"`
	Description   string                   `json:"description"`
	Channel       []*DetailChannelResponse `json:"channel"`
	Configuration []*ResponseConfiguration `json:"configuration"`
}

type DetailPlatformResponse struct {
	Id              int64      `json:"id"`
	Code            string     `json:"code,omitempty"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	ApiKey          string     `json:"apiKey,omitempty"`
	SecretKey       string     `json:"secretKey,omitempty"`
	IsActive        bool       `json:"isActive"`
	NotificationUrl string     `json:"notificationUrl"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
}
