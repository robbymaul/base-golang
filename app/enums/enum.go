package enums

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type AggregatorName string

const (
	AGGREGATOR_NAME_MIDTRANS  = "MIDTRANS"
	AGGREGATOR_NAME_ESPAY     = "ESPAY"
	AGGREGATOR_NAME_SENANGPAY = "SENANGPAY"
)

type FeeType string

const (
	FEE_TYPE_FIXED            FeeType = "fixed"
	FEE_TYPE_PERCENTAGE       FeeType = "percentage"
	FEE_TYPE_NONE             FeeType = "none"
	FEE_TYPE_FIXED_PERCENTAGE FeeType = "fixed_percentage"
)

type ActionAdminActivityLog string

const (
	ACTION_ADMIN_ACTVITY_LOG_LOGIN          ActionAdminActivityLog = "login"
	ACTION_ADMIN_ACTVITY_LOG_LOGOUT         ActionAdminActivityLog = "logout"
	ACTION_ADMIN_ACTVITY_LOG_CANCEL_PAYMENT ActionAdminActivityLog = "cancel_payment"
	ACTION_ADMIN_ACCTIVITY_CREATE           ActionAdminActivityLog = "create"
	ACTION_ADMIN_ACCTIVITY_UPDATE           ActionAdminActivityLog = "update"
	ACTION_ADMIN_ACCTIVITY_DELETE           ActionAdminActivityLog = "delete"
)

type NameAdminRoles string

const (
	NAME_ADMIN_ROLES_SUPER_ADMIN NameAdminRoles = "super_admin"
	NAME_ADMIN_ROLES_ADMIN       NameAdminRoles = "admin"
	NAME_ADMIN_ROLE_OPERATOR     NameAdminRoles = "operator"
)

type CodeAdminRole string

const (
	CODE_ADMIN_ROLES_SUPER_ADMIN CodeAdminRole = "super_admin"
	CODE_ADMIN_ROLES_ADMIN       CodeAdminRole = "admin"
	CODE_ADMIN_ROLE_OPERATOR     CodeAdminRole = "operator"
)

type CodePlatform string

const (
	CODE_PLATFORM_KNET CodePlatform = "web_knet"
	CODE_PLATFORM_SMS  CodePlatform = "web_sms"
)

type CodePaymentMethod string

const (
	CODE_PAYMENT_METHOD_BANK_TRANSFER CodePaymentMethod = "bank_transfer"
	CODE_PAYMENT_METHOD_E_WALLET      CodePaymentMethod = "e_wallet"
	CODE_PAYMENT_METHOD_CREDIT_CARD   CodePaymentMethod = "credit_card"
	CODE_PAYMENT_METHOD_VA            CodePaymentMethod = "va"
	CODE_PAYMENT_METHOD_QRIS          CodePaymentMethod = "qris"
)

type PaymentMethod string

const (
	//PAYMENT_METHOD_BCA_VA          PaymentMethod = "bca_va"
	//PAYMENT_METHOD_BNI_VA          PaymentMethod = "bni_va"
	//PAYMENT_METHOD_BRI_VA          PaymentMethod = "bri_va"
	//PAYMENT_METHOD_CIMB_VA         PaymentMethod = "cimb_va"
	//PAYMENT_METHOD_GOPAY           PaymentMethod = "gopay"
	//PAYMENT_METHOD_MANDIRI_VA      PaymentMethod = "mandiri_va"
	//PAYMENT_METHOD_PERMATA_VA      PaymentMethod = "permata_va"
	//PAYMENT_METHOD_SHOPEE_PAY      PaymentMethod = "shopee_pay"
	//PAYMENT_METHOD_VA              PaymentMethod = "va"
	PAYMENT_METHOD_SENANGPAY       PaymentMethod = "SENANGPAY"
	PAYMENT_METHOD_CREDIT_CARD     PaymentMethod = "CREDIT_CARD"
	PAYMENT_METHOD_BANK_TRANSFER   PaymentMethod = "BANK_TRANSFER"
	PAYMENT_METHOD_VIRTUAL_ACCOUNT PaymentMethod = "VIRTUAL_ACCOUNT"
	PAYMENT_METHOD_QRIS            PaymentMethod = "QRIS"
	PAYMENT_METHOD_GOPAY           PaymentMethod = "GOPAY"
	PAYMENT_METEHOD_K_WALLET       PaymentMethod = "K-WALLET"
	PAYMENT_METHOD_E_WALLET        PaymentMethod = "E-WALLET"
	PAYMENT_METHOD_RETAIL_OUTLET   PaymentMethod = "RETAIL_OUTLET"
	PAYMENT_METHOD_MULTI_PAYMENT   PaymentMethod = "MULTI_PAYMENT"
)

type PaymentType string

const (
	PAYMENT_TYPE_VA       PaymentType = "va"
	PAYMENT_TYPE_REDIRECT PaymentType = "redirect"
	PAYMENT_TYPE_QRIS     PaymentType = "qris"
	PAYMENT_BILL          PaymentType = "bill"
)

type Channel string

const (
	CHANNEL_BCA       Channel = "BCA"
	CHANNEL_BNI       Channel = "BNI"
	CHANNEL_BRI       Channel = "BRI"
	CHANNEL_CIMB      Channel = "CIMB"
	CHANNEL_MANDIRI   Channel = "MANDIRI"
	CHANNEL_PERMATA   Channel = "PERMATA"
	CHANNEL_GOPAY     Channel = "GOPAY"
	CHANNEL_SHOPEE    Channel = "SHOPEE"
	CHANNEL_SENANGPAY Channel = "SENANGPAY"
	CHANNEL_DANAMON   Channel = "DANAMON"
	CHANNEL_MAYBANK   Channel = "MAYBANK"
	CHANNEL_K_WALLET  Channel = "K-WALLET"
	CHANNEL_QRIS      Channel = "QRIS"
	CHANNEL_DANA      Channel = "DANA"
)

type CodePaymentType string

const (
	CODE_PAYMENT_TYPE_SALES_ORDER  CodePaymentType = "sales_order"
	CODE_PAYMENT_TYPE_TOPUP_TOKEN  CodePaymentType = "topup_token"
	CODE_PAYMENT_TYPE_TOPUP_WALLET CodePaymentType = "topup_wallet"
)

type PaymentStatus string

const (
	PAYMENT_STATUS_PENDING    PaymentStatus = "pending"
	PAYMENT_STATUS_PROCESSING PaymentStatus = "processing"
	PAYMENT_STATUS_SUCCESS    PaymentStatus = "success"
	PAYMENT_STATUS_FAILED     PaymentStatus = "failed"
	PAYMENT_STATUS_CANCELLED  PaymentStatus = "cancelled"
	PAYMENT_STATUS_EXPIRED    PaymentStatus = "expired"
)

type ProviderPaymentMethod string

const (
	PROVIDER_PAYMENT_METHOD_MIDTRANS  ProviderPaymentMethod = "midtrans"
	PROVIDER_PAYMENT_METHOD_SENANGPAY ProviderPaymentMethod = "senangpay"
	PROVIDER_PAYMENT_METHOD_ESPAY     ProviderPaymentMethod = "espay"
)

type ResourceAdminActivityLog string

const (
	RESOURCE_ADMIN_ACTIVITY_LOG_PAYMENT                ResourceAdminActivityLog = "payment"
	RESOURCE_ADMIN_ACTIVITY_LOG_USER                   ResourceAdminActivityLog = "admin_user"
	RESOURCE_ADMIN_ACTIVITY_LOG_PLATFORM               ResourceAdminActivityLog = "platform"
	RESOURCE_ADMIN_ACTIVITY_LOG_AGGREGATOR             ResourceAdminActivityLog = "aggregator"
	RESOURCE_ADMIN_ACTIVITY_LOG_CHANNEL                ResourceAdminActivityLog = "channel"
	RESOURCE_ADMIN_ACTIVITY_LOG_CONFIGURATION          ResourceAdminActivityLog = "configuration"
	RESOURCE_ADMIN_ACTIVITY_LOG_PLATFORM_CONFIGURATION ResourceAdminActivityLog = "platform_configuration"
	RESOURCE_ADMIN_ACTIVITY_LOG_PLATFORM_CHANNEL       ResourceAdminActivityLog = "platform_channel"
)

type Currency string

const (
	CURRENCY_IDR Currency = "IDR"
	CURRENCY_MYR Currency = "MYR"
)

type SymbolCurrency string

const (
	SYMBOL_CURRENCY_IDR SymbolCurrency = "Rp"
	SYMBOL_CURRENCY_MYR SymbolCurrency = "RM"
)

type KWalletDirection string

const (
	K_WALLET_DIRECTION_IN  KWalletDirection = "in"
	K_WALLET_DIRECTION_OUT KWalletDirection = "out"
)

type TxnStatusSenangpay string

const (
	SENANGPAY_STATUS_SUCCESS TxnStatusSenangpay = "1"
	SENANGPAY_STATUS_FAILED  TxnStatusSenangpay = "0"
)

type TxnStatusMidtrans string

const (
	MIDTRANS_STATUS_CAPTURE    TxnStatusMidtrans = "capture"
	MIDTRANS_STATUS_SETTLEMENT TxnStatusMidtrans = "settlement"
	MIDTRANS_STATUS_PENDING    TxnStatusMidtrans = "pending"
	MIDTRANS_STATUS_DENIED     TxnStatusMidtrans = "deny"
	MIDTRANS_STATUS_EXPIRE     TxnStatusMidtrans = "expire"
	MIDTRANS_STATUS_CANCEL     TxnStatusMidtrans = "cancel"
	MIDTRANS_STATUS_FAILURE    TxnStatusMidtrans = "failure"
	MIDTRANS_STATUS_AUTHORIZE  TxnStatusMidtrans = "authorize"
)

type ConfigName string

const (
	CONFIG_NAME_MIDTRANS  = "midtrans"
	CONFIG_NAME_SENANGPAY = "senangpay"
	CONFIG_NAME_ESPAY     = "espay"
)

type EnumSignatureEspay string

const (
	ENUM_SIGNATURE_ESPAY_INQUIRY           = "INQUIRY"
	ENUM_SIGNATURE_ESPAY_PAYMENTREPORT     = "PAYMENTREPORT"
	ENUM_SIGNATURE_ESPAY_CHECKSTATUS       = "CHECKSTATUS"
	ENUM_SIGNATURE_ESPAY_EXPIRETRANSACTION = "EXPIRETRANSACTION"
	ENUM_SIGNATURE_ESPAY_SENDINVOICE       = "SENDINVOICE"
)

type SandboxProduction string

const (
	SANDBOX    SandboxProduction = "sandbox"
	PRODUCTION SandboxProduction = "production"
)

type VaExpired string

const (
	VA_EXPIRED_1   VaExpired = "1"
	VA_EXPIRED_60  VaExpired = "60"
	VA_EXPIRED_120 VaExpired = "120"
	VA_EXPIRED_180 VaExpired = "180"
	VA_EXPIRED_240 VaExpired = "240"
)

type EspayService string

const (
	INQUIRY              EspayService = "24"
	PAYMENT              EspayService = "25"
	INQUIRY_STATUS       EspayService = "26"
	CREATE_VA            EspayService = "27"
	DELETE_VA            EspayService = "31"
	QRIS_QR_MPM          EspayService = "47"
	PAYMENT_HOST_TO_HOST EspayService = "54"
)

type EspayCaseCode string

const (
	ESPAY_SUCCESSFULL                          EspayCaseCode = "00"
	ESPAY_INVALID_MISSING_FIELD_FORMAT         EspayCaseCode = "01"
	ESPAY_INVALID_MANDATORY_FIELD              EspayCaseCode = "02"
	ESPAY_UNAUTHORIZED_REASON                  EspayCaseCode = "00"
	ESPAY_INVALID_TOKEN_B2B                    EspayCaseCode = "01"
	ESPAY_TOKEN_CUSTOMER_NOT_VALID             EspayCaseCode = "02"
	ESPAY_TOKEN_NOT_FOUND_B2B                  EspayCaseCode = "03"
	ESPAY_INVALID_TRANSACTION_STATUS           EspayCaseCode = "00"
	ESPAY_TRANSACTION_NOT_FOUND                EspayCaseCode = "01"
	ESPAY_INVALID_BILL_VIRITUAL_ACCOUNT_REASON EspayCaseCode = "12"
	ESPAY_INVALID_AMOUNT                       EspayCaseCode = "13"
	ESPAY_BILL_HAS_BEEN_PAID                   EspayCaseCode = "14"
	ESPAY_INVALID_BILL_VIRTUAL_ACCOUNT         EspayCaseCode = "19"
	ESPAY_CONFLICT                             EspayCaseCode = "00"
	ESPAY_DUPLICATE_PARTNER_REFERENCE_NO       EspayCaseCode = "01"
	ESPAY_GENERAL_ERROR                        EspayCaseCode = "00"
	ESPAY_INTERNAL_SERVER_ERROR                EspayCaseCode = "01"
	ESPAY_EXTERNAL_SERVER_ERROR                EspayCaseCode = "02"
	ESPAY_TIMEOUT                              EspayCaseCode = "00"
)

type EspayMessage string

const (
	ESPAY_MESSAGE_SUCCESSFULL                          EspayMessage = "Successful, Transaction successful."
	ESPAY_MESSAGE_INVALID_MISSING_FIELD_FORMAT         EspayMessage = "Invalid / Missing Field Format %v, Invalid data format in field %v"
	ESPAY_MESSAGE_INVALID_MANDATORY_FIELD              EspayMessage = "Invalid Mandatory Field %v, Some mandatory parameters are missing or have an invalid format."
	ESPAY_MESSAGE_UNAUTHORIZED_REASON                  EspayMessage = "Unauthorized. [%v], Authentication failed."
	ESPAY_MESSAGE_INVALID_TOKEN_B2B                    EspayMessage = "Invalid Token (B2B), The token used in the request is invalid or has expired."
	ESPAY_MESSAGE_TOKEN_CUSTOMER_NOT_VALID             EspayMessage = "Invalid Customer Token, Invalid or expired customer token."
	ESPAY_MESSAGE_TOKEN_NOT_FOUND_B2B                  EspayMessage = "Token Not Found (B2B), Token not found in the system."
	ESPAY_MESSAGE_INVALID_TRANSACTION_STATUS           EspayMessage = "Invalid Transaction Status, The current transaction status is not appropriate for this process."
	ESPAY_MESSAGE_TRANSACTION_NOT_FOUND                EspayMessage = "Transaction Not Found."
	ESPAY_MESSAGE_INVALID_BILL_VIRITUAL_ACCOUNT_REASON EspayMessage = "Invalid Bill / Virtual Account [%v], The Virtual Account (VA) number or bill is not found or blocked."
	ESPAY_MESSAGE_INVALID_AMOUNT                       EspayMessage = "Invalid Amount, The payment amount does not match the expected amount."
	ESPAY_MESSAGE_BILL_HAS_BEEN_PAID                   EspayMessage = "Bill has been paid, The bill has already been paid."
	ESPAY_MESSAGE_INVALID_BILL_VIRTUAL_ACCOUNT         EspayMessage = "Invalid Bill / Virtual Account, The bill or Virtual Account used is no longer active."
	ESPAY_MESSAGE_CONFLICT                             EspayMessage = "Conflict, A transaction with the same X-EXTERNAL-ID has already been processed today."
	ESPAY_MESSAGE_DUPLICATE_PARTNER_REFERENCE_NO       EspayMessage = "Duplicate partnerReferenceNo, A transaction with the same partnerReferenceNo has already been processed successfully."
	ESPAY_MESSAGE_GENERAL_ERROR                        EspayMessage = "General Error, A general error occurred while processing the transaction."
	ESPAY_MESSAGE_INTERNAL_SERVER_ERROR                EspayMessage = "Internal Server Error, An internal error occurred while processing the transaction."
	ESPAY_MESSAGE_EXTERNAL_SERVER_ERROR                EspayMessage = "External Server Error, Backend system failure occurred while processing the transaction."
	ESPAY_MESSAGE_TIMEOUT                              EspayMessage = "Timeout, Transaction requests to the bank or issuer are taking longer than usual."
)

func CreateEspayCodeResponse(httpCode int, service EspayService, caseCode EspayCaseCode, err error) (string, string) {
	log.Debug().Int("httpCode", httpCode).Interface("espay service", service).Interface("caseCode", caseCode).Err(err).Msg("create espay code response")
	message := CreateEspayMessageResponse(httpCode, caseCode, err)
	log.Debug().Str("message", message).Msg("return value create espay code response")
	return fmt.Sprintf("%v%v%v", httpCode, service, caseCode), message
}

func CreateEspayMessageResponse(httpCode int, caseCode EspayCaseCode, err error) string {
	log.Debug().Int("httpCode", httpCode).Interface("caseCode", caseCode).Err(err).Msg("create espay message response")
	mapMessage := map[int]map[EspayCaseCode]string{
		http.StatusOK: {
			ESPAY_SUCCESSFULL: string(ESPAY_MESSAGE_SUCCESSFULL),
		},
		http.StatusBadRequest: {
			ESPAY_INVALID_MISSING_FIELD_FORMAT: fmt.Sprintf(string(ESPAY_MESSAGE_INVALID_MISSING_FIELD_FORMAT), err, err),
			ESPAY_INVALID_MANDATORY_FIELD:      fmt.Sprintf(string(ESPAY_MESSAGE_INVALID_MANDATORY_FIELD), err),
		},
		http.StatusUnauthorized: {
			ESPAY_UNAUTHORIZED_REASON:      fmt.Sprintf(string(ESPAY_MESSAGE_UNAUTHORIZED_REASON), err),
			ESPAY_INVALID_TOKEN_B2B:        string(ESPAY_MESSAGE_INVALID_TOKEN_B2B),
			ESPAY_TOKEN_CUSTOMER_NOT_VALID: string(ESPAY_MESSAGE_TOKEN_CUSTOMER_NOT_VALID),
			ESPAY_TOKEN_NOT_FOUND_B2B:      string(ESPAY_MESSAGE_TOKEN_NOT_FOUND_B2B),
		},
		http.StatusNotFound: {
			ESPAY_INVALID_TRANSACTION_STATUS:           string(ESPAY_MESSAGE_INVALID_TRANSACTION_STATUS),
			ESPAY_TRANSACTION_NOT_FOUND:                string(ESPAY_MESSAGE_TRANSACTION_NOT_FOUND),
			ESPAY_INVALID_BILL_VIRITUAL_ACCOUNT_REASON: fmt.Sprintf(string(ESPAY_MESSAGE_INVALID_BILL_VIRITUAL_ACCOUNT_REASON), err),
			ESPAY_INVALID_AMOUNT:                       string(ESPAY_MESSAGE_INVALID_AMOUNT),
			ESPAY_BILL_HAS_BEEN_PAID:                   string(ESPAY_MESSAGE_BILL_HAS_BEEN_PAID),
			ESPAY_INVALID_BILL_VIRTUAL_ACCOUNT:         string(ESPAY_MESSAGE_INVALID_BILL_VIRTUAL_ACCOUNT),
		},
		http.StatusConflict: {
			ESPAY_CONFLICT:                       string(ESPAY_MESSAGE_CONFLICT),
			ESPAY_DUPLICATE_PARTNER_REFERENCE_NO: string(ESPAY_MESSAGE_DUPLICATE_PARTNER_REFERENCE_NO),
		},
		http.StatusInternalServerError: {
			ESPAY_GENERAL_ERROR:         string(ESPAY_MESSAGE_GENERAL_ERROR),
			ESPAY_INTERNAL_SERVER_ERROR: string(ESPAY_MESSAGE_INTERNAL_SERVER_ERROR),
			ESPAY_EXTERNAL_SERVER_ERROR: string(ESPAY_MESSAGE_EXTERNAL_SERVER_ERROR),
		},
		http.StatusGatewayTimeout: {
			ESPAY_TIMEOUT: string(ESPAY_MESSAGE_TIMEOUT),
		},
	}

	if cases, ok := mapMessage[httpCode]; ok {
		if message, ok := cases[caseCode]; ok {
			return message
		}
	}

	return ""
}

type StringEnum string

const (
	NULL_STRING StringEnum = ""
)

type KWalletStatus string

const (
	KWalletStatusActive   KWalletStatus = "active"
	KWalletStatusInactive KWalletStatus = "inactive"
)

type AssetType string

const (
	ASSET_LOGO AssetType = "logo"
)

type ImageSizeType string

const (
	IMAGE_SIZE_TYPE_S  ImageSizeType = "s"
	IMAGE_SIZE_TYPE_M  ImageSizeType = "m"
	IMAGE_SIZE_TYPE_L  ImageSizeType = "l"
	IMAGE_SIZE_TYPE_XL ImageSizeType = "xl"
)

type ImageGeometric string

const (
	IMAGE_ROUNDED   ImageGeometric = "rounded"
	IMAGE_SQUARE    ImageGeometric = "square"
	IMAGE_CIRCLE    ImageGeometric = "circle"
	IMAGE_RECTANGLE ImageGeometric = "rectangle"
)

type StatisValue string

const (
	STATIS_PHONE_NUMBER StatisValue = "082123456789"
	STATIS_EMAIL        StatisValue = "jhon.doe@gmail.com"
)
