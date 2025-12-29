package repositories

import (
	"context"
	"errors"
	"paymentserviceklink/app/models"
	"paymentserviceklink/pkg/pagination"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (rc *RepositoryContext) GetListKWalletRepository(ctx context.Context, pages *pagination.Pages) ([]*models.KWallet, error) {
	var (
		db       = rc.db
		kWallets []*models.KWallet
		err      error
	)

	db = db.Table("k_wallet").Model(&models.KWallet{}).WithContext(ctx)
	db = db.Limit(pages.Limit()).Offset(pages.Offset())

	query, args := rc.SearchQuery(pages.Filters, pages.JoinOperator)

	db = db.Where(query, args...)

	db = db.Order("id " + pages.Sort)

	err = db.Find(&kWallets).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get all k-wallet")
		return nil, err
	}

	return kWallets, err
}

func (rc *RepositoryContext) GetTotalCountKWalletRepository(ctx context.Context, pages *pagination.Pages, fn func(db *gorm.DB) *gorm.DB) (int64, error) {
	var (
		db         = rc.db
		totalCount int64
		err        error
	)

	db = db.Table("k_wallet").Model(&models.KWallet{}).WithContext(ctx)
	db = fn(db)

	err = db.Count(&totalCount).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get total count k-wallet")
		return 0, err
	}

	return totalCount, err
}

func (rc *RepositoryContext) GetKWalletRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) (*models.KWallet, error) {
	var (
		db      = rc.db
		kWallet *models.KWallet
		err     error
	)

	db = db.Table("k_wallet").Model(&models.KWallet{}).WithContext(ctx)
	db = fn(db)

	err = db.First(&kWallet).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get k-wallet")
		return nil, err
	}

	return kWallet, err
}

func (rc *RepositoryContext) FindKWalletRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) ([]*models.KWallet, error) {
	var (
		db      = rc.db
		kWallet []*models.KWallet
		err     error
	)

	db = db.Table("k_wallet").Model(&models.KWallet{}).WithContext(ctx)
	db = fn(db)

	err = db.Find(&kWallet).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get k-wallet")
		return nil, err
	}

	return kWallet, nil
}

func (rc *RepositoryContext) GetListKWalletTransactionRepository(ctx context.Context, pages *pagination.Pages) ([]*models.KWalletTransaction, error) {
	var (
		db                  = rc.db
		kWalletTransactions []*models.KWalletTransaction
		err                 error
	)

	db = db.Table("k_wallet_transaction").Model(&models.KWalletTransaction{}).WithContext(ctx)

	err = db.Find(&kWalletTransactions).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get all k-wallet transaction")
		return nil, err
	}

	return kWalletTransactions, nil
}

func (rc *RepositoryContext) GetTotalCountKWalletTransactionRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) (int64, error) {
	var (
		db         = rc.db
		totalCount int64
		err        error
	)

	db = db.Table("k_wallet_transaction").Model(&models.KWalletTransaction{}).WithContext(ctx)
	db = fn(db)

	err = db.Count(&totalCount).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get total count k-wallet transaction")
		return 0, err
	}

	return totalCount, nil
}

func (rc *RepositoryContext) GetLastIdKWalletRepository(ctx context.Context) (int64, error) {
	var (
		db      = rc.db
		kWallet *models.KWallet
		err     error
	)

	db = db.Table("k_wallet").Model(&models.KWallet{}).WithContext(ctx)
	db = db.Select("id").Order("id desc")

	err = db.First(&kWallet).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get last id k-wallet")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1, nil
		}

		return 0, err
	}

	return kWallet.ID + 1, nil
}

func (rc *RepositoryContext) InsertKWalletRepositoryTx(ctx context.Context, tx *gorm.DB, kWallet *models.KWallet) (*models.KWallet, error) {
	var (
		db  = rc.db
		err error
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("k_wallet").Model(&models.KWallet{}).WithContext(ctx)

	err = db.Create(kWallet).Error
	if err != nil {
		log.Error().Err(err).Msg("error query insert k-wallet")
		return nil, err
	}

	return kWallet, nil
}

func (rc *RepositoryContext) InsertBatchVirtualAccountRepositoryTx(ctx context.Context, tx *gorm.DB, va []*models.VirtualAccountKWallet) error {
	var (
		db  = rc.db
		err error
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("virtual_account_k_wallet").Model(&models.VirtualAccountKWallet{}).WithContext(ctx)

	err = db.Create(&va).Error
	if err != nil {
		log.Error().Err(err).Msg("error query insert batch virtual account")
		return err
	}

	return nil
}

func (rc *RepositoryContext) InsertKWalletTransactionRepositoryTx(ctx context.Context, tx *gorm.DB, kWalletTransaction *models.KWalletTransaction) error {
	var (
		db  = rc.db
		err error
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("k_wallet_transaction").Model(&models.KWalletTransaction{}).WithContext(ctx)

	err = db.Create(&kWalletTransaction).Error
	if err != nil {
		log.Error().Err(err).Msg("error query insert k-wallet transaction")
		return err
	}

	return nil
}

func (rc *RepositoryContext) UpdateKWalletRepositoryTx(ctx context.Context, tx *gorm.DB, fn func(db *gorm.DB) *gorm.DB) error {
	var (
		db  = rc.db
		err error
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("k_wallet").Model(&models.KWallet{}).WithContext(ctx)
	db = fn(db)

	err = db.Error
	if err != nil {
		log.Error().Err(err).Msg("error query update k-wallet")
		return err
	}

	return nil
}

func (rc *RepositoryContext) FindKWalletTransactionRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) ([]*models.KWalletTransaction, error) {
	var (
		db                  = rc.db
		kWalletTransactions []*models.KWalletTransaction
		err                 error
	)

	db = db.Table("k_wallet_transaction").Model(&models.KWalletTransaction{}).WithContext(ctx)
	db = fn(db)

	err = db.Find(&kWalletTransactions).Error
	if err != nil {
		log.Error().Err(err).Msg("error query find k-wallet transaction")
		return nil, err
	}

	return kWalletTransactions, nil
}

func (rc *RepositoryContext) FindVirtualAccountKWalletRepository(ctx context.Context, fn func(db *gorm.DB) *gorm.DB) ([]*models.VirtualAccountKWallet, error) {
	var (
		db                     = rc.db
		virtualAccountKWallets []*models.VirtualAccountKWallet
		err                    error
	)

	db = db.Table("virtual_account_k_wallet").Model(&models.VirtualAccountKWallet{}).WithContext(ctx)
	db = fn(db)

	err = db.Find(&virtualAccountKWallets).Error
	if err != nil {
		log.Error().Err(err).Msg("error query find virtual account k-wallet")
		return nil, err
	}

	return virtualAccountKWallets, nil
}

func (rc *RepositoryContext) FindMasDionVaCustPayBalRepository(ctx context.Context, f func(db *gorm.DB) *gorm.DB) ([]*models.MasDionVaCustPayBal, error) {
	var (
		db   = rc.db
		data []*models.MasDionVaCustPayBal
		err  error
	)

	db = db.Table("mas_dion_va_cust_pay_bal").Model(&models.MasDionVaCustPayBal{}).WithContext(ctx)
	db = f(db)

	err = db.Find(&data).Error
	if err != nil {
		log.Error().Err(err).Msg("error query find mas dion va cust pay bal")
		return nil, err
	}

	return data, nil
}

func (rc *RepositoryContext) UpdateMasDionVaCustPayBalRepository(ctx context.Context, tx *gorm.DB, f func(db *gorm.DB) *gorm.DB) error {
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("mas_dion_va_cust_pay_bal").Model(&models.MasDionVaCustPayBal{}).WithContext(ctx)

	db = f(db)

	err := db.Error
	if err != nil {
		log.Error().Err(err).Msg("error query update mas dion va cust pay bal")
		return err
	}

	return nil
}

func (rc *RepositoryContext) FindMasAmmarKWalletMemberRepository(ctx context.Context, f func(db *gorm.DB) *gorm.DB) ([]*models.MasAmmarKWalletMember, error) {
	var (
		db   = rc.db
		data []*models.MasAmmarKWalletMember
		err  error
	)

	db = db.Table("mas_ammar_k_wallet_member").Model(&models.MasAmmarKWalletMember{}).WithContext(ctx)
	db = f(db)

	err = db.Find(&data).Error
	if err != nil {
		log.Error().Err(err).Msg("error query find mas ammar k-wallet member")
		return nil, err
	}

	return data, nil
}

func (rc *RepositoryContext) GetMasAmmarKWalletMemberSaldoRepository(ctx context.Context, f func(db *gorm.DB) *gorm.DB) (*models.MasAmmarKWalletMemberSaldo, error) {
	var (
		db                 = rc.db
		kWalletMemberSaldo *models.MasAmmarKWalletMemberSaldo
		err                error
	)

	db = db.Table("mas_ammar_k_wallet_member_saldo").Model(&models.MasAmmarKWalletMemberSaldo{}).WithContext(ctx)
	db = f(db)

	err = db.First(&kWalletMemberSaldo).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get mas ammar k-wallet member saldo")
		return nil, err
	}

	return kWalletMemberSaldo, nil
}

func (rc *RepositoryContext) GetMasAmmarKWalletGenVaRepository(ctx context.Context, f func(db *gorm.DB) *gorm.DB) (*models.MasAmmarKWalletGenVa, error) {
	var (
		db    = rc.db
		genVa *models.MasAmmarKWalletGenVa
		err   error
	)

	db = db.Table("mas_ammar_k_wallet_gen_va").Model(&models.MasAmmarKWalletGenVa{}).WithContext(ctx)

	db = f(db)

	err = db.First(&genVa).Error
	if err != nil {
		log.Error().Err(err).Msg("error query get mas ammar k-wallet gen va")
		return nil, err
	}

	return genVa, nil
}

func (rc *RepositoryContext) FindMasAmmarKwalletGenVaDetailRepository(ctx context.Context, f func(db *gorm.DB) *gorm.DB) ([]*models.MasAmmarKWalletGenVaDetail, error) {
	var (
		db           = rc.db
		genVaDetails []*models.MasAmmarKWalletGenVaDetail
		err          error
	)

	db = db.Table("mas_ammar_k_wallet_gen_va_detail").Model(&models.MasAmmarKWalletGenVaDetail{}).WithContext(ctx)

	db = f(db)

	err = db.Find(&genVaDetails).Error
	if err != nil {
		log.Error().Err(err).Msg("error query find mas ammar k-wallet gen va detail")
		return nil, err
	}

	return genVaDetails, nil
}

func (rc *RepositoryContext) UpdateMasAmmarKWalletMemberRepository(ctx context.Context, tx *gorm.DB, f func(db *gorm.DB) *gorm.DB) error {
	var (
		db = rc.db
	)

	if tx != nil {
		db = tx
	}

	db = db.Table("mas_ammar_k_wallet_member").Model(&models.MasAmmarKWalletMember{}).WithContext(ctx)

	db = f(db)

	err := db.Error
	if err != nil {
		log.Error().Err(err).Msg("error query update mas ammar k-wallet member")
		return err
	}

	return nil
}
