package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/repositories"
	"paymentserviceklink/app/web"
	"paymentserviceklink/config"
	pkgjwt "paymentserviceklink/pkg/jwt"
	"paymentserviceklink/pkg/pagination"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ChannelService struct {
	service     *Service
	serviceName string
}

func NewChannelService(ctx context.Context, repo *repositories.RepositoryContext, cfg *config.Config) *ChannelService {
	return &ChannelService{
		service:     NewService(ctx, repo, cfg),
		serviceName: "ChannelService",
	}
}

func (s *ChannelService) AdminCreateChannelService(session *pkgjwt.JwtResponse, payload []web.CreateChannelRequest) ([]*web.DetailChannelResponse, error) {
	var err error
	log.Debug().Interface("payload", payload).
		Interface("context", s.serviceName).Msg("admin create channel service")
	//
	//aggregator, err := s.service.repository.GetAggregatorByIdRepository(s.service.ctx, aggregatorId)
	//if err != nil {
	//	log.Error().Err(err).Str("context", s.serviceName).Msg("get aggregator by id repository error")
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
	//	}
	//
	//	return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	//}
	//log.Debug().Interface("aggregator", aggregator).Interface("context", s.serviceName).Msg("result data get aggregator by id")

	//if !aggregator.IsActive {
	//	return nil, helpers.NewErrorTrace(fmt.Errorf("aggregator is not active"), s.serviceName).WithStatusCode(http.StatusForbidden)
	//}

	channels := make([]*models.Channel, 0, len(payload))

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
		for i, ch := range payload {
			log.Debug().Interface("index", i).Interface("data channel", ch).Interface("context", s.serviceName).
				Msg("iteration append data payload")

			code := s.generateCode(ch.BankName, ch.PaymentMethod)

			channel := &models.Channel{
				//AggregatorId:    aggregator.Id,
				Code:            code,
				Name:            ch.Name,
				PaymentMethod:   ch.PaymentMethod,
				TransactionType: ch.TransactionType,
				//Provider:        aggregator.Slug,
				Currency:      ch.Currency,
				FeeType:       ch.FeeType,
				FeeAmount:     ch.FeeFixed,
				FeePercentage: ch.FeePercentage,
				IsActive:      true,
				//IsEspay:         ch.IsEspay,
				ProductName: ch.ProductName,
				ProductCode: ch.ProductCode,
				BankName:    ch.BankName,
				BankCode:    ch.BankCode,
				Instruction: ch.Instruction,
			}

			for _, image := range ch.Image {
				channelImage := &models.ChannelImage{
					ChannelID:     channel.Id,
					FileName:      image.FileName,
					SizeType:      image.SizeType,
					GeometricType: image.Geometric,
				}

				channel.ChannelImage = append(channel.ChannelImage, channelImage)
			}

			channels = append(channels, channel)
		}
		log.Debug().Interface("channels data", channels).Interface("context", s.serviceName).Msg("data channels for insert database")

		channels, err = s.service.repository.InsertBatchChannelRepositoryTx(s.service.ctx, tx, channels)
		if err != nil {
			log.Error().Err(err).Msg("insert batch channel repository failed")
			return helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("channels", channels).Interface("context", s.serviceName).Msg("data channels after insert database")

		adminActivityLog := models.NewAdminActivityLogs(
			session.Id,
			enums.ACTION_ADMIN_ACCTIVITY_CREATE,
			enums.RESOURCE_ADMIN_ACTIVITY_LOG_CHANNEL,
			"",
			fmt.Sprint(enums.RESOURCE_ADMIN_ACTIVITY_LOG_CHANNEL, enums.ACTION_ADMIN_ACCTIVITY_CREATE),
			enums.NULL_STRING,
			session.Sub,
			payload,
			channels,
		)

		err = s.service.repository.InsertAdminActivityLogRepositoryTx(s.service.ctx, tx, adminActivityLog)
		if err != nil {
			log.Error().Err(err).Interface("context", s.serviceName).Msg("insert admin activity log repository tx")
		}

		return nil
	})
	if err != nil {
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	result := make([]*web.DetailChannelResponse, 0, len(channels))

	for _, ch := range channels {
		result = append(result, s.service.mapToDetailChannelResponse(ch))
	}

	return result, nil
}

func (s *ChannelService) AdminGetListChannelService(pages *pagination.Pages) (map[enums.Currency]map[enums.PaymentMethod][]*web.DetailChannelResponse, error) {
	log.Debug().Interface("pages", pages).Msg("admin get list channel service")

	// check pages filter
	filter, err := helpers.FilterColumnValidation(pages.Filters, models.AllowedFilterColumnChannel())
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("filter column validation")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

	channels, err := s.service.repository.FindChannelRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("channels.id," +
			"channels.code," +
			"channels.name," +
			"channels.currency," +
			"channels.fee_type," +
			"channels.fee_amount," +
			"channels.is_active," +
			"channels.payment_method," +
			"channels.transaction_type," +
			"channels.product_name," +
			"channels.product_code," +
			"channels.bank_code," +
			"channels.bank_name," +
			"channels.instruction",
		)
		//db = db.Limit(pages.Limit()).Offset(pages.Offset())
		query, args := s.service.repository.SearchQuery(filter, pages.JoinOperator)
		if query != "" {
			db = db.Where(query, args...)
		}
		db = db.Order("channels.name " + pages.Sort)
		db = db.Preload("ChannelImage", func(db *gorm.DB) *gorm.DB {
			db = db.Select("channel_images.id," +
				"channel_images.channel_id," +
				"channel_images.file_name," +
				"channel_images.size_type," +
				"channel_images.geometric_type")
			return db
		})
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get list payment method repository error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data channel", channels).Msg("result admin get list channel repository")

	return s.service.mapToChannelResponseWithCurrency(channels), nil
}

func (s *ChannelService) GetListChannelService(platformId int64, pages *pagination.Pages, payload *web.GetListChannelRequest) (*web.ListResponse, error) {
	filter, err := helpers.FilterColumnValidation(pages.Filters, models.AllowedFilterColumnChannel())
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("filter column validation")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

	paymentMethods, err := s.service.repository.FindChannelRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("channels.id," +
			"channels.name," +
			"channels.currency," +
			"channels.fee_type," +
			"channels.fee_amount," +
			"channels.is_active," +
			"channels.payment_method," +
			"channels.bank_name",
		)
		//db = db.Limit(pages.Limit()).Offset(pages.Offset())
		query, args := s.service.repository.SearchQuery(filter, pages.JoinOperator)
		if query != "" {
			db = db.Where(query, args...)
		}
		db = db.Joins("JOIN platform_channel on platform_channel.channel_id = channels.id")
		db = db.Where("platform_channel.platform_id = ?", platformId)
		//db = db.Where("platform_channel.platform_id = ?", platformId)
		db = db.Where("channels.currency = ?", payload.Currency)
		db = db.Order("channels.name " + pages.Sort)
		db = db.Preload("ChannelImage", func(db *gorm.DB) *gorm.DB {
			db = db.Select("channel_images.channel_id," +
				"channel_images.file_name," +
				"channel_images.size_type," +
				"channel_images.geometric_type",
			)
			return db
		})
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get list payment method repository error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	response := make([]*web.DetailChannelResponse, 0, len(paymentMethods))

	for _, paymentMethod := range paymentMethods {
		response = append(response, s.service.mapToDetailChannelResponse(paymentMethod))
	}

	return &web.ListResponse{Items: response, Metadata: pages.GetMetadata()}, nil
}

func (s *ChannelService) ClientGetListChannelService(payload *web.GetListChannelRequest) ([]*web.ChannelResponse, error) {
	s.serviceName = "ChannelService.ClientGetListChannel"

	// get client auth
	platform, err := GetClientAuth(s.service)
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("get client auth error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusUnauthorized)
	}
	log.Debug().Interface("data platform", platform).Interface("context", s.serviceName).Msg("result data platform client auth")

	// get data k-wallet member
	//kWallets, err := s.service.repository.FindKWalletRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
	//	db = db.Select("k_wallet.member_id," +
	//		"k_wallet.full_name," +
	//		"k_wallet.no_rekening," +
	//		"k_wallet.balance," +
	//		"k_wallet.symbol," +
	//		"k_wallet.status")
	//	db = db.Where("k_wallet.member_id = ? and k_wallet.currency = ?", payload.MemberId, payload.Currency)
	//	return db
	//})

	// find channels
	channels, err := s.service.repository.FindChannelRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("channels.id," +
			"channels.name," +
			"channels.currency," +
			"channels.fee_type," +
			"channels.fee_amount," +
			"channels.is_active," +
			"channels.payment_method," +
			"channels.bank_name," +
			"channels.instruction",
		)
		db = db.Joins("JOIN platform_channel on platform_channel.channel_id = channels.id")
		db = db.Where("channels.currency = ? and platform_channel.platform_id = ? ", payload.Currency, platform.Id)
		db = db.Preload("ChannelImage", func(db *gorm.DB) *gorm.DB {
			db = db.Select("channel_images.channel_id," +
				"channel_images.file_name," +
				"channel_images.size_type," +
				"channel_images.geometric_type")
			return db
		})
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get list payment method repository error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data channels", channels).Interface("context", s.serviceName).Msg("result data channel get list channel repository")

	response := make([]*web.ChannelResponse, 0, len(channels))

	for _, channel := range channels {
		response = append(response, s.service.mapToChannelResponse(channel))
	}

	//return s.service.mapToChannelResponseWithoutCurrency(channels, kWallets), nil
	return response, nil
}

func (s *ChannelService) AdminGetDetailChannelService(channelId int64) (*web.DetailChannelResponse, error) {
	log.Debug().Interface("channel_id", channelId).
		Interface("context", s.serviceName).
		Msg("admin get detail channel service")

	channel, err := s.service.repository.GetChannelRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("channels.id," +
			"channels.code," +
			"channels.name," +
			"channels.provider," +
			"channels.currency," +
			"channels.fee_type," +
			"channels.fee_amount," +
			"channels.is_active," +
			"channels.payment_method," +
			"channels.product_name," +
			"channels.product_code," +
			"channels.bank_code," +
			"channels.bank_name")
		db = db.Where("channels.id = ?", channelId)
		db = db.Preload("ChannelImage", func(db *gorm.DB) *gorm.DB {
			db = db.Select("channel_images.id," +
				"channel_images.channel_id," +
				"channel_images.file_name," +
				"channel_images.size_type," +
				"channel_images.geometric_type")
			return db
		})
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get channel by id repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("result data channel", channel).Interface("context", s.serviceName).
		Msg("result data channel get channel by id aggregator and channel id")

	return s.service.mapToDetailChannelResponse(channel), nil
}

func (s *ChannelService) AdminUpdateChannelService(session *pkgjwt.JwtResponse, payload *web.DetailChannelResponse, channelId int64) (*web.DetailChannelResponse, error) {
	log.Debug().Interface("session", session).Interface("payload", payload).
		Interface("channel_id", channelId).Interface("context", s.serviceName).Msg("admin update channel service")

	if payload.Id != channelId {
		return nil, helpers.NewErrorTrace(fmt.Errorf("channel id param = %v invalid payload id = %v", channelId, payload.Id), s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

	channel, err := s.service.repository.GetChannelRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("channels.id," +
			"channels.code," +
			"channels.name," +
			"channels.provider," +
			"channels.currency," +
			"channels.fee_type," +
			"channels.fee_amount," +
			"channels.is_active," +
			"channels.payment_method," +
			"channels.product_name," +
			"channels.product_code," +
			"channels.bank_code," +
			"channels.bank_name," +
			"channels.instruction")
		db = db.Where("channels.id = ?", channelId)
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get payment method by id repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data channel", channel).Interface("context", s.serviceName).Msg("result data channel get channel by id aggregator and id channel repository")

	s.updateChannel(&channel, payload)
	log.Debug().Interface("data update channel", channel).Interface("context", s.serviceName).Msg("update data channel")

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
		err = s.service.repository.UpdateChannelRepositoryTx(s.service.ctx, tx, func(db *gorm.DB) *gorm.DB {
			db = db.Where("channels.id = ?", channel.Id)
			updateColumn := map[string]interface{}{
				"code":       channel.Code,
				"name":       channel.Name,
				"provider":   channel.Provider,
				"currency":   channel.Currency,
				"fee_type":   channel.FeeType,
				"fee_amount": channel.FeeAmount,
				"is_active":  channel.IsActive,
				"updated_at": time.Now(),
			}
			db = db.UpdateColumns(updateColumn)
			return db
		})
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("update payment method repository error")
			return err
		}

		adminActivityLog := models.NewAdminActivityLogs(
			session.Id,
			enums.ACTION_ADMIN_ACCTIVITY_UPDATE,
			enums.RESOURCE_ADMIN_ACTIVITY_LOG_CHANNEL,
			fmt.Sprint(channel.Id),
			fmt.Sprint(enums.RESOURCE_ADMIN_ACTIVITY_LOG_CHANNEL, enums.ACTION_ADMIN_ACCTIVITY_UPDATE),
			enums.NULL_STRING,
			session.Sub,
			payload,
			channel,
		)

		err = s.service.repository.InsertAdminActivityLogRepositoryTx(s.service.ctx, tx, adminActivityLog)
		if err != nil {
			log.Error().Err(err).Interface("context", s.serviceName).Msg("insert admin activity log repository tx")
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("admin update channel service transaction error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return s.service.mapToDetailChannelResponse(channel), nil
}

//
//func (s *ChannelService) mapToDetailChannelResponse(channel *models.Channel) *web.DetailChannelResponse {
//	channelImages := make([]*web.ImageResponse, 0, len(channel.ChannelImage))
//
//	for _, channelImage := range channel.ChannelImage {
//		imageResponse := s.service.mapToChannelImageResponse(channelImage)
//		if imageResponse != nil {
//			channelImages = append(channelImages, imageResponse)
//		}
//	}
//
//	return &web.DetailChannelResponse{
//		Id:            channel.Id,
//		Code:          channel.Code,
//		Name:          channel.Name,
//		PaymentMethod: channel.PaymentMethod,
//		//TransactionType: channel.TransactionType,
//		//Provider:        channel.Provider,
//		Currency:    channel.Currency,
//		FeeType:     channel.FeeType,
//		FeeAmount:   channel.FeeAmount.IntPart(),
//		IsActive:    channel.IsActive,
//		ProductName: channel.ProductName,
//		ProductCode: channel.ProductCode,
//		BankName:    channel.BankName,
//		BankCode:    channel.BankCode,
//		Images:      channelImages,
//	}
//}

func (s *ChannelService) updateChannel(channel **models.Channel, payload *web.DetailChannelResponse) {
	//if (*channel).Code != payload.Code {
	//	(*channel).Code = payload.Code
	//}

	if (*channel).Name != payload.Name {
		(*channel).Name = payload.Name
	}

	//if (*channel).Provider != payload.Provider {
	//	(*channel).Provider = payload.Provider
	//}

	if (*channel).Currency != payload.Currency {
		(*channel).Currency = payload.Currency
	}

	if (*channel).FeeType != payload.FeeType {
		(*channel).FeeType = payload.FeeType
	}

	if !decimal.NewFromInt((*channel).FeeAmount).Equal(decimal.NewFromInt(payload.FeeFixed)) {
		(*channel).FeeAmount = payload.FeeFixed
	}

	if !decimal.NewFromFloat32((*channel).FeePercentage).Equals(decimal.NewFromFloat32(payload.FeePercentage)) {
		(*channel).FeePercentage = payload.FeePercentage
	}

	if (*channel).IsActive != payload.IsActive {
		(*channel).IsActive = payload.IsActive
	}

	if (*channel).PaymentMethod != payload.PaymentMethod {
		(*channel).PaymentMethod = payload.PaymentMethod
	}

	//if (*channel).TransactionType != payload.TransactionType {
	//	(*channel).TransactionType = payload.TransactionType
	//}

	if (*channel).BankName != payload.BankName {
		(*channel).BankName = payload.BankName
	}
}

func (s *ChannelService) generateCode(bankName enums.Channel, method enums.PaymentMethod) string {
	return fmt.Sprint(strings.ToUpper(string(bankName)), "_", strings.ToUpper(string(method)))
}

func (s *ChannelService) mapToDetailChannelKWalletResponse(ch *models.Channel, kWallet *models.KWallet) *web.DetailChannelResponse {
	channelImageResponse := make([]*web.ImageResponse, 0, len(ch.ChannelImage))

	for _, channelImage := range ch.ChannelImage {
		channelImageResponse = append(channelImageResponse, s.service.mapToChannelImageResponse(channelImage))
	}

	return &web.DetailChannelResponse{
		Id:            ch.Id,
		Code:          ch.Code,
		Name:          ch.Name,
		PaymentMethod: ch.PaymentMethod,
		//TransactionType: ch.TransactionType,
		//Provider:        ch.Provider,
		Currency:      ch.Currency,
		FeeType:       ch.FeeType,
		FeeFixed:      ch.FeeAmount,
		FeePercentage: ch.FeePercentage,
		IsActive:      ch.IsActive,
		MemberID:      kWallet.MemberID,
		FullName:      kWallet.FullName,
		NoRekening:    kWallet.NoRekening,
		GenVa:         kWallet.GenVA,
		Balance:       kWallet.Balance.IntPart(),
		Symbol:        kWallet.Symbol,
		Status:        kWallet.Status,
		ProductName:   ch.ProductName,
		ProductCode:   ch.ProductCode,
		BankName:      ch.BankName,
		BankCode:      ch.BankCode,
		Images:        channelImageResponse,
	}
}

//func (s *ChannelService) mapToChannelImageResponse(channelImage *models.ChannelImage) *web.ImageResponse {
//	assetPath, err := s.assetService.GetAssetPath(enums.ASSET_LOGO, channelImage.FileName)
//	if err != nil {
//		return nil
//	}
//
//	return &web.ImageResponse{
//		FileUrl:   s.assetService.urlFile(assetPath),
//		SizeType:  channelImage.SizeType,
//		Geometric: channelImage.GeometricType,
//	}
//}
