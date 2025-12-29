package repositories

import (
	"context"
	"fmt"
	"paymentserviceklink/app/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (rc *RepositoryContext) InsertBatchPlatformConfigurationRepositoryTx(ctx context.Context, tx *gorm.DB, platformConfiguration []*models.PlatformConfiguration) ([]*models.PlatformConfiguration, error) {
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("platform_configuration").Model(&models.PlatformConfiguration{}).WithContext(ctx)
	err := db.Create(&platformConfiguration).Error
	if err != nil {
		log.Error().Err(err).Msg("error query insert batch platform configuration")
		return nil, err
	}

	return platformConfiguration, err
}

func (rc *RepositoryContext) GetPlatformConfigurationByConfigurationIdAndPlatformIdRepository(ctx context.Context, configurationId int64, platformId int64) (*models.PlatformConfiguration, error) {
	var (
		db                    = rc.db
		platformConfiguration *models.PlatformConfiguration
		err                   error
	)

	db = db.Table("platform_configuration").Model(&models.PlatformConfiguration{}).WithContext(ctx)

	db = db.Where("configuration_id = ? and platform_id = ?", configurationId, platformId)

	err = db.First(&platformConfiguration).Error
	if err != nil {
		log.Error().Err(err).
			Msg(fmt.Sprintf("error query get platform configuration where configuration_id = %v and platform_id = %v", configurationId, platformId))
		return nil, err
	}

	return platformConfiguration, nil
}

func (rc *RepositoryContext) GetPlatformConfigurationRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) (*models.PlatformConfiguration, error) {
	var (
		db                    = rc.db
		platformConfiguration *models.PlatformConfiguration
		err                   error
	)

	db = db.Table("platform_configuration").Model(&models.PlatformConfiguration{}).WithContext(ctx)
	db = fn(db)

	err = db.First(&platformConfiguration).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get first platform configuration")
		return nil, err
	}

	return platformConfiguration, nil
}

func (rc *RepositoryContext) DeleteINPlatformConfigurationRepositoryTx(ctx context.Context, tx *gorm.DB, id []int64) error {
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("platform_configuration").Model(&models.PlatformConfiguration{}).WithContext(ctx)

	db = db.Exec("delete from platform_configuration where id in ?", id)

	return db.Error
}

func (rc *RepositoryContext) GetExistsPlatformConfigurationAggregatorRepository(ctx context.Context, configuration *models.Configuration, platform *models.Platforms) (bool, error) {
	var (
		db                     = rc.db
		platformConfigurations []*models.PlatformConfiguration
		err                    error
	)

	db = db.Table("platform_configuration").Model(&models.PlatformConfiguration{}).WithContext(ctx)

	db = db.Joins("JOIN configurations on configurations.id = platform_configuration.configuration_id")

	db = db.Where("configurations.aggregator_id = ? and platform_configuration.platform_id = ?", configuration.AggregatorId, platform.Id)

	err = db.Find(&platformConfigurations).Error
	if err != nil {
		log.Error().Err(err).Msg("error query find platform configuration data exists")
		return false, err
	}

	if len(platformConfigurations) == 0 {
		return false, nil
	}

	return true, nil
}
