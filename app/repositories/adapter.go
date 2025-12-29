package repositories

import (
	"fmt"
	"paymentserviceklink/app/client/espay"
	"paymentserviceklink/app/client/midtrans"
	restyclient "paymentserviceklink/app/client/resty"
	clientsenangpay "paymentserviceklink/app/client/senangpay"
	"paymentserviceklink/app/strategy"
	"paymentserviceklink/config"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Adapter struct {
	JakartaLoc *time.Location
	Senangpay  *clientsenangpay.Senangpay
	Midtrans   *midtrans.Midtrans
	Espay      *espay.Espay
	HttpClient *restyclient.RestyClient
	Strategy   *strategy.Strategy
	Minio      *minio.Client
}

func NewRepositoryAdapter(cfg *config.Config) (*Adapter, error) {
	adapter := new(Adapter)

	location, err := time.LoadLocation(cfg.ServerTimeZone)
	if err != nil {
		return nil, err
	}

	adapter.JakartaLoc = location

	adapter.HttpClient = restyclient.NewRestyClient()

	//adapter.Senangpay = senangpay.NewSenangpay(cfg.SenangpayUrl, cfg.SenangpaySecretKey, cfg.SenangpayMerchantId, adapter.HttpClient)
	//adapter.Midtrans = midtrans.NewMidtrans(adapter.HttpClient, cfg)
	//adapter.Strategy = strategy.NewStrategy(adapter.Senangpay, adapter.Midtrans, nil)

	minioClient, err := minio.New(
		cfg.S3Endpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(cfg.S3AccessKeyId, cfg.S3SecretAccessKey, ""),
			Secure: cfg.S3UseSSL,
		},
	)
	if err != nil {
		return nil, newError(fmt.Sprintf("error set up minio client:"), err.Error())
	}

	adapter.Minio = minioClient

	return adapter, nil
}
