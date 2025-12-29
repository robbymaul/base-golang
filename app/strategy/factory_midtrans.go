package strategy

import (
	"context"
	"fmt"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/web"
)

type MidtransPayment interface {
	Pay(ctx context.Context, req any) (map[string]interface{}, error)
	ClientResponse(channel *models.Channel, payment *models.Payments) (*web.PaymentResponse, error)
}

type MidtransStrategy struct {
	BCAVirtualAccount     MidtransPayment
	BNIVirtualAccount     MidtransPayment
	BRIVirtualAccount     MidtransPayment
	CIMBVa                MidtransPayment
	Gopay                 MidtransPayment
	MandiriBillPayment    MidtransPayment
	PermataVirtualAccount MidtransPayment
	Qris                  MidtransPayment
	ShopeePay             MidtransPayment
	Dana                  MidtransPayment
}

func NewMidtransStrategy(
	BCAVirtualAccount MidtransPayment,
	BNIVirtualAccount MidtransPayment,
	BRIVirtualAccount MidtransPayment,
	CIMBVa MidtransPayment,
	Gopay MidtransPayment,
	MandiriBillPayment MidtransPayment,
	PermataViratualAccount MidtransPayment,
	Qris MidtransPayment,
	ShopeePay MidtransPayment,
	Dana MidtransPayment,
) *MidtransStrategy {
	return &MidtransStrategy{
		BCAVirtualAccount:     BCAVirtualAccount,
		BNIVirtualAccount:     BNIVirtualAccount,
		BRIVirtualAccount:     BRIVirtualAccount,
		CIMBVa:                CIMBVa,
		Gopay:                 Gopay,
		MandiriBillPayment:    MandiriBillPayment,
		PermataVirtualAccount: PermataViratualAccount,
		Qris:                  Qris,
		ShopeePay:             ShopeePay,
		Dana:                  Dana,
	}
}

//func (s *MidtransStrategy) GetStrategy(method enums.PaymentMethod) (MidtransPayment, error) {
//	switch method {
//	case enums.PAYMENT_METHOD_BCA_VA:
//		return s.BCAVirtualAccount, nil
//	case enums.PAYMENT_METHOD_BNI_VA:
//		return s.BNIVirtualAccount, nil
//	case enums.PAYMENT_METHOD_BRI_VA:
//		return s.BRIVirtualAccount, nil
//	case enums.PAYMENT_METHOD_CIMB_VA:
//		return s.CIMBVa, nil
//	case enums.PAYMENT_METHOD_GOPAY:
//		return s.Gopay, nil
//	case enums.PAYMENT_METHOD_MANDIRI_VA:
//		return s.MandiriBillPayment, nil
//	case enums.PAYMENT_METHOD_PERMATA_VA:
//		return s.PermataVirtualAccount, nil
//	case "qris":
//		return s.Qris, nil
//	case enums.PAYMENT_METHOD_SHOPEE_PAY:
//		return s.ShopeePay, nil
//	default:
//		return nil, errors.New("unknown payment method")
//	}
//}

func (s *MidtransStrategy) GetStrategy(method enums.PaymentMethod, channel enums.Channel) (MidtransPayment, error) {
	//switch method {
	//case enums.PAYMENT_METHOD_BCA_VA:
	//	return s.BCAVirtualAccount, nil
	//case enums.PAYMENT_METHOD_BNI_VA:
	//	return s.BNIVirtualAccount, nil
	//case enums.PAYMENT_METHOD_BRI_VA:
	//	return s.BRIVirtualAccount, nil
	//case enums.PAYMENT_METHOD_CIMB_VA:
	//	return s.CIMBVa, nil
	//case enums.PAYMENT_METHOD_GOPAY:
	//	return s.Gopay, nil
	//case enums.PAYMENT_METHOD_MANDIRI_VA:
	//	return s.MandiriBillPayment, nil
	//case enums.PAYMENT_METHOD_PERMATA_VA:
	//	return s.PermataVirtualAccount, nil
	//case "qris":
	//	return s.Qris, nil
	//case enums.PAYMENT_METHOD_SHOPEE_PAY:
	//	return s.ShopeePay, nil
	//default:
	//	return nil, errors.New("unknown payment method")
	//}

	if method == enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT && channel == enums.CHANNEL_BCA {
		return s.BCAVirtualAccount, nil
	}

	if method == enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT && channel == enums.CHANNEL_BNI {
		return s.BNIVirtualAccount, nil
	}

	if method == enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT && channel == enums.CHANNEL_BRI {
		return s.BRIVirtualAccount, nil
	}

	if method == enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT && channel == enums.CHANNEL_CIMB {
		return s.CIMBVa, nil
	}

	if method == enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT && channel == enums.CHANNEL_MANDIRI {
		return s.MandiriBillPayment, nil
	}

	if method == enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT && channel == enums.CHANNEL_PERMATA {
		return s.PermataVirtualAccount, nil
	}

	if channel == enums.CHANNEL_GOPAY {
		return s.Gopay, nil
	}

	if channel == enums.CHANNEL_SHOPEE {
		return s.ShopeePay, nil
	}

	if method == enums.PAYMENT_METHOD_QRIS {
		return s.Qris, nil
	}

	if channel == enums.CHANNEL_DANA {
		return s.Dana, nil
	}

	return nil, fmt.Errorf("payment method or channel not recognized. Please choose another payment option or contact our support team")
}
