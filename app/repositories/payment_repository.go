package repositories

import (
	"context"
	"fmt"
	"net/http"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/web"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (rc *RepositoryContext) InsertPaymentRepository(ctx context.Context, payment *models.Payments) (*models.Payments, error) {
	var (
		db = rc.db
	)

	db = db.Table("payments").Model(&models.Payments{}).WithContext(ctx)

	err := db.Create(payment).Error
	if err != nil {
		return nil, err
	}

	err = db.Where("id = ?", payment.Id).Preload("Platform").Preload("Channel").First(&payment).Error
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (rc *RepositoryContext) UpdatePaymentAfterTryRepository(ctx context.Context, payment *models.Payments) (*models.Payments, error) {
	var (
		db = rc.db
	)

	db = db.Table("payments").Model(&models.Payments{}).WithContext(ctx)

	db = db.Where("transaction_id", payment.TransactionId)

	updateColumn := map[string]interface{}{
		"aggregator_id": payment.AggregatorId,
		"updated_at":    time.Now(),
	}

	err := db.UpdateColumns(updateColumn).Error
	if err != nil {
		log.Error().Err(err).Msg("error query update payment after try repository")
		return nil, err
	}

	return payment, nil
}

func (rc *RepositoryContext) UpdatePaymentAfterTryRepositoryTx(ctx context.Context, tx *gorm.DB, payment *models.Payments) (*models.Payments, error) {
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("payments").Model(&models.Payments{}).WithContext(ctx)

	db = db.Where("transaction_id", payment.TransactionId)

	updateColumn := map[string]interface{}{
		"aggregator_id": payment.AggregatorId,
		"updated_at":    time.Now(),
	}

	err := db.UpdateColumns(updateColumn).Error
	if err != nil {
		log.Error().Err(err).Msg("error query update payment after try repository")
		return nil, err
	}

	return payment, nil
}

func (rc *RepositoryContext) GetPaymentByTransactionIdRepository(ctx context.Context, transactionId string) (payment *models.Payments, err error) {
	var (
		db = rc.db
	)

	db = db.Table("payments").Model(&models.Payments{}).WithContext(ctx)

	err = db.First(&payment).Error

	return
}

func (rc *RepositoryContext) GetPaymentByOrderIdRepository(ctx context.Context, orderId string) (payment *models.Payments, err error) {
	var (
		db = rc.db
	)

	db = db.Table("payments").Model(&models.Payments{}).WithContext(ctx)

	db = db.Where("order_id = ? and status = ?", orderId, enums.PAYMENT_STATUS_PENDING)

	db = db.Preload("Platform").Preload("Channel")

	err = db.First(&payment).Error

	return
}

func (rc *RepositoryContext) GetPaymentByTransactionIdAndPlatformIdRepository(ctx context.Context, orderId string, platformId int64) (payment *models.Payments, err error) {
	var (
		db = rc.db
	)

	db = db.Table("payments").Model(&models.Payments{}).WithContext(ctx)

	db = db.Where("order_id = ? and platform_id = ?", orderId, platformId)

	db = db.Preload("Platform").Preload("Channel").Preload("Aggregator")

	err = db.First(&payment).Error

	return
}

func (rc *RepositoryContext) UpdatePaymentRepositoryTx(ctx context.Context, tx *gorm.DB, fn func(db *gorm.DB) *gorm.DB) error {
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("payments").Model(&models.Payments{}).WithContext(ctx)
	db = fn(db)

	err := db.Error
	if err != nil {
		log.Error().Err(err).Msg("error update payment repository transaction")
	}

	return nil
}

func (rc *RepositoryContext) UpdatePaymentRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) error {
	var (
		db = rc.db
	)

	db = db.Table("payments").Model(&models.Payments{}).WithContext(ctx)
	db = fn(db)

	err := db.Error
	if err != nil {
		log.Error().Err(err).Msg("error update payment repository transaction")
	}

	return nil
}

func (rc *RepositoryContext) UpdatePaymentAfterGetResponsePaymentGatewayRepository(ctx context.Context, payment *models.Payments) error {
	var (
		db = rc.db
	)

	db = db.Table("payments").Model(&models.Payments{}).WithContext(ctx)

	db = db.Where("id = ?", payment.Id).Omit("id")

	//updateColumn := map[string]interface{}{
	//	"amount":           payment.Amount,
	//	"fee_amount":       payment.FeeAmount,
	//	"total_amount":     payment.TotalAmount,
	//	"gateway_response": payment.GatewayResponse,
	//	"updated_at":       time.Now(),
	//}
	err := db.Updates(payment).Error
	if err != nil {
		log.Error().Err(err).Msg("error update payment after get response gateway repository")
		return err
	}

	//return db.UpdateColumns(updateColumn).Error
	return nil
}

func (rc *RepositoryContext) UpdatePaymentAfterGetResponsePaymentGatewayRepositoryTx(ctx context.Context, tx *gorm.DB, payment *models.Payments) error {
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("payments").Model(&models.Payments{}).WithContext(ctx)

	db = db.Where("id = ?", payment.Id).Omit("id")

	//updateColumn := map[string]interface{}{
	//	"amount":           payment.Amount,
	//	"fee_amount":       payment.FeeAmount,
	//	"total_amount":     payment.TotalAmount,
	//	"gateway_response": payment.GatewayResponse,
	//	"updated_at":       time.Now(),
	//}
	err := db.Updates(payment).Error
	if err != nil {
		log.Error().Err(err).Msg("error update payment after get response gateway repository")
		return err
	}

	//return db.UpdateColumns(updateColumn).Error
	return nil
}

func (rc *RepositoryContext) CountGetListPaymentRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) (int64, error) {
	var (
		db    = rc.db
		count int64
		err   error
	)
	db = db.Table("payments").Model(&models.Payments{}).WithContext(ctx)
	db = fn(db)

	err = db.Count(&count).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get count payments list")
		return 0, err
	}

	return count, nil
}

func (rc *RepositoryContext) FindPaymentRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) ([]*models.Payments, error) {
	var (
		db       = rc.db
		payments []*models.Payments
		err      error
	)

	db = db.Table("payments").Model(&models.Payments{}).WithContext(ctx)
	db = fn(db)

	err = db.Find(&payments).Error
	if err != nil {
		log.Error().Err(err).Msg("error query find payment error")
		return nil, err
	}

	return payments, nil
}

func (rc *RepositoryContext) GetPaymentByIdRepository(ctx context.Context, id int64) (*models.Payments, error) {
	var (
		db      = rc.db
		payment *models.Payments
		err     error
	)

	db = db.Table("payments").Model(&models.Payments{}).WithContext(ctx)

	db = db.Where("id = ? ", id)

	db = db.Preload("Platform").Preload("Channel")

	err = db.First(&payment).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get payment by id")
		return nil, err
	}

	return payment, nil
}

func (rc *RepositoryContext) InsertPaymentRepositoryTx(ctx context.Context, tx *gorm.DB, payment *models.Payments) (*models.Payments, error) {
	var (
		db  = rc.db
		err error
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("payments").Model(&models.Payments{}).WithContext(ctx)

	err = db.Create(payment).Error
	if err != nil {
		log.Error().Err(err).Msg("error query insert payment")
		return nil, err
	}

	return payment, nil
}

func (rc *RepositoryContext) CallbackFunctionRepository(ctx context.Context, url string, callback *web.PaymentCallback) error {
	client := rc.HttpClient.Client.R()

	client.SetContext(ctx)
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("Accept", "application/json")
	client.SetBody(callback)
	response, err := client.Post(url)
	if err != nil {
		log.Error().Err(err).Msg("error query callback function")
		return err
	}

	if response.StatusCode() != http.StatusOK {
		log.Error().Msg("error query callback function")
		return helpers.NewErrorTrace(fmt.Errorf("error query callback function"), "").WithStatusCode(http.StatusInternalServerError)
	}

	return nil
}

func (rc *RepositoryContext) UpdatePaymentNotificationRepository(ctx context.Context, payment *models.Payments) error {
	var (
		db = rc.db
	)

	db = db.Table("payments").Model(&models.Payments{}).WithContext(ctx)

	db = db.Where("transaction_id", payment.TransactionId)

	updateColumn := map[string]interface{}{
		"notification_callback": true,
		"updated_at":            time.Now(),
	}

	err := db.UpdateColumns(updateColumn).Error
	if err != nil {
		log.Error().Err(err).Msg("error query update payment after try repository")
		return err
	}

	return nil
}

func (rc *RepositoryContext) GetPaymentRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) (*models.Payments, error) {
	var (
		db      = rc.db
		payment *models.Payments
		err     error
	)

	db = db.Table("payments").Model(&models.Payments{}).WithContext(ctx)
	db = fn(db)

	err = db.First(&payment).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get payments")
		return nil, err
	}

	return payment, err
}

//
//func (rc *RepositoryContext) FindPaymentRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) ([]*models.Payments, error) {
//	var (
//		db       = rc.db
//		payments []*models.Payments
//		err      error
//	)
//
//	db = db.Table("payments").Model(&models.Payments{}).WithContext(ctx)
//	db = fn(db)
//
//	err = db.Find(&payments).Error
//	if err != nil {
//		log.Error().Err(err).Msg("error query find payments")
//		return nil, err
//	}
//
//	return payments, nil
//}
