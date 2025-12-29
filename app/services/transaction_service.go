package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/repositories"
	"paymentserviceklink/app/web"
	"paymentserviceklink/config"
	"paymentserviceklink/pkg/pagination"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type TransactionService struct {
	service     *Service
	serviceName string
}

func NewTransactionService(ctx context.Context, repo *repositories.RepositoryContext, cfg *config.Config) *TransactionService {
	return &TransactionService{
		service:     NewService(ctx, repo, cfg),
		serviceName: "transaction_service",
	}
}

func (s TransactionService) AdminGetListTransactionService(pages *pagination.Pages, payload *web.ListTransactionRequest) (any, error) {
	s.serviceName = "TransactionService.AdminGetListTransactionService"
	var starDate time.Time
	var endDate time.Time
	var transactions []*models.Payments
	var totalCount int64
	var err error

	log.Debug().Interface("pages", pages).Interface("payload", payload).Interface("context", s.serviceName).Msg("admin get list transaction service")

	// check payload start date
	if payload.StartDate != "" {
		starDate, _ = time.Parse(time.DateOnly, payload.StartDate)
		log.Debug().Interface("starDate", starDate).Interface("context", s.serviceName).Msg("set start date")
		payload.StartDate = starDate.Format(time.DateOnly)
	} else {
		starDate = time.Now()
		payload.StartDate = starDate.Format(time.DateOnly)
	}

	// check payload end date
	if payload.EndDate != "" {
		endDate, _ = time.Parse(time.DateOnly, payload.EndDate)
		log.Debug().Interface("endDate", endDate).Interface("context", s.serviceName).Msg("set end date")

		payload.EndDate = endDate.AddDate(0, 0, 1).Format(time.DateOnly)
		log.Debug().Interface("endDate", payload.EndDate).Interface("context", s.serviceName).Msg("payload end date")
	} else {
		endDate = time.Now().AddDate(0, 0, 1)

		payload.EndDate = endDate.Format(time.DateOnly)
	}

	// check pages filter
	_, err = helpers.FilterColumnValidation(pages.Filters, models.AllowedFilterColumnPayment())
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("filter column validation")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

	// ctx with cancel
	ctx, cancelFn := context.WithCancel(s.service.ctx)
	eg := errgroup.Group{}

	eg.Go(func() error {
		transactions, err = s.service.repository.FindPaymentRepository(ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Select("payments.id," +
				"payments.transaction_id," +
				"payments.platform_id," +
				"payments.payment_method_id," +
				"payments.amount," +
				"payments.fee_amount," +
				"payments.total_amount," +
				"payments.currency," +
				"payments.status," +
				"payments.customer_id," +
				"payments.customer_name," +
				"payments.customer_email," +
				"payments.customer_phone," +
				"payments.reference_id," +
				"payments.reference_type," +
				"payments.gateway_transaction_id," +
				"payments.gateway_reference," +
				"payments.callback_url," +
				"payments.return_url," +
				"payments.expired_at," +
				"payments.paid_at," +
				"payments.order_id," +
				"payments.expired_time," +
				"payments.aggregator_id," +
				"payments.notification_callback",
			)
			query, args := s.service.repository.SearchQuery(pages.Filters, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			if payload.Currency != "" {
				db = db.Where("currency = ?", payload.Currency)
			}
			if payload.StartDate != "" {
				db = db.Where("created_at >= ?", payload.StartDate)
			}
			if payload.EndDate != "" {
				db = db.Where("created_at <= ?", payload.EndDate)
			}
			db = db.Limit(pages.Limit()).Offset(pages.Offset())
			db = db.Order("payments.id " + pages.Sort)
			db = db.Preload("Platform", func(db *gorm.DB) *gorm.DB {
				db = db.Select("platforms.id," +
					"platforms.name",
				)
				return db
			})
			db = db.Preload("Channel", func(db *gorm.DB) *gorm.DB {
				db = db.Select("channels.id," +
					"channels.name," +
					"channels.bank_name",
				)
				return db
			})
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
			cancelFn()
			log.Error().Err(err).Interface("context", s.serviceName).Msg("get list payment repository")
			return err
		}
		log.Debug().Interface("data transactions", transactions).Interface("context", s.serviceName).Msg("result data get list payment repository")

		return nil
	})

	eg.Go(func() error {
		totalCount, err = s.service.repository.CountGetListPaymentRepository(ctx, func(db *gorm.DB) *gorm.DB {
			query, args := s.service.repository.SearchQuery(pages.Filters, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			if payload.Currency != "" {
				db = db.Where("currency = ?", payload.Currency)
			}
			if payload.StartDate != "" {
				db = db.Where("created_at >= ?", payload.StartDate)
			}
			if payload.EndDate != "" {
				db = db.Where("created_at <= ?", payload.EndDate)
			}
			return db
		})
		if err != nil {
			cancelFn()
			log.Error().Err(err).Interface("context", s.serviceName).Msg("count get list payment repository")
			return err
		}
		log.Debug().Interface("count", totalCount).Interface("context", s.serviceName).Msg("result count get list payment repository")

		return nil
	})

	err = eg.Wait()
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("err sync group get count and get list transaction")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	response := make([]*web.DetailPaymentResponse, 0)
	pages.TotalCount = int(totalCount)

	for _, transaction := range transactions {
		response = append(response, s.service.mapToDetailPaymentResponse(transaction))
	}

	return web.ListResponse{Items: response, Metadata: pages.GetMetadata()}, nil
}

func (s TransactionService) AdminGetDetailTransaction(transactionId string) (any, error) {
	s.serviceName = "TransactionService.AdminGetDetailTransaction"
	log.Debug().Interface("transactionId", transactionId).Interface("context", s.serviceName).Msg("admin get detail transaction service")

	transaction, err := s.service.repository.GetPaymentRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("payments.id," +
			"payments.transaction_id," +
			"payments.platform_id," +
			"payments.payment_method_id," +
			"payments.amount," +
			"payments.fee_amount," +
			"payments.total_amount," +
			"payments.currency," +
			"payments.status," +
			"payments.customer_id," +
			"payments.customer_name," +
			"payments.customer_email," +
			"payments.customer_phone," +
			"payments.reference_id," +
			"payments.reference_type," +
			"payments.gateway_transaction_id," +
			"payments.gateway_reference," +
			"payments.callback_url," +
			"payments.return_url," +
			"payments.expired_at," +
			"payments.paid_at," +
			"payments.order_id," +
			"payments.expired_time," +
			"payments.aggregator_id," +
			"payments.notification_callback",
		)
		db = db.Where("payments.transaction_id = ?", transactionId)
		db = db.Preload("Platform", func(db *gorm.DB) *gorm.DB {
			db = db.Select("platforms.id," +
				"platforms.name",
			)
			return db
		})
		db = db.Preload("Channel", func(db *gorm.DB) *gorm.DB {
			db = db.Select("channels.id," +
				"channels.name," +
				"channels.bank_name",
			)
			return db
		})
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
		log.Error().Err(err).Interface("context", s.serviceName).Msg("get payment by id repository")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(fmt.Errorf("transaction, %v", err), s.serviceName).WithStatusCode(http.StatusNotFound)
		}

		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data transaction", transaction).Interface("context", s.serviceName).Msg("result data get payment repository")

	return s.service.mapToDetailPaymentResponse(transaction), nil
}
