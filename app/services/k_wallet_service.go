package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"paymentserviceklink/app/client/espay"
	clientsenangpay "paymentserviceklink/app/client/senangpay"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/helpers"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/repositories"
	"paymentserviceklink/app/web"
	"paymentserviceklink/config"
	"paymentserviceklink/pkg/pagination"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type KWalletService struct {
	service     *Service
	serviceName string
}

type IKWallet interface {
	create(payload *web.CreateKWalletRequest, lastId int64) (*web.KWalletResponse, error)
	createTopup(kWallet *models.KWallet, topupTransaction *models.TopupTransaction) (*web.TopupTransaction, error)
}

func NewKWalletService(ctx context.Context, repo *repositories.RepositoryContext, cfg *config.Config) *KWalletService {
	return &KWalletService{
		service:     NewService(ctx, repo, cfg),
		serviceName: "KWalletService",
	}
}

func (s *KWalletService) AdminGetListKWalletService(pages *pagination.Pages) (*web.ListResponse, error) {
	log.Debug().Interface("pages", pages).Interface("context", s.serviceName).Msg("admin get list k-wallet service")

	var kWallets []*models.KWallet
	var totalCount int64
	var err error

	filter, err := helpers.FilterColumnValidation(pages.Filters, models.AllowedFilterColumnKWallet())
	if err != nil {
		log.Error().Err(err).Msg("filter column validation error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

	ctx, cancelFn := context.WithCancel(s.service.ctx)
	eg := errgroup.Group{}

	eg.Go(func() error {
		kWallets, err = s.service.repository.FindKWalletRepository(ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Select("k_wallet.id," +
				"k_wallet.member_id," +
				"k_wallet.full_name," +
				"k_wallet.no_rekening," +
				"k_wallet.gen_va," +
				"k_wallet.balance," +
				"k_wallet.currency," +
				"k_wallet.symbol," +
				"k_wallet.status," +
				"k_wallet.is_active",
			)
			query, args := s.service.repository.SearchQuery(filter, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			db = db.Limit(pages.Limit()).Offset(pages.Offset())
			db = db.Order("id " + pages.Sort)
			return db
		})
		if err != nil {
			cancelFn()
			log.Error().Err(err).Msg("get list k-wallet repository error")
			return err
		}

		return nil
	})

	eg.Go(func() error {
		totalCount, err = s.service.repository.GetTotalCountKWalletRepository(ctx, pages, func(db *gorm.DB) *gorm.DB {
			query, args := s.service.repository.SearchQuery(pages.Filters, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			return db
		})
		if err != nil {
			cancelFn()
			log.Error().Err(err).Msg("get total count k-wallet repository error")
			return err
		}

		return nil
	})

	err = eg.Wait()
	if err != nil {
		log.Error().Err(err).Msg("async group error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	response := make([]*web.KWalletResponse, 0, len(kWallets))

	for _, kWallet := range kWallets {
		response = append(response, s.service.mapToKWalletResponse(kWallet))
	}

	pages.TotalCount = int(totalCount)

	return &web.ListResponse{Items: response, Metadata: pages.GetMetadata()}, nil
}

func (s *Service) mapToKWalletResponse(kWallet *models.KWallet) *web.KWalletResponse {
	virtualAccounts := make([]*web.VirtualAccountKWallet, 0, len(kWallet.VirtualAccount))

	for _, va := range kWallet.VirtualAccount {
		virtualAccounts = append(virtualAccounts, s.mapToVirtualAccount(va))
	}

	return &web.KWalletResponse{
		MemberID:       kWallet.MemberID,
		FullName:       kWallet.FullName,
		NoRekening:     kWallet.NoRekening,
		GenVa:          kWallet.GenVA,
		Balance:        kWallet.Balance.IntPart(),
		Currency:       kWallet.Currency,
		Symbol:         kWallet.Symbol,
		Status:         kWallet.Status,
		IsActive:       kWallet.IsActive,
		VirtualAccount: virtualAccounts,
	}
}

func (s *KWalletService) AdminGetDetailKWalletService(noRekening string) (*web.KWalletResponse, error) {
	var kWallet *models.KWallet
	var err error

	kWallet, err = s.service.repository.GetKWalletRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("k_wallet.id," +
			"k_wallet.member_id," +
			"k_wallet.full_name," +
			"k_wallet.no_rekening," +
			"k_wallet.gen_va," +
			"k_wallet.balance," +
			"k_wallet.currency," +
			"k_wallet.symbol," +
			"k_wallet.status," +
			"k_wallet.is_active",
		)
		db = db.Where("k_wallet.no_rekening = ?", noRekening)
		db = db.Preload("VirtualAccount", func(db *gorm.DB) *gorm.DB {
			db = db.Select("virtual_account_k_wallet.id," +
				"virtual_account_k_wallet.k_wallet_id," +
				"virtual_account_k_wallet.virtual_account," +
				"virtual_account_k_wallet.bank",
			)
			return db
		})
		return db
	})
	if err != nil {
		log.Error().Err(err).Msg("get k-wallet repository error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return s.service.mapToKWalletResponse(kWallet), nil
}

func (s *KWalletService) AdminGetListKWalletTransactionService(noRekening string, pages *pagination.Pages) (interface{}, error) {
	var kWalletTransactions []*models.KWalletTransaction
	var totalCount int64
	var err error

	pages.Filters, err = models.TransformFilterColumnKWalletTransaction(pages.Filters, models.AllowedFilterColumnKWalletTransaction())
	if err != nil {
		log.Error().Err(err).Msg("transform filter column k-wallet transaction error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

	filter, err := helpers.FilterColumnValidation(pages.Filters, models.AllowedFilterColumnKWalletTransaction())
	if err != nil {
		log.Error().Err(err).Msg("filter column validation error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusBadRequest)
	}

	ctx, cancelFn := context.WithCancel(s.service.ctx)
	eg := errgroup.Group{}

	eg.Go(func() error {
		kWalletTransactions, err = s.service.repository.FindKWalletTransactionRepository(ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Select("k_wallet_transaction.id," +
				"k_wallet_transaction.k_wallet_id," +
				"k_wallet_transaction.k_wallet_type_transaction_id," +
				"k_wallet_transaction.payment_id," +
				"k_wallet_transaction.title," +
				"k_wallet_transaction.payment_code," +
				"k_wallet_transaction.transaction_code," +
				"k_wallet_transaction.transaction_type," +
				"k_wallet_transaction.direction," +
				"k_wallet_transaction.counterparty_name," +
				"k_wallet_transaction.counterparty_bank," +
				"k_wallet_transaction.payment_channel," +
				"k_wallet_transaction.description," +
				"k_wallet_transaction.balance," +
				"k_wallet_transaction.debit," +
				"k_wallet_transaction.credit," +
				"k_wallet_transaction.amount," +
				"k_wallet_transaction.currency," +
				"k_wallet_transaction.symbol," +
				"k_wallet_transaction.status," +
				"k_wallet_transaction.month," +
				"k_wallet_transaction.year," +
				"k_wallet_transaction.date," +
				"k_wallet_transaction.time," +
				"k_wallet_transaction.datetime",
			)
			db = db.Joins("JOIN k_wallet on k_wallet.id = k_wallet_transaction.k_wallet_id")
			db = db.Where("k_wallet.no_rekening = ?", noRekening)
			query, args := s.service.repository.SearchQuery(filter, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			db = db.Limit(pages.Limit()).Offset(pages.Offset())
			db = db.Order("k_wallet_transaction.id " + pages.Sort)
			return db
		})
		if err != nil {
			cancelFn()
			log.Error().Err(err).Msg("get list k-wallet transaction repository error")
			return err
		}

		return nil
	})

	eg.Go(func() error {
		totalCount, err = s.service.repository.GetTotalCountKWalletTransactionRepository(ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Joins("JOIN k_wallet on k_wallet.id = k_wallet_transaction.k_wallet_id")
			db = db.Where("k_wallet.no_rekening = ?", noRekening)
			query, args := s.service.repository.SearchQuery(pages.Filters, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			return db
		})
		if err != nil {
			cancelFn()
			log.Error().Err(err).Msg("get total count k-wallet transaction repository error")
			return err
		}

		return nil
	})

	err = eg.Wait()
	if err != nil {
		log.Error().Err(err).Msg("async group error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	pages.TotalCount = int(totalCount)

	response := make([]*web.KWalletTransaction, 0, len(kWalletTransactions))

	for _, kWalletTransaction := range kWalletTransactions {
		response = append(response, s.mapToKWalletTransactionResponse(kWalletTransaction))
	}

	return web.ListResponse{Items: response, Metadata: pages.GetMetadata()}, nil
}

func (s *KWalletService) mapToKWalletTransactionResponse(kWalletTransaction *models.KWalletTransaction) *web.KWalletTransaction {
	return &web.KWalletTransaction{
		ID:                       kWalletTransaction.ID,
		KWalletTypeTransactionID: kWalletTransaction.KWalletID,
		PaymentID:                kWalletTransaction.PaymentID,
		Title:                    kWalletTransaction.Title,
		PaymentCode:              kWalletTransaction.PaymentCode,
		TransactionCode:          kWalletTransaction.TransactionCode,
		TransactionType:          kWalletTransaction.TransactionType,
		Direction:                kWalletTransaction.Direction,
		CounterpartyName:         kWalletTransaction.CounterpartyName,
		CounterpartyBank:         kWalletTransaction.CounterpartyBank,
		PaymentChannel:           kWalletTransaction.PaymentChannel,
		Description:              kWalletTransaction.Description,
		Balance:                  kWalletTransaction.Balance.IntPart(),
		Debit:                    kWalletTransaction.Debit.IntPart(),
		Credit:                   kWalletTransaction.Credit.IntPart(),
		Amount:                   kWalletTransaction.Amount.IntPart(),
		Currency:                 kWalletTransaction.Currency,
		Symbol:                   kWalletTransaction.Symbol,
		Date:                     kWalletTransaction.Date.Format(time.DateOnly),
		Time:                     kWalletTransaction.Time,
		DateTime:                 kWalletTransaction.DateTime.Format(time.DateTime),
	}
}

func (s *KWalletService) AdminGetListTopupKWalletService(pages *pagination.Pages) (*web.ListResponse, error) {
	var topups []*models.TopupTransaction
	var totalCount int64
	var err error

	ctx, cancelFn := context.WithCancel(s.service.ctx)
	eg := errgroup.Group{}

	eg.Go(func() error {
		topups, err = s.service.repository.FindTopupTransactionRepository(ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Select("topup_transaction.id," +
				"topup_transaction.k_wallet_id," +
				"topup_transaction.member_id," +
				"topup_transaction.channel_id," +
				"topup_transaction.aggregator," +
				"topup_transaction.merchant," +
				"topup_transaction.amount," +
				"topup_transaction.fee_admin," +
				"topup_transaction.currency," +
				"topup_transaction.symbol," +
				"topup_transaction.reference_id," +
				"topup_transaction.status," +
				"topup_transaction.completed_at," +
				"topup_transaction.description")
			db = db.Joins("JOIN k_wallet on k_wallet.id = topup_transaction.k_wallet_id")
			query, args := s.service.repository.SearchQuery(pages.Filters, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			db = db.Preload("KWallet", func(db *gorm.DB) *gorm.DB {
				db = db.Select("k_wallet.id," +
					"k_wallet.member_id," +
					"k_wallet.full_name," +
					"k_wallet.no_rekening",
				)
				return db
			})
			db = db.Limit(pages.Limit()).Offset(pages.Offset())
			db = db.Order("topup_transaction.id " + pages.Sort)
			return db
		})
		if err != nil {
			cancelFn()
			log.Error().Err(err).Msg("get list topup transaction repository error")
			return err
		}

		return nil
	})

	eg.Go(func() error {
		totalCount, err = s.service.repository.GetTotalCountTopupTransactionRepository(ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Joins("JOIN k_wallet on k_wallet.id = topup_transaction.k_wallet_id")
			query, args := s.service.repository.SearchQuery(pages.Filters, pages.JoinOperator)
			if query != "" {
				db = db.Where(query, args...)
			}
			return db
		})
		if err != nil {
			cancelFn()
			log.Error().Err(err).Msg("get total count topup transaction repository error")
			return err
		}

		return nil
	})

	err = eg.Wait()
	if err != nil {
		log.Error().Err(err).Msg("async group error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	pages.TotalCount = int(totalCount)

	response := make([]*web.TopupTransaction, 0, len(topups))

	for _, topup := range topups {
		response = append(response, s.mapToTopupTransaction(topup))
	}

	return &web.ListResponse{Items: response, Metadata: pages.GetMetadata()}, nil
}

func (s *KWalletService) mapToTopupTransaction(topup *models.TopupTransaction) *web.TopupTransaction {
	var kWallet *web.KWalletResponse
	if topup.KWallet != nil {
		kWallet = s.service.mapToKWalletResponse(topup.KWallet)
	}
	return &web.TopupTransaction{
		ID:          topup.ID,
		KWalletID:   topup.KWalletID,
		MemberID:    topup.MemberId,
		ChannelId:   topup.ChannelID,
		Aggregator:  topup.Aggregator,
		Merchant:    topup.Merchant,
		Amount:      topup.Amount,
		FeeAdmin:    topup.FeeAdmin,
		Currency:    topup.Currency,
		Symbol:      topup.Symbol,
		ReferenceID: topup.ReferenceID,
		Status:      topup.Status,
		CompletedAt: topup.CompletedAt.Format(time.DateTime),
		Description: topup.Description,
		KWallet:     kWallet,
	}
}

type KWalletIDR struct {
	kWalletService *KWalletService
}

func NewKWalletIDR(kWalletService *KWalletService) IKWallet {
	return &KWalletIDR{
		kWalletService: kWalletService,
	}
}

func (k *KWalletIDR) create(payload *web.CreateKWalletRequest, lastId int64) (*web.KWalletResponse, error) {

	//filter := make([]*pagination.Filter, 0)
	//[]pagination.Filter{
	//	//{
	//	//	ID:       "",
	//	//	Value:    nil,
	//	//	Variant:  "",
	//	//	Operator: "",
	//	//	FilterID: "",
	//	//},
	//}

	//filter, err = helpers.FilterColumnValidation(filter, models.AllowedFilterColumnConfiguration())
	//if err != nil {
	//	log.Error().Err(err).Interface("context", s.serviceName).Msg("filter column validation error")
	//	return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	//}
	//
	//configuration, err := s.service.repository.GetConfigurationRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
	//	query, args := s.service.repository.SearchQuery(filter, "and")
	//	if query != "" {
	//		db = db.Where(query, args...)
	//	}
	//	return db
	//})
	//if err != nil {
	//	log.Error().Err(err).Interface("context", s.serviceName).Msg("get configuration repository error")
	//	return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	//}
	//log.Debug().Interface("data configuration", configuration).Interface("context", s.serviceName).Msg("result data configuration get configuration repository")

	configValue := enums.SANDBOX
	if k.kWalletService.service.config.EspayIsProduction {
		configValue = enums.PRODUCTION
	}

	configuration := &models.Configuration{
		Id:           0,
		AggregatorId: 0,
		ConfigKey:    "",
		ConfigValue:  configValue,
		ConfigName:   "",
		ConfigJson: models.ConfigJson{
			SandboxBaseUrl:       k.kWalletService.service.config.EspayTopupBaseUrl,
			ProductionBaseUrl:    k.kWalletService.service.config.EspayTopupBaseUrl,
			SandboxMerchantId:    "",
			ProductionMerchantId: "",
			SandboxMerchantCode:  helpers.EncryptAES(k.kWalletService.service.config.EspayTopupMerchantCode),
			//SandboxMerchantCode:          helpers.EncryptAES("SGWKLINKTOPUP01"),
			ProductionMerchantCode: helpers.EncryptAES(k.kWalletService.service.config.EspayTopupMerchantCode),
			//SandboxMerchantName:          helpers.EncryptAES("KLINKTOPUP01"),
			SandboxMerchantName:    helpers.EncryptAES(k.kWalletService.service.config.EspayTopupMerchantName),
			ProductionMerchantName: helpers.EncryptAES(k.kWalletService.service.config.EspayTopupMerchantName),
			//SandboxApiKey:                helpers.EncryptAES("c5c52819439156857a63a74d23f534ee"),
			SandboxApiKey:       helpers.EncryptAES(k.kWalletService.service.config.EspayTopupApiKey),
			ProductionApiKey:    helpers.EncryptAES(k.kWalletService.service.config.EspayTopupApiKey),
			SandboxServerKey:    "",
			ProductionServerKey: "",
			SandboxSecretKey:    "",
			ProductionSecretKey: "",
			SandboxClientKey:    "",
			ProductionClientKey: "",
			//SandboxSignatureKey:          helpers.EncryptAES("sc429jdqy5sgd6tc"),
			SandboxSignatureKey:    helpers.EncryptAES(k.kWalletService.service.config.EspayTopupSignatureKey),
			ProductionSignatureKey: helpers.EncryptAES(k.kWalletService.service.config.EspayTopupSignatureKey),
			//SandboxCredentialPassword:    helpers.EncryptAES("KQYORQKV"),
			SandboxCredentialPassword:    helpers.EncryptAES(k.kWalletService.service.config.EspayTopupCredentialPassword),
			ProductionCredentialPassword: helpers.EncryptAES(k.kWalletService.service.config.EspayTopupCredentialPassword),
			PublicKey:                    helpers.EncryptAES(k.kWalletService.service.config.EspayTopupPublicKey),
			PrivateKey:                   helpers.EncryptAES(k.kWalletService.service.config.EspayTopupPrivateKey),
			ReturnUrl:                    "",
		},
		IsActive: true,
		Aggregator: &models.Aggregator{
			Id:          0,
			Name:        enums.AGGREGATOR_NAME_ESPAY,
			Slug:        enums.ProviderPaymentMethod(strings.ToLower(fmt.Sprint(enums.AGGREGATOR_NAME_ESPAY))),
			Description: "",
			IsActive:    true,
			Currency:    enums.CURRENCY_IDR,
		},
	}

	espayStrategy := espay.NewEspay(k.kWalletService.service.repository.HttpClient, configuration)

	orderId := time.Now().Unix()
	orderString1 := strconv.Itoa(int(lastId))
	orderString2 := strconv.Itoa(int(orderId))
	var last8 string

	if len(orderString1) >= 8 {
		last8 = orderString1[len(orderString1)-8:]
	} else {
		last8 = fmt.Sprintf("%08s", orderString1)
	}

	//last8 = "00000001"

	espayRequest := espay.PaymentRequest{
		RQUUID:        uuid.New().String(),
		RQDateTime:    time.Now(),
		OrderID:       last8,
		Amount:        "",
		CCY:           enums.CURRENCY_IDR,
		Method:        enums.PAYMENT_METHOD_VIRTUAL_ACCOUNT,
		CustomerPhone: payload.CustomerPhone,
		CustomerName:  payload.CustomerName,
		CustomerEmail: payload.CustomerEmail,
		BankCode:      "",
		VaExpired:     "999999",
	}

	vaStatic, err := espayStrategy.CreateVaStatic(k.kWalletService.service.ctx, espayRequest)
	if err != nil {
		log.Error().Err(err).Interface("context", k.kWalletService.serviceName).Msg("pay espay strategy error")
		return nil, helpers.NewErrorTrace(err, k.kWalletService.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	var response *web.KWalletResponse

	err = k.kWalletService.service.repository.WithTransaction(func(tx *gorm.DB) error {
		// k-wallet data
		kWallet := &models.KWallet{
			MemberID:   payload.MemberID,
			FullName:   payload.CustomerName,
			NoRekening: fmt.Sprintf("%s%s", orderString2, last8),
			GenVA:      last8,
			Balance:    decimal.NewFromInt(0),
			Currency:   enums.CURRENCY_IDR,
			Symbol:     enums.SYMBOL_CURRENCY_IDR,
			Status:     enums.KWalletStatusActive,
			IsActive:   true,
		}

		kWallet, err = k.kWalletService.service.repository.InsertKWalletRepositoryTx(k.kWalletService.service.ctx, tx, kWallet)
		if err != nil {
			log.Error().Err(err).Interface("context", k.kWalletService.serviceName).Msg("insert k-wallet repository error")
			return helpers.NewErrorTrace(err, k.kWalletService.serviceName).WithStatusCode(http.StatusInternalServerError)
		}

		// virtual account k-wallet
		virtualAccountKWallets := make([]*models.VirtualAccountKWallet, 0)

		for key, value := range vaStatic {
			channel, errCh := k.kWalletService.service.repository.GetChannelRepository(k.kWalletService.service.ctx, func(db *gorm.DB) *gorm.DB {
				db = db.Select("channels.id," +
					"channels.bank_name," +
					"channels.bank_code",
				)
				db = db.Where("channels.bank_code = ?", key)
				return db
			})
			if errCh != nil {
				log.Error().Err(err).Msg("get channel repository error")
				continue
			}

			virtualAccountKWallets = append(virtualAccountKWallets, &models.VirtualAccountKWallet{
				KWalletId:      kWallet.ID,
				VirtualAccount: value.VaNumber,
				Bank:           channel.BankName,
				BankCode:       channel.BankCode,
			})
		}

		err = k.kWalletService.service.repository.InsertBatchVirtualAccountRepositoryTx(k.kWalletService.service.ctx, tx, virtualAccountKWallets)
		if err != nil {
			log.Error().Err(err).Interface("context", k.kWalletService.serviceName).Msg("insert batch virtual account repository error")
			return err
		}

		kWallet.VirtualAccount = virtualAccountKWallets

		response = k.kWalletService.service.mapToKWalletResponse(kWallet)

		return nil
	})
	if err != nil {
		log.Error().Err(err).Interface("context", k.kWalletService.serviceName).Msg("with transaction error")
		return nil, helpers.NewErrorTrace(err, k.kWalletService.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return response, nil
}

func (k *KWalletIDR) createTopup(kWallet *models.KWallet, toptupTransaction *models.TopupTransaction) (*web.TopupTransaction, error) {
	return nil, helpers.NewErrorTrace(fmt.Errorf("not implemented"), k.kWalletService.serviceName).WithStatusCode(http.StatusInternalServerError)
}

type KWalletMYR struct {
	kWalletService *KWalletService
}

func NewKWalletMYR(kWalletService *KWalletService) IKWallet {
	return &KWalletMYR{
		kWalletService: kWalletService,
	}
}

func (k *KWalletMYR) create(payload *web.CreateKWalletRequest, lastId int64) (*web.KWalletResponse, error) {
	var err error
	orderId := time.Now().Unix()
	orderString1 := strconv.Itoa(int(lastId))
	orderString2 := strconv.Itoa(int(orderId))
	var last8 string

	if len(orderString1) >= 8 {
		last8 = orderString1[len(orderString1)-8:]
	} else {
		last8 = fmt.Sprintf("%08s", orderString1)
	}

	var response *web.KWalletResponse

	err = k.kWalletService.service.repository.WithTransaction(func(tx *gorm.DB) error {
		// k-wallet data
		kWallet := &models.KWallet{
			MemberID:   payload.MemberID,
			FullName:   payload.CustomerName,
			NoRekening: fmt.Sprintf("%s%s", orderString2, last8),
			GenVA:      last8,
			Balance:    decimal.NewFromInt(0),
			Currency:   enums.CURRENCY_MYR,
			Symbol:     enums.SYMBOL_CURRENCY_MYR,
			Status:     enums.KWalletStatusActive,
			IsActive:   true,
		}

		kWallet, err = k.kWalletService.service.repository.InsertKWalletRepositoryTx(k.kWalletService.service.ctx, tx, kWallet)
		if err != nil {
			log.Error().Err(err).Interface("context", k.kWalletService.serviceName).Msg("insert k-wallet repository error")
			return helpers.NewErrorTrace(err, k.kWalletService.serviceName).WithStatusCode(http.StatusInternalServerError)
		}

		response = k.kWalletService.service.mapToKWalletResponse(kWallet)

		return nil
	})
	if err != nil {
		log.Error().Err(err).Interface("context", k.kWalletService.serviceName).Msg("with transaction error")
		return nil, helpers.NewErrorTrace(err, k.kWalletService.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	return response, nil
}

func (k *KWalletMYR) createTopup(kWallet *models.KWallet, topupTransaction *models.TopupTransaction) (*web.TopupTransaction, error) {

	configuration := &models.Configuration{}

	senangpay := clientsenangpay.NewSenangpay(k.kWalletService.service.repository.HttpClient, configuration)

	paymentRequest := &clientsenangpay.PaymentRequest{
		OrderID: topupTransaction.ReferenceID,
		Amount:  topupTransaction.Amount.String(),
		Detail:  fmt.Sprintf("Topup k-wallet : %v", kWallet.NoRekening),
		Name:    kWallet.FullName,
		Email:   string(enums.STATIS_EMAIL),
		Phone:   string(enums.STATIS_PHONE_NUMBER),
	}

	_ = senangpay.GeneratePaymentURL(paymentRequest)
	return nil, helpers.NewErrorTrace(fmt.Errorf("not implemented"), k.kWalletService.serviceName).WithStatusCode(http.StatusInternalServerError)
}

func (s *KWalletService) CreateKWalletService(payload *web.CreateKWalletRequest) (*web.KWalletResponse, error) {
	platform, err := GetClientAuth(s.service)
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("get client auth error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusUnauthorized)
	}
	log.Debug().Interface("data platform", platform).Interface("context", s.serviceName).Msg("result data platform client auth")

	// get last id in table k-wallet for gen va
	lastId, err := s.service.repository.GetLastIdKWalletRepository(s.service.ctx)
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("get last id repository error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}
	log.Debug().Interface("data last id", lastId).Interface("context", s.serviceName).Msg("result data last id get last id repository")

	var kWallet IKWallet
	if payload.Currency == enums.CURRENCY_IDR {
		kWallet = NewKWalletIDR(s)
	} else if payload.Currency == enums.CURRENCY_MYR {
		kWallet = NewKWalletMYR(s)
	} else {
		return nil, helpers.NewErrorTrace(fmt.Errorf("currency not implemented"), s.serviceName).WithStatusCode(http.StatusNotImplemented)
	}

	response, err := kWallet.create(payload, lastId)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *Service) mapToVirtualAccount(va *models.VirtualAccountKWallet) *web.VirtualAccountKWallet {

	return &web.VirtualAccountKWallet{
		ID:             va.ID,
		Bank:           va.Bank,
		BankCode:       va.BankCode,
		VirtualAccount: va.VirtualAccount,
	}
}

func (s *KWalletService) GetKWalletMemberService(payload *web.GetKWalletRequest) ([]*web.KWalletResponse, error) {
	kWallets, err := s.service.repository.FindKWalletRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("k_wallet.member_id," +
			"k_wallet.full_name," +
			"k_wallet.no_rekening," +
			"k_wallet.balance," +
			"k_wallet.currency," +
			"k_wallet.symbol," +
			"k_wallet.status," +
			"k_wallet.is_active",
		)
		db = db.Where("k_wallet.member_id = ?", payload.MemberId)
		return db
	})
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("get k-wallet repository error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	response := make([]*web.KWalletResponse, 0, len(kWallets))

	for _, kWallet := range kWallets {
		response = append(response, s.service.mapToKWalletResponse(kWallet))
	}

	return response, nil
}

func (s *KWalletService) GetListKWalletTransactionMemberService(payload *web.GetListKWalletTransactionRequest) ([]*web.KWalletTransaction, error) {
	response := make([]*web.KWalletTransaction, 0)
	timeNow := time.Now()
	var err error

	filterColumnKWalletTransaction := make([]*pagination.Filter, 0)

	if payload.FromDate != "" && payload.ToDate != "" {
		filterColumnKWalletTransaction = append(filterColumnKWalletTransaction, []*pagination.Filter{
			&pagination.Filter{
				ID:       "date",
				Value:    payload.FromDate,
				Variant:  "date",
				Operator: "gte",
				FilterID: "",
			},
			&pagination.Filter{
				ID:       "date",
				Value:    payload.ToDate,
				Variant:  "date",
				Operator: "lte",
				FilterID: "",
			},
		}...)

		fromDate, errF := time.Parse(time.DateOnly, payload.FromDate)
		if errF != nil {
			log.Error().Err(errF).Msg("error parse from date")
			return nil, helpers.NewErrorTrace(fmt.Errorf("from date format must be %v error %v", time.DateOnly, errF), s.serviceName).WithStatusCode(http.StatusBadRequest)
		}

		toDate, errT := time.Parse(time.DateOnly, payload.ToDate)
		if errT != nil {
			log.Error().Err(errT).Msg("error parse to date")
			return nil, helpers.NewErrorTrace(fmt.Errorf("to date format must be %v error %v", time.DateOnly, errT), s.serviceName).WithStatusCode(http.StatusBadRequest)
		}

		message, validationDuration := helpers.ValidationDuration(fromDate, toDate, 30*24*time.Hour)
		if !validationDuration {
			return nil, helpers.NewErrorTrace(fmt.Errorf(message), s.serviceName).WithStatusCode(http.StatusBadRequest)
		}
	} else if payload.Month != 0 {
		filterColumnKWalletTransaction = append(filterColumnKWalletTransaction, []*pagination.Filter{
			&pagination.Filter{
				ID:       "month",
				Value:    payload.Month,
				Variant:  "number",
				Operator: "eq",
				FilterID: "",
			},
		}...)
	} else {
		filterColumnKWalletTransaction = append(filterColumnKWalletTransaction, []*pagination.Filter{
			&pagination.Filter{
				ID:       "month",
				Value:    timeNow.Month(),
				Variant:  "number",
				Operator: "eq",
				FilterID: "",
			},
		}...)
	}
	log.Debug().Interface("filter column k-wallet transaction", filterColumnKWalletTransaction).Msg("data filter column k-wallet transaction")

	filter, err := helpers.FilterColumnValidation(filterColumnKWalletTransaction, models.AllowedFilterColumnKWalletTransaction())
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("filter column validation error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	kWalletTransactions, err := s.service.repository.FindKWalletTransactionRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("k_wallet_transaction.id," +
			"k_wallet_transaction.k_wallet_id," +
			"k_wallet_transaction.k_wallet_type_transaction_id," +
			"k_wallet_transaction.payment_id," +
			"k_wallet_transaction.title," +
			"k_wallet_transaction.payment_code," +
			"k_wallet_transaction.transaction_code," +
			"k_wallet_transaction.transaction_type," +
			"k_wallet_transaction.direction," +
			"k_wallet_transaction.counterparty_name," +
			"k_wallet_transaction.counterparty_bank," +
			"k_wallet_transaction.payment_channel," +
			"k_wallet_transaction.description," +
			"k_wallet_transaction.balance," +
			"k_wallet_transaction.debit," +
			"k_wallet_transaction.credit," +
			"k_wallet_transaction.amount," +
			"k_wallet_transaction.currency," +
			"k_wallet_transaction.symbol," +
			"k_wallet_transaction.status," +
			"k_wallet_transaction.month," +
			"k_wallet_transaction.year," +
			"k_wallet_transaction.date," +
			"k_wallet_transaction.time," +
			"k_wallet_transaction.datetime",
		)
		db = db.Joins("JOIN k_wallet on k_wallet.id = k_wallet_transaction.k_wallet_id")
		db = db.Where("k_wallet.member_id = ? and k_wallet.no_rekening = ?", payload.MemberId, payload.NoRekening)
		query, args := s.service.repository.SearchQuery(filter, "and")
		if query != "" {
			db = db.Where(query, args...)
		}
		db = db.Order("k_wallet_transaction.id desc")
		return db
	})
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("get k-wallet transaction repository error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	for _, kWalletTransaction := range kWalletTransactions {
		response = append(response, s.mapToKWalletTransactionResponse(kWalletTransaction))
	}

	return response, nil
}

func (s *KWalletService) GetVirtualAccountKWalletService(payload *web.GetVirtualAccountKWalletRequest) ([]*web.VirtualAccountKWallet, error) {
	response := make([]*web.VirtualAccountKWallet, 0)

	virtualAccountKWallets, err := s.service.repository.FindVirtualAccountKWalletRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("virtual_account_k_wallet.k_wallet_id," +
			"virtual_account_k_wallet.virtual_account," +
			"virtual_account_k_wallet.bank," +
			"virtual_account_k_wallet.bank_code",
		)
		db = db.Joins("JOIN k_wallet on k_wallet.id = virtual_account_k_wallet.k_wallet_id")
		db = db.Where("k_wallet.member_id = ? and k_wallet.no_rekening = ?", payload.MemberId, payload.NoRekening)
		db = db.Order("virtual_account_k_wallet.bank asc")
		return db
	})
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("get virtual account k-wallet repository error")
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	for _, virtualAccountKWallet := range virtualAccountKWallets {
		response = append(response, s.service.mapToVirtualAccount(virtualAccountKWallet))
	}

	return response, nil
}

func (s *KWalletService) CreateTopupKWalletService(payload *web.CreateTopupKWalletRequest) (*web.TopupTransaction, error) {
	kWallet, err := s.service.repository.GetKWalletRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
		db = db.Select("k_wallet.member_id," +
			"k_wallet.full_name," +
			"k_wallet.no_rekening," +
			"k_wallet.balance," +
			"k_wallet.currency," +
			"k_wallet.symbol," +
			"k_wallet.status," +
			"k_wallet.is_active",
		)
		db = db.Where("k_wallet.member_id = ?", payload.MemberId)
		return db
	})
	if err != nil {
		log.Error().Err(err).Interface("context", s.serviceName).Msg("get k-wallet repository error")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusNotFound)
		}
		return nil, helpers.NewErrorTrace(err, s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	if kWallet.Status != enums.KWalletStatusActive {
		return nil, helpers.NewErrorTrace(fmt.Errorf("k-wallet is %v", kWallet.Status), s.serviceName).WithStatusCode(http.StatusNotFound)
	}

	var iKWallet IKWallet

	if kWallet.Currency == enums.CURRENCY_IDR {
		iKWallet = NewKWalletIDR(s)
	} else if kWallet.Currency == enums.CURRENCY_MYR {
		iKWallet = NewKWalletMYR(s)
	} else {
		return nil, helpers.NewErrorTrace(fmt.Errorf("currency not implemented"), s.serviceName).WithStatusCode(http.StatusInternalServerError)
	}

	topupTransaction := &models.TopupTransaction{}

	response, err := iKWallet.createTopup(kWallet, topupTransaction)
	if err != nil {
		return nil, err
	}

	return response, nil
}
