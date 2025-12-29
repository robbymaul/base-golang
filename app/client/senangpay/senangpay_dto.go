package clientsenangpay

type CallbackRequest struct {
	StatusID      string `form:"status_id"`
	OrderID       string `form:"order_id"`
	TransactionID string `form:"transaction_id"`
	Msg           string `form:"msg"`
	Hash          string `form:"hash"`
}

type PaymentRequest struct {
	OrderID string
	Amount  string
	Detail  string
	Name    string
	Email   string
	Phone   string
}

type VerifyPayment struct {
	StatusID      string `form:"txn_status" json:"txn_status"`
	OrderID       string `form:"order_id" json:"order_id"`
	TransactionID string `form:"transaction_id" json:"transaction_id"`
	Message       string `form:"txn_msg" json:"txn_msg"`
	ReceiveHash   string `form:"hashed_value" json:"hashed_value"`
	Amount        string `form:"amount_paid" json:"amount_paid"`
}

type CheckStatusPaymentSenangpayResponse struct {
	Status int64                                      `json:"status"`
	Msg    string                                     `json:"msg"`
	Data   []*DataCheckStatusPaymentSenangpayResponse `json:"data"`
}

type DataCheckStatusPaymentSenangpayResponse struct {
	BuyerContact    BuyerContact    `json:"buyer_contact"`
	DeliveryAddress DeliveryAddress `json:"delivery_address"`
	OrderDetail     OrderDetail     `json:"order_detail"`
	PaymentInfo     PaymentInfo     `json:"payment_info"`
	DiscountInfo    DiscountInfo    `json:"discount_info"`
	NetworkInfo     NetworkInfo     `json:"network_info"`
	Product         Product         `json:"product"`
}

type DiscountInfo struct {
	Code          string `json:"code"`
	Amount        string `json:"amount"`
	Type          string `json:"type"`
	ValidityStart string `json:"validity_start"`
	ValidityEnd   string `json:"validity_end"`
}

type PaymentInfo struct {
	TransactionReference string `json:"transaction_reference"`
	TransactionDate      string `json:"transaction_date"`
	PaymentMode          string `json:"payment_mode"`
	Status               string `json:"status"`
}

type OrderDetail struct {
	Quantity       string `json:"quantity"`
	UnitPrice      string `json:"unit_price"`
	DeliveryCharge string `json:"delivery_charge"`
	Gst            string `json:"gst"`
	GrandTotal     string `json:"grand_total"`
}

type DeliveryAddress struct {
	Address1 string `json:"address_1"`
	Address2 string `json:"address_2"`
	City     string `json:"city"`
	Postcode string `json:"postcode"`
	State    string `json:"state"`
	Country  string `json:"country"`
}

type BuyerContact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Note  string `json:"Note"`
}

type NetworkInfo struct {
	Referer string `json:"referer"`
	Ip      string `json:"ip"`
	City    string `json:"city"`
	Country string `json:"country"`
	Browser string `json:"browser"`
}

type Product struct {
	ProductName string `json:"product_name"`
	Description string `json:"description"`
}

type CheckStatusPaymentRequest struct {
	TransactionId string `json:"transactionId"`
	OrderId       string `json:"orderId"`
}
