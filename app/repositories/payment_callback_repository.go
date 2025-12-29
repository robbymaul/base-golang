package repositories

import (
	"context"
	"gorm.io/gorm"
	"paymentserviceklink/app/models"
)

func (rc *RepositoryContext) InsertPaymentCallbackRepositoryTx(ctx context.Context, tx *gorm.DB, callback *models.PaymentCallbacks) error {
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("payment_callbacks").Model(&models.PaymentCallbacks{}).WithContext(ctx)

	return db.Create(&callback).Error
}
