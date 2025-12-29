package clientsenangpay

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	restyclient "paymentserviceklink/app/client/resty"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/web"
	"strings"

	"github.com/rs/zerolog/log"
)

type Senangpay struct {
	SenangpayUrl string
	SecretKey    string
	MerchantId   string
	serviceName  string
	httpClient   *restyclient.RestyClient
}

type SenangpayResponse struct {
	RedirectUrl string `json:"redirectUrl"`
}

func NewSenangpay(httpClient *restyclient.RestyClient, config *models.Configuration) *Senangpay {
	log.Debug().Interface("config", config).Msg("new senangpay")
	if config.ConfigValue == enums.PRODUCTION {
		return &Senangpay{
			SenangpayUrl: config.ConfigJson.ProductionBaseUrl,
			SecretKey:    helpers.DecryptAES(config.ConfigJson.ProductionSecretKey),
			MerchantId:   helpers.DecryptAES(config.ConfigJson.ProductionMerchantId),
			serviceName:  "senangpay",
			httpClient:   httpClient,
		}
	}

	return &Senangpay{
		SenangpayUrl: config.ConfigJson.SandboxBaseUrl,
		SecretKey:    helpers.DecryptAES(config.ConfigJson.SandboxSecretKey),
		MerchantId:   helpers.DecryptAES(config.ConfigJson.SandboxMerchantId),
		serviceName:  "senangpay",
		httpClient:   httpClient,
	}
}

func (s *Senangpay) Pay(ctx context.Context, req any) (map[string]interface{}, error) {
	log.Debug().Interface("req", req).Msg("senangpay pay")
	var generateUrl string
	log.Debug().Str("generateUrl", generateUrl).Msg("generate url")
	if request, ok := req.(PaymentRequest); ok {
		generateUrl = s.GeneratePaymentURL(&request)
		log.Debug().Str("generateUrl", generateUrl).Msg("generate url")
	} else {
		return nil, errors.New("invalid payment request request type")
	}

	// Contoh request
	_, err := s.httpClient.Client.R().
		SetContext(ctx).
		Get(generateUrl)
	if err != nil {
		return nil, err
	}

	// Debug responsenya
	//println(resp.Status())
	//println(resp.String())

	return map[string]interface{}{
		"redirectUrl": generateUrl,
	}, nil
}

func (s *Senangpay) MapCheckStatusPayment(payment *models.Payments, statusPayment any) (*web.CheckStatusPaymentResponse, error) {
	data, ok := statusPayment.(CheckStatusPaymentSenangpayResponse)
	if !ok {
		return nil, helpers.NewErrorTrace(errors.New("invalid status payment response type"), s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	status := enums.TxnStatusSenangpay(fmt.Sprint(data.Status))

	if status == enums.SENANGPAY_STATUS_FAILED {
		return nil, helpers.NewErrorTrace(fmt.Errorf(data.Msg), s.serviceName).WithStatusCode(http.StatusPaymentRequired)
	}

	currentStatus := enums.PAYMENT_STATUS_SUCCESS

	return &web.CheckStatusPaymentResponse{
		Id:            payment.Id,
		TransactionId: payment.TransactionId,
		OrderId:       payment.OrderId,
		Status:        currentStatus,
		Amount:        payment.TotalAmount.IntPart(),
		Currency:      payment.Currency,
	}, nil
}

func (s *Senangpay) MapResponsePayment(channel *models.Channel, payment *models.Payments) (*web.PaymentResponse, error) {
	var senangpayResponse SenangpayResponse

	err := json.Unmarshal(payment.GatewayResponse, &senangpayResponse)
	if err != nil {
		return nil, err
	}

	return &web.PaymentResponse{
		Id:            payment.Id,
		TransactionId: payment.TransactionId,
		OrderId:       payment.OrderId,
		Status:        payment.Status,
		Amount:        payment.Amount.IntPart(),
		FeeAdmin:      payment.FeeAmount.IntPart(),
		TotalAmount:   payment.TotalAmount.IntPart(),
		Currency:      payment.Currency,
		PaymentMethod: payment.Channel.PaymentMethod,
		PaymentType:   enums.PAYMENT_TYPE_REDIRECT,
		PaymentDetail: web.PaymentDetail{
			Bank: enums.CHANNEL_SENANGPAY,
			Url: []*web.Actions{
				{
					Name:   "url payment",
					Method: http.MethodGet,
					Url:    senangpayResponse.RedirectUrl,
				},
			},
			VaNumber:        "",
			TransactionTime: "",
			ExpireTime:      "",
			Instruction:     channel.Instruction,
		},
		Customer: web.Customer{
			MemberId: payment.CustomerId,
			Name:     payment.CustomerName,
			Email:    payment.CustomerEmail,
			Phone:    payment.CustomerPhone,
		},
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	}, nil
}

func (s *Senangpay) Send(ctx context.Context, url string) error {
	// Contoh request
	resp, err := s.httpClient.Client.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return err
	}

	// Debug responsenya
	println(resp.Status())
	println(resp.String())

	return nil
}

func (s *Senangpay) CheckStatusPayment(ctx context.Context, request any) (any, error) {
	checkStatusPaymentRequest, ok := request.(CheckStatusPaymentRequest)
	if !ok {
		return nil, helpers.NewErrorTrace(errors.New("invalid request type"), s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

	subscribe := NewSubscribe(s, &checkStatusPaymentRequest)
	nonSubscribe := NewNonSubscribe(s, &checkStatusPaymentRequest)

	subscribeStrategy := NewSubscribeStrategy(subscribe, nonSubscribe)

	setSubscribe := subscribeStrategy.SetSubscribe(checkStatusPaymentRequest.TransactionId)

	generateUrl := setSubscribe.GenerateUrl()

	result, err := setSubscribe.CheckStatusPayment(ctx, generateUrl)
	if err != nil {
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return result, nil

	//url := s.GenerateStatusPaymentURL(checkStatusPaymentRequest.TransactionId)

	//generateUrl := s.GenerateQueryOrderStatusURL(checkStatusPaymentRequest.OrderId)
	//
	//log.Debug().Str("generateUrl", generateUrl).Msg("url check status payment")
	//
	//resp, err := client.R().
	//	SetContext(ctx).
	//	Get(generateUrl)
	//if err != nil {
	//	return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	//}
	//
	//Ambil raw response body
	//body := resp.Body()
	//bodyStr := string(body)
	//
	//log.Debug().Str("raw_body", bodyStr).Msg("raw response body")

	// Perbaiki casing field "Note" → "note"
	//bodyStr = strings.ReplaceAll(bodyStr, `"Note":`, `"note":`)

	// Parse kembali ke struct
	//var result CheckStatusPaymentSenangpayResponse
	//if err = json.Unmarshal([]byte(bodyStr), &result); err != nil {
	//	log.Error().Err(err).Msg("failed to unmarshal senangpay response")
	//	return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	//}
	//
	//log.Debug().Interface("response", result).Msg("response from check status payment senangpay")

	//var name string
	//var email string
	//var phone string
	//var amount string
	//var message string
	//var orderId string
	//name = result.Data[0].BuyerContact.Name
	//email = result.Data[0].BuyerContact.Email
	//phone = result.Data[0].BuyerContact.Phone
	//amount = result.Data[0].OrderDetail.GrandTotal
	//message = result.Msg
	//orderId = result.Data[0].PaymentInfo.TransactionReference

	//return &models.CallbackData{
	//	Name:          name,
	//	Email:         email,
	//	Phone:         phone,
	//	AmountPaid:    amount,
	//	TxnStatus:     enums.TxnStatusSenangpay(fmt.Sprint(result.Status)),
	//	TxnMessage:    message,
	//	OrderId:       "",
	//	TransactionId: orderId,
	//	HashedValue:   "",
	//}, nil

	//return result, nil
}

func (s *Senangpay) CheckKey() (any, error) {

	return nil, nil
}

func (s *Senangpay) GenerateHashMD5(args ...string) string {
	data := strings.Join(args, "")
	sum := md5.Sum([]byte(data))
	return hex.EncodeToString(sum[:])
}

func (s *Senangpay) VerifyHashMD5(ctx context.Context, payload *VerifyPayment) bool {
	data := s.SecretKey + payload.StatusID + payload.OrderID + payload.TransactionID + payload.Message

	log.Debug().
		Str("secretKey", s.SecretKey).
		Str("statusID", payload.StatusID).
		Str("orderID", payload.OrderID).
		Str("transactionID", payload.TransactionID).
		Str("message", payload.Message).
		Str("data", data).
		Msg("String before MD5")

	sum := md5.Sum([]byte(data))
	expectedHash := hex.EncodeToString(sum[:])

	log.Debug().
		Str("expectedHash", expectedHash).
		Str("receivedHash", payload.ReceiveHash).
		Msg("hash compare")

	return strings.EqualFold(expectedHash, payload.ReceiveHash)
}

func (s *Senangpay) GenerateHash(args ...string) string {
	data := strings.Join(args, "")
	h := hmac.New(sha256.New, []byte(s.SecretKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *Senangpay) VerifyHash(ctx context.Context, payload *VerifyPayment) bool {
	// Format data sesuai dokumentasi SenangPay:
	// secretKey + statusID + orderID + transactionID + message
	data := s.SecretKey + payload.StatusID + payload.OrderID + payload.TransactionID + payload.Message

	log.Debug().
		Str("secretKey", s.SecretKey).
		Str("statusID", payload.StatusID).
		Str("orderID", payload.OrderID).
		Str("transactionID", payload.TransactionID).
		Str("message", payload.Message).
		Str("data", data).
		Msg("String before HMAC SHA256")

	// Hitung HMAC-SHA256
	// Note: SecretKey digunakan sebagai kunci HMAC DAN termasuk dalam data yang di-hash
	h := hmac.New(sha256.New, []byte(s.SecretKey))
	h.Write([]byte(data))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	log.Debug().
		Str("expectedHash", expectedHash).
		Str("receivedHash", payload.ReceiveHash).
		Msg("hash compare")

	// Bandingkan dengan case-insensitive
	return strings.EqualFold(expectedHash, payload.ReceiveHash)
}

func (s *Senangpay) GeneratePaymentURL(payload *PaymentRequest) string {
	log.Debug().
		Str("MerchantId", s.MerchantId).
		Str("SecretKey", s.SecretKey).
		Str("Detail", payload.Detail).
		Str("Amount", payload.Amount).
		Str("OrderID", payload.OrderID).
		Str("Name", payload.Name).
		Str("Email", payload.Email).
		Str("Phone", payload.Phone).
		Msg("GeneratePaymentURL")
	hash := s.GenerateHash(s.SecretKey, payload.Detail, payload.Amount, payload.OrderID)
	log.Debug().
		Str("hash", hash).
		Msg("GenerateHash")

	params := url.Values{}
	params.Set("detail", payload.Detail)
	params.Set("amount", payload.Amount)
	params.Set("order_id", payload.OrderID)
	params.Set("name", payload.Name)
	params.Set("email", payload.Email)
	params.Set("phone", payload.Phone)
	params.Set("hash", hash)

	return fmt.Sprintf("%s%s", s.SenangpayUrl, fmt.Sprintf("/payment/%s?%s", s.MerchantId, params.Encode()))
}

func (s *Senangpay) GeneratePaymentURLMD5(payload *PaymentRequest) string {
	hash := s.GenerateHashMD5(payload.Detail, payload.Amount, payload.OrderID)

	params := url.Values{}
	params.Set("detail", payload.Detail)
	params.Set("amount", payload.Amount)
	params.Set("order_id", payload.OrderID)
	params.Set("name", payload.Name)
	params.Set("email", payload.Email)
	params.Set("phone", payload.Phone)
	params.Set("hash", hash)

	return fmt.Sprintf(s.SenangpayUrl, s.MerchantId, params.Encode())
}

func (s *Senangpay) GenerateStatusPaymentURL(paymentId string) string {
	hash := s.GenerateHash(s.MerchantId, s.SecretKey, paymentId)

	params := url.Values{}
	params.Add("merchant_id", s.MerchantId)
	params.Add("transaction_reference", paymentId)
	params.Add("hash", hash)

	return fmt.Sprintf("%v/apiv1/query_transaction_status?%s", s.SenangpayUrl, params.Encode())
}

func (s *Senangpay) GenerateQueryOrderStatusURL(orderId string) string {
	hash := s.GenerateHashMD5(s.MerchantId, s.SecretKey, orderId)

	params := url.Values{}
	params.Add("merchant_id", s.MerchantId)
	params.Add("order_id", orderId)
	params.Add("hash", hash)

	return fmt.Sprintf("%v/apiv1/query_order_status?%s", s.SenangpayUrl, params.Encode())
}
