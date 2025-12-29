package clientsenangpay

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/rs/zerolog/log"
)

type ISubscribeStrategy interface {
	CheckStatusPayment(ctx context.Context, generateUrl string) (any, error)
	GenerateUrl() string
}

type SubscribeStrategy struct {
	Subscribe    ISubscribeStrategy
	NonSubscribe ISubscribeStrategy
}

func NewSubscribeStrategy(subscribe ISubscribeStrategy, nonSubscribe ISubscribeStrategy) *SubscribeStrategy {
	return &SubscribeStrategy{Subscribe: subscribe, NonSubscribe: nonSubscribe}
}

func (s *SubscribeStrategy) SetSubscribe(transactionId string) ISubscribeStrategy {
	if transactionId == "" {
		return s.Subscribe
	}
	return s.NonSubscribe
}

type Subscribe struct {
	senangpay *Senangpay
	request   *CheckStatusPaymentRequest
}

func NewSubscribe(s *Senangpay, request *CheckStatusPaymentRequest) *Subscribe {
	return &Subscribe{
		senangpay: s,
		request:   request,
	}
}

func (s *Subscribe) CheckStatusPayment(ctx context.Context, generateUrl string) (any, error) {
	client := s.senangpay.httpClient.Client

	resp, err := client.R().
		SetContext(ctx).
		Get(generateUrl)
	if err != nil {
		return nil, err
	}

	// Ambil raw response body
	body := resp.Body()
	bodyStr := string(body)

	log.Debug().Str("raw_body", bodyStr).Msg("raw response body")

	// Perbaiki casing field "Note" → "note"
	bodyStr = strings.ReplaceAll(bodyStr, `"Note":`, `"note":`)

	// Parse kembali ke struct
	var result CheckStatusPaymentSenangpayResponse
	if err = json.Unmarshal([]byte(bodyStr), &result); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal senangpay response")
		return nil, err
	}

	log.Debug().Interface("response", result).Msg("response from check status payment senangpay")
	return result, nil
}

func (s *Subscribe) GenerateUrl() string {
	return s.senangpay.GenerateQueryOrderStatusURL(s.request.OrderId)
}

type NonSubscribe struct {
	senangpay *Senangpay
	request   *CheckStatusPaymentRequest
}

func NewNonSubscribe(s *Senangpay, request *CheckStatusPaymentRequest) *NonSubscribe {
	return &NonSubscribe{
		senangpay: s,
		request:   request,
	}
}

func (s *NonSubscribe) CheckStatusPayment(ctx context.Context, generateUrl string) (any, error) {
	client := s.senangpay.httpClient.Client

	resp, err := client.R().
		SetContext(ctx).
		Get(generateUrl)
	if err != nil {
		return nil, err
	}

	// Ambil raw response body
	body := resp.Body()
	bodyStr := string(body)

	log.Debug().Str("raw_body", bodyStr).Msg("raw response body")

	// Perbaiki casing field "Note" → "note"
	bodyStr = strings.ReplaceAll(bodyStr, `"Note":`, `"note":`)

	// Parse kembali ke struct
	var result CheckStatusPaymentSenangpayResponse
	if err = json.Unmarshal([]byte(bodyStr), &result); err != nil {
		log.Error().Err(err).Msg("failed to unmarshal senangpay response")
		return nil, err
	}

	log.Debug().Interface("response", result).Msg("response from check status payment senangpay")
	return result, nil
}

func (s *NonSubscribe) GenerateUrl() string {
	return s.senangpay.GenerateStatusPaymentURL(s.request.TransactionId)
}
