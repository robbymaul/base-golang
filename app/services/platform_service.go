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
	"paymentserviceklink/pkg/util"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type PlatformService struct {
	service     *Service
	serviceName string
}

func NewPlatformService(ctx context.Context, repo *repositories.RepositoryContext, cfg *config.Config) *PlatformService {
	return &PlatformService{
		service:     NewService(ctx, repo, cfg),
		serviceName: "PlatformService",
	}
}

func (s *PlatformService) AdminCreatePlatformService(session *pkgjwt.JwtResponse, payload *web.CreatePlatformRequest) (*web.DetailPlatformResponse, error) {
	s.serviceName = "PlatformService.AdminCreatePlatform"
	log.Debug().Interface("payload", payload).Interface("context", s.serviceName).Msg("create platform service")

	var err error

	// generate random string api key
	apiKey := util.GenerateRandomString(30)
	log.Debug().Str("api_key", apiKey).Interface("context", s.serviceName).Msg("generate random string api key")

	// generate random string secret key
	secretKey := util.GenerateRandomString(50)
	log.Debug().Str("secret_key", secretKey).Interface("context", s.serviceName).Msg("generate random string secret key")

	// encrypted api key string
	//	encryptApiKey := helpers.EncryptAES(apiKey)
	//	log.Debug().Str("api_key", apiKey).Msg("encrypt api key string")

	// encrypted secret key string
	//	encryptSecretKey := helpers.EncryptAES(secretKey)
	//	log.Debug().Str("api_key", apiKey).Msg("encrypt secret key string")

	// platform code
	platformCode := strings.Join(strings.Split(payload.Name, " "), "_")

	platforms := &models.Platforms{
		Code:        platformCode,
		Name:        payload.Name,
		Description: payload.Description,
		ApiKey:      apiKey,
		SecretKey:   secretKey,
		IsActive:    true,
	}
	log.Debug().Interface("data platforms", platforms).Interface("context", s.serviceName).Msg("data platforms for insert")

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
		platforms, err = s.service.repository.InsertPlatformRepositoryTx(s.service.ctx, tx, platforms)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("insert platform repository error")
			return err
		}
		log.Debug().Interface("data platforms", platforms).Interface("context", s.serviceName).Msg("result data insert platform repository")

		if len(payload.Channel) > 0 {
			log.Debug().Interface("payload channel", payload.Channel).Msg("process payload channel for platform channel")
			platformChannel := make([]*models.PlatformChannel, 0, len(payload.Channel))

			for i, ch := range payload.Channel {
				log.Debug().Interface("index", i).Interface("data", ch).Msg("loop payload channel process")

				platformChannel = append(platformChannel, &models.PlatformChannel{
					PlatformId: platforms.Id,
					ChannelId:  ch.Id,
				})
			}

			_, errChannel := s.service.repository.InsertBatchPlatformChannelRepositoryTx(s.service.ctx, tx, platformChannel)
			if errChannel != nil {
				log.Error().Err(errChannel).Msg("insert batch platform channel repository tx error")
				return errChannel
			}
		}

		if len(payload.Configuration) > 0 {
			log.Debug().Interface("payload configuration", payload.Configuration).Msg("process payload channel for platform configuration")
			platformConfiguration := make([]*models.PlatformConfiguration, 0, len(payload.Configuration))

			for i, conf := range payload.Configuration {
				log.Debug().Interface("index", i).Interface("data", conf).Msg("loop payload configuration process")

				platformConfiguration = append(platformConfiguration, &models.PlatformConfiguration{
					ConfigurationId: conf.Id,
					PlatformId:      platforms.Id,
				})
			}

			_, errConf := s.service.repository.InsertBatchPlatformConfigurationRepositoryTx(s.service.ctx, tx, platformConfiguration)
			if errConf != nil {
				log.Debug().Err(err).Msg("insert batch platform configuration repository tx error")
				return errConf
			}
		}

		adminActivityLog := models.NewAdminActivityLogs(
			session.Id,
			enums.ACTION_ADMIN_ACCTIVITY_CREATE,
			enums.RESOURCE_ADMIN_ACTIVITY_LOG_PLATFORM,
			fmt.Sprint(platforms.Id),
			fmt.Sprint(enums.RESOURCE_ADMIN_ACTIVITY_LOG_PLATFORM, enums.ACTION_ADMIN_ACCTIVITY_CREATE),
			enums.NULL_STRING,
			session.Sub,
			payload,
			platforms,
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

	return s.service.mapToDetailPlatformResponse(platforms), nil
}

func (s *PlatformService) AdminGetListPlatformService(pages *pagination.Pages) (*web.ListResponse, error) {
	var platforms []*models.Platforms
	var totalCount int64
	var err error

	s.serviceName = "PlatformService.AdminGetListPlatform"

	// check pages filter
	_, err = helpers.FilterColumnValidation(pages.Filters, models.AllowedFilterColumnPlatform())
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("filter column validation")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

	log.Debug().Interface("pages", pages).Interface("context", s.serviceName).
		Msg("admin get list platform service")

	ctx, cancelFn := context.WithCancel(s.service.ctx)
	eg := errgroup.Group{}

	eg.Go(func() error {
		platforms, err = s.service.repository.FindPlatformRepository(ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Select("platforms.id," +
				"platforms.code," +
				"platforms.name," +
				"platforms.description," +
				"platforms.is_active," +
				"platforms.notification_url",
			)
			query, args := s.service.repository.SearchQuery(pages.Filters, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			db = db.Limit(pages.Limit()).Offset(pages.Offset())
			db = db.Order("platforms.id " + pages.Sort)
			return db
		})
		if err != nil {
			cancelFn()
			log.Error().Err(err).Str("context", s.serviceName).Msg("get list platform repository error")
			return err
		}
		log.Debug().Interface("data platforms", platforms).Interface("context", s.serviceName).Msg("result data get list platform repository")

		return nil
	})

	eg.Go(func() error {
		totalCount, err = s.service.repository.GetTotalCountPlatformRepository(ctx, pages, func(db *gorm.DB) *gorm.DB {
			query, args := s.service.repository.SearchQuery(pages.Filters, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			return db
		})
		if err != nil {
			cancelFn()
			log.Error().Err(err).Str("context", s.serviceName).Msg("get total count platform repository error")
			return err
		}
		log.Debug().Interface("total count", totalCount).Interface("context", s.serviceName).Msg("result data get total count platform repository")

		return nil
	})

	err = eg.Wait()
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("err sync group error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	pages.TotalCount = int(totalCount)

	response := make([]*web.DetailPlatformResponse, 0, len(platforms))

	for _, platform := range platforms {
		response = append(response, s.service.mapToDetailPlatformResponse(platform))
	}

	return &web.ListResponse{Items: response, Metadata: pages.GetMetadata()}, nil
}

func (s *PlatformService) AdminGetDetailPlatformService(platformId int64) (*web.DetailPlatformResponse, error) {
	s.serviceName = "PlatformService.AdminGetDetailPlatform"
	log.Debug().Interface("platform_id", platformId).Interface("context", s.serviceName).Msg("admin get detail platform service")

	platform, err := s.service.repository.GetPlatformRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("platforms.id," +
			"platforms.code," +
			"platforms.name," +
			"platforms.description," +
			"platforms.api_key," +
			"platforms.secret_key," +
			"platforms.is_active," +
			"platforms.notification_url",
		)
		db = db.Where("platforms.id = ?", platformId)
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get platform by id repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data platform", platform).Msg("result data get platform repository")

	return s.service.mapToDetailPlatformResponse(platform), nil
}

func (s *PlatformService) AdminUpdatePlatformService(session *pkgjwt.JwtResponse, platformId int64, payload *web.DetailPlatformResponse) (*web.DetailPlatformResponse, error) {
	s.serviceName = "PlatformService.AdminUpdatePlatform"

	log.Debug().Interface("platform_id", platformId).Interface("payload", payload).Interface("context", s.serviceName).
		Msg("admin update platform service")

	if platformId != payload.Id {
		return nil, helpers.NewErrorTrace(fmt.Errorf("invalid platform id request update %v", payload.Name), s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

	platform, err := s.service.repository.GetPlatformRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("platforms.id," +
			"platforms.name," +
			"platforms.description," +
			"platforms.is_active," +
			"platforms.notification_url",
		)
		db = db.Where("platforms.id = ?", platformId)
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get platform by id repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data platform", platform).Interface("context", s.serviceName).Msg("result data get platform repository")

	s.updatePlatform(&platform, payload)
	log.Debug().Interface("update data platform", platform).Interface("context", s.serviceName).Msg("result update platform")

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
		err = s.service.repository.UpdatePlatformRepository(s.service.ctx, platform, func(db *gorm.DB) *gorm.DB {
			db = db.Where("platforms.id = ? ", platform.Id)
			updateColumn := map[string]interface{}{
				//"code":        platform.Code,
				"name":        platform.Name,
				"description": platform.Description,
				"is_active":   platform.IsActive,
				"updated_at":  time.Now(),
			}
			db = db.UpdateColumns(updateColumn)
			return db
		})
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("update platform repository error")
			return err
		}

		adminActivityLog := models.NewAdminActivityLogs(
			session.Id,
			enums.ACTION_ADMIN_ACCTIVITY_UPDATE,
			enums.RESOURCE_ADMIN_ACTIVITY_LOG_PLATFORM,
			fmt.Sprint(platformId),
			fmt.Sprint(enums.RESOURCE_ADMIN_ACTIVITY_LOG_PLATFORM, enums.ACTION_ADMIN_ACCTIVITY_UPDATE),
			enums.NULL_STRING,
			session.Sub,
			payload,
			platform,
		)

		err = s.service.repository.InsertAdminActivityLogRepositoryTx(s.service.ctx, tx, adminActivityLog)
		if err != nil {
			log.Error().Err(err).Interface("context", s.serviceName).Msg("insert admin activity log repository tx")
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("update platform transaction error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return s.service.mapToDetailPlatformResponse(platform), nil
}

func (s *PlatformService) AdminUpdatePlatformSecretKeyService(session *pkgjwt.JwtResponse, platformId int64) (*web.DetailPlatformResponse, error) {
	s.serviceName = "PlatformService.AdminUpdatePlatformSecretKey"

	log.Debug().Interface("session", session).Interface("platform_id", platformId).Interface("context", s.serviceName).
		Msg("admin update platform secret key service")

	platform, err := s.service.repository.GetPlatformRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("platforms.id," +
			"platforms.name," +
			"platforms.description," +
			"platforms.api_key," +
			"platforms.secret_key," +
			"platforms.is_active," +
			"platforms.notification_url",
		)
		db = db.Where("platforms.id = ?", platformId)
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get platform by id repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data platform", platform).Interface("context", s.serviceName).Msg("result data get platform by id repository")

	s.updateSecretKey(&platform)
	log.Debug().Interface("update data platform", platform).Interface("context", s.serviceName).Msg("result update platform secret key")

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {

		err = s.service.repository.UpdatePlatformRepository(s.service.ctx, platform, func(db *gorm.DB) *gorm.DB {
			db = db.Where("platforms.id = ? ", platform.Id)
			updateColumn := map[string]interface{}{
				"api_key":    platform.ApiKey,
				"secret_key": platform.SecretKey,
				"updated_at": time.Now(),
			}
			db = db.UpdateColumns(updateColumn)
			return db
		})
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("update platform repository error")
			return err
		}

		adminActivityLog := models.NewAdminActivityLogs(
			session.Id,
			enums.ACTION_ADMIN_ACCTIVITY_UPDATE,
			enums.RESOURCE_ADMIN_ACTIVITY_LOG_PLATFORM,
			fmt.Sprint(platformId),
			fmt.Sprint(enums.RESOURCE_ADMIN_ACTIVITY_LOG_PLATFORM, enums.ACTION_ADMIN_ACCTIVITY_UPDATE),
			enums.NULL_STRING,
			session.Sub,
			nil,
			platform,
		)

		err = s.service.repository.InsertAdminActivityLogRepositoryTx(s.service.ctx, tx, adminActivityLog)
		if err != nil {
			log.Error().Err(err).Interface("context", s.serviceName).Msg("insert admin activity log repository tx")
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("update platform transaction")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return s.service.mapToDetailPlatformResponse(platform), nil
}

//func (s *PlatformService) mapToDetailPlatformResponse(platform *models.Platforms) *web.DetailPlatformResponse {
//	return &web.DetailPlatformResponse{
//		Id:              platform.Id,
//		Code:            platform.Code,
//		Name:            platform.Name,
//		Description:     platform.Description,
//		ApiKey:          platform.ApiKey,
//		SecretKey:       platform.SecretKey,
//		IsActive:        platform.IsActive,
//		NotificationUrl: platform.NotificationURL,
//		CreatedAt:       platform.CreatedAt,
//		UpdatedAt:       platform.UpdatedAt,
//	}
//}

func (s *PlatformService) updatePlatform(platform **models.Platforms, payload *web.DetailPlatformResponse) {
	//if (*platform).Code != payload.Code {
	//	(*platform).Code = payload.Code
	//}

	if (*platform).Name != payload.Name {
		(*platform).Name = payload.Name
	}

	if (*platform).Description != payload.Description {
		(*platform).Description = payload.Description
	}

	if (*platform).IsActive != payload.IsActive {
		(*platform).IsActive = payload.IsActive
	}

	if (*platform).NotificationURL != payload.NotificationUrl {
		(*platform).NotificationURL = payload.NotificationUrl
	}
}

func (s *PlatformService) updateSecretKey(platform **models.Platforms) {
	secretKey := util.GenerateRandomString(50)
	(*platform).SecretKey = secretKey
}

func (s *PlatformService) AdminAssignmentPlatformConfigurationService(session *pkgjwt.JwtResponse, platformId int64, payload []web.ResponseConfiguration) (any, error) {
	s.serviceName = "PlatformService.AdminAssignmentPlatformConfiguration"

	log.Debug().Interface("session", session).Interface("platform_id", platformId).Interface("payload", payload).
		Interface("context", s.serviceName).Msg("admin assignment platform configuration service")

	platform, err := s.service.repository.GetPlatformRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("platforms.id," +
			"platforms.name",
		)
		db = db.Where("platforms.id = ?", platformId)
		return db
	})
	if err != nil {
		log.Debug().Err(err).Interface("context", s.serviceName).Msg("get platform by id repository")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(fmt.Errorf("platform, %v", err), s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data platform", platform).Interface("context", s.serviceName).Msg("result data get platform by id repository")

	platformConfiguration := make([]*models.PlatformConfiguration, 0)

	for idx1, conf := range payload {
		log.Debug().Interface("conf", conf).Interface("index", idx1).Interface("context", s.serviceName).
			Msg("iteration payload configuration")

		configuration, errConf := s.service.repository.GetConfigurationRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Select("configurations.id," +
				"configurations.config_name," +
				"configurations.aggregator_id",
			)
			db = db.Where("configurations.id = ?", conf.Id)
			db = db.Preload("Aggregator", func(db *gorm.DB) *gorm.DB {
				db = db.Select("aggregators.id," +
					"aggregators.name",
				)
				return db
			})
			return db
		})
		if errConf != nil {
			log.Error().Err(errConf).Str("context", s.serviceName).Msg("get configuration by id repository error")
			if errors.Is(errConf, gorm.ErrRecordNotFound) {
				return nil, helpers.NewErrorTrace(fmt.Errorf("configuration %v, %v", conf.ConfigName, errConf), s.serviceName).WithStatusCode(http.StatusNotFound)
			}

			return nil, helpers.NewErrorTrace(errConf, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("data configuration", configuration).Interface("context", s.serviceName).
			Msg("result data configuration get configuration by id repository")

		exists, errExist := s.service.repository.GetExistsPlatformConfigurationAggregatorRepository(s.service.ctx, configuration, platform)
		if errExist != nil {
			log.Debug().Err(err).Interface("context", s.serviceName).Msg("get exists platform configuration aggregator repository")
			return nil, helpers.NewErrorTrace(errExist, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("exists", exists).Interface("context", s.serviceName).
			Msg("result data exist get platform configuration aggregator")

		if exists {
			return nil, helpers.NewErrorTrace(fmt.Errorf("the platform get exists configuration on aggregator %v, please check data %v", configuration.Aggregator.Name, platform.Name), s.serviceName).
				WithStatusCode(http.StatusConflict)
		}

		platformConfiguration = append(platformConfiguration, &models.PlatformConfiguration{
			ConfigurationId: configuration.Id,
			PlatformId:      platform.Id,
		})
	}
	log.Debug().Interface("data platform configuration", platformConfiguration).Interface("context", s.serviceName).
		Msg("result data platform configuration after iteration")

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
		platformConfiguration, err = s.service.repository.InsertBatchPlatformConfigurationRepositoryTx(s.service.ctx, tx, platformConfiguration)
		if err != nil {
			log.Debug().Err(err).Interface("context", s.serviceName).Msg("insert batch platform configuration repository tx error")
			return err
		}

		adminActivityLog := models.NewAdminActivityLogs(
			session.Id,
			enums.ACTION_ADMIN_ACCTIVITY_CREATE,
			enums.RESOURCE_ADMIN_ACTIVITY_LOG_PLATFORM_CONFIGURATION,
			fmt.Sprint(),
			fmt.Sprint(enums.RESOURCE_ADMIN_ACTIVITY_LOG_PLATFORM_CONFIGURATION, enums.ACTION_ADMIN_ACCTIVITY_CREATE),
			enums.NULL_STRING,
			session.Sub,
			payload,
			platformConfiguration,
		)

		err = s.service.repository.InsertAdminActivityLogRepositoryTx(s.service.ctx, tx, adminActivityLog)
		if err != nil {
			log.Error().Err(err).Interface("context", s.serviceName).Msg("insert admin activity log repository tx")
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("assignment platform configuration transaction error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return platformConfiguration, nil
}

func (s *PlatformService) AdminRemovalPlatformConfigurationService(session *pkgjwt.JwtResponse, platformId int64, payload []web.ResponseConfiguration) (any, error) {
	s.serviceName = "PlatformService.AdminAssignmentPlatformConfiguration"

	log.Debug().Interface("session", session).Interface("platform_id", platformId).Interface("payload", payload).
		Interface("context", s.serviceName).Msg("admin assignment platform configuration service")

	platform, err := s.service.repository.GetPlatformRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("platforms.id," +
			"platforms.name",
		)
		db = db.Where("platforms.id = ?", platformId)
		return db
	})
	if err != nil {
		log.Debug().Err(err).Interface("context", s.serviceName).Msg("get platform by id repository")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(fmt.Errorf("platform, %v", err), s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data platform", platform).Interface("context", s.serviceName).Msg("result data get platform by id repository")

	platformConfigurationId := make([]int64, 0)

	for idx1, conf := range payload {
		log.Debug().Interface("conf", conf).Interface("index", idx1).Interface("context", s.serviceName).
			Msg("iteration payload configuration")

		configuration, errConf := s.service.repository.GetConfigurationRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Select("configurations.id," +
				"configurations.config_name",
			)
			db = db.Where("configurations.id = ?", conf.Id)
			return db
		})
		if errConf != nil {
			log.Error().Err(errConf).Str("context", s.serviceName).Msg("get configuration by id repository error")
			if errors.Is(errConf, gorm.ErrRecordNotFound) {
				return nil, helpers.NewErrorTrace(fmt.Errorf("configuration %v, %v", conf.ConfigName, errConf), s.serviceName).WithStatusCode(http.StatusNotFound)
			}

			return nil, helpers.NewErrorTrace(errConf, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("data configuration", configuration).Interface("context", s.serviceName).
			Msg("result data configuration get configuration by id repository")

		platformConfiguration, errPg := s.service.repository.GetPlatformConfigurationRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Where("platform_configuration.configuration_id = ? and platform_configuration.platform_id = ?", configuration.Id, platform.Id)
			return db
		})
		if errPg != nil {
			log.Error().Err(errPg).Interface("context", s.serviceName).Msg("get platform configuration by configuration id and platform id repository error")
			if errors.Is(errPg, gorm.ErrRecordNotFound) {
				return nil, helpers.NewErrorTrace(fmt.Errorf("%v platfom has not been entered into aggregator %v configuration %v", platform.Name, configuration.Aggregator.Name, configuration.ConfigName), s.serviceName).
					WithStatusCode(http.StatusNotFound)
			}

			return nil, helpers.NewErrorTrace(errPg, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("data platform configuration", platformConfiguration).Interface("context", s.serviceName).
			Msg("result data platform configuration by configuration id and platform id repository")

		platformConfigurationId = append(platformConfigurationId, platformConfiguration.Id)
	}
	log.Debug().Interface("data platform configuration", platformConfigurationId).Interface("context", s.serviceName).
		Msg("result data platform configuration after iteration")

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
		err = s.service.repository.DeleteINPlatformConfigurationRepositoryTx(s.service.ctx, tx, platformConfigurationId)
		if err != nil {
			log.Error().Err(err).Interface("context", s.serviceName).Msg("delete platform configuration repository tx error")
			return err
		}

		adminActivityLog := models.NewAdminActivityLogs(
			session.Id,
			enums.ACTION_ADMIN_ACCTIVITY_DELETE,
			enums.RESOURCE_ADMIN_ACTIVITY_LOG_PLATFORM_CONFIGURATION,
			fmt.Sprint(),
			fmt.Sprint(enums.RESOURCE_ADMIN_ACTIVITY_LOG_PLATFORM_CONFIGURATION, enums.ACTION_ADMIN_ACCTIVITY_DELETE),
			enums.NULL_STRING,
			session.Sub,
			payload,
			platformConfigurationId,
		)

		err = s.service.repository.InsertAdminActivityLogRepositoryTx(s.service.ctx, tx, adminActivityLog)
		if err != nil {
			log.Error().Err(err).Interface("context", s.serviceName).Msg("insert admin activity log repository tx")
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("removal platform configuration transaction error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return nil, nil
}

func (s *PlatformService) AdminGetListConfigurationInPlatformService(platformId int64, pages *pagination.Pages) (any, error) {
	s.serviceName = "PlatformService.AdminGetListConfigurationInPlatform"

	log.Debug().Interface("platform_id", platformId).Interface("pages", pages).
		Interface("context", s.serviceName).Msg("admin get list configuration in platform service")

	_, err := helpers.FilterColumnValidation(pages.Filters, models.AllowedFilterColumnConfiguration())
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("filter column validation failed")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

	configurations, err := s.service.repository.GetListConfigurationInPlatformRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("configurations.id," +
			"configurations.config_name," +
			"configurations.config_value," +
			"configurations.config_key," +
			"configurations.config_json," +
			"configurations.is_active," +
			"configurations.aggregator_id",
		)
		db = db.Joins("JOIN platform_configuration on platform_configuration.configuration_id = configurations.id")
		query, args := s.service.repository.SearchQuery(pages.Filters, pages.JoinOperator)
		if query != "" {
			db = db.Where(query, args...)
		}
		db = db.Where("platform_configuration.platform_id = ?", platformId)
		db = db.Preload("Aggregator", func(db *gorm.DB) *gorm.DB {
			db = db.Select("aggregators.id," +
				"aggregators.name," +
				"aggregators.currency",
			)
			return db
		})
		return db
	})
	if err != nil {
		log.Debug().Err(err).Interface("context", s.serviceName).Msg("get list platform in configuration repository error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data configurations", configurations).Interface("context", s.serviceName).Msg("result data get list configuration in platform repository")

	response := make([]*web.ResponseConfiguration, 0)

	for _, configuration := range configurations {
		response = append(response, s.service.mapToConfigurationResponse(configuration))
	}

	return web.ListResponse{Items: response, Metadata: pages.GetMetadata()}, nil
}

func (s *PlatformService) AdminAssignmentPlatformChannelService(session *pkgjwt.JwtResponse, platformId int64, payload []web.DetailChannelResponse) (any, error) {
	s.serviceName = "PlatformService.AdminAssignmentPlatformChannel"

	log.Debug().Interface("session", session).Interface("platform_id", platformId).Interface("payload", payload).
		Interface("context", s.serviceName).Msg("admin assignment platform channel service")

	platform, err := s.service.repository.GetPlatformRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("platforms.id," +
			"platforms.name",
		)
		db = db.Where("platforms.id = ?", platformId)
		return db
	})
	if err != nil {
		log.Debug().Err(err).Interface("context", s.serviceName).Msg("get platform by id repository")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(fmt.Errorf("platform, %v", err), s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data platform", platform).Interface("context", s.serviceName).Msg("result data get platform by id repository")

	platformChannel := make([]*models.PlatformChannel, 0)

	for idx1, ch := range payload {
		log.Debug().Interface("ch", ch).Interface("index", idx1).Interface("context", s.serviceName).
			Msg("iteration payload channel")

		channel, errCh := s.service.repository.GetChannelRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Select("channels.id," +
				"channels.name",
			)
			db = db.Where("channels.id = ?", ch.Id)
			return db
		})
		if errCh != nil {
			log.Error().Err(errCh).Str("context", s.serviceName).Msg("get configuration by id repository error")
			if errors.Is(errCh, gorm.ErrRecordNotFound) {
				return nil, helpers.NewErrorTrace(fmt.Errorf("channel %v, %v", ch.Name, errCh), s.serviceName).WithStatusCode(http.StatusNotFound)
			}

			return nil, helpers.NewErrorTrace(errCh, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("data channel", channel).Interface("context", s.serviceName).
			Msg("result data channel get channel by id repository")

		exists, errExist := s.service.repository.GetExistsPlatformChannelRepository(s.service.ctx, channel, platform)
		if errExist != nil {
			log.Debug().Err(err).Interface("context", s.serviceName).Msg("get exists platform channel repository error")
			return nil, helpers.NewErrorTrace(errExist, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("exists", exists).Interface("context", s.serviceName).
			Msg("result data exist get platform channel")

		if exists {
			return nil, helpers.NewErrorTrace(fmt.Errorf("the platform get exists channel %v, please check data %v", channel.Name, platform.Name), s.serviceName).
				WithStatusCode(http.StatusConflict)
		}

		platformChannel = append(platformChannel, &models.PlatformChannel{
			ChannelId:  channel.Id,
			PlatformId: platform.Id,
		})
	}
	log.Debug().Interface("data platform channel", platformChannel).Interface("context", s.serviceName).
		Msg("result data platform channel after iteration")

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
		platformChannel, err = s.service.repository.InsertBatchPlatformChannelRepositoryTx(s.service.ctx, tx, platformChannel)
		if err != nil {
			log.Debug().Err(err).Interface("context", s.serviceName).Msg("insert batch platform channel repository tx error")
			return err
		}

		adminActivityLog := models.NewAdminActivityLogs(
			session.Id,
			enums.ACTION_ADMIN_ACCTIVITY_CREATE,
			enums.RESOURCE_ADMIN_ACTIVITY_LOG_PLATFORM_CHANNEL,
			fmt.Sprint(),
			fmt.Sprint(enums.RESOURCE_ADMIN_ACTIVITY_LOG_PLATFORM_CHANNEL, enums.ACTION_ADMIN_ACCTIVITY_CREATE),
			enums.NULL_STRING,
			session.Sub,
			payload,
			platformChannel,
		)

		err = s.service.repository.InsertAdminActivityLogRepositoryTx(s.service.ctx, tx, adminActivityLog)
		if err != nil {
			log.Error().Err(err).Interface("context", s.serviceName).Msg("insert admin activity log repository tx")
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("assignment platform configuration transaction error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return platformChannel, nil
}

func (s *PlatformService) AdminAssignmentPlatformListChannelService(platformId int64) (map[enums.Currency]map[enums.PaymentMethod][]*web.DetailChannelResponse, error) {
	s.serviceName = "PlatformService.AdminAssignmentPlatformListChannel"

	log.Debug().Interface("platform_id", platformId).Interface("context", s.serviceName).Msg("admin assignment platform channel service")

	platform, err := s.service.repository.GetPlatformRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("platforms.id," +
			"platforms.name",
		)
		db = db.Where("platforms.id = ?", platformId)
		return db
	})
	if err != nil {
		log.Debug().Err(err).Interface("context", s.serviceName).Msg("get platform by id repository")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(fmt.Errorf("platform, %v", err), s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data platform", platform).Interface("context", s.serviceName).Msg("result data get platform by id repository")

	//configurations, err := s.service.repository.GetConfigurationByPlatformIdRepository(s.service.ctx, platformId)
	//if err != nil {
	//	log.Error().Err(err).Interface("context", s.serviceName).Msg("get platform configuration by platform id repository error")
	//	return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	//}
	//
	//if len(configurations) < 1 {
	//	return nil, helpers.NewErrorTrace(fmt.Errorf("you have not configuration for add channel payment"), s.serviceName).WithStatusCode(http.StatusNotFound)
	//}

	platformChannels, err := s.service.repository.GetPlatformChannelByPlatformIdRepository(s.service.ctx, platformId)
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("get platform channel by platform id repository")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	//aggregatorId := make([]int64, 0)
	channelId := make([]int64, 0)

	//for _, conf := range configurations {
	//	aggregatorId = append(aggregatorId, conf.AggregatorId)
	//}

	for _, pc := range platformChannels {
		channelId = append(channelId, pc.ChannelId)
	}

	channels, err := s.service.repository.GetChannelINByRepository(s.service.ctx, channelId)
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("get channel IN by aggregator id repository error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	response := s.service.mapToChannelResponseWithCurrency(channels)

	return response, nil
}

func (s *PlatformService) AdminRemovalPlatformChannelService(session *pkgjwt.JwtResponse, platformId int64, payload []web.DetailChannelResponse) (any, error) {
	s.serviceName = "PlatformService.AdminRemovalPlatformChannel"

	log.Debug().Interface("session", session).Interface("platform_id", platformId).Interface("payload", payload).
		Interface("context", s.serviceName).Msg("admin removal platform channel service")

	platform, err := s.service.repository.GetPlatformRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("platforms.id," +
			"platforms.name",
		)
		db = db.Where("platforms.id = ?", platformId)
		return db
	})
	if err != nil {
		log.Debug().Err(err).Interface("context", s.serviceName).Msg("get platform by id repository")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(fmt.Errorf("platform, %v", err), s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data platform", platform).Interface("context", s.serviceName).Msg("result data get platform by id repository")

	platformChannelId := make([]int64, 0)

	for idx1, ch := range payload {
		log.Debug().Interface("ch", ch).Interface("index", idx1).Interface("context", s.serviceName).
			Msg("iteration payload configuration")

		channel, errCh := s.service.repository.GetChannelRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Select("channels.id," +
				"channels.name",
			)
			db = db.Where("channels.id = ?", ch.Id)
			return db
		})
		if errCh != nil {
			log.Error().Err(errCh).Str("context", s.serviceName).Msg("get configuration by id repository error")
			if errors.Is(errCh, gorm.ErrRecordNotFound) {
				return nil, helpers.NewErrorTrace(fmt.Errorf("channel %v, %v", ch.Name, errCh), s.serviceName).WithStatusCode(http.StatusNotFound)
			}

			return nil, helpers.NewErrorTrace(errCh, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("data channel", channel).Interface("context", s.serviceName).
			Msg("result data channel get channel by id repository")

		platformChannel, errPh := s.service.repository.GetPlatformChannelRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Where("channel_id = ? and platform_id = ?", channel.Id, platform.Id)
			return db
		})
		if errPh != nil {
			log.Error().Err(errPh).Interface("context", s.serviceName).Msg("get platform channel by channel id and platform id repository error")
			if errors.Is(errPh, gorm.ErrRecordNotFound) {
				return nil, helpers.NewErrorTrace(fmt.Errorf("%v platfom has not been entered into channel %v", platform.Name, channel.Name), s.serviceName).
					WithStatusCode(http.StatusNotFound)
			}

			return nil, helpers.NewErrorTrace(errPh, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("data platform channel", platformChannel).Interface("context", s.serviceName).
			Msg("result data platform channel by channel id and platform id repository")

		platformChannelId = append(platformChannelId, platformChannel.Id)
	}
	log.Debug().Interface("data platform channel", platformChannelId).Interface("context", s.serviceName).
		Msg("result data platform channel after iteration")

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
		err = s.service.repository.DeleteINPlatformChannelRepositoryTx(s.service.ctx, tx, platformChannelId)
		if err != nil {
			log.Error().Err(err).Interface("context", s.serviceName).Msg("delete in platform channel repository tx error")
			return err
		}

		adminActivityLog := models.NewAdminActivityLogs(
			session.Id,
			enums.ACTION_ADMIN_ACCTIVITY_DELETE,
			enums.RESOURCE_ADMIN_ACTIVITY_LOG_PLATFORM_CHANNEL,
			fmt.Sprint(),
			fmt.Sprint(enums.RESOURCE_ADMIN_ACTIVITY_LOG_PLATFORM_CHANNEL, enums.ACTION_ADMIN_ACCTIVITY_DELETE),
			enums.NULL_STRING,
			session.Sub,
			payload,
			nil,
		)

		err = s.service.repository.InsertAdminActivityLogRepositoryTx(s.service.ctx, tx, adminActivityLog)
		if err != nil {
			log.Error().Err(err).Interface("context", s.serviceName).Msg("insert admin activity log repository tx")
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("removal platform channel transaction error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return nil, nil
}

func (s *PlatformService) AdminGetListChannelInPlatformService(platformId int64, pages *pagination.Pages) (map[enums.Currency]map[enums.PaymentMethod][]*web.DetailChannelResponse, error) {
	s.serviceName = "PlatformService.AdminGetListChannelInPlatform"

	log.Debug().Interface("platform_id", platformId).Interface("pages", pages).
		Interface("context", s.serviceName).Msg("admin get list channel in platform service")

	// check pages filter
	_, err := helpers.FilterColumnValidation(pages.Filters, models.AllowedFilterColumnChannel())
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("filter column validation")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

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
		query, args := s.service.repository.SearchQuery(pages.Filters, pages.JoinOperator)
		if query != "" {
			db = db.Where(query, args...)
		}
		db = db.Where("platform_channel.platform_id = ?", platformId)
		return db
	})
	if err != nil {
		log.Debug().Err(err).Interface("context", s.serviceName).Msg("get list platform in channel repository error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data channels", channels).Interface("context", s.serviceName).Msg("result data get list channel in platform repository")

	response := make([]*web.DetailChannelResponse, 0)

	responseData2 := map[enums.Currency]map[enums.PaymentMethod][]*web.DetailChannelResponse{
		enums.CURRENCY_MYR: map[enums.PaymentMethod][]*web.DetailChannelResponse{
			enums.PAYMENT_METHOD_MULTI_PAYMENT:   []*web.DetailChannelResponse{},
			enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT: []*web.DetailChannelResponse{},
		},
		enums.CURRENCY_IDR: map[enums.PaymentMethod][]*web.DetailChannelResponse{
			enums.PAYMENT_METHOD_MULTI_PAYMENT:   []*web.DetailChannelResponse{},
			enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT: []*web.DetailChannelResponse{},
		},
	}

	for _, ch := range channels {
		channel := s.service.mapToDetailChannelResponse(ch)

		response = append(response, s.service.mapToDetailChannelResponse(ch))

		if _, ok := responseData2[ch.Currency]; ok {
			//responseData2[ch.Currency] = map[enums.PaymentMethod][]*web.DetailChannelResponse{
			//	enums.PAYMENT_METHOD_MULTI_PAYMENT:   []*web.DetailChannelResponse{},
			//	enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT: []*web.DetailChannelResponse{},
			//}

			switch ch.PaymentMethod {
			case enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT:
				responseData2[ch.Currency][enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT] = append(responseData2[ch.Currency][enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT], channel)
			default:
				responseData2[ch.Currency][enums.PAYMENT_METHOD_MULTI_PAYMENT] = append(responseData2[ch.Currency][enums.PAYMENT_METHOD_MULTI_PAYMENT], channel)
			}

		}
	}

	return responseData2, nil
}
