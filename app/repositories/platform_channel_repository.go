package repositories

import (
	"context"
	"fmt"
	"paymentserviceklink/app/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (rc *RepositoryContext) InsertBatchPlatformChannelRepositoryTx(ctx context.Context, tx *gorm.DB, platformChannel []*models.PlatformChannel) ([]*models.PlatformChannel, error) {
	log.Debug().Interface("data platform channel", platformChannel).Msg("insert batch platform channel repository")
	var (
		db  = rc.db
		err error
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("platform_channel").Model(&models.PlatformChannel{}).WithContext(ctx)

	err = db.Create(&platformChannel).Error
	if err != nil {
		log.Error().Err(err).Msg("error insert batch platform channel repository")
		return nil, err
	}

	return platformChannel, err
}

func (rc *RepositoryContext) GetExistsPlatformChannelRepository(ctx context.Context, channel *models.Channel, platform *models.Platforms) (bool, error) {
	var (
		db              = rc.db
		platformChannel []*models.PlatformChannel
		err             error
	)

	db = db.Table("platform_channel").Model(&models.PlatformChannel{}).WithContext(ctx)

	db = db.Where("channel_id = ? and platform_id = ?", channel.Id, platform.Id)

	err = db.Find(&platformChannel).Error
	if err != nil {
		log.Error().Err(err).Msg("error query find platform channel data exists")
		return false, err
	}

	if len(platformChannel) == 0 {
		return false, nil
	}

	return true, nil
}

func (rc *RepositoryContext) GetPlatformChannelRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) (*models.PlatformChannel, error) {
	var (
		db              = rc.db
		platformChannel *models.PlatformChannel
		err             error
	)

	db = db.Table("platform_channel").Model(&models.PlatformChannel{}).WithContext(ctx)
	db = fn(db)

	err = db.First(&platformChannel).Error
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("error query get first platform channel"))
		return nil, err
	}

	return platformChannel, nil
}

func (rc *RepositoryContext) DeleteINPlatformChannelRepositoryTx(ctx context.Context, tx *gorm.DB, id []int64) error {
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("platform_channel").Model(&models.PlatformChannel{}).WithContext(ctx)

	db = db.Exec("delete from platform_channel where id in ?", id)

	return db.Error
}

func (rc *RepositoryContext) GetPlatformChannelByPlatformIdRepository(ctx context.Context, platformId int64) ([]*models.PlatformChannel, error) {
	var (
		db               = rc.db
		platformChannels []*models.PlatformChannel
		err              error
	)

	db = db.Table("platform_channel").Model(&models.PlatformChannel{}).WithContext(ctx)
	db = db.Where("platform_id = ?", platformId)
	err = db.Find(&platformChannels).Error
	if err != nil {
		log.Error().Err(err).Msg("error query find platform channel by platform id")
		return nil, err
	}

	return platformChannels, nil
}
