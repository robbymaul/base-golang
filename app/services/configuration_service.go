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
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type ConfigurationService struct {
	service     *Service
	serviceName string
}

func NewConfigurationService(ctx context.Context, repo *repositories.RepositoryContext, cfg *config.Config) *ConfigurationService {
	return &ConfigurationService{
		service:     NewService(ctx, repo, cfg),
		serviceName: "platform configuration",
	}
}

func (s *ConfigurationService) AdminCreateConfigurationService(session *pkgjwt.JwtResponse, payload []web.CreateConfigurationRequest) (any, error) {
	log.Debug().Interface("session", session).Interface("payload", payload).Msg("create platform configuration service")
	var err error

	configurations := make([]*models.Configuration, 0, len(payload))

	for idx1, conf := range payload {
		log.Debug().Interface("index", idx1).Interface("conf", conf).
			Msg(fmt.Sprintf("iteration index %v paylod create configurtion request", idx1))

		aggregator, errAg := s.service.repository.GetAggregatorByIdRepository(s.service.ctx, conf.Aggregator.Id)
		if errAg != nil {
			log.Error().Err(errAg).Interface("index", idx1).Interface("conf data", conf).
				Str("context", s.serviceName).Msg("get aggregator by id repository error")

			if errors.Is(errAg, gorm.ErrRecordNotFound) {
				return nil, helpers.NewErrorTrace(fmt.Errorf("aggregator, %v on index %v", errAg, idx1), s.serviceName).WithStatusCode(http.StatusNotFound)
			}

			return nil, helpers.NewErrorTrace(errAg, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("index", idx1).Interface("aggregator", aggregator).Msg("get aggregator by id repository")

		configJson := models.ConfigJson{
			SandboxBaseUrl:               conf.ConfigJson.SandboxBaseUrl,
			ProductionBaseUrl:            conf.ConfigJson.ProductionBaseUrl,
			SandboxMerchantId:            helpers.EncryptAES(conf.ConfigJson.SandboxMerchantId),
			ProductionMerchantId:         helpers.EncryptAES(conf.ConfigJson.ProductionMerchantId),
			SandboxMerchantCode:          helpers.EncryptAES(conf.ConfigJson.SandboxMerchantCode),
			ProductionMerchantCode:       helpers.EncryptAES(conf.ConfigJson.ProductionMerchantCode),
			SandboxMerchantName:          helpers.EncryptAES(conf.ConfigJson.SandboxMerchantName),
			ProductionMerchantName:       helpers.EncryptAES(conf.ConfigJson.ProductionMerchantName),
			SandboxApiKey:                helpers.EncryptAES(conf.ConfigJson.SandboxApiKey),
			ProductionApiKey:             helpers.EncryptAES(conf.ConfigJson.ProductionApiKey),
			SandboxServerKey:             helpers.EncryptAES(conf.ConfigJson.SandboxServerKey),
			ProductionServerKey:          helpers.EncryptAES(conf.ConfigJson.ProductionServerKey),
			SandboxSecretKey:             helpers.EncryptAES(conf.ConfigJson.SandboxSecretKey),
			ProductionSecretKey:          helpers.EncryptAES(conf.ConfigJson.ProductionSecretKey),
			SandboxClientKey:             helpers.EncryptAES(conf.ConfigJson.SandboxClientKey),
			ProductionClientKey:          helpers.EncryptAES(conf.ConfigJson.ProductionClientKey),
			SandboxSignatureKey:          helpers.EncryptAES(conf.ConfigJson.SandboxSignatureKey),
			ProductionSignatureKey:       helpers.EncryptAES(conf.ConfigJson.ProductionSignatureKey),
			SandboxCredentialPassword:    helpers.EncryptAES(conf.ConfigJson.SandboxCredentialPassword),
			ProductionCredentialPassword: helpers.EncryptAES(conf.ConfigJson.ProductionCredentialPassword),
			ReturnUrl:                    conf.ConfigJson.ReturnUrl,
		}
		log.Debug().Interface("index", idx1).Interface("configJson", configJson).Interface("context", s.serviceName).
			Msg("config json")

		configuration := &models.Configuration{
			AggregatorId: aggregator.Id,
			ConfigKey:    uuid.New().String(),
			ConfigValue:  conf.ConfigValue,
			ConfigName:   conf.ConfigName,
			IsActive:     true,
			ConfigJson:   configJson,
		}
		log.Debug().Interface("index", idx1).Interface("configuration", configuration).Interface("context", s.serviceName).
			Msg("platform configuration")

		configurations = append(configurations, configuration)
	}

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
		configurations, err = s.service.repository.InsertBatchConfigurationRepositoryTx(s.service.ctx, tx, configurations)
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("create platform configuration repository error")

			return err
		}

		adminActivityLog := models.NewAdminActivityLogs(
			session.Id,
			enums.ACTION_ADMIN_ACCTIVITY_CREATE,
			enums.RESOURCE_ADMIN_ACTIVITY_LOG_CONFIGURATION,
			fmt.Sprint(),
			fmt.Sprint(enums.RESOURCE_ADMIN_ACTIVITY_LOG_CONFIGURATION, enums.ACTION_ADMIN_ACCTIVITY_CREATE),
			enums.NULL_STRING,
			session.Sub,
			payload,
			configurations,
		)

		err = s.service.repository.InsertAdminActivityLogRepositoryTx(s.service.ctx, tx, adminActivityLog)
		if err != nil {
			log.Error().Err(err).Interface("context", s.serviceName).Msg("insert admin activity log repository tx")
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("create configuration with transaction error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	response := make([]*web.ResponseConfiguration, 0, len(configurations))

	for _, conf := range configurations {
		response = append(response, s.mapToConfigurationResponse(conf))
	}

	return response, nil
}

func (s *ConfigurationService) AdminUpdateConfigurationService(session *pkgjwt.JwtResponse, configurationId int64, payload *web.ResponseConfiguration) (any, error) {
	log.Debug().Interface("session", session).Interface("configuration_id", configurationId).Interface("payload", payload).
		Msg("admin update configuration service request")

	configuration, err := s.service.repository.GetConfigurationRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("configurations.id," +
			"configurations.config_key," +
			"configurations.config_value," +
			"configurations.config_name," +
			"configurations.is_active," +
			"configurations.config_json," +
			"configurations.aggregator_id")
		db = db.Where("configurations.id=?", configurationId)
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get platform configuration by id repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data configuration", configuration).Interface("context", s.serviceName).
		Msg("result data configuration get configuration by id")

	s.updateConfiguration(&configuration, payload)
	log.Debug().Interface("data update configuration", configuration).Interface("context", s.serviceName).
		Msg("result data update configuration")

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
		err = s.service.repository.UpdateConfigurationRepositoryTx(s.service.ctx, tx, func(db *gorm.DB) *gorm.DB {
			db = db.Where("configurations.id = ?", configuration.Id)
			updateColumn := map[string]any{
				"config_key":    configuration.ConfigKey,
				"config_value":  configuration.ConfigValue,
				"config_name":   configuration.ConfigName,
				"config_json":   configuration.ConfigJson,
				"is_active":     configuration.IsActive,
				"aggregator_id": configuration.AggregatorId,
				"updated_at":    time.Now(),
			}
			db = db.UpdateColumns(&updateColumn)
			return db
		})
		if err != nil {
			log.Error().Err(err).Str("context", s.serviceName).Msg("update platform configuration repository tx error")
			return err
		}

		adminActivityLog := models.NewAdminActivityLogs(
			session.Id,
			enums.ACTION_ADMIN_ACCTIVITY_UPDATE,
			enums.RESOURCE_ADMIN_ACTIVITY_LOG_CONFIGURATION,
			fmt.Sprint(),
			fmt.Sprint(enums.RESOURCE_ADMIN_ACTIVITY_LOG_CONFIGURATION, enums.ACTION_ADMIN_ACCTIVITY_UPDATE),
			enums.NULL_STRING,
			session.Sub,
			payload,
			configuration,
		)

		err = s.service.repository.InsertAdminActivityLogRepositoryTx(s.service.ctx, tx, adminActivityLog)
		if err != nil {
			log.Error().Err(err).Interface("context", s.serviceName).Msg("insert admin activity log repository tx")
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("update configuration with transaction error ")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return s.mapToConfigurationResponse(configuration), nil
}

func (s *ConfigurationService) mapToConfigurationResponse(conf *models.Configuration) *web.ResponseConfiguration {
	var aggregator *web.AggregatorResponse

	if conf.Aggregator != nil {
		aggregator = s.service.mapToAggregatorResponse(conf.Aggregator)
	}

	return &web.ResponseConfiguration{
		Id:           conf.Id,
		AggregatorId: conf.AggregatorId,
		Aggregator:   aggregator,
		ConfigName:   conf.ConfigName,
		ConfigValue:  conf.ConfigValue,
		IsActive:     conf.IsActive,
		ConfigJson:   s.service.mapToConfigJsonResponse(conf.ConfigJson),
		CreatedAt:    conf.CreatedAt,
		UpdatedAt:    conf.UpdatedAt,
	}
}

//
//func (s *ConfigurationService) mapToConfigJsonResponse(json models.ConfigJson) *web.ConfigJson {
//	return &web.ConfigJson{
//		SandboxBaseUrl:               json.SandboxBaseUrl,
//		ProductionBaseUrl:            json.ProductionBaseUrl,
//		SandboxMerchantId:            helpers.DecryptAES(json.SandboxMerchantId),
//		ProductionMerchantId:         helpers.DecryptAES(json.ProductionMerchantId), // json.ProductionMerchantId,
//		SandboxMerchantCode:          helpers.DecryptAES(json.SandboxMerchantCode),
//		ProductionMerchantCode:       helpers.DecryptAES(json.ProductionMerchantCode),
//		SandboxMerchantName:          helpers.DecryptAES(json.SandboxMerchantName),
//		ProductionMerchantName:       helpers.DecryptAES(json.ProductionMerchantName),
//		SandboxApiKey:                helpers.DecryptAES(json.SandboxApiKey),
//		ProductionApiKey:             helpers.DecryptAES(json.ProductionApiKey),
//		SandboxServerKey:             helpers.DecryptAES(json.SandboxServerKey),       // json.SandboxServerKey,
//		ProductionServerKey:          helpers.DecryptAES(json.ProductionServerKey),    // json.ProductionServerKey,
//		SandboxSecretKey:             helpers.DecryptAES(json.SandboxSecretKey),       // json.SandboxSecretKey,
//		ProductionSecretKey:          helpers.DecryptAES(json.ProductionSecretKey),    // json.ProductionSecretKey,
//		SandboxClientKey:             helpers.DecryptAES(json.SandboxClientKey),       // json.SandboxClientKey,
//		ProductionClientKey:          helpers.DecryptAES(json.ProductionClientKey),    // json.ProductionClientKey,
//		SandboxSignatureKey:          helpers.DecryptAES(json.SandboxSignatureKey),    // json.SandboxSignatureKey,
//		ProductionSignatureKey:       helpers.DecryptAES(json.ProductionSignatureKey), // json.ProductionSignatureKey,
//		SandboxCredentialPassword:    helpers.DecryptAES(json.SandboxCredentialPassword),
//		ProductionCredentialPassword: helpers.DecryptAES(json.ProductionCredentialPassword),
//		ReturnUrl:                    json.ReturnUrl,
//	}
//}

func (s *ConfigurationService) updateConfiguration(config **models.Configuration, payload *web.ResponseConfiguration) {
	if (*config).ConfigValue != payload.ConfigValue {
		(*config).ConfigValue = payload.ConfigValue
	}

	if (*config).ConfigName != payload.ConfigName {
		(*config).ConfigName = payload.ConfigName
	}

	if (*config).IsActive != payload.IsActive {
		(*config).IsActive = payload.IsActive
	}

	(*config).ConfigJson = models.ConfigJson{
		SandboxBaseUrl:         payload.ConfigJson.SandboxBaseUrl,
		ProductionBaseUrl:      payload.ConfigJson.ProductionBaseUrl,
		SandboxMerchantId:      helpers.EncryptAES(payload.ConfigJson.SandboxMerchantId),
		ProductionMerchantId:   helpers.EncryptAES(payload.ConfigJson.ProductionMerchantId),
		SandboxServerKey:       helpers.EncryptAES(payload.ConfigJson.SandboxServerKey),
		ProductionServerKey:    helpers.EncryptAES(payload.ConfigJson.ProductionServerKey),
		SandboxSecretKey:       helpers.EncryptAES(payload.ConfigJson.SandboxSecretKey),
		ProductionSecretKey:    helpers.EncryptAES(payload.ConfigJson.ProductionSecretKey),
		SandboxClientKey:       helpers.EncryptAES(payload.ConfigJson.SandboxClientKey),
		ProductionClientKey:    helpers.EncryptAES(payload.ConfigJson.ProductionClientKey),
		SandboxSignatureKey:    helpers.EncryptAES(payload.ConfigJson.SandboxSignatureKey),
		ProductionSignatureKey: helpers.EncryptAES(payload.ConfigJson.ProductionSignatureKey),
		ReturnUrl:              payload.ConfigJson.ReturnUrl,
	}
	if (*config).AggregatorId != payload.Aggregator.Id {
		(*config).AggregatorId = payload.Aggregator.Id
	}
}

func (s *ConfigurationService) AdminGetListConfigurationService(pages *pagination.Pages, aggregator string) (any, error) {
	var configurations []*models.Configuration
	var totalCount int64
	var err error

	log.Debug().Interface("pages", pages).Interface("context", s.serviceName).
		Msg("admin get list configuration service")

	filter, err := helpers.FilterColumnValidation(pages.Filters, models.AllowedFilterColumnConfiguration())
	if err != nil {
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

	ctx, cancelFn := context.WithCancel(s.service.ctx)
	eg := errgroup.Group{}

	eg.Go(func() error {
		totalCount, err = s.service.repository.GetCountConfigurationRepository(ctx, func(db *gorm.DB) *gorm.DB {
			query, args := s.service.repository.SearchQuery(filter, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			if aggregator != "" {
				db = db.Joins("JOIN aggregators on aggregators.id = configurations.aggregator_id")
				db = db.Where("aggregators.name = ?", aggregator)
			}
			return db
		})
		if err != nil {
			cancelFn()
			log.Error().Err(err).Str("context", s.serviceName).Msg("get count configuration repository error")
			return err
		}
		log.Debug().Int64("count", totalCount).Interface("context", s.serviceName).Msg("result data count configuration")

		return nil
	})

	eg.Go(func() error {
		configurations, err = s.service.repository.FindConfigurationRepository(ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Select("configurations.id," +
				"configurations.config_key," +
				"configurations.config_value," +
				"configurations.config_name," +
				"configurations.config_json," +
				"configurations.is_active," +
				"configurations.aggregator_id",
			)
			db = db.Limit(pages.Limit()).Offset(pages.Offset())
			query, args := s.service.repository.SearchQuery(filter, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			if aggregator != "" {
				db = db.Joins("JOIN aggregators on aggregators.id = configurations.aggregator_id")
				db = db.Where("aggregators.name = ?", aggregator)
			}
			db = db.Order("id " + pages.Sort)
			db = db.Preload("Aggregator", func(db *gorm.DB) *gorm.DB {
				db.Select("aggregators.id," +
					"aggregators.name," +
					"aggregators.is_active",
				)
				return db
			})
			return db
		})
		if err != nil {
			cancelFn()
			log.Error().Err(err).Str("context", s.serviceName).Msg("get list configuration repository error")
			return err
		}
		log.Debug().Interface("data configuration", configurations).Interface("context", s.serviceName).Msg("result data get list configuration")

		return nil
	})

	err = eg.Wait()
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("err sync group error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	pages.TotalCount = int(totalCount)

	configsResponse := make([]*web.ResponseConfiguration, 0)
	for _, platformConfig := range configurations {
		configsResponse = append(configsResponse, s.mapToConfigurationResponse(platformConfig))
	}

	return &web.ListResponse{Items: configsResponse, Metadata: pages.GetMetadata()}, nil
}

func (s *ConfigurationService) AdminGetConfigurationService(configurationId int64) (*web.ResponseConfiguration, error) {
	log.Debug().Int64("configuration_id", configurationId).Interface("context", s.serviceName).
		Msg("admin get platform configuration service")

	configuration, err := s.service.repository.GetConfigurationRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("configurations.id," +
			"configurations.config_key," +
			"configurations.config_value," +
			"configurations.config_name," +
			"configurations.config_json," +
			"configurations.is_active," +
			"configurations.aggregator_id",
		)
		db = db.Where("configurations.id=?", configurationId)
		db = db.Preload("Aggregator", func(db *gorm.DB) *gorm.DB {
			db = db.Select("aggregators.id," +
				"aggregators.name," +
				"aggregators.is_active",
			)
			return db
		})
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get configuration by id repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data configuration", configuration).Interface("context", s.serviceName).
		Msg("result data configuration get configuration by id repository")

	return s.mapToConfigurationResponse(configuration), nil
}

func (s *ConfigurationService) AdminAssignmentConfigurationToPlatformService(session *pkgjwt.JwtResponse, configurationId int64, payload []web.DetailPlatformResponse) (any, error) {
	s.serviceName = "ConfigurationService.AdminAssignmentConfigurationPlatform"

	log.Debug().Interface("session", session).Interface("configuration_id", configurationId).Interface("payload", payload).
		Interface("context", s.serviceName).Msg("admin assignment configuration platform")

	// get configuration repository
	configuration, err := s.service.repository.GetConfigurationRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("configurations.id," +
			"configurations.config_name," +
			"configurations.aggregator_id",
		)
		db = db.Where("configurations.id = ?", configurationId)
		db = db.Preload("Aggregator", func(db *gorm.DB) *gorm.DB {
			db = db.Select("aggregators.id," +
				"aggregators.name",
			)
			return db
		})
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get configuration by id repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data configuration", configuration).Interface("context", s.serviceName).
		Msg("result data configuration get configuration by id repository")

	platformConfiguration := make([]*models.PlatformConfiguration, 0)

	for idx1, pf := range payload {
		log.Debug().Interface("pf", pf).Interface("index", idx1).Interface("context", s.serviceName).
			Msg("iteration payload platform")

		// get platform configuration
		platform, errPf := s.service.repository.GetPlatformRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Select("platforms.id," +
				"platforms.name",
			)
			db = db.Where("platforms.id = ?", pf.Id)
			return db
		})
		if errPf != nil {
			log.Debug().Err(errPf).Interface("context", s.serviceName).Msg("get platform by id repository")
			if errors.Is(errPf, gorm.ErrRecordNotFound) {
				return nil, helpers.NewErrorTrace(fmt.Errorf("platform %v, %v", pf.Name, errPf), s.serviceName).WithStatusCode(http.StatusNotFound)
			}

			return nil, helpers.NewErrorTrace(errPf, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("data platform", platform).Interface("context", s.serviceName).Msg("result data get platform by id repository")

		// check if exists platform configuration aggregator
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

	return nil, nil
}

func (s *ConfigurationService) AdminRemovalPlatformFromConfigurationService(session *pkgjwt.JwtResponse, configurationId int64, payload []web.DetailPlatformResponse) (any, error) {
	s.serviceName = "ConfigurationService.AdminRemovalPlatformFromConfiguration"

	configuration, err := s.service.repository.GetConfigurationRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("configurations.id," +
			"configurations.config_name," +
			"configurations.aggregator_id",
		)
		db = db.Where("configurations.id = ?", configurationId)
		db = db.Preload("Aggregator", func(db *gorm.DB) *gorm.DB {
			db = db.Select("aggregators.id," +
				"aggregators.name")
			return db
		})
		return db
	})
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("get configuration by id repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data configuration", configuration).Interface("context", s.serviceName).
		Msg("result data configuration get configuration by id repository")

	platformConfigurationId := make([]int64, 0)

	for idx1, pf := range payload {
		log.Debug().Interface("platform", pf).Interface("index", idx1).Interface("context", s.serviceName).
			Msg("iteration payload platform")

		platform, errPf := s.service.repository.GetPlatformRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Select("platforms.id," +
				"platforms.name",
			)
			db = db.Where("platforms.id = ?", pf.Id)
			return db
		})
		if errPf != nil {
			log.Debug().Err(errPf).Interface("context", s.serviceName).Msg("get platform by id repository")
			if errors.Is(errPf, gorm.ErrRecordNotFound) {
				return nil, helpers.NewErrorTrace(fmt.Errorf("platform %v, %v", pf.Name, errPf), s.serviceName).WithStatusCode(http.StatusNotFound)
			}

			return nil, helpers.NewErrorTrace(errPf, s.serviceName).WithStatusCode(http.StatusInternalServerError)
		}
		log.Debug().Interface("data platform", platform).Interface("context", s.serviceName).Msg("result data get platform by id repository")

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

func (s *ConfigurationService) AdminGetListPlatformInConfigurationService(configurationId int64, pages *pagination.Pages) (any, error) {
	s.serviceName = "ConfigurationService.AdminGetListPlatformInConfiguration"

	log.Debug().Interface("configuration_id", configurationId).Interface("pages", pages).
		Interface("context", s.serviceName).Msg("admin get list platform in configuration service")

	filter, err := helpers.FilterColumnValidation(pages.Filters, models.AllowedFilterColumnPlatform())
	if err != nil {
		log.Error().Err(err).Msg("filter column validation error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

	platforms, err := s.service.repository.FindPlatformInConfigurationRepository(s.service.ctx, configurationId, pages, func(db *gorm.DB) *gorm.DB {
		db = db.Select("platforms.id," +
			"platforms.code," +
			"platforms.name," +
			"platforms.description," +
			"platforms.is_active",
		)
		db = db.Joins("JOIN platform_configuration on platform_configuration.platform_id = platforms.id")
		query, args := s.service.repository.SearchQuery(filter, pages.JoinOperator)
		if query != "" {
			db = db.Where(query, args...)
		}
		db = db.Where("platform_configuration.configuration_id = ?", configurationId)
		return db
	})
	if err != nil {
		log.Debug().Err(err).Interface("context", s.serviceName).Msg("get list platform in configuration repository error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data platforms", platforms).Interface("context", s.serviceName).Msg("result data get list platform in configuration repository")

	response := make([]*web.DetailPlatformResponse, 0)

	for _, platform := range platforms {
		response = append(response, s.service.mapToDetailPlatformResponse(platform))
	}

	return web.ListResponse{Items: response, Metadata: pages.GetMetadata()}, nil
}
