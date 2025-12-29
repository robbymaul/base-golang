package repositories

import (
	"context"
	"gorm.io/gorm"
	"paymentserviceklink/app/models"
)

func (rc *RepositoryContext) InsertAdminActivityLogRepositoryTx(ctx context.Context, tx *gorm.DB, adminActivityLog *models.AdminActivityLogs) error {
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("admin_activity_logs").Model(&models.AdminActivityLogs{}).WithContext(ctx)

	return db.Create(&adminActivityLog).Error
}
