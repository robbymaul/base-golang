package repositories

import (
	"context"
	"paymentserviceklink/app/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (rc *RepositoryContext) FindTopupTransactionRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) ([]*models.TopupTransaction, error) {
	var (
		db     = rc.db
		topups []*models.TopupTransaction
		err    error
	)

	db = db.Table("topup_transaction").Model(&models.TopupTransaction{}).WithContext(ctx)
	db = fn(db)

	err = db.Find(&topups).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get all topup transaction")
		return nil, err
	}

	return topups, nil
}

func (rc *RepositoryContext) GetTotalCountTopupTransactionRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) (int64, error) {
	var (
		db         = rc.db
		totalCount int64
		err        error
	)

	db = db.Table("topup_transaction").Model(&models.TopupTransaction{}).WithContext(ctx)
	db = fn(db)

	err = db.Count(&totalCount).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get total count topup transaction")
		return 0, err
	}

	return totalCount, nil
}

func (rc *RepositoryContext) InsertTopupTransactionRepositoryTx(ctx context.Context, tx *gorm.DB, topupTransaction models.TopupTransaction) error {
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("topup_transaction").Model(&models.TopupTransaction{}).WithContext(ctx)

	err := db.Create(&topupTransaction).Error
	if err != nil {
		log.Error().Err(err).Msg("error query insert topup transaction")
		return err
	}

	return nil
}
