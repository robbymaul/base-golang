package espay

import (
	"context"
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
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Midtrans implements PaymentStrategy interface
type Espay struct {
	BaseUrl          string
	MerchantCode     string
	MerchantName     string
	ApiKey           string
	SignatureKey     string
	PublicKey        string
	PrivateKey       string
	MerchantId       string
	MerchantPassword string
	serviceName      string
	httpClient       *restyclient.RestyClient
	EspayStrategy    *strategy.EspayStrategy
	VirtualAccount   *VirtualAccount
	CreditCard       *CreditCard
	QRIS             *QRIS
	PaymentLink      *PaymentLink
}

func NewEspay(
	http *restyclient.RestyClient,
	config *models.Configuration,
	// espayStrategy *strategy.EspayStrategy,
) *Espay {

	var espay *Espay

	if config.ConfigValue == enums.PRODUCTION {
		espay = &Espay{
			BaseUrl:          config.ConfigJson.ProductionBaseUrl,
			MerchantCode:     helpers.DecryptAES(config.ConfigJson.ProductionMerchantCode),
			MerchantName:     helpers.DecryptAES(config.ConfigJson.ProductionMerchantName),
			ApiKey:           helpers.DecryptAES(config.ConfigJson.ProductionApiKey),
			SignatureKey:     helpers.DecryptAES(config.ConfigJson.ProductionSignatureKey),
			MerchantId:       helpers.DecryptAES(config.ConfigJson.ProductionMerchantId),
			MerchantPassword: helpers.DecryptAES(config.ConfigJson.ProductionCredentialPassword),
			PublicKey:        helpers.DecryptAES(config.ConfigJson.PublicKey),
			PrivateKey:       helpers.DecryptAES(config.ConfigJson.PrivateKey),
			httpClient:       http,
			serviceName:      "espay",
		}
	} else {
		//		publicKey := `-----BEGIN PUBLIC KEY-----
		//MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2O9xDMTBiZ5oOy3LBVn6
		//TerxWMHEwxl6gr0SX1dRt4be5vq2voFMoCHokeowqpeU5ZQi0EM36W7Q1K8hH6KR
		//jdNqhdIHyMh7X0yhVJTQ3Fz9QcjBfeMwoovmIYHP+U08GKz7j99VojSSriYvzT1m
		//PdwvTuAdFT3QEXfgdMLKQCjtXF/eyg2Q+xCYJALv+zeaPlsu00RO3TM5NGaCSbFC
		//oF/xa4IOfV+215beBvl1fUhW6mkEo7gdhK8T0ddk5bInEJs3YzDwQNtAutLEFVot
		//EKX2ETqIk8S1H7Pou7tSo73O0fFGaSBhG610bKIb9lLTXCQYJKk8bygPaL3aoT+5
		//QwIDAQAB
		//-----END PUBLIC KEY-----`
		//
		//		privateKey := `-----BEGIN RSA PRIVATE KEY-----
		//MIIEpAIBAAKCAQEA2O9xDMTBiZ5oOy3LBVn6TerxWMHEwxl6gr0SX1dRt4be5vq2
		//voFMoCHokeowqpeU5ZQi0EM36W7Q1K8hH6KRjdNqhdIHyMh7X0yhVJTQ3Fz9QcjB
		//feMwoovmIYHP+U08GKz7j99VojSSriYvzT1mPdwvTuAdFT3QEXfgdMLKQCjtXF/e
		//yg2Q+xCYJALv+zeaPlsu00RO3TM5NGaCSbFCoF/xa4IOfV+215beBvl1fUhW6mkE
		//o7gdhK8T0ddk5bInEJs3YzDwQNtAutLEFVotEKX2ETqIk8S1H7Pou7tSo73O0fFG
		//aSBhG610bKIb9lLTXCQYJKk8bygPaL3aoT+5QwIDAQABAoIBAQDTSPIcc43kUWpH
		//KSSxQ59sQEVsIt1W//u4VhoMzekDDNMQuGNATIKq/Bud8jAQFq6oo4z8tltAefPf
		//Eer6+sU1ExKO369BOTIf8Wy4CnEaD1+CsNrzl1EJH6S2Qc6jizva9K/WwriO0RGD
		//mCG6jfCEk21oLxNkWt3KBa2RSx7dOLO+ct07jtRbfYCVCAezyx6fWxLJ6eVmGZXM
		//kOhAr9tQ6IC3v/iQgA00LNPXR+X12obcmNXtcng5uHffeZNr6tmpLpXTYLdwZlwl
		//FINuTGpPjp1yy6q6GQYphF51ywRFN17g8NoVHLXDAfnrmB1lgtbC3nSiAvqEq2c6
		//XQkAIZbBAoGBAPgQDtG7RJ/Wdo5ra9HMgceVqDQgrdY4vw4cnV4NVGSBGn8jNhk7
		//YrJ8siJbLxqi5cPwJzu7xS8krKyt3vBY8AFKvVJ9yZ06VVL4d2LvWr50ym/zshnC
		//w1WlKhcuyaqP6MCiC6pZNA5LR1AN6hK2B1ZnmSrvDkg+MZtGTAxJrF6TAoGBAN/g
		//aVtaHuw1Zh2ixRfUjjQ4YMSxt/68DnmAJemmWQysvFsTZLfy87KLenmLABnG9qke
		//sOLD/vC7h5s5G9+vN4JMbmTYGBYp0VW5wWaC7Nw8cskgsmb7BZ+K7HsQbmtxh9Nu
		//BeQqdmQHZvLQ6wgY+0QTy/1KTUPwxLztyJttGjiRAoGAQEkpDgFSD3osz0vXbU9q
		//cqa+KIQviMy79pRD1BPwQvuSOlCNvIw/T7IxF+Y5ltWQZe7evAQ1XbpLZZTJqc/i
		//ovMTjUU78psjcZUim2kcQy9RJyIojbSDmrZq6gceDC2vS/yyuTrU2r93g6+XcbHq
		//xOGkOBQrx10Wzf6xxp1xJjECgYBfSk6t4nsdAVGYtap8jS2GDqUps5dkZrkmgCQj
		//AnoOygtWHLgXD+MokPOtfjupvSVKMNULgG8oGjoLGNDDcfoHjO7EH7KI5H3Epk8q
		//ifm1eElHUJJ/AMOQ9/nWG9VUCDvPA5qgVm6T/w6TtdcEWFXC0UZXZmPi0j17SR7F
		//AThS8QKBgQCCyPFJzwGIP99PcakQ38oFcoU8u/ahb0ghgJfSgK+K/ChXSyfbq5zt
		//jRkj6UWLa3plYX3po9h0Yp6f2IxnbOa3VK6fPkcSvBxhgK3RrugPerUJzFEPd3k4
		//GTqOBXtXO6N7zEMYxZxv0SgrV24LPfPz0aPObDeH6F0kuzXjanopIw==
		//-----END RSA PRIVATE KEY-----`
		espay = &Espay{
			BaseUrl:          config.ConfigJson.SandboxBaseUrl,
			MerchantCode:     helpers.DecryptAES(config.ConfigJson.SandboxMerchantCode),
			MerchantName:     helpers.DecryptAES(config.ConfigJson.SandboxMerchantName),
			ApiKey:           helpers.DecryptAES(config.ConfigJson.SandboxApiKey),
			SignatureKey:     helpers.DecryptAES(config.ConfigJson.SandboxSignatureKey),
			MerchantId:       helpers.DecryptAES(config.ConfigJson.SandboxMerchantId),
			MerchantPassword: helpers.DecryptAES(config.ConfigJson.SandboxCredentialPassword),
			PublicKey:        helpers.DecryptAES(config.ConfigJson.PublicKey),
			PrivateKey:       helpers.DecryptAES(config.ConfigJson.PrivateKey),
			httpClient:       http,
			serviceName:      "espay",
		}
	}

	espay.VirtualAccount = NewVirtualAccount(espay)
	espay.CreditCard = NewCreditCard(espay)
	espay.QRIS = NewQRIS(espay)
	espay.PaymentLink = NewPaymentLink(espay)

	espay.EspayStrategy = strategy.NewEspayStrategy(espay.VirtualAccount, espay.CreditCard, espay.QRIS, espay.PaymentLink)

	return espay
}

// Pay implements PaymentStrategy.Pay
func (e *Espay) Pay(ctx context.Context, req any) (map[string]interface{}, error) {
	paymentReq, ok := req.(PaymentRequest)
	if !ok {
		return nil, errors.New("invalid payment request type")
	}

	// get strategy
	strategyPayment, err := e.EspayStrategy.GetStrategy(paymentReq.Method)
	if err != nil {
		return nil, err
	}

	result, err := strategyPayment.Pay(ctx, paymentReq)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (e *Espay) MapResponsePayment(channel *models.Channel, payment *models.Payments) (*web.PaymentResponse, error) {
	log.Debug().Str("context", e.serviceName).Interface("payment", payment).Msg("map espay payment response")
	strategyPayment, err := e.EspayStrategy.GetStrategy(payment.Channel.PaymentMethod)
	if err != nil {
		return nil, err
	}

	response, err := strategyPayment.ClientResponse(payment, channel)
	if err != nil {
		payment.Channel.IsActive = false

		log.Error().Err(err).Str("context", e.serviceName).Msg("failed to map espay payment response")
		return nil, err
	}

	return response, nil
}

func (e *Espay) MapCheckStatusPayment(payment *models.Payments, statusPayment any) (*web.CheckStatusPaymentResponse, error) {

	return nil, nil
}

// CheckStatusPayment implements PaymentStrategy.CheckStatusPayment
func (e *Espay) CheckStatusPayment(ctx context.Context, req any) (any, error) {
	request, ok := req.(CheckStatusPaymentRequest)
	if !ok {
		return nil, errors.New("invalid check status payment request type")
	}

	// get strategy
	signature := NewSignature(e.SignatureKey, e.MerchantCode)

	asymmetricSignature, _ := signature.AsymmetricSignatureGenerate(nil, "")

	return e.CheckStatus(ctx, request, asymmetricSignature)
}

// check client key and server key correct
func (e *Espay) CheckKey() (any, error) {

	return nil, fmt.Errorf("not implemented")
}

//func (e *Espay)

func (e *Espay) SendInvoicePaymentVa(ctx context.Context, chargeReq map[string]string) (map[string]interface{}, error) {
	log.Debug().Interface("chargeReq", chargeReq).Msg("send invoice payment va")

	// Make API request to Espay
	url := fmt.Sprintf("%s%s", e.BaseUrl, PathSendInvoicePaymentVa)
	log.Debug().Str("url", url).Msg("send invoice payment va url")

	client := e.httpClient.Client.SetRetryCount(0).R()

	resp, err := client.SetContext(ctx).
		SetHeader("Accept", "*/*").
		//SetHeader("Authorization", e.GetBase64Authorization()).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(chargeReq).
		Post(url)
	if err != nil {
		log.Error().Err(err).Str("context", e.serviceName).Msg("failed to create espay payment")
		return nil, err
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		errMsg := fmt.Sprintf("espay API error: %v", string(resp.Body()))
		log.Error().Str("error", errMsg).Msg("espay payment failed")
		return nil, errors.New(errMsg)
	}

	log.Info().Interface("response", string(resp.Body())).Str("context", e.serviceName).Msg("espay payment success")

	// Parse response
	var result map[string]interface{}
	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		log.Error().Err(err).Str("context", e.serviceName).Msg("failed to parse espay response")
		return nil, err
	}

	// Check for error response
	if resp.StatusCode() != http.StatusOK {
		errMsg := fmt.Sprintf("espay API error: %v", result["status_message"])
		log.Error().Str("error", errMsg).Msg("espay payment failed")
		return nil, errors.New(errMsg)
	}

	errorCode, _ := result["error_code"]

	if errorCode.(any) != "0000" {
		errMsg := fmt.Sprintf("espay API error: %v", result["status_message"])
		return nil, errors.New(errMsg)
	}

	return result, nil
}

func (e *Espay) CheckStatus(ctx context.Context, payload CheckStatusPaymentRequest, signature string) (any, error) {
	//log.Debug().Interface("payload", payload).Msg("check status payment espay")
	//url := fmt.Sprintf("%s/apimerchant/v1.0/transfer-va/status", e.BaseUrl)
	//
	//e.httpClient.Client.R().SetContext(ctx).SetHeader("Content-Type", "application/json").
	//	SetHeader("X-TIMESTAMP", time.Now().String()).
	//	SetHeader("X-SIGNATURE", signature).
	//	SetHeader("X-EXTERNAL-ID", "uuid").
	//	SetHeader("X-PARTNER-ID", e.MerchantCode).
	//	SetHeader("CHANNEL-ID", " ESPAY")

	return nil, nil
}

type ResultVaStatic struct {
	ErrorCode    string            `json:"error_code"`
	ErrorMessage string            `json:"error_message"`
	RqUuid       string            `json:"rq_uuid"`
	RsDatetime   string            `json:"rs_datetime"`
	VaList       map[string]VaList `json:"va_list"`
}

type VaList struct {
	Amount         int    `json:"amount"`
	BankCode       string `json:"bank_code"`
	Description    string `json:"description"`
	ErrorCode      string `json:"error_code"`
	ErrorMessage   string `json:"error_message"`
	ExpiryDateTime string `json:"expiry_date_time"`
	Fee            string `json:"fee"`
	TotalAmount    string `json:"total_amount"`
	VaNumber       string `json:"va_number"`
}

func (e *Espay) CreateVaStatic(ctx context.Context, payload PaymentRequest) (map[string]VaList, error) {
	// signature
	signature := NewSignature(e.SignatureKey, e.MerchantCode)

	signatureSendInvoice := signature.HashBashSignatureGenerate(payload, enums.ENUM_SIGNATURE_ESPAY_SENDINVOICE)
	log.Debug().Interface("signature", signatureSendInvoice).Msg("signature")

	payloadVaStatic := e.payloadVaStatic(payload, signatureSendInvoice)

	result, err := e.SendInvoicePaymentVaStatic(ctx, payloadVaStatic)
	if err != nil {
		log.Error().Err(err).Msg("send invoice payment va static")
		return nil, fmt.Errorf("send invoice payment va static failed")
	}

	return result.VaList, nil
}

func (e *Espay) payloadVaStatic(payload PaymentRequest, signature string) map[string]string {
	return map[string]string{
		"rq_uuid":     payload.RQUUID,                           // uuid -> clear
		"rq_datetime": payload.RQDateTime.Format(time.DateTime), // time transaction -> clear
		"order_id":    payload.OrderID,                          // order id -> clear
		"amount":      "",
		"ccy":         string(payload.CCY),
		"comm_code":   e.MerchantCode,
		"remark1":     string(enums.STATIS_PHONE_NUMBER),
		"remark2":     payload.CustomerName,
		"remark3":     string(enums.STATIS_EMAIL),
		"remark4":     e.MerchantName,
		"update":      "N",
		"bank_code":   payload.BankCode,
		"va_expired":  fmt.Sprint(payload.VaExpired),
		"signature":   signature,
	}
}

func (e *Espay) SendInvoicePaymentVaStatic(ctx context.Context, chargeReq map[string]string) (*ResultVaStatic, error) {
	log.Debug().Interface("chargeReq", chargeReq).Msg("send invoice payment va")

	// Make API request to Espay
	url := fmt.Sprintf("%s%s", e.BaseUrl, PathSendInvoicePaymentVa)
	log.Debug().Str("url", url).Msg("send invoice payment va url")

	client := e.httpClient.Client.R()

	resp, err := client.SetContext(ctx).
		SetHeader("Accept", "*/*").
		//SetHeader("Authorization", e.GetBase64Authorization()).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(chargeReq).
		Post(url)
	if err != nil {
		log.Error().Err(err).Str("context", e.serviceName).Msg("failed to create espay payment")
		return nil, err
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		errMsg := fmt.Sprintf("espay API error: %v", string(resp.Body()))
		log.Error().Str("error", errMsg).Msg("espay payment failed")
		return nil, fmt.Errorf("send espay invoice payment va static failed")
	}

	log.Info().Interface("response", string(resp.Body())).Str("context", e.serviceName).Msg("espay payment success")

	// Parse response
	var result *ResultVaStatic
	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		log.Error().Err(err).Str("context", e.serviceName).Msg("failed to parse espay response")
		return nil, err
	}

	// Check for error response
	if resp.StatusCode() != http.StatusOK {
		errMsg := fmt.Sprintf("espay API error: %v", result.ErrorMessage)
		log.Error().Str("error", errMsg).Msg("espay payment failed")
		return nil, errors.New(errMsg)
	}

	errorCode := result.ErrorCode

	if errorCode != "0000" {
		errMsg := fmt.Sprintf("espay API error: %v", result.ErrorMessage)
		return nil, fmt.Errorf("send espay invoice payment va static failed %v", errMsg)
	}

	return result, nil
}

func (e *Espay) SendPaymentHostToHost(ctx context.Context, chargeReq CreditCardRequest, signature string) (map[string]interface{}, error) {
	log.Debug().Interface("chargeReq", chargeReq).Msg("send payment host to host")

	// Make API request to Espay
	url := fmt.Sprintf("%s%s", e.BaseUrl, PathPaymentHostToHost)
	log.Debug().Str("url", url).Msg("send payment host to host url")

	client := e.httpClient.Client.SetRetryCount(0).R()
	//
	//for key, value := range chargeReq.Headers(time.Now(), chargeReq.Signature, uuid.New().String(), e.MerchantCode) {
	//
	//}

	client = client.SetHeaders(chargeReq.Headers(time.Now(), signature, uuid.New().String(), e.MerchantCode))

	resp, err := client.SetContext(ctx).
		//SetHeader("Accept", "*/*").
		////SetHeader("Authorization", e.GetBase64Authorization()).
		SetHeader("Content-Type", "application/json").
		SetBody(chargeReq).
		Post(url)
	if err != nil {
		log.Error().Err(err).Str("context", e.serviceName).Msg("failed to create espay payment")
		return nil, err
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		errMsg := fmt.Sprintf("espay API error: %v", string(resp.Body()))
		log.Error().Str("error", errMsg).Msg("espay payment failed")
		return nil, errors.New(errMsg)
	}

	log.Info().Interface("response", string(resp.Body())).Str("context", e.serviceName).Msg("espay payment success")

	// Parse response
	var result map[string]interface{}
	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		log.Error().Err(err).Str("context", e.serviceName).Msg("failed to parse espay response")
		return nil, err
	}

	// Check for error response
	if resp.StatusCode() != http.StatusOK {
		errMsg := fmt.Sprintf("espay API error: %v", result["status_message"])
		log.Error().Str("error", errMsg).Msg("espay payment failed")
		return nil, errors.New(errMsg)
	}

	responseCode, _ := result["responseCode"]

	if responseCode.(any) != "2005400" {
		errMsg := fmt.Sprintf("espay API error: %v", result["responseMessage"])
		return nil, errors.New(errMsg)
	}

	return result, nil
}

func (e *Espay) SendPaymentQRIS(ctx context.Context, chargeReq QRISRequest, signature string) (map[string]interface{}, error) {
	log.Debug().Interface("chargeReq", chargeReq).Interface("signature", signature).Msg("send payment qris")

	// url qris
	url := fmt.Sprintf("%s%s", e.BaseUrl, PathPaymentQRIS)
	log.Debug().Interface("url", url).Msg("send payment qris url")

	client := e.httpClient.Client.R()

	// set headers
	client = client.SetHeaders(chargeReq.Headers(time.Now(), signature, uuid.New().String(), e.MerchantCode))

	resp, err := client.SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(chargeReq).
		Post(url)
	if err != nil {
		log.Error().Err(err).Str("context", e.serviceName).Msg("failed to create espay payment")
		return nil, err
	}
	log.Debug().Interface("body response", string(resp.Body())).Msg("response body espay qris")

	// Parse response
	var result map[string]interface{}
	if err = json.Unmarshal(resp.Body(), &result); err != nil {
		log.Error().Err(err).Str("context", e.serviceName).Msg("failed to parse espay response")
		return nil, err
	}

	// Check for error response
	if resp.StatusCode() != http.StatusOK {
		errMsg := fmt.Sprintf("espay API error: %v", result["responseMessage"])
		log.Error().Str("error", errMsg).Msg("espay payment failed")
		return nil, errors.New(errMsg)
	}

	responseCode, _ := result["responseCode"]

	if responseCode.(any) != "2004700" {
		errMsg := fmt.Sprintf("espay API error: %v", result["responseMessage"])
		return nil, errors.New(errMsg)
	}

	return result, nil
}

func (e *Espay) SendPaymentPaymentLink(ctx context.Context, req PaymentLinkRequest, generate string) (map[string]interface{}, error) {
	return nil, errors.New("not implemented")
}
