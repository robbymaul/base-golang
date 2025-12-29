package repositories

import (
	"context"
	"errors"
	"fmt"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (rc *RepositoryContext) GetConfigurationRepository(ctx context.Context, f func(db *gorm.DB) *gorm.DB) (*models.Configuration, error) {
	var (
		db     = rc.db
		config *models.Configuration
		err    error
	)

	db = db.Table("configurations").Model(&models.Configuration{}).WithContext(ctx)
	db = f(db)

	err = db.First(&config).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get platform configuration")
		return nil, err
	}

	return config, err
}

func (rc *RepositoryContext) InsertConfigurationRepository(ctx context.Context, platformConfiguration *models.Configuration) (*models.Configuration, error) {
	var (
		db = rc.db
	)

	db = db.Table("configurations").Model(&models.Configuration{}).WithContext(ctx)

	err := db.Create(&platformConfiguration).Error
	if err != nil {
		log.Error().Err(err).Msg("error query insert platform configuration")
		return nil, err
	}

	return platformConfiguration, err
}

func (rc *RepositoryContext) InsertBatchConfigurationRepository(ctx context.Context, configurations []*models.Configuration) ([]*models.Configuration, error) {
	log.Debug().Interface("configurations", configurations).Msg("insert batch configuration repository")

	var (
		db = rc.db
	)

	db = db.Table("configurations").Model(&models.Configuration{}).WithContext(ctx)

	err := db.Create(&configurations).Error
	if err != nil {
		log.Error().Err(err).Msg("error query insert bulking configuration")
		return nil, err
	}

	return configurations, err
}

func (rc *RepositoryContext) InsertBatchConfigurationRepositoryTx(ctx context.Context, tx *gorm.DB, configurations []*models.Configuration) ([]*models.Configuration, error) {
	log.Debug().Interface("configurations", configurations).Msg("insert batch configuration repository")

	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("configurations").Model(&models.Configuration{}).WithContext(ctx)

	err := db.Create(&configurations).Error
	if err != nil {
		log.Error().Err(err).Msg("error query insert bulking configuration")
		return nil, err
	}

	return configurations, err
}

func (rc *RepositoryContext) UpdateConfigurationRepository(ctx context.Context, platformConfiguration *models.Configuration) (*models.Configuration, error) {
	var (
		db = rc.db
	)

	db = db.Table("configurations").Model(&models.Configuration{}).WithContext(ctx)

	err := db.Save(&platformConfiguration).Error
	if err != nil {
		log.Error().Err(err).Msg("error query update platform configuration")
		return nil, err
	}

	return platformConfiguration, err
}

func (rc *RepositoryContext) UpdateConfigurationRepositoryTx(ctx context.Context, tx *gorm.DB, fn func(db *gorm.DB) *gorm.DB) error {
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("configurations").Model(&models.Configuration{}).WithContext(ctx)

	db = fn(db)

	return db.Error
}

func (rc *RepositoryContext) GetConfigurationByIdRepository(ctx context.Context, id int64) (*models.Configuration, error) {
	var (
		db     = rc.db
		config *models.Configuration
		err    error
	)

	db = db.Table("configurations").Model(&models.Configuration{}).WithContext(ctx)

	db = db.Where("id = ?", id)

	db = db.Preload("Aggregator")

	err = db.First(&config).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get platform configuration by id")
		return nil, err
	}

	return config, nil
}

func (rc *RepositoryContext) GetExistsConfigurationByPlatformIdAndAggregatorIdRepository(ctx context.Context, idPlatform int64, idAggregator int64) (bool, error) {
	var (
		db     = rc.db
		config *models.Configuration
		err    error
	)

	db = db.Table("configurations").Model(&models.Configuration{}).WithContext(ctx)

	db = db.Where("platform_id = ? and aggregator_id = ?", idPlatform, idAggregator)

	err = db.First(&config).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		log.Error().Err(err).Msg("error query get platform configuration by platform id and aggregator id")
		return false, err
	}

	return true, nil
}

func (rc *RepositoryContext) GetCountConfigurationRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) (int64, error) {
	var (
		db    = rc.db
		count int64
		err   error
	)

	db = db.Table("configurations").Model(&models.Configuration{}).WithContext(ctx)
	db = fn(db)

	err = db.Count(&count).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get count platform configuration")
		return 0, err
	}

	return count, err
}

func (rc *RepositoryContext) FindConfigurationRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) ([]*models.Configuration, error) {
	var (
		db      = rc.db
		configs []*models.Configuration
		err     error
	)

	db = db.Table("configurations").Model(&models.Configuration{}).WithContext(ctx)
	db = fn(db)

	err = db.Find(&configs).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get list platform configuration")
		return nil, err
	}

	return configs, nil
}

func (rc *RepositoryContext) GetListConfigurationInPlatformRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) ([]*models.Configuration, error) {
	var (
		db             = rc.db
		configurations []*models.Configuration
		err            error
	)

	db = db.Table("configurations").Model(&models.Configuration{}).WithContext(ctx)
	db = fn(db)

	err = db.Find(&configurations).Error
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("error query find configurations join platform configuration"))
		return nil, err
	}

	return configurations, nil
}

func (rc *RepositoryContext) GetConfigurationByPlatformIdRepository(ctx context.Context, platformId int64) ([]*models.Configuration, error) {
	var (
		db             = rc.db
		configurations []*models.Configuration
		err            error
	)

	db = db.Table("configurations").Model(&models.Configuration{}).WithContext(ctx)

	db = db.Joins("JOIN platform_configuration on platform_configuration.configuration_id = configurations.id")
	db = db.Where("platform_configuration.platform_id = ?", platformId)

	err = db.Find(&configurations).Error
	if err != nil {
		log.Error().Err(err).Msg("error query find configurations by platform_id")
		return nil, err
	}

	return configurations, nil
}

func (rc *RepositoryContext) GetConfigurationByPlatformIdAndCurrencyRepository(ctx context.Context, platformId int64, currency enums.Currency) ([]*models.Configuration, error) {
	var (
		db            = rc.db
		configuration []*models.Configuration
		err           error
	)

	db = db.Table("configurations").Model(&models.Configuration{}).WithContext(ctx)
	db = db.Joins("JOIN platform_configuration on platform_configuration.configuration_id = configurations.id")
	db = db.Joins("JOIN aggregators on aggregators.id = configurations.aggregator_id")
	db = db.Where("aggregators.currency = ? and platform_configuration.platform_id = ? ", currency, platformId)

	db = db.Preload("Aggregator")

	err = db.Find(&configuration).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get configuration by platform id and aggregator id")
		return nil, err
	}

	return configuration, nil
}

func (rc *RepositoryContext) GetConfigurationByPlatformIdAndAggregatorIdRepository(ctx context.Context, platformId int64, aggregatorId int64) (*models.Configuration, error) {
	var (
		db            = rc.db
		configuration *models.Configuration
		err           error
	)

	db = db.Table("configurations").Model(&models.Configuration{}).WithContext(ctx)
	db = db.Joins("JOIN platform_configuration on platform_configuration.configuration_id = configurations.id")
	db = db.Where("configurations.aggregator_id = ? and platform_configuration.platform_id = ? ", aggregatorId, platformId)

	db = db.Preload("Aggregator")

	err = db.First(&configuration).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get configuration by platform id and aggregator id")
		return nil, err
	}

	return configuration, nil
}
