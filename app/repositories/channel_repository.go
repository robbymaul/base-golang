package repositories

import (
	"context"
	"fmt"
	"paymentserviceklink/app/models"
	"paymentserviceklink/pkg/pagination"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (rc *RepositoryContext) InsertChannelRepository(ctx context.Context, channel *models.Channel) (*models.Channel, error) {
	var (
		db = rc.db
	)

	db = db.Table("channels").Model(&models.Channel{}).WithContext(ctx)

	err := db.Create(&channel).Error
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (rc *RepositoryContext) InsertBatchChannelRepositoryTx(ctx context.Context, tx *gorm.DB, channel []*models.Channel) ([]*models.Channel, error) {
	log.Debug().Interface("channel data", channel).Msg("insert batch channel repository tx")
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("channels").Model(&models.Channel{}).WithContext(ctx)

	err := db.Create(&channel).Error
	if err != nil {
		log.Debug().Err(err).Msg("error query insert batch channel tx")
		return nil, err
	}

	return channel, nil
}

func (rc *RepositoryContext) FindChannelForPlatformRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) ([]*models.Channel, error) {
	var (
		db      = rc.db
		channel []*models.Channel
		err     error
	)

	db = db.Table("channels").Model(&models.Channel{}).WithContext(ctx)
	db = fn(db)

	err = db.Find(&channel).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get list channel for platform repository")
		return nil, err
	}

	return channel, nil
}

//
//func (rc *RepositoryContext) FindChannelForPlatformRepository(ctx context.Context, filter []pagination.Filter, selectJoin func(db *gorm.DB) *gorm.DB) ([]*models.Channel, error) {
//	var (
//		db       = rc.db
//		channels []*models.Channel
//		err      error
//	)
//
//	db = db.Table("channels").Model(&models.Channel{}).WithContext(ctx)
//
//	db = selectJoin(db)
//
//	query, args := rc.SearchQuery(filter, "and")
//
//	db = db.Where(query, args...)
//
//	db = db.Order("name " + "asc")
//
//	err = db.Find(&channels).Error
//	if err != nil {
//		log.Error().Err(err).Msg("error query find channel for platform repository")
//		return nil, err
//	}
//
//	return channels, nil
//}

func (rc *RepositoryContext) GetChannelByIdChannelRepository(ctx context.Context, channelId int64) (paymentMethod *models.Channel, err error) {
	var (
		db = rc.db
	)

	db = db.Table("channels").Model(&models.Channel{}).WithContext(ctx)

	db = db.Where("id = ?", channelId)

	err = db.First(&paymentMethod).Error

	return
}

func (rc *RepositoryContext) GetChannelByChannelIdAndPlatformIdRepository(ctx context.Context, channelId int64, platformId int64) (*models.Channel, error) {
	var (
		db      = rc.db
		channel *models.Channel
		err     error
	)

	db = db.Table("channels").Model(&models.Channel{}).WithContext(ctx)
	db = db.Joins("JOIN platform_channel on platform_channel.channel_id = channels.id")

	db = db.Where("channels.id = ? and platform_channel.platform_id = ?", channelId, platformId)

	err = db.First(&channel).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get channel by channel id and platform id")
		return nil, err
	}

	return channel, nil
}

func (rc *RepositoryContext) UpdateChannelRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) error {
	var (
		db = rc.db
	)

	db = db.Table("channels").Model(&models.Channel{}).WithContext(ctx)
	db = fn(db)

	err := db.Error
	if err != nil {
		log.Error().Err(err).Msg("error update channel repository")
		return err
	}

	return nil
}

func (rc *RepositoryContext) UpdateChannelRepositoryTx(ctx context.Context, tx *gorm.DB, fn func(db *gorm.DB) *gorm.DB) error {
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("channels").Model(&models.Channel{}).WithContext(ctx)
	db = fn(db)

	err := db.Error
	if err != nil {
		log.Error().Err(err).Msg("error update channel repository")
		return err
	}

	return nil
}

func (rc *RepositoryContext) FindChannelRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) ([]*models.Channel, error) {
	var (
		db       = rc.db
		channels []*models.Channel
		err      error
	)

	db = db.Table("channels").Model(&models.Channel{}).WithContext(ctx)
	db = fn(db)

	err = db.Find(&channels).Error
	if err != nil {
		log.Error().Err(err).Msg("query error find channel")
		return nil, err
	}

	return channels, nil
}

func (rc *RepositoryContext) GetChannelByIdRepository(ctx context.Context, id int64) (*models.Channel, error) {
	var (
		db      = rc.db
		channel *models.Channel
		err     error
	)

	db = db.Table("channels").Model(&models.Channel{}).WithContext(ctx)
	db = db.Where("id = ?", id)

	err = db.First(&channel).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get channel by id")
		return nil, err
	}

	return channel, nil
}

func (rc *RepositoryContext) GetChannelINByRepository(ctx context.Context, channelId []int64) ([]*models.Channel, error) {
	var (
		db       = rc.db
		channels []*models.Channel
		err      error
	)

	db = db.Table("channels").Model(&models.Channel{}).WithContext(ctx)

	//db = db.Where("aggregator_id IN ?", aggregator)

	if len(channelId) > 0 {
		db = db.Not("id", channelId)
	}

	err = db.Find(&channels).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get channel in by aggregator id")
		return nil, err
	}

	return channels, nil
}

func (rc *RepositoryContext) GetListChannelInPlatformRepository(ctx context.Context, platformId int64, pages *pagination.Pages) ([]*models.Channel, error) {
	var (
		db       = rc.db
		channels []*models.Channel
		err      error
	)

	db = db.Table("channels").Model(&models.Channel{}).WithContext(ctx)
	db = db.Joins("JOIN platform_channel on platform_channel.channel_id = channels.id")

	query, args := rc.SearchQuery(pages.Filters, pages.JoinOperator)

	db = db.Where(query, args...)
	db = db.Where("platform_channel.platform_id = ?", platformId)

	err = db.Find(&channels).Error
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("error query find channels join platform channel where platform_channel.platform_id = %v", platformId))
		return nil, err
	}

	return channels, nil
}

func (rc *RepositoryContext) GetChannelRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) (*models.Channel, error) {
	var (
		db      = rc.db
		channel *models.Channel
		err     error
	)

	db = db.Table("channels").Model(&models.Channel{}).WithContext(ctx)
	db = fn(db)

	err = db.First(&channel).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get channel by filter")
		return nil, err
	}

	return channel, nil
}
