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
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type AggregatorService struct {
	service     *Service
	serviceName string
}

func NewAggregatorService(ctx context.Context, repo *repositories.RepositoryContext, cfg *config.Config) *AggregatorService {
	return &AggregatorService{
		service:     NewService(ctx, repo, cfg),
		serviceName: "AggregatorService",
	}
}

func (s *AggregatorService) AdminCreateAggregatorService(session *pkgjwt.JwtResponse, payload *web.CreateAggregatorRequest) (any, error) {
	log.Debug().Interface("payload", payload).Interface("context", s.serviceName).Msg("create aggregator service")
	exist, err := s.service.repository.GetExistAggregatorByNameRepository(s.service.ctx, payload.Name)
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("error get exist aggregator by name")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("aggregator exists", exist).Interface("context", s.serviceName).Msg("get exist aggregator by name repository")

	if exist {
		return nil, helpers.NewErrorTrace(fmt.Errorf("aggregator already exist %v", payload.Name), s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	slug := strings.ToLower(string(payload.Name))

	aggregator := &models.Aggregator{
		Name:        payload.Name,
		Slug:        enums.ProviderPaymentMethod(slug),
		Description: payload.Description,
		IsActive:    true,
	}

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
		aggregator, err = s.service.repository.InsertAggregatorRepositoryTx(s.service.ctx, tx, aggregator)
		if err != nil {
			log.Error().Err(err).Interface("context", s.serviceName).Msg("insert aggregator repository errir")
			return err
		}
		log.Debug().Interface("aggregator", aggregator).Interface("context", s.serviceName).Msg("aggregator result value insert aggregator repository")

		var wg sync.WaitGroup
		var mu sync.Mutex

		//if len(payload.Channel) > 0 {
		//	channel := make([]*models.Channel, 0)
		//
		//	for i, ch := range payload.Channel {
		//		log.Debug().Interface("index", i).Interface("channel data", ch).Interface("context", s.serviceName).Msg("iteration channel proses append")
		//		wg.Add(1)
		//
		//		go func(ch *web.CreateChannelRequest) {
		//			defer wg.Done()
		//
		//			mu.Lock()
		//			channel = append(channel, &models.Channel{
		//				//AggregatorId:    aggregator.Id,
		//				//Code:            ch.Code,
		//				Name:            ch.Name,
		//				PaymentMethod:   ch.PaymentMethod,
		//				TransactionType: ch.TransactionType,
		//				Provider:        aggregator.Slug,
		//				Currency:        ch.Currency,
		//				FeeType:         ch.FeeType,
		//				FeeAmount:       ch.FeeAmount,
		//				IsActive:        true,
		//				//IsEspay:         ch.IsEspay,
		//				ProductName: ch.ProductName,
		//				ProductCode: ch.ProductCode,
		//				BankName:    ch.BankName,
		//				BankCode:    ch.BankCode,
		//			})
		//			mu.Unlock()
		//
		//		}(&ch)
		//	}
		//
		//	wg.Wait()
		//
		//	_, err = s.service.repository.InsertBatchChannelRepositoryTx(s.service.ctx, tx, channel)
		//	if err != nil {
		//		log.Debug().Err(err).Interface("context", s.serviceName).Interface("context", s.serviceName).Msg("insert batch channel repository tx")
		//		return err
		//	}
		//}

		if len(payload.Configuration) > 0 {
			configurations := make([]*models.Configuration, 0, len(payload.Configuration))

			for i, conf := range payload.Configuration {
				log.Debug().Interface("index", i).Interface("config data", conf).
					Interface("context", s.serviceName).Msg("iteration configuration proses append")

				wg.Add(1)

				go func(conf *web.CreateConfigurationRequest) {
					defer wg.Done()
					configKey := uuid.New()

					mu.Lock()
					configurations = append(configurations, &models.Configuration{
						AggregatorId: aggregator.Id,
						ConfigKey:    configKey.String(),
						ConfigValue:  conf.ConfigValue,
						ConfigName:   conf.ConfigName,
						ConfigJson: models.ConfigJson{
							SandboxBaseUrl:               conf.ConfigJson.SandboxBaseUrl,
							ProductionBaseUrl:            conf.ConfigJson.ProductionBaseUrl,
							SandboxMerchantId:            helpers.DecryptAES(conf.ConfigJson.SandboxMerchantId),
							ProductionMerchantId:         helpers.DecryptAES(conf.ConfigJson.ProductionMerchantId),
							SandboxMerchantCode:          helpers.DecryptAES(conf.ConfigJson.SandboxMerchantCode),
							ProductionMerchantCode:       helpers.DecryptAES(conf.ConfigJson.ProductionMerchantCode),
							SandboxMerchantName:          helpers.DecryptAES(conf.ConfigJson.SandboxMerchantName),
							ProductionMerchantName:       helpers.DecryptAES(conf.ConfigJson.ProductionMerchantName),
							SandboxApiKey:                helpers.DecryptAES(conf.ConfigJson.SandboxApiKey),
							ProductionApiKey:             helpers.DecryptAES(conf.ConfigJson.ProductionApiKey),
							SandboxServerKey:             helpers.DecryptAES(conf.ConfigJson.SandboxServerKey),
							ProductionServerKey:          helpers.DecryptAES(conf.ConfigJson.ProductionServerKey),
							SandboxSecretKey:             helpers.DecryptAES(conf.ConfigJson.SandboxSecretKey),
							ProductionSecretKey:          helpers.DecryptAES(conf.ConfigJson.ProductionSecretKey),
							SandboxClientKey:             helpers.DecryptAES(conf.ConfigJson.SandboxClientKey),
							ProductionClientKey:          helpers.DecryptAES(conf.ConfigJson.ProductionClientKey),
							SandboxSignatureKey:          helpers.DecryptAES(conf.ConfigJson.SandboxSignatureKey),
							ProductionSignatureKey:       helpers.DecryptAES(conf.ConfigJson.ProductionSignatureKey),
							SandboxCredentialPassword:    helpers.DecryptAES(conf.ConfigJson.SandboxCredentialPassword),
							ProductionCredentialPassword: helpers.DecryptAES(conf.ConfigJson.ProductionCredentialPassword),
							ReturnUrl:                    "",
						},
						IsActive: true,
					})
				}(&conf)
				mu.Unlock()
			}

			wg.Wait()

			_, err = s.service.repository.InsertBatchConfigurationRepositoryTx(s.service.ctx, tx, configurations)
			if err != nil {
				log.Debug().Err(err).Interface("context", s.serviceName).Interface("context", s.serviceName).
					Msg("insert batch configuration repository")
				return err
			}
		}

		adminActivityLog := models.NewAdminActivityLogs(
			session.Id,
			enums.ACTION_ADMIN_ACCTIVITY_CREATE,
			enums.RESOURCE_ADMIN_ACTIVITY_LOG_AGGREGATOR,
			fmt.Sprint(aggregator.Id),
			fmt.Sprint(enums.RESOURCE_ADMIN_ACTIVITY_LOG_AGGREGATOR, enums.ACTION_ADMIN_ACCTIVITY_CREATE),
			enums.NULL_STRING,
			session.Sub,
			payload,
			aggregator,
		)

		err = s.service.repository.InsertAdminActivityLogRepositoryTx(s.service.ctx, tx, adminActivityLog)
		if err != nil {
			log.Error().Err(err).Interface("context", s.serviceName).Msg("insert admin activity log repository tx")
		}

		return nil
	})
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("repository with transaction error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return s.service.mapToAggregatorResponse(aggregator), nil
}

func (s *AggregatorService) AdminGetAggregatorService(id int64) (any, error) {
	log.Debug().Interface("id", id).Interface("context", s.serviceName).Msg("admin get aggregator service")

	aggregator, err := s.service.repository.GetAggregatorRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("aggregators.id," +
			"aggregators.name," +
			"aggregators.slug," +
			"aggregators.description," +
			"aggregators.is_active," +
			"aggregators.currency")
		db = db.Where("aggregators.id = ?", id)
		return db
	})
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("error get aggregator repository")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("aggregator data", aggregator).Interface("context", s.serviceName).Msg("result data get aggregator by id")

	return s.service.mapToAggregatorResponse(aggregator), nil
}

func (s *AggregatorService) AdminGetAllAggregatorService(pages *pagination.Pages) (any, error) {
	var aggregators []*models.Aggregator
	var totalCount int64
	var err error

	log.Debug().Interface("pages", pages).Interface("context", s.serviceName).
		Msg("admin get all aggregator service")

	filter, err := helpers.FilterColumnValidation(pages.Filters, models.AllowedFilterColumnAggregator())
	if err != nil {
		log.Error().Err(err).Str("context", s.serviceName).Msg("filter column validation error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

	ctx, cancelFn := context.WithCancel(s.service.ctx)
	eg := errgroup.Group{}

	eg.Go(func() error {
		aggregators, err = s.service.repository.FindAggregatorRepository(ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Select("aggregators.id," +
				"aggregators.name," +
				"aggregators.slug," +
				"aggregators.description," +
				"aggregators.is_active," +
				"aggregators.currency",
			)
			query, args := s.service.repository.SearchQuery(filter, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			db = db.Limit(pages.Limit()).Offset(pages.Offset())
			db = db.Order(fmt.Sprintf("aggregators.id %s", pages.Sort))
			return db
		})
		if err != nil {
			cancelFn()
			log.Error().Err(err).Interface("context", s.serviceName).Msg("error get all aggregator repository")
			return err
		}
		log.Debug().Interface("aggregator list data", aggregators).Interface("context", s.serviceName).Msg("get all aggregator repository")

		return nil
	})

	eg.Go(func() error {
		totalCount, err = s.service.repository.GetCountAggregatorRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
			query, args := s.service.repository.SearchQuery(filter, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			return db
		})
		if err != nil {
			cancelFn()
			log.Error().Err(err).Interface("context", s.serviceName).Msg("error get count aggregator repository")
			return err
		}
		log.Debug().Interface("count data", totalCount).Interface("context", s.serviceName).Msg("get count aggregator repository")

		return nil
	})

	err = eg.Wait()
	if err != nil {
		log.Error().Err(err).Msg("async group error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	pages.TotalCount = int(totalCount)

	response := make([]*web.AggregatorResponse, 0)
	for _, aggregator := range aggregators {
		response = append(response, s.service.mapToAggregatorResponse(aggregator))
	}

	return web.ListResponse{Items: response, Metadata: pages.GetMetadata()}, nil
}

func (s *AggregatorService) AdminUpdateAggregatorService(session *pkgjwt.JwtResponse, id int64, payload *web.AggregatorResponse) (any, error) {
	aggregator, err := s.service.repository.GetAggregatorRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("aggregators.id," +
			"aggregators.name," +
			"aggregators.description," +
			"aggregators.is_active," +
			"aggregators.currency",
		)
		db = db.Where("aggregators.id = ?", id)
		return db
	})
	if err != nil {
		log.Error().Err(err).Msg("error get aggregator repository")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("aggregator data", aggregator).Interface("context", s.serviceName).Msg("result data get aggregator data by id")

	s.updateAggregator(&aggregator, payload)
	log.Debug().Interface("aggregator update", aggregator).Interface("context", s.serviceName).Msg("result data update aggregator")

	err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
		err = s.service.repository.UpdateAggregatorRepositoryTx(s.service.ctx, tx, func(db *gorm.DB) *gorm.DB {
			db = db.Where("aggregators.id = ?", aggregator.Id)
			updateColumn := map[string]interface{}{
				"name":        aggregator.Name,
				"description": aggregator.Description,
				"is_active":   aggregator.IsActive,
				"currency":    aggregator.Currency,
				"updated_at":  time.Now(),
			}
			db = db.UpdateColumns(updateColumn)
			return db
		})
		if err != nil {
			log.Error().Err(err).Msg("error update aggregator repository")
			return err
		}
		log.Debug().Interface("aggregator update result", aggregator).Interface("context", s.serviceName).Msg("result data update aggregator repository")

		adminActivityLog := models.NewAdminActivityLogs(
			session.Id,
			enums.ACTION_ADMIN_ACCTIVITY_UPDATE,
			enums.RESOURCE_ADMIN_ACTIVITY_LOG_AGGREGATOR,
			fmt.Sprint(aggregator.Id),
			fmt.Sprint(enums.RESOURCE_ADMIN_ACTIVITY_LOG_AGGREGATOR, enums.ACTION_ADMIN_ACCTIVITY_UPDATE),
			enums.NULL_STRING,
			session.Sub,
			payload,
			aggregator,
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

	return s.service.mapToAggregatorResponse(aggregator), nil
}

//func (s *AggregatorService) mapToAggregatorResponse(aggregator *models.Aggregator) *web.AggregatorResponse {
//	return &web.AggregatorResponse{
//		Id:          aggregator.Id,
//		Name:        aggregator.Name,
//		Slug:        aggregator.Slug,
//		Description: aggregator.Description,
//		IsActive:    aggregator.IsActive,
//		Currency:    aggregator.Currency,
//		CreatedAt:   aggregator.CreatedAt,
//		UpdatedAt:   aggregator.UpdatedAt,
//	}
//}

func (s *AggregatorService) updateAggregator(aggregator **models.Aggregator, payload *web.AggregatorResponse) {
	if (*aggregator).Name != payload.Name {
		(*aggregator).Name = payload.Name
	}

	if (*aggregator).Description != payload.Description {
		(*aggregator).Description = payload.Description
	}

	if (*aggregator).IsActive != payload.IsActive {
		(*aggregator).IsActive = payload.IsActive
	}

	if (*aggregator).Currency != payload.Currency {
		(*aggregator).Currency = payload.Currency
	}
}
