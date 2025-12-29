package midtrans

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	restyclient "paymentserviceklink/app/client/resty"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/strategy"
	"paymentserviceklink/app/web"

	"github.com/rs/zerolog/log"
)

// Midtrans implements PaymentStrategy interface
type Midtrans struct {
	BaseUrl            string
	ServerKey          string
	ClientKey          string
	serviceName        string
	MerchantId         string
	Http               *restyclient.RestyClient
	MidtransStrategy   *strategy.MidtransStrategy
	Gopay              *Gopay
	PermataVa          *PermataVa
	BcaVa              *BcaVa
	BniVa              *BniVa
	BriVa              *BriVa
	CIMBVa             *CIMBVa
	MandiriBillPayment *MandiriBillPayment
	Qris               *Qris
	ShopeePay          *ShopeePay
	Dana               *Dana
}

func NewMidtrans(httpClient *restyclient.RestyClient, config *models.Configuration) *Midtrans {
	var midtrans *Midtrans

	if config.ConfigValue == enums.PRODUCTION {
		midtrans = &Midtrans{
			BaseUrl:     config.ConfigJson.ProductionBaseUrl,
			ServerKey:   helpers.DecryptAES(config.ConfigJson.ProductionServerKey),
			ClientKey:   helpers.DecryptAES(config.ConfigJson.ProductionClientKey),
			MerchantId:  helpers.DecryptAES(config.ConfigJson.SandboxMerchantId),
			serviceName: "midtrans",
			Http:        httpClient,
		}
	} else {
		midtrans = &Midtrans{
			BaseUrl:     config.ConfigJson.SandboxBaseUrl,
			ServerKey:   helpers.DecryptAES(config.ConfigJson.SandboxServerKey),
			ClientKey:   helpers.DecryptAES(config.ConfigJson.SandboxClientKey),
			MerchantId:  helpers.DecryptAES(config.ConfigJson.SandboxMerchantId),
			serviceName: "midtrans",
			Http:        httpClient,
		}
	}
	log.Debug().Interface("base_url", midtrans.BaseUrl).Interface("server_key", midtrans.ServerKey).
		Interface("client_key", midtrans.ClientKey).Interface("merchant_id", midtrans.MerchantId).
		Interface("service_name", midtrans.serviceName).Msg("new midtrans configuration set up")

	midtrans.Gopay = NewGopay(midtrans)
	midtrans.PermataVa = NewPermataVa(midtrans)
	midtrans.BcaVa = NewBcaVa(midtrans)
	midtrans.BniVa = NewBniVa(midtrans)
	midtrans.BriVa = NewBriVa(midtrans)
	midtrans.CIMBVa = NewCIMBVa(midtrans)
	midtrans.MandiriBillPayment = NewMandiriBillPayment(midtrans)
	midtrans.Qris = NewQris(midtrans)
	midtrans.ShopeePay = NewShopeePay(midtrans)
	midtrans.Dana = NewDana(midtrans)

	midtrans.MidtransStrategy = strategy.NewMidtransStrategy(
		midtrans.BcaVa,
		midtrans.BniVa,
		midtrans.BriVa,
		midtrans.CIMBVa,
		midtrans.Gopay,
		midtrans.MandiriBillPayment,
		midtrans.PermataVa,
		midtrans.Qris,
		midtrans.ShopeePay,
		midtrans.Dana,
	)

	return midtrans
}

// Pay implements PaymentStrategy.Pay
func (m *Midtrans) Pay(ctx context.Context, req any) (map[string]interface{}, error) {
	paymentReq, ok := req.(PaymentRequest)
	if !ok {
		return nil, errors.New("invalid payment request type")
	}

	// get strategy
	strategyPayment, err := m.MidtransStrategy.GetStrategy(paymentReq.Method, paymentReq.Channel)
	if err != nil {
		return nil, err
	}

	// pay
	result, err := strategyPayment.Pay(ctx, paymentReq)
	if err != nil {
		return nil, err
	}
	log.Debug().Interface("result", result).Msg("result strategy payment pay()")

	return result, nil
}

func (m *Midtrans) MapResponsePayment(channel *models.Channel, payment *models.Payments) (*web.PaymentResponse, error) {
	log.Debug().Str("context", m.serviceName).Interface("payment", payment).Msg("map midtrans payment response")

	strategyPayment, err := m.MidtransStrategy.GetStrategy(channel.PaymentMethod, channel.BankName)
	if err != nil {
		return nil, err
	}

	response, err := strategyPayment.ClientResponse(channel, payment)
	if err != nil {
		channel.IsActive = false

		log.Error().Err(err).Str("context", m.serviceName).Msg("failed to map midtrans payment response")
		return nil, err
	}

	return response, nil
}

func (m *Midtrans) MapCheckStatusPayment(payment *models.Payments, statusPayment any) (*web.CheckStatusPaymentResponse, error) {
	data, ok := statusPayment.(CheckStatusPaymentResponse)
	if !ok {
		return nil, helpers.NewErrorTrace(errors.New("invalid status payment response type"), m.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	status := enums.TxnStatusMidtrans(data.TransactionStatus)

	if status == enums.MIDTRANS_STATUS_AUTHORIZE {
		return nil, helpers.NewErrorTrace(fmt.Errorf(data.StatusMessage), m.serviceName).WithStatusCode(http.StatusPaymentRequired)
	}

	currentStatus := payment.Status

	if status == enums.MIDTRANS_STATUS_CAPTURE || status == enums.MIDTRANS_STATUS_SETTLEMENT {
		currentStatus = enums.PAYMENT_STATUS_SUCCESS
	} else if status == enums.MIDTRANS_STATUS_PENDING {
		currentStatus = enums.PAYMENT_STATUS_PENDING
	} else if status == enums.MIDTRANS_STATUS_CANCEL {
		currentStatus = enums.PAYMENT_STATUS_CANCELLED
	} else if status == enums.MIDTRANS_STATUS_EXPIRE {
		currentStatus = enums.PAYMENT_STATUS_EXPIRED
	} else if status == enums.MIDTRANS_STATUS_DENIED {
		currentStatus = enums.PAYMENT_STATUS_FAILED
	} else if status == enums.MIDTRANS_STATUS_FAILURE {
		currentStatus = enums.PAYMENT_STATUS_FAILED
	} else {
		currentStatus = enums.PAYMENT_STATUS_PENDING
	}

	return &web.CheckStatusPaymentResponse{
		Id:            payment.Id,
		TransactionId: payment.TransactionId,
		OrderId:       payment.OrderId,
		Status:        currentStatus,
		Amount:        payment.TotalAmount.IntPart(),
		Currency:      payment.Currency,
	}, nil
}

func (m *Midtrans) CreatePayment(ctx context.Context, chargeReq map[string]interface{}) (map[string]interface{}, error) {
	// Make API request to Midtrans
	url := fmt.Sprintf("%s/v2/charge", m.BaseUrl)
	log.Info().Str("url", url).Interface("charge request", chargeReq).Str("context", m.serviceName).Msg("midtrans create payment va")
	resp, err := m.Http.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", m.GetBase64Authorization()).
		SetHeader("Content-Type", "application/json").
		SetBody(chargeReq).
		Post(url)
	if err != nil {
		log.Error().Err(err).Str("context", m.serviceName).Msg("failed to create midtrans payment")
		return nil, err
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		errMsg := fmt.Sprintf("midtrans API error: %v", string(resp.Body()))
		log.Error().Str("error", errMsg).Msg("midtrans payment failed")
		return nil, errors.New(errMsg)
	}

	log.Info().Interface("response", string(resp.Body())).Str("context", m.serviceName).Msg("midtrans payment success")

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		log.Error().Err(err).Str("context", m.serviceName).Msg("failed to parse midtrans response")
		return nil, err
	}

	statusCode, ok := result["status_code"]
	if ok {
		if statusCode != fmt.Sprint(http.StatusCreated) {
			return nil, fmt.Errorf("error payment " + fmt.Sprint(result["status_message"]))
		}
	}
	log.Debug().Interface("result midtrans", result).Msg("result midtrans payment")

	// Check for error response
	if resp.StatusCode() != http.StatusOK {
		errMsg := fmt.Sprintf("midtrans API error: %v", result["status_message"])
		log.Error().Str("error", errMsg).Msg("midtrans payment failed")
		return nil, errors.New(errMsg)
	}

	return result, nil
}

func (m *Midtrans) CreateSnapPayment(ctx context.Context, chargeReq map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/snap/v1/transactions", m.BaseUrl)

	log.Info().Str("url", url).Interface("charge request", chargeReq).Str("context", m.serviceName).Msg("midtrans create payment va")
	resp, err := m.Http.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", m.GetBase64Authorization()).
		SetHeader("Content-Type", "application/json").
		SetBody(chargeReq).
		Post(url)
	if err != nil {
		log.Error().Err(err).Str("context", m.serviceName).Msg("failed to create midtrans payment")
		return nil, err
	}

	log.Info().Interface("status code", resp.StatusCode()).Interface("response", string(resp.Body())).Str("context", m.serviceName).Msg("midtrans payment success")

	if resp.StatusCode() != http.StatusCreated {
		return nil, fmt.Errorf("error payment " + fmt.Sprint(resp.Body()))
	}

	var result map[string]interface{}

	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		log.Error().Err(err).Msg("create snap payment unmarshal response body error")
		return nil, err
	}

	return result, nil
}

// CheckStatusPayment implements PaymentStrategy.CheckStatusPayment
func (m *Midtrans) CheckStatusPayment(ctx context.Context, req any) (any, error) {
	request, ok := req.(CheckStatusPaymentRequest)
	if !ok {
		return nil, errors.New("invalid transaction ID type")
	}

	// Delegate to the existing CheckStatus method
	return m.CheckStatus(ctx, request.TransactionId)
}

// CheckStatus checks the status of a payment
func (m *Midtrans) CheckStatus(ctx context.Context, orderId string) (any, error) {
	log.Debug().Interface("context", ctx).Str("orderId", orderId).Msg("midtrans check payment status")
	url := fmt.Sprintf("%s/v2/%s/status", m.BaseUrl, orderId)

	resp, err := m.Http.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Basic "+m.GetBase64Authorization()).
		Get(url)
	if err != nil {
		log.Error().Err(err).Str("context", m.serviceName).Msg("failed to check midtrans payment status")
		return nil, helpers.NewErrorTrace(err, m.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("response", resp.Body()).Str("context", m.serviceName).Msg("midtrans check payment status success")

	// Parse response
	var result CheckStatusPaymentResponse
	if errParse := json.Unmarshal(resp.Body(), &result); errParse != nil {
		log.Error().Err(errParse).Str("context", m.serviceName).Msg("failed to parse midtrans status response")
		return nil, helpers.NewErrorTrace(errParse, m.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	// Map Midtrans response to CallbackData
	return result, nil
}

// check client key and server key correct
func (m *Midtrans) CheckKey() (any, error) {
	log.Debug().Str("base url", m.BaseUrl).Str("server key", m.ServerKey).Str("client key", m.ClientKey).Msg("check midtrans key")

	url := fmt.Sprintf("%s/v2/charge/status", m.BaseUrl)
	ctx := context.Background()

	resp, err := m.Http.Client.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", m.GetBase64Authorization()).
		Get(url)
	if err != nil {
		log.Error().Err(err).Str("context", m.serviceName).Msg("failed to check midtrans payment status")
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		log.Error().Err(err).Msg("failed to decode midtrans response")
		return nil, err
	}

	return result, nil
}

// createChargeRequest converts our PaymentRequest to Midtrans charge request format
func (m *Midtrans) createChargeRequest(req PaymentRequest) map[string]interface{} {
	return map[string]interface{}{
		"payment_type": req.PaymentType,
		"transaction_details": map[string]interface{}{
			"order_id":     req.OrderID,
			"gross_amount": req.Amount,
		},
		"customer_details": map[string]interface{}{
			"first_name": req.CustomerName,
			"email":      req.CustomerEmail,
			"phone":      req.CustomerPhone,
		},
		"item_details": []map[string]interface{}{
			{
				"id":       req.OrderID,
				"price":    req.Amount,
				"quantity": 1,
				"name":     req.Description,
			},
		},
	}
}

// mapToCallbackData converts Midtrans response to our standard CallbackData
func (m *Midtrans) mapToCallbackData(data map[string]interface{}) *models.CallbackData {
	// Get customer details
	customer := getMap(data, "customer_details")

	// Get transaction details
	transaction := getMap(data, "transaction_details")

	callback := &models.CallbackData{
		Name:          getString(customer, "first_name"),
		Email:         getString(customer, "email"),
		Phone:         getString(customer, "phone"),
		AmountPaid:    getString(transaction, "gross_amount"),
		TxnMessage:    getString(data, "status_message"),
		OrderId:       getString(transaction, "order_id"),
		TransactionId: getString(data, "transaction_id"),
		HashedValue:   "", // Midtrans doesn't provide a direct equivalent to hashed_value
	}

	// Map status to our standard status
	switch getString(data, "transaction_status") {
	case "capture", "settlement":
		callback.TxnStatus = enums.SENANGPAY_STATUS_SUCCESS
	case "pending":
		callback.TxnStatus = enums.SENANGPAY_STATUS_FAILED
	case "deny", "expire", "cancel":
		callback.TxnStatus = enums.SENANGPAY_STATUS_FAILED
	}

	return callback
}

// Helper function to safely get a map from a map
func getMap(m map[string]interface{}, key string) map[string]interface{} {
	if val, ok := m[key]; ok {
		if m, ok := val.(map[string]interface{}); ok {
			return m
		}
	}
	return make(map[string]interface{})
}

// Helper function to safely get string from map
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (m *Midtrans) GetBase64Authorization() string {
	s := fmt.Sprintf("%s:", m.ServerKey)

	fmt.Println(s)

	return fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(s)))
}
