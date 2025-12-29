package repositories

import (
	"context"
	"gorm.io/gorm"
	"paymentserviceklink/app/models"
)

func (rc *RepositoryContext) InsertPaymentStatusHistoryRepositoryTx(ctx context.Context, tx *gorm.DB, history *models.PaymentStatusHistory) error {
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("payment_status_history").Model(&models.PaymentStatusHistory{}).WithContext(ctx)

	return db.Create(&history).Error
}
