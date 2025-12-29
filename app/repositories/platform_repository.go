package repositories

import (
	"context"
	"fmt"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/web"
	"paymentserviceklink/pkg/pagination"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (rc *RepositoryContext) InsertPlatformRepository(ctx context.Context, platforms *models.Platforms) (*models.Platforms, error) {
	var (
		db = rc.db
	)

	db = db.Table("platforms").Model(&models.Platforms{}).WithContext(ctx)

	err := db.Create(platforms).Error
	if err != nil {
		return nil, err
	}

	return platforms, nil
}

func (rc *RepositoryContext) InsertPlatformRepositoryTx(ctx context.Context, tx *gorm.DB, platforms *models.Platforms) (*models.Platforms, error) {
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("platforms").Model(&models.Platforms{}).WithContext(ctx)

	err := db.Create(platforms).Error
	if err != nil {
		return nil, err
	}

	return platforms, nil
}

func (rc *RepositoryContext) FindPlatformRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) ([]*models.Platforms, error) {
	var (
		db        = rc.db
		platforms []*models.Platforms
		err       error
	)

	db = db.Table("platforms").Model(&models.Platforms{}).WithContext(ctx)
	db = fn(db)

	err = db.Find(&platforms).Error
	if err != nil {
		log.Error().Err(err).Msg("error query find platform")
		return nil, err
	}

	return platforms, nil
}

func (rc *RepositoryContext) GetTotalCountPlatformRepository(ctx context.Context, pages *pagination.Pages, fn func(db *gorm.DB) *gorm.DB) (int64, error) {
	var (
		db         = rc.db
		totalCount int64
		err        error
	)

	db = db.Table("platforms").Model(&models.Platforms{}).WithContext(ctx)
	db = fn(db)

	err = db.Count(&totalCount).Error
	if err != nil {
		log.Error().Err(err).Msg("error query total count platform")
		return 0, err
	}

	return totalCount, nil
}

func (rc *RepositoryContext) GetPlatformRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) (*models.Platforms, error) {
	var (
		db       = rc.db
		platform *models.Platforms
		err      error
	)

	db = db.Table("platforms").Model(&models.Platforms{}).WithContext(ctx)
	db = fn(db)

	err = db.First(&platform).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get first platform")
		return nil, err
	}

	return platform, nil
}

func (rc *RepositoryContext) UpdatePlatformRepository(ctx context.Context, platform *models.Platforms, fn func(db *gorm.DB) *gorm.DB) error {
	var (
		db = rc.db
	)

	db = db.Table("platforms").Model(&models.Platforms{}).WithContext(ctx)
	db = fn(db)

	return db.Error
}

func (rc *RepositoryContext) GetPlatformByApiKeyAndSecretKeyRepository(ctx context.Context, payload *web.ClientRequest) (platform *models.Platforms, err error) {
	var (
		db = rc.db
	)

	db = db.Table("platforms").Model(&models.Platforms{}).WithContext(ctx)

	db = db.Where("api_key = ? and secret_key = ?", payload.ApiKey, payload.SecretKey)

	err = db.First(&platform).Error

	return
}

func (rc *RepositoryContext) FindPlatformInConfigurationRepository(ctx context.Context, configurationId int64, pages *pagination.Pages, fn func(db *gorm.DB) *gorm.DB) ([]*models.Platforms, error) {
	var (
		db        = rc.db
		platforms []*models.Platforms
		err       error
	)

	db = db.Table("platforms").Model(&models.Platforms{}).WithContext(ctx)
	db = fn(db)

	err = db.Find(&platforms).Error
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("error query find platforms join platform configuration where platform_configuration.configuration_id = %v", configurationId))
		return nil, err
	}

	return platforms, nil
}
