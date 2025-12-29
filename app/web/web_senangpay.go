package web

import "paymentserviceklink/app/enums"

type WebhookCallbackSenangpay struct {
	Name          string                   `form:"name"`
	Email         string                   `form:"email"`
	Phone         string                   `form:"phone"`
	AmountPaid    string                   `form:"amount_paid"`
	TxnStatus     enums.TxnStatusSenangpay `form:"txn_status"`
	TxnMessage    string                   `form:"txn_msg"`
	OrderId       string                   `form:"order_id"`
	TransactionId string                   `form:"transaction_id"`
	HashedValue   string                   `form:"hashed_value"`
}
