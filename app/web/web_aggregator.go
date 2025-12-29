package web

import (
	"paymentserviceklink/app/enums"
	"time"
)

type CreateAggregatorRequest struct {
	Name        enums.AggregatorName `json:"name"`
	Description string               `json:"description"`
	//Channel       []CreateChannelRequest       `json:"channel"`
	Configuration []CreateConfigurationRequest `json:"configuration"`
	Currency      enums.Currency               `json:"currency"`
}

type AggregatorResponse struct {
	Id          int64                       `json:"id,omitempty"`
	Name        enums.AggregatorName        `json:"name,omitempty"`
	Slug        enums.ProviderPaymentMethod `json:"slug,omitempty"`
	Description string                      `json:"description,omitempty"`
	IsActive    bool                        `json:"isActive"`
	Currency    enums.Currency              `json:"currency,omitempty"`
	CreatedAt   *time.Time                  `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time                  `json:"updatedAt,omitempty"`
	DeletedAt   *time.Time                  `json:"deletedAt,omitempty"`
}
