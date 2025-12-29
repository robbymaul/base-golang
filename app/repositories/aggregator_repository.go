package repositories

import (
	"context"
	"errors"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/models"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (rc *RepositoryContext) InsertAggregatorRepository(ctx context.Context, aggregator *models.Aggregator) (*models.Aggregator, error) {
	var (
		db = rc.db
	)

	db = db.Table("aggregators").Model(&models.Aggregator{}).WithContext(ctx)

	if err := db.Create(aggregator).Error; err != nil {
		log.Error().Err(err).Msg("error query insert aggregator")
		return nil, err
	}

	return aggregator, nil
}

func (rc *RepositoryContext) InsertAggregatorRepositoryTx(ctx context.Context, tx *gorm.DB, aggregator *models.Aggregator) (*models.Aggregator, error) {
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("aggregators").Model(&models.Aggregator{}).WithContext(ctx)

	if err := db.Create(aggregator).Error; err != nil {
		log.Error().Err(err).Msg("error query insert aggregator")
		return nil, err
	}

	return aggregator, nil
}

func (rc *RepositoryContext) GetAggregatorRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) (*models.Aggregator, error) {
	var (
		db         = rc.db
		aggregator *models.Aggregator
		err        error
	)

	db = db.Table("aggregators").Model(&models.Aggregator{}).WithContext(ctx)
	db = fn(db)

	err = db.First(&aggregator).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get aggregator")
		return nil, err
	}

	return aggregator, err
}

func (rc *RepositoryContext) UpdateAggregatorRepository(ctx context.Context, aggregator *models.Aggregator) (*models.Aggregator, error) {
	var (
		db = rc.db
	)

	db = db.Table("aggregators").Model(&models.Aggregator{}).WithContext(ctx)

	db = db.Where("id = ?", aggregator.Id)

	updateColumn := map[string]interface{}{
		"name":        aggregator.Name,
		"description": aggregator.Description,
		"is_active":   aggregator.IsActive,
		"updated_at":  time.Now(),
	}

	if err := db.UpdateColumns(updateColumn).Error; err != nil {
		log.Error().Err(err).Msg("error query update aggregator")
		return nil, err
	}

	return aggregator, nil
}

func (rc *RepositoryContext) UpdateAggregatorRepositoryTx(ctx context.Context, tx *gorm.DB, fn func(db *gorm.DB) *gorm.DB) error {
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("aggregators").Model(&models.Aggregator{}).WithContext(ctx)
	db = fn(db)

	return db.Error
}

func (rc *RepositoryContext) GetAggregatorByNameRepository(ctx context.Context, name string) (*models.Aggregator, error) {
	var (
		db         = rc.db
		aggregator *models.Aggregator
		err        error
	)

	db = db.Table("aggregators").Model(&models.Aggregator{}).WithContext(ctx)

	db = db.Where("name = ?", name)

	err = db.First(&aggregator).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get aggregator by name")
		return nil, err
	}

	return aggregator, err
}

func (rc *RepositoryContext) FindAggregatorRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) ([]*models.Aggregator, error) {
	var (
		db          = rc.db
		aggregators []*models.Aggregator
		err         error
	)

	db = db.Table("aggregators").Model(&models.Aggregator{}).WithContext(ctx)
	db = fn(db)

	err = db.Find(&aggregators).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get all aggregator")
		return nil, err
	}

	return aggregators, nil
}

func (rc *RepositoryContext) GetCountAggregatorRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) (int64, error) {
	var (
		db    = rc.db
		count int64
		err   error
	)

	db = db.Table("aggregators").Model(&models.Aggregator{}).WithContext(ctx)
	db = fn(db)

	err = db.Count(&count).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get count aggregator")
		return 0, err
	}

	return count, nil
}

func (rc *RepositoryContext) GetExistAggregatorByNameRepository(ctx context.Context, name enums.AggregatorName) (bool, error) {
	var (
		db         = rc.db
		aggregator *models.Aggregator
		err        error
	)

	db = db.Table("aggregators").Model(&models.Aggregator{}).WithContext(ctx)

	db = db.Where("name = ?", name)

	err = db.First(&aggregator).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		log.Error().Err(err).Msg("error query get aggregator by name")
		return false, err
	}

	return true, nil
}

func (rc *RepositoryContext) GetAggregatorByIdRepository(ctx context.Context, id int64) (*models.Aggregator, error) {
	var (
		db         = rc.db
		aggregator *models.Aggregator
		err        error
	)

	db = db.Table("aggregators").Model(&models.Aggregator{}).WithContext(ctx)

	db = db.Where("id = ?", id)

	err = db.First(&aggregator).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get aggregator by id")
		return nil, err
	}

	return aggregator, err
}
