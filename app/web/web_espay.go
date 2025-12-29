package web

import "time"

type EspayPaymentNotificationRequest struct {
	PartnerServiceId string                            `json:"partnerServiceId"`
	CustomerNo       string                            `json:"customerNo"`
	VirtualAccountNo string                            `json:"virtualAccountNo"`
	TrxId            string                            `json:"trxId"`
	PaymentRequestId string                            `json:"paymentRequestId"`
	PaidAmount       PaymentNotificationPaidAmount     `json:"paidAmount"`
	TotalAmount      PaymentNotificationTotalAmount    `json:"totalAmount"`
	TrxDateTime      time.Time                         `json:"trxDateTime"`
	AdditionalInfo   PaymentNotificationAdditionalInfo `json:"additionalInfo"`
}

type PaymentNotificationTotalAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type PaymentNotificationPaidAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type PaymentNotificationAdditionalInfo struct {
	TransactionStatus string                                   `json:"transactionStatus"`
	MemberCode        string                                   `json:"memberCode"`
	DebitFrom         string                                   `json:"debitFrom"`
	DebitFromName     string                                   `json:"debitFromName"`
	DebitFromBank     string                                   `json:"debitFromBank"`
	CreditTo          string                                   `json:"creditTo"`
	CreditToName      string                                   `json:"creditToName"`
	CreditToBank      string                                   `json:"creditToBank"`
	ProductCode       string                                   `json:"productCode"`
	ProductValue      string                                   `json:"productValue"`
	Message           PaymentNotificationAdditionalInfoMessage `json:"message"`
	FeeType           string                                   `json:"feeType"`
	TxFee             string                                   `json:"txFee"`
	PaymentRef        string                                   `json:"paymentRef"`
	PaymentRemark     string                                   `json:"paymentRemark"`
	Rrn               string                                   `json:"rrn"`
	ApprovalCode      string                                   `json:"approvalCode"`
	Token             string                                   `json:"token"`
	UserId            string                                   `json:"userId"`
}
type PaymentNotificationAdditionalInfoMessage struct {
	CHANNELFLAG interface{} `json:"CHANNEL_FLAG"`
}

type EspayPaymentNotificationResponse struct {
	ResponseCode       string `json:"responseCode"`
	ResponseMessage    string `json:"responseMessage"`
	VirtualAccountData struct {
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
	} `json:"virtualAccountData"`
	AdditionalInfo struct {
		ReconcileId       string    `json:"reconcileId"`
		ReconcileDatetime time.Time `json:"reconcileDatetime"`
	} `json:"additionalInfo"`
}

type EspayInquiryRequest struct {
	PartnerServiceId string `json:"partnerServiceId"`
	CustomerNo       string `json:"customerNo"`
	VirtualAccountNo string `json:"virtualAccountNo"`
	TrxDateInit      string `json:"trxDateInit"`
	InquiryRequestId string `json:"inquiryRequestId"`
}

type EspayInquiryResponse struct {
	ResponseCode       string                   `json:"responseCode"`
	ResponseMessage    string                   `json:"responseMessage"`
	VirtualAccountData *EspayVirtualAccountData `json:"virtualAccountData,omitempty"`
}

type EspayVirtualAccountData struct {
	PartnerServiceId    string                             `json:"partnerServiceId"`
	CustomerNo          string                             `json:"customerNo"`
	VirtualAccountNo    string                             `json:"virtualAccountNo"`
	VirtualAccountName  string                             `json:"virtualAccountName"`
	VirtualAccountEmail string                             `json:"virtualAccountEmail"`
	VirtualAccountPhone string                             `json:"virtualAccountPhone"`
	InquiryRequestId    string                             `json:"inquiryRequestId"`
	TotalAmount         *EspayVirtualAccountTotalAmount    `json:"totalAmount"`
	BillDetails         []*EspayVirtualAccountBillDetails  `json:"billDetails"`
	AdditionalInfo      *EspayVirtualAccountAdditionalInfo `json:"additionalInfo,omitempty"`
}

type EspayVirtualAccountTotalAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type EspayVirtualAccountBillDetailsDescription struct {
	English   string `json:"english"`
	Indonesia string `json:"indonesia"`
}

type EspayVirtualAccountBillDetails struct {
	BillDescription *EspayVirtualAccountBillDetailsDescription `json:"billDescription,omitempty"`
}

type EspayVirtualAccountAdditionalInfo struct {
	Token           string                                            `json:"token,omitempty"`
	TransactionDate time.Time                                         `json:"transactionDate"`
	DataMerchant    *EspayVirtualAccountAdditionalInfoDataMerchant    `json:"dataMerchant,omitempty"`
	ShippingAddress *EspayVirtualAccountAdditionalInfoShippingAddress `json:"shippingAddress,omitempty"`
	Items           []*EspayVirtualAccountAdditionalInfoItems         `json:"items"`
}

type EspayVirtualAccountAdditionalInfoDataMerchant struct {
	KodeCa         string                                                    `json:"kodeCa"`
	KodeSubCa      string                                                    `json:"kodeSubCa"`
	NoKontrak      string                                                    `json:"noKontrak"`
	NamaPelanggan  string                                                    `json:"namaPelanggan"`
	AngsuranKe     string                                                    `json:"angsuranKe"`
	JmlBayarExcAdm float64                                                   `json:"jmlBayarExcAdm"`
	Denda          float64                                                   `json:"denda"`
	FeeCa          float64                                                   `json:"feeCa"`
	FeeSwitcher    float64                                                   `json:"feeSwitcher"`
	TotalAdmin     float64                                                   `json:"totalAdmin"`
	JumlahBayar    float64                                                   `json:"jumlahBayar"`
	MinimumAmount  float64                                                   `json:"minimumAmount"`
	TotalAngsuran  int                                                       `json:"totalAngsuran"`
	Customer       *EspayVirtualAccountAdditionalInfoDataMerchantCustomer    `json:"customer,omitempty"`
	JatuhTempo     string                                                    `json:"jatuhTempo"`
	ListBills      []*EspayVirtualAccountAdditionalInfoDataMerchantListBills `json:"listBills"`
}

type EspayVirtualAccountAdditionalInfoShippingAddress struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Address     string `json:"address"`
	City        string `json:"city"`
	PostalCode  string `json:"postalCode"`
	PhoneNumber string `json:"phoneNumber"`
	CountryCode string `json:"countryCode"`
}

type EspayVirtualAccountAdditionalInfoItems struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Type     string `json:"type"`
	Url      string `json:"url"`
	Quantity string `json:"quantity"`
}

type EspayVirtualAccountAdditionalInfoDataMerchantCustomer struct {
	Email string `json:"email"`
}

type EspayVirtualAccountAdditionalInfoDataMerchantListBills struct {
	BillCode   interface{} `json:"billCode"`
	BillName   string      `json:"billName"`
	BillAmount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"billAmount"`
	BillSubCompany string `json:"billSubCompany"`
}

//
//type EspayTopupNotificationRequest struct {
//	RQUUID             string `json:"rq_uuid" form:"rq_uuid"`
//	RQDateTime         string `json:"rq_datetime" form:"rq_datetime"`
//	SenderId           string `json:"sender_id" form:"sender_id"`
//	ReceiverId         string `json:"receiver_id" form:"receiver_id"`
//	Password           string `json:"password" form:"password"`
//	CommCode           string `json:"comm_code" form:"comm_code"`
//	MemberCode         string `json:"member_code" form:"member_code"`
//	MemberCustomerId   string `json:"member_cust_id" form:"member_cust_id"`
//	MemberCustomerName string `json:"member_cust_name" form:"member_cust_name"`
//	CCY                string `json:"ccy" form:"ccy"`
//	Amount             string `json:"amount" form:"amount"`
//	DebitFrom          string `json:"debit_from" form:"debit_from"`
//	DebitFromName      string `json:"debit_from_name" form:"debit_from_name"`
//	DebitFromBank      string `json:"debit_from_bank" form:"debit_from_bank"`
//	CreditTo           string `json:"credit_to" form:"credit_to"`
//	CreditToName       string `json:"credit_to_name" form:"credit_to_name"`
//	CreditToBank       string `json:"credit_to_bank" form:"credit_to_bank"`
//	PaymentDatetime    string `json:"payment_datetime" form:"payment_datetime"`
//	PaymentRef         string `json:"payment_ref" form:"payment_ref"`
//	PaymentRemark      string `json:"payment_remark" form:"payment_remark"`
//	OrderId            string `json:"order_id" form:"order_id"`
//	ProductCode        string `json:"product_code" form:"product_code"`   // product code bank espay
//	ProductValue       string `json:"product_value" form:"product_value"` // number virtual account
//	Message            string `json:"message" form:"message"`
//	Status             string `json:"status" form:"status"`
//	Token              string `json:"token" form:"token"`
//	TotalAmount        string `json:"total_amount" form:"total_amount"`  // jumlah topup
//	TxKey              string `json:"tx_key" form:"tx_key"`              // kode pembayaran transaksi key di espay
//	FeeType            string `json:"fee_type" form:"fee_type"`          // biaya admin
//	TxStatus           string `json:"tx_status" form:"tx_status"`        // status pembayaran
//	UserId             string `json:"user_id" form:"user_id"`            // gen va k-wallet
//	ReferenceId        string `json:"reference_id" form:"reference_id" ` // merchant name espay
//	IsSnap             string `json:"is_snap" form:"is_snap"`
//	MemberId           string `json:"member_id" form:"member_id"` // merchant code espay
//	Signature          string `json:"signature" form:"signature"`
//}

type EspayTopupNotificationRequest struct {
	PartnerServiceId string `json:"partnerServiceId"`
	CustomerNo       string `json:"customerNo"`
	VirtualAccountNo string `json:"virtualAccountNo"`
	TrxId            string `json:"trxId"`
	PaymentRequestId string `json:"paymentRequestId"`
	PaidAmount       struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"paidAmount"`
	TotalAmount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"totalAmount"`
	TrxDateTime    time.Time `json:"trxDateTime"`
	AdditionalInfo struct {
		TransactionStatus string `json:"transactionStatus"`
		MemberCode        string `json:"memberCode"`
		DebitFrom         string `json:"debitFrom"`
		DebitFromName     string `json:"debitFromName"`
		DebitFromBank     string `json:"debitFromBank"`
		CreditTo          string `json:"creditTo"`
		CreditToName      string `json:"creditToName"`
		CreditToBank      string `json:"creditToBank"`
		ProductCode       string `json:"productCode"`
		ProductValue      string `json:"productValue"`
		Message           struct {
			CHANNELFLAG interface{} `json:"CHANNEL_FLAG"`
		} `json:"message"`
		FeeType       string `json:"feeType"`
		TxFee         string `json:"txFee"`
		PaymentRef    string `json:"paymentRef"`
		PaymentRemark string `json:"paymentRemark"`
		Rrn           string `json:"rrn"`
		ApprovalCode  string `json:"approvalCode"`
		Token         string `json:"token"`
		UserId        string `json:"userId"`
	} `json:"additionalInfo"`
	BillDetails []struct {
		BillNo string `json:"billNo"`
	} `json:"billDetails"`
}
