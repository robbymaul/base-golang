package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"paymentserviceklink/app/client/espay"
	"paymentserviceklink/app/client/midtrans"
	clientsenangpay "paymentserviceklink/app/client/senangpay"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/repositories"
	"paymentserviceklink/app/strategy"
	"paymentserviceklink/app/web"
	"paymentserviceklink/config"
	"paymentserviceklink/pkg/pagination"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type IPaymentStrategy interface {
	CreatePayment(ctx context.Context) (*web.PaymentResponse, error)
}

type PaymentStrategyKWallet struct {
	paymentService *PaymentService
	payload        *web.CreatePaymentRequest
	payment        *web.Payment
	channel        *models.Channel
	platform       *models.Platforms
	noRekening     string
	strategy       string
}

type PaymentStrategyPG struct {
	paymentService *PaymentService
	channel        *models.Channel
	platform       *models.Platforms
	payload        *web.CreatePaymentRequest
	payment        *web.Payment
	strategy       string
}

func (s *PaymentService) StrategyPayment(payload *web.CreatePaymentRequest, payment *web.Payment, platform *models.Platforms, channel *models.Channel, noRekening string) (IPaymentStrategy, error) {
	switch channel.PaymentMethod {
	case enums.PAYMENT_METEHOD_K_WALLET:
		return NewPaymentStrategyKWallet(s, payload, payment, channel, platform, noRekening), nil
	//case enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT:
	//	return NewPaymentStrategyVa(s, payload, payment, platform, channel), nil
	//case enums.PAYMENT_METHOD_CREDIT_CARD:
	//	return NewPaymentStrategyCreditCard(s, payload, payment, platform, channel), nil
	default:
		//return nil, fmt.Errorf("payment method %s not implemented", channel.PaymentMethod)
		return NewPaymentStrategyPG(s, payload, payment, platform, channel), nil
	}
}

type PaymentService struct {
	service     *Service
	serviceName string
}

func NewPaymentService(ctx context.Context, repo *repositories.RepositoryContext, cfg *config.Config) *PaymentService {
	return &PaymentService{
		service:     NewService(ctx, repo, cfg),
		serviceName: "PaymentService",
	}
}

func (s *PaymentService) CreatePaymentService(payload *web.CreatePaymentRequest) ([]*web.PaymentResponse, error) {
	var err error
	webResponses := make([]*web.PaymentResponse, 0)
	s.serviceName = "PaymentService.CreatePaymentService"

	log.Debug().Interface("payload", payload).Interface("context", s.serviceName).Msg("create payment service")

	platform, err := GetClientAuth(s.service)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Interface("context", s.serviceName).Msg("get client auth error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusUnauthorized)
	}
	log.Debug().Interface("platform", platform).Interface("context", s.serviceName).Msg("get client auth")

	for _, pay := range payload.Payment {
		totalRequested := decimal.NewFromInt(pay.Amount)
		totalProcessed := decimal.NewFromInt(0)

		for _, ch := range pay.Channel {
			// get channel
			channel, err := s.service.repository.GetChannelByChannelIdAndPlatformIdRepository(s.service.ctx, ch.Id, platform.Id)
			if err != nil {
				log.Error().Err(err).Str("context", s.serviceName).Msg("get payment method by id repository error")
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, helpers.NewErrorTrace(fmt.Errorf("%v, chnnel", err), s.serviceName).WithStatusCode(http.StatusNotFound)
				}

				return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
			}
			log.Debug().Interface("channel", channel).Msg("get channel by id repository")

			// Default amount (kalau tidak dikirim di payload)
			amount := decimal.NewFromInt(0)
			//if decimal.NewFromInt(ch.Amount) != decimal.Zero {
			//	amount = decimal.NewFromInt(ch.Amount)
			//} else {
			//	// fallback: bagi rata atau total sisa
			//	amount = totalRequested.Sub(totalProcessed)
			//}

			if totalProcessed != totalRequested {
				amount = decimal.NewFromInt(ch.Amount)
			}

			if totalProcessed == totalRequested {
				break
			}

			switch channel.PaymentMethod {
			case enums.PAYMENT_METEHOD_K_WALLET:
				pay.Amount = amount.IntPart()
				strategyKWallet := NewPaymentStrategyKWallet(s, payload, pay, channel, platform, ch.NoRekening)

				webResponse, err := strategyKWallet.CreatePayment(s.service.ctx)
				if err != nil {
					return nil, err
				}

				totalProcessed = totalProcessed.Add(amount)
				webResponses = append(webResponses, webResponse)
			default:
				pay.Amount = amount.IntPart()
				// strategy payment channel method
				strategy, err := s.StrategyPayment(payload, pay, platform, channel, ch.NoRekening)
				if err != nil {
					log.Error().Err(err).Str("context", s.serviceName).Msg("strategy payment channel method error")
					return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
				}
				log.Debug().Interface("strategy", strategy).Msg("strategy payment channel method")

				webResponse, err := strategy.CreatePayment(s.service.ctx)
				if err != nil {
					return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
				}

				totalProcessed = amount
				webResponses = append(webResponses, webResponse)
			}

		}
	}

	return webResponses, nil
}

func (s *PaymentService) CheckStatusPaymentService(payload *web.CheckStatusPaymentRequest) (*web.CheckStatusPaymentResponse, error) {
	s.serviceName = "PaymentService.CheckStatusPayment"

	log.Debug().Interface("payload", payload).Msg("check status payment request")

	platform, err := GetClientAuth(s.service)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get client auth error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusUnauthorized)
	}
	log.Debug().Interface("platform", platform).Msg("get client auth")

	payment, err := s.service.repository.GetPaymentByTransactionIdAndPlatformIdRepository(s.service.ctx, payload.OrderId, platform.Id)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get payment by transaction id repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(fmt.Errorf("%v, payment transaction", err), s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("payment", payment).Msg("get payment by transaction id repository")

	if payment.Status == enums.PAYMENT_STATUS_SUCCESS || payment.Status == enums.PAYMENT_STATUS_EXPIRED || payment.Status == enums.PAYMENT_STATUS_FAILED {
		return s.service.mapToCheckStatusPaymentResponse(payment), nil
	}

	if payment.Channel == nil {
		return nil, helpers.NewErrorTrace(fmt.Errorf("payment method is missing, payment notification"), s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	configuration, err := s.service.repository.GetConfigurationByPlatformIdAndAggregatorIdRepository(s.service.ctx, platform.Id, *payment.AggregatorId)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get configuration by platform id and aggregator id repository")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(fmt.Errorf("%v, platform configuration", err), s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("configuration", configuration).Interface("context", s.serviceName).Msg("get platform repository")

	s.service.repository.SetConfigurationPayment(configuration)

	// get strategy
	strategy, err := s.service.repository.Strategy.GetStrategy(configuration)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get strategy payment error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	//var responseData *models.CallbackData
	var request any

	if payment.Aggregator.Slug == enums.PROVIDER_PAYMENT_METHOD_SENANGPAY {
		//responseData, err = s.service.repository.CheckStatusPaymentSenangpayRepository(s.service.ctx, payment.ReferenceId)
		//if err != nil {
		//	return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		//}

		//s.updatePaymentSenangpay(&payment, &web.WebhookCallbackSenangpay{
		//	Name:          responseData.Name,
		//	Email:         responseData.Email,
		//	Phone:         responseData.Phone,
		//	AmountPaid:    responseData.AmountPaid,
		//	TxnStatus:     responseData.TxnStatus,
		//	TxnMessage:    responseData.TxnMessage,
		//	OrderId:       responseData.OrderId,
		//	TransactionId: responseData.TransactionId,
		//	HashedValue:   responseData.HashedValue,
		//})
		request = clientsenangpay.CheckStatusPaymentRequest{
			TransactionId: payment.ReferenceId,
			OrderId:       payment.OrderId,
		}
	} else if payment.Aggregator.Slug == enums.PROVIDER_PAYMENT_METHOD_MIDTRANS {
		request = midtrans.CheckStatusPaymentRequest{
			TransactionId: payment.OrderId,
			OrderId:       payment.OrderId,
		}

	} else {
		return nil, helpers.NewErrorTrace(fmt.Errorf("payment method provider is missing, payment notification"), s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("request", request).Interface("context", s.serviceName).Msg("check status payment")

	statusPayment, err := strategy.CheckStatusPayment(s.service.ctx, request)
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("check status payment error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("statusPayment", statusPayment).Interface("context", s.serviceName).Msg("check status payment")

	//if responseData == nil {
	//	return nil, helpers.NewErrorTrace(fmt.Errorf("response data is missing, payment notification"), s.serviceName).WithStatusCode(http.StatusInternalServerError)
	//}

	dataTypesJson, err := helpers.ConvertAnyToDatatypeJson(statusPayment)
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("convert amy to datatype json error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("dataTypesJson", dataTypesJson).Msg("convert amy to datatype json")

	result, err := strategy.MapCheckStatusPayment(payment, statusPayment)
	if err != nil {
		log.Error().Err(err).Msg("map check status payment error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("result", result).Interface("context", s.serviceName).Msg("map check status payment")

	s.updatePaymentStatus(&payment, result)

	paymentCallback := models.PaymentCallbacks{
		PaymentId:    payment.Id,
		GatewayName:  payment.GatewayReference,
		CallbackData: dataTypesJson,
		ResponseData: nil,
		Status:       payment.Status,
	}

	paymentStatusHistory := models.PaymentStatusHistory{
		PaymentId: payment.Id,
		Status:    payment.Status,
		Notes:     "",
		CreatedBy: models.CreatedBy{
			ID:       platform.ApiKey,
			Name:     platform.SecretKey,
			Role:     "platform",
			Platform: platform.Name,
		},
	}

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
		err = s.service.repository.UpdatePaymentRepositoryTx(s.service.ctx, tx, func(db *gorm.DB) *gorm.DB {
			db = db.Where("payments.id = ?", payment.Id)
			updateColumn := map[string]interface{}{
				"status":       payment.Status,
				"paid_at":      payment.PaidAt,
				"reference_id": payment.ReferenceId,
				"updated_at":   time.Now(),
			}
			db = db.UpdateColumns(updateColumn)
			return db
		})
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("update payment repository error")
			return helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}

		err = s.service.repository.InsertPaymentCallbackRepositoryTx(s.service.ctx, tx, &paymentCallback)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("insert payment callback repository error")
			return helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}

		err = s.service.repository.InsertPaymentStatusHistoryRepositoryTx(s.service.ctx, tx, &paymentStatusHistory)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("insert payment status history repository error")
			return helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}

		return nil
	})

	return result, err
}

func (s *PaymentService) CallbackSenangpayNotificationService(payload *web.WebhookCallbackSenangpay) error {
	payment, err := s.service.repository.GetPaymentByOrderIdRepository(s.service.ctx, payload.OrderId)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get payment by order id repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helpers.NewErrorTrace(fmt.Errorf("%v, payment notification", err), s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	s.updatePaymentSenangpay(&payment, payload)

	paymentCallback := models.PaymentCallbacks{
		PaymentId:   payment.Id,
		GatewayName: payment.GatewayReference,
		//CallbackData: models.CallbackData{
		//	Name:          payload.Name,
		//	Email:         payload.Email,
		//	Phone:         payload.Phone,
		//	AmountPaid:    payload.AmountPaid,
		//	TxnStatus:     payload.TxnStatus,
		//	TxnMessage:    payload.TxnMessage,
		//	OrderId:       payload.OrderId,
		//	TransactionId: payload.TransactionId,
		//	HashedValue:   payload.HashedValue,
		//},
		CallbackData: nil,
		ResponseData: nil,
		Status:       payment.Status,
	}

	paymentStatusHistory := models.PaymentStatusHistory{
		PaymentId: payment.Id,
		Status:    payment.Status,
		Notes:     "",
		CreatedBy: models.CreatedBy{
			ID:       "system",
			Name:     "callback notification senangpay",
			Role:     "system",
			Platform: "callback_senangpay",
		},
	}

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
		err = s.service.repository.UpdatePaymentRepositoryTx(s.service.ctx, tx, func(db *gorm.DB) *gorm.DB {
			db = db.Where("payments.id = ?", payment.Id)
			updateColumn := map[string]interface{}{
				"status":       payment.Status,
				"paid_at":      payment.PaidAt,
				"reference_id": payment.ReferenceId,
				"updated_at":   time.Now(),
			}
			db = db.UpdateColumns(updateColumn)
			return db
		})
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("update payment repository error")
			return err
		}

		err = s.service.repository.InsertPaymentCallbackRepositoryTx(s.service.ctx, tx, &paymentCallback)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("insert payment callback repository error")
			return err
		}

		err = s.service.repository.InsertPaymentStatusHistoryRepositoryTx(s.service.ctx, tx, &paymentStatusHistory)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("insert payment status history repository error")
			return err
		}

		return nil
	})

	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("update payment status history error")
		return helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	callback := web.PaymentCallback{
		TransactionId: payment.TransactionId,
		OrderId:       payment.OrderId,
		Status:        payment.Status,
		Amount:        payment.TotalAmount.IntPart(),
		Currency:      payment.Currency,
	}

	defer func() {
		err = s.service.repository.CallbackFunctionRepository(s.service.ctx, payment.Platform.NotificationURL, &callback)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("callback function repository error")
			return
		}

		err = s.service.repository.UpdatePaymentNotificationRepository(s.service.ctx, payment)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("update payment repository error")
			return
		}
	}()

	return nil
}

func (s *PaymentService) CheckKeyMidtransService() (any, error) {
	configuration := models.Configuration{
		Id:           0,
		AggregatorId: 0,
		ConfigKey:    "",
		ConfigValue:  "",
		ConfigName:   "",
		ConfigJson:   models.ConfigJson{},
		IsActive:     false,
		Aggregator: &models.Aggregator{
			Id:          0,
			Name:        "",
			Slug:        "midtrans",
			Description: "",
			IsActive:    false,
			Currency:    "",
			BaseField:   models.BaseField{},
		},
		BaseField: models.BaseField{},
	}

	strategy, err := s.service.repository.Strategy.GetStrategy(&configuration)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get strategy error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return strategy.CheckKey()

}

func (s *PaymentService) EspayValidationInquiryService(payload *web.EspayInquiryRequest) *web.EspayInquiryResponse {
	log.Debug().Interface("payload", payload).Msg("payment service espay validation inquiry service")

	// get transaction with order id
	orderId := payload.VirtualAccountNo
	payment, err := s.service.repository.GetPaymentByOrderIdRepository(s.service.ctx, orderId)
	if err != nil {
		log.Error().Err(err).Msg("get payment by order id repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.mapToEspayInquiryResponseError(http.StatusNotFound, enums.INQUIRY, enums.ESPAY_TRANSACTION_NOT_FOUND, err)
		}

		return s.mapToEspayInquiryResponseError(http.StatusInternalServerError, enums.INQUIRY, enums.ESPAY_INTERNAL_SERVER_ERROR, err)
	}

	return s.mapToEspayInquiryResponseSuccess(payment, payload)
}

func (s *PaymentService) EspayPaymentNotificationService(payload *web.EspayPaymentNotificationRequest) *web.EspayPaymentNotificationResponse {
	log.Debug().Interface("payload", payload).Msg("payment service espay payment notification service")
	payment, err := s.service.repository.GetPaymentByOrderIdRepository(s.service.ctx, payload.VirtualAccountNo)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get payment by order id repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.mapToEspayNotificationResponseError(payload, http.StatusNotFound, enums.ESPAY_INTERNAL_SERVER_ERROR, err)
		}

		return s.mapToEspayNotificationResponseError(payload, http.StatusInternalServerError, enums.ESPAY_INTERNAL_SERVER_ERROR, err)
	}

	statusPayment := &web.CheckStatusPaymentResponse{
		Id:            payment.Id,
		TransactionId: payment.TransactionId,
		OrderId:       payment.OrderId,
		Status:        enums.PAYMENT_STATUS_SUCCESS,
		Amount:        payment.Amount.IntPart(),
		Currency:      payment.Currency,
	}

	s.updatePaymentStatus(&payment, statusPayment)

	jsonPayload, _ := json.Marshal(payload)

	paymentCallback := models.PaymentCallbacks{
		PaymentId:    payment.Id,
		GatewayName:  payment.GatewayReference,
		CallbackData: jsonPayload,
		ResponseData: nil,
		Status:       payment.Status,
	}

	paymentStatusHistory := models.PaymentStatusHistory{
		PaymentId: payment.Id,
		Status:    payment.Status,
		Notes:     "",
		CreatedBy: models.CreatedBy{
			ID:       "system",
			Name:     "callback notification espay",
			Role:     "system",
			Platform: "callback_espay",
		},
	}

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
		err = s.service.repository.UpdatePaymentRepositoryTx(s.service.ctx, tx, func(db *gorm.DB) *gorm.DB {
			db = db.Where("payments.id = ?", payment.Id)
			updateColumn := map[string]interface{}{
				"status":       payment.Status,
				"paid_at":      payment.PaidAt,
				"reference_id": payment.ReferenceId,
				"updated_at":   time.Now(),
			}
			db = db.UpdateColumns(updateColumn)
			return db
		})
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("update payment repository error")
			return err
		}

		err = s.service.repository.InsertPaymentCallbackRepositoryTx(s.service.ctx, tx, &paymentCallback)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("insert payment callback repository error")
			return err
		}

		err = s.service.repository.InsertPaymentStatusHistoryRepositoryTx(s.service.ctx, tx, &paymentStatusHistory)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("insert payment status history repository error")
			return err
		}

		return nil
	})
	if err != nil {
		return s.mapToEspayNotificationResponseError(payload, http.StatusInternalServerError, enums.ESPAY_INTERNAL_SERVER_ERROR, err)
	}

	callback := web.PaymentCallback{
		TransactionId: payment.TransactionId,
		OrderId:       payment.OrderId,
		Status:        payment.Status,
		Amount:        payment.TotalAmount.IntPart(),
		Currency:      payment.Currency,
	}

	defer func() {
		err = s.service.repository.CallbackFunctionRepository(s.service.ctx, payment.Platform.NotificationURL, &callback)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("callback function repository error")
			return
		}

		err = s.service.repository.UpdatePaymentNotificationRepository(s.service.ctx, payment)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("update payment repository error")
			return
		}
	}()

	return s.mapToEspayNotificationResponseSuccess(payload, payment)
}

func (s *PaymentService) EspayTopupNotificationService(payload *web.EspayTopupNotificationRequest) *web.EspayPaymentNotificationResponse {
	log.Debug().Interface("payload", payload).Msg("payment service espay payment notification service")
	//filterKwallet := []*pagination.Filter{
	//	{
	//		ID:       "gen_va",
	//		Value:    payload.AdditionalInfo.UserId,
	//		Variant:  "string",
	//		Operator: "eq",
	//		FilterID: "",
	//	},
	//}

	//filterKWallet, err := helpers.FilterColumnValidation(filterKwallet, models.AllowedFilterColumnKWallet())
	//if err != nil {
	//	log.Error().Err(err).Str("context", s.serviceName).Msg("filter column validation error")
	//	return s.mapToEspayTopupNotificationResponseError(payload, http.StatusInternalServerError, enums.ESPAY_INTERNAL_SERVER_ERROR, err)
	//}

	kWallet, err := s.service.repository.GetKWalletRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Where("k_wallet.gen_va = ?", payload.AdditionalInfo.UserId).Or("k_wallet.gen_va = ?", payload.AdditionalInfo.DebitFrom).
			Or("k_wallet.gen_va = ?", payload.AdditionalInfo.DebitFromName).
			Or("k_wallet.gen_va = ?", payload.AdditionalInfo.ProductValue)
		//query, args := s.service.repository.SearchQuery(filterKWallet, "and")
		//if query != "" {
		//	db = db.Where(query, args)
		//}
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get kwallet repository error")
		return s.mapToEspayTopupNotificationResponseError(payload, http.StatusInternalServerError, enums.ESPAY_INTERNAL_SERVER_ERROR, err)
	}

	//filterChannel := []*pagination.Filter{
	//	{
	//		ID:       "bank_code",
	//		Value:    payload.AdditionalInfo.DebitFromBank,
	//		Variant:  "string",
	//		Operator: "eq",
	//		FilterID: "",
	//	},
	//	{
	//		ID:       "payment_method",
	//		Value:    enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT,
	//		Variant:  "string",
	//		Operator: "eq",
	//		FilterID: "",
	//	},
	//}
	//
	//filterCh, err := helpers.FilterColumnValidation(filterChannel, models.AllowedFilterColumnChannel())
	//if err != nil {
	//	log.Error().Err(err).Str("context", s.serviceName).Msg("filter column validation error")
	//	return s.mapToEspayTopupNotificationResponseError(payload, http.StatusInternalServerError, enums.ESPAY_INTERNAL_SERVER_ERROR, err)
	//}

	channel, err := s.service.repository.GetChannelRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Where(
			"channels.bank_code = ? and channels.product_code = ? and channels.payment_method = ?",
			payload.AdditionalInfo.DebitFromBank,
			payload.AdditionalInfo.ProductCode,
			enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT,
		)
		//query, args := s.service.repository.SearchQuery(filterCh, "and")
		//if query != "" {
		//	db = db.Where(query, args)
		//}
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get channel repository error")
		return s.mapToEspayTopupNotificationResponseError(payload, http.StatusInternalServerError, enums.ESPAY_INTERNAL_SERVER_ERROR, err)
	}

	amount, err := decimal.NewFromString(payload.TotalAmount.Value)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("new from string error")
		return s.mapToEspayTopupNotificationResponseError(payload, http.StatusInternalServerError, enums.ESPAY_INTERNAL_SERVER_ERROR, err)
	}
	//totalAmount, _ := strconv.ParseFloat(payload.TotalAmount, 32)
	txFee, err := decimal.NewFromString(payload.AdditionalInfo.TxFee)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("new from string error")
		return s.mapToEspayTopupNotificationResponseError(payload, http.StatusInternalServerError, enums.ESPAY_INTERNAL_SERVER_ERROR, err)
	}

	_, txFee = s.calculatePayment(amount, channel)

	//s.updateTopupKWallet(&kWallet, payload)

	kWallet.AddBalance(amount.Sub(txFee))

	totalAmount := amount.Sub(txFee)

	topupTransaction := models.TopupTransaction{
		KWalletID:   kWallet.ID,
		MemberId:    kWallet.MemberID,
		ChannelID:   channel.Id,
		Aggregator:  enums.AGGREGATOR_NAME_ESPAY,
		Merchant:    payload.AdditionalInfo.MemberCode,
		Amount:      amount,
		FeeAdmin:    txFee,
		Currency:    enums.Currency(payload.PaidAmount.Currency),
		Symbol:      enums.SYMBOL_CURRENCY_IDR,
		ReferenceID: payload.TrxId,
		Status:      enums.PAYMENT_STATUS_SUCCESS,
		CompletedAt: time.Now(),
		Description: "",
	}

	kWalletTransaction := &models.KWalletTransaction{
		KWalletID:                kWallet.ID,
		KWalletTypeTransactionID: 1,
		PaymentID:                payload.TrxId,
		PaymentCode:              fmt.Sprint(time.Now().Unix()),
		TransactionCode:          fmt.Sprint(time.Now().Unix()),
		TransactionType:          "",
		CounterpartyName:         kWallet.FullName,
		CounterpartyBank:         channel.BankName,
		PaymentChannel:           kWallet.NoRekening,
		Description:              fmt.Sprint(payload.AdditionalInfo.Message),
		Currency:                 enums.CURRENCY_IDR,
		Symbol:                   enums.SYMBOL_CURRENCY_IDR,
		Status:                   enums.PAYMENT_STATUS_SUCCESS,
	}

	kWalletTransaction.TopupTransaction(kWallet.Balance, totalAmount)

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
		// insert topup transaction
		err = s.service.repository.InsertTopupTransactionRepositoryTx(s.service.ctx, tx, topupTransaction)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("insert topup transaction repository error")
			return err
		}

		// insert k-wallet transaction
		err = s.service.repository.InsertKWalletTransactionRepositoryTx(s.service.ctx, tx, kWalletTransaction)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("insert k-wallet transaction repository error")
			return err
		}

		// update k-wallet balance
		err = s.service.repository.UpdateKWalletRepositoryTx(s.service.ctx, tx, func(db *gorm.DB) *gorm.DB {
			db = db.Where("k_wallet.id = ?", kWallet.ID)
			updateColumn := map[string]interface{}{
				"balance":    kWallet.Balance,
				"updated_at": time.Now(),
			}
			db = db.UpdateColumns(updateColumn)
			return db
		})
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("update k-wallet repository error")
			return err
		}

		return nil
	})
	if err != nil {
		return s.mapToEspayTopupNotificationResponseError(payload, http.StatusInternalServerError, enums.ESPAY_INTERNAL_SERVER_ERROR, err)
	}

	return s.mapToEspayTopupNotificationResponseSuccess(payload, kWalletTransaction)
}

func (s *PaymentService) calculatePayment(amount decimal.Decimal, channel *models.Channel) (totalAmount decimal.Decimal, feeAmount decimal.Decimal) {
	switch channel.FeeType {
	case enums.FEE_TYPE_PERCENTAGE:
		return helpers.CalculatePercentage(amount, decimal.NewFromFloat32(channel.FeePercentage))
	case enums.FEE_TYPE_FIXED:
		return amount.Add(decimal.NewFromInt(channel.FeeAmount)), decimal.NewFromInt(channel.FeeAmount)
	case enums.FEE_TYPE_FIXED_PERCENTAGE:
		amountPercentage, feePercentage := helpers.CalculatePercentage(amount, decimal.NewFromFloat32(channel.FeePercentage))
		feeFix := decimal.NewFromInt(channel.FeeAmount)
		return amountPercentage.Add(feeFix), feePercentage.Add(feeFix)
	default:
		return amount, decimal.NewFromInt(channel.FeeAmount)
	}
}

func (s *PaymentService) updatePaymentSenangpay(payment **models.Payments, payload *web.WebhookCallbackSenangpay) enums.PaymentStatus {
	if payload.TxnStatus == enums.SENANGPAY_STATUS_SUCCESS {
		(*payment).Status = enums.PAYMENT_STATUS_SUCCESS
	} else if payload.TxnStatus == enums.SENANGPAY_STATUS_FAILED {
		(*payment).Status = enums.PAYMENT_STATUS_FAILED
	}

	(*payment).ReferenceId = payload.TransactionId

	now := time.Now()
	(*payment).PaidAt = &now

	return (*payment).Status
}

func (s *PaymentService) updatePaymentStatus(payment **models.Payments, result *web.CheckStatusPaymentResponse) {
	log.Debug().Interface("result", result).Msg("update payment status")
	paidAt := time.Now()
	(*payment).Status = result.Status
	(*payment).PaidAt = &paidAt
}

func (s *PaymentService) mapToEspayInquiryResponseSuccess(payment *models.Payments, payload *web.EspayInquiryRequest) *web.EspayInquiryResponse {
	espayCode, message := enums.CreateEspayCodeResponse(http.StatusOK, enums.INQUIRY, enums.ESPAY_SUCCESSFULL, nil)

	billsData := make([]*web.EspayVirtualAccountBillDetails, 0)

	billsData = append(billsData, &web.EspayVirtualAccountBillDetails{
		BillDescription: &web.EspayVirtualAccountBillDetailsDescription{
			English:   fmt.Sprintf("Invoice No %v", payment.OrderId),
			Indonesia: fmt.Sprintf("Tagihan No %v", payment.OrderId),
		},
	})

	return &web.EspayInquiryResponse{
		ResponseCode:    espayCode,
		ResponseMessage: message,
		VirtualAccountData: &web.EspayVirtualAccountData{
			PartnerServiceId:    payload.PartnerServiceId,
			CustomerNo:          payload.CustomerNo,
			VirtualAccountNo:    payment.OrderId,
			VirtualAccountName:  payment.CustomerName,
			VirtualAccountEmail: payment.CustomerEmail,
			VirtualAccountPhone: payment.CustomerPhone,
			InquiryRequestId:    payload.InquiryRequestId,
			TotalAmount: &web.EspayVirtualAccountTotalAmount{
				Value:    fmt.Sprint(payment.Amount),
				Currency: string(payment.Currency),
			},
			BillDetails:    billsData,
			AdditionalInfo: nil,
		},
	}
}

func (s *PaymentService) mapToEspayInquiryResponseError(httpCode int, serviceCode enums.EspayService, caseCode enums.EspayCaseCode, err error) *web.EspayInquiryResponse {
	espayCode, message := enums.CreateEspayCodeResponse(httpCode, serviceCode, caseCode, err)

	return &web.EspayInquiryResponse{
		ResponseCode:       espayCode,
		ResponseMessage:    message,
		VirtualAccountData: nil,
	}
}

func (s *PaymentService) mapToEspayNotificationResponseError(payload *web.EspayPaymentNotificationRequest, httpCode int, caseCode enums.EspayCaseCode, err error) *web.EspayPaymentNotificationResponse {
	espayCode, message := enums.CreateEspayCodeResponse(httpCode, enums.PAYMENT, caseCode, err)
	return &web.EspayPaymentNotificationResponse{
		ResponseCode:    espayCode,
		ResponseMessage: message,
	}
}

func (s *PaymentService) mapToEspayNotificationResponseSuccess(payload *web.EspayPaymentNotificationRequest, payment *models.Payments) *web.EspayPaymentNotificationResponse {
	espayCode, message := enums.CreateEspayCodeResponse(http.StatusOK, enums.PAYMENT, enums.ESPAY_SUCCESSFULL, nil)

	return &web.EspayPaymentNotificationResponse{
		ResponseCode:    espayCode,
		ResponseMessage: message,
		VirtualAccountData: struct {
			PartnerServiceId   string `json:"partnerServiceId"`
			CustomerNo         string `json:"customerNo"`
			VirtualAccountNo   string `json:"virtualAccountNo"`
			VirtualAccountName string `json:"virtualAccountName"`
			PaymentRequestId   string `json:"paymentRequestId"`
			TotalAmount        struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			} `json:"totalAmount"`
			BillDetails []struct {
				BillDescription struct {
					English   string `json:"english"`
					Indonesia string `json:"indonesia"`
				} `json:"billDescription"`
			} `json:"billDetails"`
		}{
			PartnerServiceId:   payload.PartnerServiceId,
			CustomerNo:         payload.CustomerNo,
			VirtualAccountNo:   payload.VirtualAccountNo,
			VirtualAccountName: payment.CustomerName,
			PaymentRequestId:   payload.PaymentRequestId,
			TotalAmount: struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			}{
				Value:    payload.TotalAmount.Value,
				Currency: payload.TotalAmount.Currency,
			},
			BillDetails: nil,
		},
		AdditionalInfo: struct {
			ReconcileId       string    `json:"reconcileId"`
			ReconcileDatetime time.Time `json:"reconcileDatetime"`
		}{
			ReconcileId:       payment.TransactionId,
			ReconcileDatetime: *payment.PaidAt,
		},
	}
}

func (s *PaymentService) mapToEspayTopupNotificationResponseError(payload *web.EspayTopupNotificationRequest, httpCode int, caseCode enums.EspayCaseCode, err error) *web.EspayPaymentNotificationResponse {
	espayCode, message := enums.CreateEspayCodeResponse(httpCode, enums.PAYMENT, caseCode, err)
	return &web.EspayPaymentNotificationResponse{
		ResponseCode:    espayCode,
		ResponseMessage: message,
	}
}

//func (s *PaymentService) updateTopupKWallet(kWallet **models.KWallet, payload *web.EspayTopupNotificationRequest) {
//	(*kWallet).Balance += float64(payload.Amount)
//	//(*kWallet).UpdatedAt = time.Now()
//}

func (s *PaymentService) mapToEspayTopupNotificationResponseSuccess(payload *web.EspayTopupNotificationRequest, kWalletTransaction *models.KWalletTransaction) *web.EspayPaymentNotificationResponse {
	espayCode, message := enums.CreateEspayCodeResponse(http.StatusOK, enums.PAYMENT, enums.ESPAY_SUCCESSFULL, nil)

	return &web.EspayPaymentNotificationResponse{
		ResponseCode:    espayCode,
		ResponseMessage: message,
		VirtualAccountData: struct {
			PartnerServiceId   string `json:"partnerServiceId"`
			CustomerNo         string `json:"customerNo"`
			VirtualAccountNo   string `json:"virtualAccountNo"`
			VirtualAccountName string `json:"virtualAccountName"`
			PaymentRequestId   string `json:"paymentRequestId"`
			TotalAmount        struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			} `json:"totalAmount"`
			BillDetails []struct {
				BillDescription struct {
					English   string `json:"english"`
					Indonesia string `json:"indonesia"`
				} `json:"billDescription"`
			} `json:"billDetails"`
		}{
			PartnerServiceId:   payload.PartnerServiceId,
			CustomerNo:         payload.CustomerNo,
			VirtualAccountNo:   payload.VirtualAccountNo,
			VirtualAccountName: payload.AdditionalInfo.DebitFromName,
			PaymentRequestId:   payload.PaymentRequestId,
			TotalAmount: struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			}{
				Value:    payload.TotalAmount.Value,
				Currency: payload.TotalAmount.Currency,
			},
			BillDetails: nil,
		},
		AdditionalInfo: struct {
			ReconcileId       string    `json:"reconcileId"`
			ReconcileDatetime time.Time `json:"reconcileDatetime"`
		}{
			ReconcileId:       kWalletTransaction.TransactionCode,
			ReconcileDatetime: kWalletTransaction.DateTime,
		},
	}
}

func NewPaymentStrategyKWallet(paymentService *PaymentService, payload *web.CreatePaymentRequest, payment *web.Payment, channel *models.Channel, platform *models.Platforms, noRekning string) IPaymentStrategy {
	return &PaymentStrategyKWallet{
		paymentService: paymentService,
		payload:        payload,
		payment:        payment,
		channel:        channel,
		platform:       platform,
		noRekening:     noRekning,
		strategy:       "k-wallet",
	}
}

func (s *PaymentStrategyKWallet) CreatePayment(ctx context.Context) (*web.PaymentResponse, error) {
	var err error
	// get k-wallet member
	//filterKwallet := []*pagination.Filter{
	//	{
	//		ID:       "member_id",
	//		Value:    s.payload.CustomerId,
	//		Variant:  "string",
	//		Operator: "eq",
	//		FilterID: "",
	//	},
	//	{
	//		ID:       "currency",
	//		Value:    s.channel.Currency,
	//		Variant:  "string",
	//		Operator: "eq",
	//		FilterID: "",
	//	},
	//	{
	//		ID:       "no_rekening",
	//		Value:    s.noRekening,
	//		Variant:  "string",
	//		Operator: "eq",
	//		FilterID: "",
	//	},
	//}
	//
	//filter, err := helpers.FilterColumnValidation(filterKwallet, models.AllowedFilterColumnKWallet())
	//if err != nil {
	//	log.Error().Err(err).Str("strategy", s.strategy).Msg("filter column validation error")
	//	return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
	//}

	kWallet, err := s.paymentService.service.repository.GetKWalletRepository(ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Where("k_wallet.member_id = ? and k_wallet.no_rekening = ? and k_wallet.currency = ?", s.payload.CustomerId, s.noRekening, s.channel.Currency)
		//query, args := s.paymentService.service.repository.SearchQuery(filter, "and")
		//if query != "" {
		//	db = db.Where(query, args...)
		//}
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("strategy", s.strategy).Msg("get k-wallet repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(fmt.Errorf("can't continue payment, member id %v has not registered k-wallet", s.payload.CustomerId), s.strategy).WithStatusCode(http.StatusPaymentRequired)
		}

		return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	amountDecimal := decimal.NewFromInt(s.payment.Amount)

	// check k-wallet
	if kWallet.Balance.LessThan(amountDecimal) && len(s.payment.Channel) < 1 {
		return nil, helpers.NewErrorTrace(fmt.Errorf("can't continue payment, member id %v has not enough balance", s.payload.CustomerId), s.strategy).WithStatusCode(http.StatusPaymentRequired)
	}

	if kWallet.Balance.GreaterThan(amountDecimal) && len(s.payment.Channel) > 1 {
		return nil, helpers.NewErrorTrace(fmt.Errorf("can't continue payment, member id %v have enough balance and cannot split payment", s.payload.CustomerId), s.strategy).WithStatusCode(http.StatusBadRequest)
	}

	transactionId := uuid.New().String()
	totalAmount, feeAmount := s.paymentService.calculatePayment(amountDecimal, s.channel)
	timeNow := time.Now()

	// create payment
	payment := &models.Payments{
		TransactionId:        transactionId,
		OrderId:              s.payment.OrderId,
		PlatformId:           s.platform.Id,
		PaymentMethodId:      s.channel.Id,
		Amount:               amountDecimal,
		FeeAmount:            feeAmount,
		TotalAmount:          totalAmount,
		Currency:             s.channel.Currency,
		Status:               enums.PAYMENT_STATUS_SUCCESS,
		CustomerId:           s.payload.CustomerId,
		CustomerName:         s.payload.CustomerName,
		CustomerEmail:        s.payload.CustomerEmail,
		CustomerPhone:        s.payload.CustomerPhone,
		ReferenceId:          s.payload.ReferenceId,
		ReferenceType:        s.payload.ReferenceType,
		GatewayTransactionId: s.platform.Code,
		GatewayReference:     s.platform.Name,
		GatewayResponse:      nil,
		CallbackUrl:          "",
		ReturnUrl:            "",
		ExpiredAt:            &timeNow,
		ExpiredTime:          "0",
		PaidAt:               &timeNow,
	}

	// update k-wallet balance
	kWallet.SubBalance(totalAmount)

	// create k-wallet transaction
	kWalletTransaction := &models.KWalletTransaction{
		KWalletID:                kWallet.ID,
		KWalletTypeTransactionID: 1,
		PaymentID:                s.payment.OrderId,
		PaymentCode:              transactionId,
		TransactionCode:          transactionId,
		TransactionType:          "",
		CounterpartyName:         kWallet.FullName,
		CounterpartyBank:         "",
		PaymentChannel:           s.platform.Name,
		Description:              fmt.Sprintf("Pembayaran Transaksi %v", s.payment.OrderId),
		Currency:                 kWallet.Currency,
		Symbol:                   kWallet.Symbol,
		Status:                   enums.PAYMENT_STATUS_SUCCESS,
	}

	kWalletTransaction.PaymentTransaction(kWallet.Balance, totalAmount)

	log.Debug().Interface("kwallet", kWallet).Interface("kwallet transaction", kWalletTransaction).Msg("data kwallet and kwallet transaction")

	gatewayResponse, err := json.Marshal(kWallet)
	if err != nil {
		log.Error().Err(err).Msg("gateway marshal kwallet error")
		return nil, err
	}
	payment.GatewayResponse = gatewayResponse

	err = s.paymentService.service.repository.WithTransaction(func(tx *gorm.DB) error {
		// insert payment
		payment, err = s.paymentService.service.repository.InsertPaymentRepositoryTx(ctx, tx, payment)
		if err != nil {
			log.Error().Err(err).Str("strategy", s.strategy).Msg("insert payment repository tx error")
			return err
		}

		// update k-wallet
		err = s.paymentService.service.repository.UpdateKWalletRepositoryTx(ctx, tx, func(db *gorm.DB) *gorm.DB {
			db = db.Where("k_wallet.id = ?", kWallet.ID)
			updateColumn := map[string]interface{}{
				"balance":    kWallet.Balance,
				"updated_at": time.Now(),
			}
			db = db.UpdateColumns(updateColumn)
			return db
		})
		if err != nil {
			log.Error().Err(err).Str("strategy", s.strategy).Msg("update k-wallet repository tx error")
			return err
		}

		// insert k-wallet transaction
		err = s.paymentService.service.repository.InsertKWalletTransactionRepositoryTx(ctx, tx, kWalletTransaction)
		if err != nil {
			log.Error().Err(err).Str("strategy", s.strategy).Msg("insert k-wallet transaction repository tx error")
			return err
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Str("strategy", s.strategy).Msg("with transaction error")
		return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	response := &web.PaymentResponse{
		Id:            payment.Id,
		TransactionId: payment.TransactionId,
		OrderId:       payment.OrderId,
		Status:        payment.Status,
		Amount:        payment.Amount.IntPart(),
		FeeAdmin:      payment.FeeAmount.IntPart(),
		TotalAmount:   payment.TotalAmount.IntPart(),
		Currency:      payment.Currency,
		PaymentMethod: s.channel.PaymentMethod,
		PaymentType:   enums.PAYMENT_TYPE_VA,
		PaymentDetail: web.PaymentDetail{
			Bank:            "",
			Url:             nil,
			VaNumber:        kWallet.NoRekening,
			BillKey:         "",
			BIllCode:        "",
			TransactionTime: payment.CreatedAt.Format(time.DateTime),
			ExpireTime:      payment.ExpiredTime,
		},
		Customer: web.Customer{
			MemberId: payment.CustomerId,
			Name:     payment.CustomerName,
			Email:    payment.CustomerEmail,
			Phone:    payment.CustomerPhone,
		},
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	}

	return response, nil
}

func NewPaymentStrategyPG(paymentService *PaymentService, payload *web.CreatePaymentRequest, payment *web.Payment, platform *models.Platforms, channel *models.Channel) IPaymentStrategy {
	return &PaymentStrategyPG{
		paymentService: paymentService,
		channel:        channel,
		platform:       platform,
		payload:        payload,
		payment:        payment,
		strategy:       "va",
	}
}

func (s *PaymentStrategyPG) CreatePayment(ctx context.Context) (*web.PaymentResponse, error) {
	configurations, err := s.paymentService.service.repository.GetConfigurationByPlatformIdAndCurrencyRepository(ctx, s.platform.Id, s.channel.Currency)
	if err != nil {
		log.Error().Err(err).Str("context", s.paymentService.serviceName).Msg("get configuration by platform id and aggregator id repository")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(fmt.Errorf("%v, platform configuration", err), s.paymentService.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("configuration", configurations).Msg("get platform configuration repository")

	try := 0
	transactionId := uuid.New().String()
	amountDecimal := s.payment.Amount

	payments := &models.Payments{
		TransactionId:   transactionId,
		OrderId:         s.payment.OrderId,
		PlatformId:      s.platform.Id,
		PaymentMethodId: s.channel.Id,
		//AggregatorId:         &configuration.AggregatorId,
		Amount: decimal.NewFromInt(amountDecimal),
		//FeeAmount:            fee,
		//TotalAmount:          totalAmount,
		Currency:             s.channel.Currency,
		Status:               enums.PAYMENT_STATUS_PENDING,
		CustomerId:           s.payload.CustomerId,
		CustomerName:         s.payload.CustomerName,
		CustomerEmail:        s.payload.CustomerEmail,
		CustomerPhone:        s.payload.CustomerPhone,
		ReferenceId:          s.payload.ReferenceId,
		ReferenceType:        s.payload.ReferenceType,
		GatewayTransactionId: s.platform.Code,
		GatewayReference:     s.platform.Name,
		GatewayResponse:      nil,
		CallbackUrl:          "",
		ReturnUrl:            "",
		ExpiredAt:            nil,
		PaidAt:               nil,
	}
	log.Debug().Interface("payment", payments).Msg("model payment")

	for _, configuration := range configurations {
		s.paymentService.service.repository.SetConfigurationPayment(configuration)

		totalAmount, fee := s.paymentService.calculatePayment(decimal.NewFromInt(amountDecimal), s.channel)
		log.Debug().Interface("totalAmount", totalAmount).Interface("fee", fee).Msg("calculate payment")

		payments.AggregatorId = &configuration.AggregatorId
		payments.FeeAmount = fee
		payments.TotalAmount = totalAmount

		if try == 0 {
			payments, err = s.paymentService.service.repository.InsertPaymentRepository(ctx, payments)
			if err != nil {
				log.Error().Err(err).Str("context", s.paymentService.serviceName).Msg("insert payment repository error")
				return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
			}
			log.Debug().Interface("payments", payments).Msg("return value insert payment repository")
		} else {
			payments, err = s.paymentService.service.repository.UpdatePaymentAfterTryRepository(ctx, payments)
			if err != nil {
				log.Error().Err(err).Msg("update payment after get response payment gateway repository")
				return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
			}
		}

		try++

		// get strategy
		strategyPayment, errStrategy := s.paymentService.service.repository.Strategy.GetStrategy(configuration)
		if errStrategy != nil {
			log.Error().Err(errStrategy).Str("context", s.paymentService.serviceName).Msg("get strategy payment error")
			return nil, helpers.NewErrorTrace(errStrategy, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("strategy", strategyPayment).Msg("get strategy payment error")

		//var gatewayResponse *models.GatewayResponse

		var paymentRequest any

		if configuration.Aggregator.Slug == enums.PROVIDER_PAYMENT_METHOD_SENANGPAY {
			paymentRequest = clientsenangpay.PaymentRequest{
				OrderID: s.payment.OrderId,
				Amount:  fmt.Sprint(totalAmount),
				Detail:  fmt.Sprintf("Shopping_id_%v", s.payment.OrderId),
				Name:    s.payload.CustomerName,
				//Email:   s.payload.CustomerEmail,
				Email: "email.statis@gmail.com",
				//Phone:   s.payload.CustomerPhone,
				Phone: "021111111111",
			}

			//gatewayResponse, err = s.service.repository.SenangpayPaymentRedirectUrlRepository(s.service.ctx, &paymentRequest)
			//if err != nil {
			//	return nil, err
			//}
		} else if configuration.Aggregator.Slug == enums.PROVIDER_PAYMENT_METHOD_MIDTRANS {
			paymentRequest = midtrans.PaymentRequest{
				OrderID:      s.payment.OrderId,
				Amount:       totalAmount,
				PaymentType:  s.channel.TransactionType,
				Method:       s.channel.PaymentMethod,
				Channel:      s.channel.BankName,
				Description:  fmt.Sprintf("Shopping_id_%v", s.payment.OrderId),
				CustomerName: s.payload.CustomerName,
				//CustomerEmail: s.payload.CustomerEmail,
				CustomerEmail: "email.statis@gmail.com",
				//CustomerPhone: s.payload.CustomerPhone,
				CustomerPhone: "021111111111",
			}
		} else if configuration.Aggregator.Slug == enums.PROVIDER_PAYMENT_METHOD_ESPAY {
			paymentRequest = espay.PaymentRequest{
				RQUUID:     uuid.New().String(),
				RQDateTime: time.Now(),
				OrderID:    s.payment.OrderId,
				Amount:     fmt.Sprint(totalAmount),
				FeeAmount:  fmt.Sprint(fee),
				CCY:        s.channel.Currency,
				//CommCode:      "",
				Method:     s.channel.PaymentMethod,
				CustomerID: s.payload.CustomerId,
				//CustomerPhone: s.payload.CustomerPhone,
				CustomerPhone: "021111111111",
				CustomerName:  s.payload.CustomerName,
				//CustomerEmail: s.payload.CustomerEmail,
				CustomerEmail: "email.statis@gmail.com",
				//Description:   payload.Description,
				BankCode:    s.channel.BankCode,
				ProductCode: s.channel.ProductCode,
				ProductName: s.channel.ProductName,
				VaExpired:   enums.VA_EXPIRED_180,
				ReturnUrl:   s.payload.ReturnUrl,
				//ReturnUrl: "https://google.com",
			}
		} else {
			return nil, helpers.NewErrorTrace(fmt.Errorf("provider %v not found", configuration.Aggregator.Slug), s.paymentService.serviceName).WithStatusCode(http.StatusNotFound)
		}

		log.Debug().Interface("paymentRequest", paymentRequest).Msg("payment request")

		response, errPay := strategyPayment.Pay(ctx, paymentRequest)
		if errPay != nil {
			log.Error().Err(errPay).Interface("response", response).Str("context", s.paymentService.serviceName).Msg("pay strategy error")
			//return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
			continue
		}
		log.Debug().Interface("response", response).Msg("pay strategy")

		gatewayResponse, errM := json.Marshal(response)
		if errM != nil {
			return nil, helpers.NewErrorTrace(errM, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("gatewayResponse", gatewayResponse).Msg("response marshal gateway response")

		payments.GatewayResponse = gatewayResponse
		log.Debug().Interface("payment.GatewayResponse", payments.GatewayResponse).Msg("assignment payment gateway response")

		webResponse, err := strategyPayment.MapResponsePayment(s.channel, payments)
		if err != nil {
			log.Error().Err(err).Str("context", s.paymentService.serviceName).Msg("map response strategy error")
			return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("webResponse", webResponse).Msg("map response strategy error")

		payments.SetFixPayment(decimal.NewFromInt(webResponse.TotalAmount), decimal.NewFromInt(webResponse.FeeAdmin), webResponse.PaymentDetail.ExpireTime)
		webResponse.Amount = payments.Amount.IntPart()
		//update local payment after get response on payment gateway
		err = s.paymentService.service.repository.UpdatePaymentAfterGetResponsePaymentGatewayRepository(ctx, payments)
		if err != nil {
			log.Error().Err(err).Msg("update payment after get response payment gateway repository")
			return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
		}

		//payment.ExpiredTime = webResponse.PaymentDetail.ExpireTime

		//return s.mapToDetailPaymentResponse(payments), nil
		return webResponse, nil
	}

	err = s.paymentService.service.repository.UpdatePaymentRepository(ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Where("payments.id = ?", payments.Id)
		updateColumn := map[string]interface{}{
			"status":     enums.PAYMENT_STATUS_FAILED,
			"updated_at": time.Now(),
		}
		db = db.UpdateColumns(updateColumn)
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.strategy).Msg("failed to update payment method")
		return nil, helpers.NewErrorTrace(err, s.strategy).WithStatusCode(http.StatusInternalServerError)
	}

	err = s.paymentService.service.repository.UpdateChannelRepository(ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Where("channels.id = ?", s.channel.Id)
		updateColumn := map[string]interface{}{
			//"code":       channel.Code,
			//"name":       channel.Name,
			//"provider":   channel.Provider,
			//"currency":   channel.Currency,
			//"fee_type":   channel.FeeType,
			//"fee_amount": channel.FeeAmount,
			"is_active":  false,
			"updated_at": time.Now(),
		}
		db = db.UpdateColumns(updateColumn)
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.strategy).Msg("failed to update payment method")
		return nil, helpers.NewErrorTrace(err, s.strategy).WithStatusCode(http.StatusInternalServerError)
	}

	return nil, helpers.NewErrorTrace(fmt.Errorf("payment failed, please call support team payment"), s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
}

type PaymentStrategyCC struct {
	paymentService *PaymentService
	channel        *models.Channel
	platform       *models.Platforms
	payload        *web.CreatePaymentRequest
	payment        *web.Payment
	strategy       string
}

func NewPaymentStrategyCreditCard(s *PaymentService, payload *web.CreatePaymentRequest, payment *web.Payment, platform *models.Platforms, channel *models.Channel) IPaymentStrategy {
	return &PaymentStrategyCC{
		paymentService: s,
		payload:        payload,
		payment:        payment,
		platform:       platform,
		channel:        channel,
		strategy:       "credit_card",
	}
}

func (s *PaymentStrategyCC) CreatePayment(ctx context.Context) (*web.PaymentResponse, error) {
	// filter configuration
	filterConfiguration := []*pagination.Filter{
		{
			ID:       "aggregators.slug",
			Value:    enums.PROVIDER_PAYMENT_METHOD_ESPAY,
			Variant:  "string",
			Operator: "eq",
			FilterID: "",
		},
		{
			ID:       "aggregators.currency",
			Value:    s.channel.Currency,
			Variant:  "string",
			Operator: "eq",
			FilterID: "",
		},
		{
			ID:       "platform_configuration.platform_id",
			Value:    s.platform.Id,
			Variant:  "number",
			Operator: "eq",
			FilterID: "",
		},
	}

	_, err := helpers.FilterColumnValidation(filterConfiguration, models.AllowedFilterColumnConfiguration())
	if err != nil {
		log.Error().Err(err).Str("context", s.paymentService.serviceName).Msg("filter column validation")
		return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	configuration, err := s.paymentService.service.repository.GetConfigurationRepository(ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Joins("join aggregators on aggregators.id = configurations.aggregator_id")
		db = db.Joins("join platform_configuration on platform_configuration.configuration_id = configurations.id")
		query, args := s.paymentService.service.repository.SearchQuery(filterConfiguration, "and")
		if query != "" {
			db = db.Where(query, args...)
		}
		db = db.Preload("Aggregator")

		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.paymentService.serviceName).Msg("get configuration repository")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(fmt.Errorf("%v, platform configuration", err), s.paymentService.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("configuration", configuration).Msg("get configuration repository")

	transactionId := uuid.New().String()
	amountDecimal := s.payment.Amount

	s.paymentService.service.repository.SetConfigurationPayment(configuration)

	totalAmount, fee := s.paymentService.calculatePayment(decimal.NewFromInt(amountDecimal), s.channel)
	log.Debug().Interface("totalAmount", totalAmount).Interface("fee", fee).Msg("calculate payment")

	payments := &models.Payments{
		TransactionId:   transactionId,
		OrderId:         s.payment.OrderId,
		PlatformId:      s.platform.Id,
		PaymentMethodId: s.channel.Id,
		AggregatorId:    &configuration.AggregatorId,
		Amount:          decimal.NewFromInt(amountDecimal),
		FeeAmount:       fee,
		TotalAmount:     totalAmount,
		Currency:        s.channel.Currency,
		Status:          enums.PAYMENT_STATUS_PENDING,
		CustomerId:      s.payload.CustomerId,
		CustomerName:    s.payload.CustomerName,
		//CustomerEmail:        s.payload.CustomerEmail,
		CustomerEmail: "email.statis@gmail.com",
		//CustomerPhone:        s.payload.CustomerPhone,
		CustomerPhone:        "021111111111",
		ReferenceId:          s.payload.ReferenceId,
		ReferenceType:        s.payload.ReferenceType,
		GatewayTransactionId: s.platform.Code,
		GatewayReference:     s.platform.Name,
		GatewayResponse:      nil,
		CallbackUrl:          "",
		ReturnUrl:            "",
		ExpiredAt:            nil,
		PaidAt:               nil,
	}
	log.Debug().Interface("payment", payments).Msg("model payment")

	payments, err = s.paymentService.service.repository.InsertPaymentRepository(ctx, payments)
	if err != nil {
		log.Error().Err(err).Str("context", s.paymentService.serviceName).Msg("insert payment repository error")
		return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("payments", payments).Msg("return value insert payment repository")

	// get strategy
	strategy, errStrategy := s.paymentService.service.repository.Strategy.GetStrategy(configuration)
	if errStrategy != nil {
		log.Error().Err(errStrategy).Str("context", s.paymentService.serviceName).Msg("get strategy payment error")
		return nil, helpers.NewErrorTrace(errStrategy, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("strategy", strategy).Msg("get strategy payment error")

	//var gatewayResponse *models.GatewayResponse

	var paymentRequest any

	if configuration.Aggregator.Slug == enums.PROVIDER_PAYMENT_METHOD_SENANGPAY {
		paymentRequest = clientsenangpay.PaymentRequest{
			OrderID: s.payment.OrderId,
			Amount:  fmt.Sprint(totalAmount),
			Detail:  fmt.Sprintf("Shopping_id_%v", s.payment.OrderId),
			Name:    s.payload.CustomerName,
			Email:   s.payload.CustomerEmail,
			Phone:   s.payload.CustomerPhone,
		}

		//gatewayResponse, err = s.service.repository.SenangpayPaymentRedirectUrlRepository(s.service.ctx, &paymentRequest)
		//if err != nil {
		//	return nil, err
		//}
	} else if configuration.Aggregator.Slug == enums.PROVIDER_PAYMENT_METHOD_MIDTRANS {
		paymentRequest = midtrans.PaymentRequest{
			OrderID:       s.payment.OrderId,
			Amount:        totalAmount,
			PaymentType:   s.channel.TransactionType,
			Method:        s.channel.PaymentMethod,
			Channel:       s.channel.BankName,
			Description:   fmt.Sprintf("Shopping_id_%v", s.payment.OrderId),
			CustomerName:  s.payload.CustomerName,
			CustomerEmail: s.payload.CustomerEmail,
			CustomerPhone: s.payload.CustomerPhone,
		}
	} else if configuration.Aggregator.Slug == enums.PROVIDER_PAYMENT_METHOD_ESPAY {
		paymentRequest = espay.PaymentRequest{
			RQUUID:     uuid.New().String(),
			RQDateTime: time.Now(),
			OrderID:    s.payment.OrderId,
			Amount:     fmt.Sprint(payments.Amount),
			FeeAmount:  fmt.Sprint(fee),
			CCY:        s.channel.Currency,
			//CommCode:      "",
			Method:        s.channel.PaymentMethod,
			CustomerID:    s.payload.CustomerId,
			CustomerPhone: s.payload.CustomerPhone,
			CustomerName:  s.payload.CustomerName,
			CustomerEmail: s.payload.CustomerEmail,
			//Description:   payload.Description,
			BankCode:    s.channel.BankCode,
			ProductCode: s.channel.ProductCode,
			ProductName: s.channel.ProductName,
			VaExpired:   enums.VA_EXPIRED_180,
			ReturnUrl:   s.payload.ReturnUrl,
		}
	} else {
		return nil, helpers.NewErrorTrace(fmt.Errorf("provider %v not found", configuration.Aggregator.Slug), s.paymentService.serviceName).WithStatusCode(http.StatusNotFound)
	}

	log.Debug().Interface("paymentRequest", paymentRequest).Msg("payment request")

	response, errPay := strategy.Pay(ctx, paymentRequest)
	if errPay != nil {
		log.Error().Err(errPay).Str("context", s.paymentService.serviceName).Msg("pay strategy error")
		return nil, helpers.NewErrorTrace(errPay, s.strategy).WithStatusCode(http.StatusInternalServerError)
		//continue
	}
	log.Debug().Interface("response", response).Msg("pay strategy")

	gatewayResponse, errM := json.Marshal(response)
	if errM != nil {
		return nil, helpers.NewErrorTrace(errM, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("gatewayResponse", gatewayResponse).Msg("response marshal gateway response")

	payments.GatewayResponse = gatewayResponse
	log.Debug().Interface("payment.GatewayResponse", string(payments.GatewayResponse)).Msg("assignment payment gateway response")

	//update local payment after get response on payment gateway
	err = s.paymentService.service.repository.UpdatePaymentAfterGetResponsePaymentGatewayRepository(ctx, payments)
	if err != nil {
		log.Error().Err(err).Msg("update payment after get response payment gateway repository")
		return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	webResponse, err := strategy.MapResponsePayment(s.channel, payments)
	if err != nil {
		//errUpdate := s.service.repository.UpdateChannelRepository(s.service.ctx, payments.Channel)
		//if errUpdate != nil {
		//	log.Error().Err(errUpdate).Str("context", s.serviceName).Msg("failed to update payment method")
		//	return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		//}
		log.Error().Err(err).Str("context", s.paymentService.serviceName).Msg("map response strategy error")
		return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("webResponse", webResponse).Msg("map response strategy error")

	//payment.ExpiredTime = webResponse.PaymentDetail.ExpireTime

	//return s.mapToDetailPaymentResponse(payments), nil
	return webResponse, nil
}

type PaymentStrategyQris struct {
	paymentService *PaymentService
	payment        *models.Payments
	platform       *models.Platforms
	channel        *models.Channel
	payload        *web.CreatePaymentRequest
}

func NewPaymentStrategyQris(paymentService *PaymentService, payment *models.Payments, platform *models.Platforms, channel *models.Channel, payload *web.CreatePaymentRequest) *PaymentStrategyQris {
	return &PaymentStrategyQris{
		paymentService: paymentService,
		payment:        payment,
		platform:       platform,
		channel:        channel,
		payload:        payload,
	}
}

func (s *PaymentStrategyQris) CreatePayment(ctx context.Context) (*web.PaymentResponse, error) {
	configurations, err := s.paymentService.service.repository.GetConfigurationByPlatformIdAndCurrencyRepository(ctx, s.platform.Id, s.channel.Currency)
	if err != nil {
		log.Error().Err(err).Str("context", s.paymentService.serviceName).Msg("get configuration by platform id and aggregator id repository")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(fmt.Errorf("%v, platform configuration", err), s.paymentService.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("configuration", configurations).Msg("get platform configuration repository")

	try := 0
	transactionId := uuid.New().String()
	amountDecimal := s.payment.Amount

	for _, configuration := range configurations {
		s.paymentService.service.repository.SetConfigurationPayment(configuration)

		totalAmount, fee := s.paymentService.calculatePayment(amountDecimal, s.channel)
		log.Debug().Interface("totalAmount", totalAmount).Interface("fee", fee).Msg("calculate payment")

		payments := &models.Payments{
			TransactionId:        transactionId,
			OrderId:              s.payment.OrderId,
			PlatformId:           s.platform.Id,
			PaymentMethodId:      s.channel.Id,
			AggregatorId:         &configuration.AggregatorId,
			Amount:               amountDecimal,
			FeeAmount:            fee,
			TotalAmount:          totalAmount,
			Currency:             s.channel.Currency,
			Status:               enums.PAYMENT_STATUS_PENDING,
			CustomerId:           s.payload.CustomerId,
			CustomerName:         s.payload.CustomerName,
			CustomerEmail:        s.payload.CustomerEmail,
			CustomerPhone:        s.payload.CustomerPhone,
			ReferenceId:          s.payload.ReferenceId,
			ReferenceType:        s.payload.ReferenceType,
			GatewayTransactionId: s.platform.Code,
			GatewayReference:     s.platform.Name,
			GatewayResponse:      nil,
			CallbackUrl:          "",
			ReturnUrl:            "",
			ExpiredAt:            nil,
			PaidAt:               nil,
		}
		log.Debug().Interface("payment", payments).Msg("model payment")

		if try == 0 {
			payments, err = s.paymentService.service.repository.InsertPaymentRepository(ctx, payments)
			if err != nil {
				log.Error().Err(err).Str("context", s.paymentService.serviceName).Msg("insert payment repository error")
				return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
			}
			log.Debug().Interface("payments", payments).Msg("return value insert payment repository")
		} else {
			payments, err = s.paymentService.service.repository.UpdatePaymentAfterTryRepository(ctx, payments)
			if err != nil {
				log.Error().Err(err).Msg("update payment after get response payment gateway repository")
				return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
			}
		}

		try++

		// get strategy
		strategy, errStrategy := s.paymentService.service.repository.Strategy.GetStrategy(configuration)
		if errStrategy != nil {
			log.Error().Err(errStrategy).Str("context", s.paymentService.serviceName).Msg("get strategy payment error")
			return nil, helpers.NewErrorTrace(errStrategy, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("strategy", strategy).Msg("get strategy payment error")

		//var gatewayResponse *models.GatewayResponse

		var paymentRequest any

		if configuration.Aggregator.Slug == enums.PROVIDER_PAYMENT_METHOD_SENANGPAY {
			paymentRequest = clientsenangpay.PaymentRequest{
				OrderID: s.payment.OrderId,
				Amount:  fmt.Sprint(totalAmount),
				Detail:  fmt.Sprintf("Shopping_id_%v", s.payment.OrderId),
				Name:    s.payload.CustomerName,
				//Email:   s.payload.CustomerEmail,
				Email: "email.statis@gmail.com",
				//Phone:   s.payload.CustomerPhone,
				Phone: "021111111111",
			}

			//gatewayResponse, err = s.service.repository.SenangpayPaymentRedirectUrlRepository(s.service.ctx, &paymentRequest)
			//if err != nil {
			//	return nil, err
			//}
		} else if configuration.Aggregator.Slug == enums.PROVIDER_PAYMENT_METHOD_MIDTRANS {
			paymentRequest = midtrans.PaymentRequest{
				OrderID:      s.payment.OrderId,
				Amount:       totalAmount,
				PaymentType:  s.channel.TransactionType,
				Method:       s.channel.PaymentMethod,
				Channel:      s.channel.BankName,
				Description:  fmt.Sprintf("Shopping_id_%v", s.payment.OrderId),
				CustomerName: s.payload.CustomerName,
				//CustomerEmail: s.payload.CustomerEmail,
				CustomerEmail: "email.statis@gmail.com",
				//CustomerPhone: s.payload.CustomerPhone,
				CustomerPhone: "021111111111",
			}
		} else if configuration.Aggregator.Slug == enums.PROVIDER_PAYMENT_METHOD_ESPAY {
			paymentRequest = espay.PaymentRequest{
				RQUUID:     uuid.New().String(),
				RQDateTime: time.Now(),
				OrderID:    s.payment.OrderId,
				Amount:     fmt.Sprint(payments.TotalAmount),
				FeeAmount:  fmt.Sprint(fee),
				CCY:        s.channel.Currency,
				//CommCode:      "",
				Method:     s.channel.PaymentMethod,
				CustomerID: s.payload.CustomerId,
				//CustomerPhone: s.payload.CustomerPhone,
				CustomerPhone: "021111111111",
				CustomerName:  s.payload.CustomerName,
				//CustomerEmail: s.payload.CustomerEmail,
				CustomerEmail: "email.statis@gmail.com",
				//Description:   payload.Description,
				BankCode:    s.channel.BankCode,
				ProductCode: s.channel.ProductCode,
				ProductName: s.channel.ProductName,
				VaExpired:   enums.VA_EXPIRED_180,
				ReturnUrl:   s.payload.ReturnUrl,
				//ReturnUrl: "https://google.com",
			}
		} else {
			return nil, helpers.NewErrorTrace(fmt.Errorf("provider %v not found", configuration.Aggregator.Slug), s.paymentService.serviceName).WithStatusCode(http.StatusNotFound)
		}

		log.Debug().Interface("paymentRequest", paymentRequest).Msg("payment request")

		response, errPay := strategy.Pay(ctx, paymentRequest)
		if errPay != nil {
			log.Error().Err(errPay).Str("context", s.paymentService.serviceName).Msg("pay strategy error")
			//return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
			continue
		}
		log.Debug().Interface("response", response).Msg("pay strategy")

		gatewayResponse, errM := json.Marshal(response)
		if errM != nil {
			return nil, helpers.NewErrorTrace(errM, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("gatewayResponse", gatewayResponse).Msg("response marshal gateway response")

		payments.GatewayResponse = gatewayResponse
		log.Debug().Interface("payment.GatewayResponse", payments.GatewayResponse).Msg("assignment payment gateway response")

		//update local payment after get response on payment gateway
		err = s.paymentService.service.repository.UpdatePaymentAfterGetResponsePaymentGatewayRepository(ctx, payments)
		if err != nil {
			log.Error().Err(err).Msg("update payment after get response payment gateway repository")
			return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
		}

		webResponse, err := strategy.MapResponsePayment(s.channel, payments)
		if err != nil {
			//errUpdate := s.service.repository.UpdateChannelRepository(s.service.ctx, payments.Channel)
			//if errUpdate != nil {
			//	log.Error().Err(errUpdate).Str("context", s.serviceName).Msg("failed to update payment method")
			//	return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
			//}
			log.Error().Err(err).Str("context", s.paymentService.serviceName).Msg("map response strategy error")
			return nil, helpers.NewErrorTrace(err, s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("webResponse", webResponse).Msg("map response strategy error")

		//payment.ExpiredTime = webResponse.PaymentDetail.ExpireTime

		//return s.mapToDetailPaymentResponse(payments), nil
		return webResponse, nil
	}

	return nil, helpers.NewErrorTrace(fmt.Errorf("payment failed, please call support team payment"), s.paymentService.serviceName).WithStatusCode(http.StatusInternalServerError)
}

func (s *PaymentService) MidtransPaymentNotificationService(payload *midtrans.CheckStatusPaymentResponse) (interface{}, error) {
	payment, err := s.service.repository.GetPaymentRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Where("payments.order_id = ? and payments.status = ?", payload.OrderId, enums.PAYMENT_STATUS_PENDING)
		db = db.Preload("Platform").Preload("Channel").Preload("Aggregator")
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get payment by transaction id repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(fmt.Errorf("%v, payment transaction", err), s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("payment", payment).Msg("get payment by transaction id repository")

	if payment.Status == enums.PAYMENT_STATUS_SUCCESS || payment.Status == enums.PAYMENT_STATUS_EXPIRED || payment.Status == enums.PAYMENT_STATUS_FAILED {
		return nil, helpers.NewErrorTrace(fmt.Errorf("notification has been accepted"), s.serviceName).WithStatusCode(http.StatusAccepted)
	}

	if payment.Channel == nil {
		return nil, helpers.NewErrorTrace(fmt.Errorf("payment method is missing, payment notification"), s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	configuration, err := s.service.repository.GetConfigurationByPlatformIdAndAggregatorIdRepository(s.service.ctx, payment.PlatformId, *payment.AggregatorId)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get configuration by platform id and aggregator id repository")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(fmt.Errorf("%v, platform configuration", err), s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("configuration", configuration).Interface("context", s.serviceName).Msg("get platform repository")

	dataTypesJson, err := helpers.ConvertAnyToDatatypeJson(payload)
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("convert amy to datatype json error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("dataTypesJson", dataTypesJson).Msg("convert amy to datatype json")

	s.service.repository.SetConfigurationPayment(configuration)

	// get strategy
	strategy, err := s.service.repository.Strategy.GetStrategy(configuration)
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get strategy payment error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	result, err := strategy.MapCheckStatusPayment(payment, *payload)
	if err != nil {
		log.Error().Err(err).Msg("map check status payment error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("result", result).Interface("context", s.serviceName).Msg("map check status payment")

	s.updatePaymentStatus(&payment, result)

	paymentCallback := models.PaymentCallbacks{
		PaymentId:    payment.Id,
		GatewayName:  payment.GatewayReference,
		CallbackData: dataTypesJson,
		ResponseData: nil,
		Status:       payment.Status,
	}

	paymentStatusHistory := models.PaymentStatusHistory{
		PaymentId: payment.Id,
		Status:    payment.Status,
		Notes:     "",
		CreatedBy: models.CreatedBy{
			ID:       "Midtrans Callback",
			Name:     "Callback Notification",
			Role:     "platform",
			Platform: "Mdtrans",
		},
	}

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
		err = s.service.repository.UpdatePaymentRepositoryTx(s.service.ctx, tx, func(db *gorm.DB) *gorm.DB {
			db = db.Where("payments.id = ?", payment.Id)
			updateColumn := map[string]interface{}{
				"status":       payment.Status,
				"paid_at":      payment.PaidAt,
				"reference_id": payment.ReferenceId,
				"updated_at":   time.Now(),
			}
			db = db.UpdateColumns(updateColumn)
			return db
		})
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("update payment repository error")
			return helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}

		err = s.service.repository.InsertPaymentCallbackRepositoryTx(s.service.ctx, tx, &paymentCallback)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("insert payment callback repository error")
			return helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}

		err = s.service.repository.InsertPaymentStatusHistoryRepositoryTx(s.service.ctx, tx, &paymentStatusHistory)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("insert payment status history repository error")
			return helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}

		return nil
	})

	callback := web.PaymentCallback{
		TransactionId: payment.TransactionId,
		OrderId:       payment.OrderId,
		Status:        payment.Status,
		Amount:        payment.TotalAmount.IntPart(),
		Currency:      payment.Currency,
	}

	defer func() {
		err = s.service.repository.CallbackFunctionRepository(s.service.ctx, payment.Platform.NotificationURL, &callback)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("callback function repository error")
			return
		}

		err = s.service.repository.UpdatePaymentNotificationRepository(s.service.ctx, payment)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("update payment repository error")
			return
		}
	}()

	return nil, nil
}

func (s *PaymentService) GetDetailPaymentService(payload *web.GetDetailPaymentRequest) ([]*web.PaymentResponse, error) {
	response := make([]*web.PaymentResponse, 0)

	payments, err := s.service.repository.FindPaymentRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Where("payments.order_id = ? and ( payments.status = ? or payments.status = ?)", payload.OrderId, enums.PAYMENT_STATUS_PENDING, enums.PAYMENT_STATUS_SUCCESS)
		db = db.Preload("Channel").Preload("Aggregator")
		return db
	})
	if err != nil {
		log.Error().Err(err).Msg("find payment repository error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("payments", payments).Msg("result data find payment repository")

	for _, payment := range payments {
		if payment.Channel.PaymentMethod == enums.PAYMENT_METEHOD_K_WALLET {
			log.Debug().Interface("gateway response", payment.GatewayResponse).Msg("data gateway response")
			var kWallet models.KWallet
			err = json.Unmarshal(payment.GatewayResponse, &kWallet)
			if err != nil {
				log.Error().Err(err).Msg("unmarshal gateway response error")
				return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
			}
			log.Debug().Interface("k-wallet", kWallet).Msg("data k-wallet")

			response = append(response, &web.PaymentResponse{
				Id:            payment.Id,
				TransactionId: payment.TransactionId,
				OrderId:       payment.OrderId,
				Status:        payment.Status,
				Amount:        payment.Amount.IntPart(),
				FeeAdmin:      payment.FeeAmount.IntPart(),
				TotalAmount:   payment.TotalAmount.IntPart(),
				Currency:      payment.Currency,
				PaymentMethod: payment.Channel.PaymentMethod,
				PaymentType:   enums.PAYMENT_TYPE_VA,
				PaymentDetail: web.PaymentDetail{
					Bank:            "",
					Url:             nil,
					VaNumber:        kWallet.NoRekening,
					BillKey:         "",
					BIllCode:        "",
					TransactionTime: payment.CreatedAt.Format(time.DateTime),
					ExpireTime:      payment.ExpiredTime,
				},
				Customer: web.Customer{
					MemberId: payment.CustomerId,
					Name:     payment.CustomerName,
					Email:    payment.CustomerEmail,
					Phone:    payment.CustomerPhone,
				},
				CreatedAt: payment.CreatedAt,
				UpdatedAt: payment.UpdatedAt,
			})
		} else {
			var paymentStrategy strategy.PaymentStrategy
			if payment.Aggregator.Slug == enums.PROVIDER_PAYMENT_METHOD_MIDTRANS {
				paymentStrategy = midtrans.NewMidtrans(s.service.repository.HttpClient, &models.Configuration{})
				//mapResponsePayment, err := strategy.MapResponsePayment(payment.Channel, payment)
				//if err != nil {
				//	log.Error().Err(err).Msg("map response payment error")
				//	return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
				//}
				//response = append(response, mapResponsePayment)
			} else if payment.Aggregator.Slug == enums.PROVIDER_PAYMENT_METHOD_ESPAY {
				paymentStrategy = espay.NewEspay(s.service.repository.HttpClient, &models.Configuration{})
			} else if payment.Aggregator.Slug == enums.PROVIDER_PAYMENT_METHOD_SENANGPAY {
				paymentStrategy = clientsenangpay.NewSenangpay(s.service.repository.HttpClient, &models.Configuration{})
			} else {
				return nil, helpers.NewErrorTrace(fmt.Errorf("provider payment method not found"), s.serviceName).WithStatusCode(http.StatusNotFound)
			}

			responsePayment, err := paymentStrategy.MapResponsePayment(payment.Channel, payment)
			if err != nil {
				log.Error().Err(err).Msg("map response payment error")
				return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
			}
			response = append(response, responsePayment)
		}

	}

	return response, nil
}
