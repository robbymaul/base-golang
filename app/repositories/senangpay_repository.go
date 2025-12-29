package repositories

import (
	"context"
	clientsenangpay "paymentserviceklink/app/client/senangpay"
	"paymentserviceklink/app/models"

	"github.com/rs/zerolog/log"
)

func (rc *RepositoryContext) SenangpayPaymentRedirectUrlRepository(ctx context.Context, paymentRequest *clientsenangpay.PaymentRequest) (*models.GatewayResponse, error) {
	senangpayUrl := rc.Senangpay.GeneratePaymentURL(paymentRequest)

	log.Debug().Str("senangpay url", senangpayUrl).Msg("value generate payment url senangpay")

	err := rc.Senangpay.Send(ctx, senangpayUrl)
	if err != nil {
		log.Error().Err(err).Msg("failed to send senangpay generate payment url")
		return nil, err
	}

	return &models.GatewayResponse{
		RedirectUrl: senangpayUrl,
	}, nil
}

func (rc *RepositoryContext) CheckStatusPaymentSenangpayRepository(ctx context.Context, transactionId string) (*models.CallbackData, error) {
	//result, err := rc.Senangpay.CheckStatusPayment(ctx, transactionId)
	//if err != nil {
	//	log.Error().Err(err).Msg("failed to check status payment")
	//	return nil, err
	//}
	////
	////var name string
	////var email string
	////var phone string
	////var amount string
	////var message string
	////var orderId string
	////name = result.Data[0].BuyerContact.Name
	////email = result.Data[0].BuyerContact.Email
	////phone = result.Data[0].BuyerContact.Phone
	////amount = result.Data[0].OrderDetail.GrandTotal
	////message = result.Msg
	////orderId = result.Data[0].PaymentInfo.TransactionReference
	////
	////return &models.CallbackData{
	////	Name:          name,
	////	Email:         email,
	////	Phone:         phone,
	////	AmountPaid:    amount,
	////	TxnStatus:     enums.TxnStatusSenangpay(fmt.Sprint(result.Status)),
	////	TxnMessage:    message,
	////	OrderId:       "",
	////	TransactionId: orderId,
	////	HashedValue:   "",
	////}, nil
	//
	//return result, nil

	return nil, nil
}
