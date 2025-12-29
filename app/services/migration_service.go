package services

import (
	"context"
	"errors"
	"fmt"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/repositories"
	"paymentserviceklink/config"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type MigrationService struct {
	service *Service
}

func NewMigration(ctx context.Context, rc *repositories.RepositoryContext, cfg *config.Config) *MigrationService {
	return &MigrationService{
		service: NewService(ctx, rc, cfg),
	}
}

func (s *MigrationService) MigrasiKWalletMasDion() error {
	var err error
	// batch size
	lastId := 0
	batchSize := 10
	//offSet := 0
	totalExecute := 0

	for {
		var data []*models.MasDionVaCustPayBal

		data, err = s.service.repository.FindMasDionVaCustPayBalRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Limit(batchSize)
			if lastId > 0 {
				db = db.Where("id > ?", lastId)
			}
			db = db.Where("migration = ?", false)
			db = db.Order("id asc")

			return db
		})

		if len(data) == 0 {
			fmt.Println("migrasi k-wallet mas dion selesai")
			break
		}
		err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
			for _, d := range data {
				log.Info().Interface("memberId", d.Dfno).Interface("gen_va", d.Novac).Msg("migrasi k-wallet mas dion")
				orderId := time.Now().Unix()
				orderString := fmt.Sprint(orderId)
				noRekening := fmt.Sprintf("%s%08s", orderString, strconv.Itoa(d.Id))

				status := enums.KWalletStatusInactive
				isActive := false

				if d.Status == "1" {
					status = enums.KWalletStatusActive
					isActive = true
				}

				kWallet := &models.KWallet{
					MemberID:   d.Dfno,
					FullName:   d.Fullnm,
					NoRekening: noRekening,
					GenVA:      d.Novac,
					Balance:    d.Amount,
					Currency:   enums.CURRENCY_IDR,
					Symbol:     enums.SYMBOL_CURRENCY_IDR,
					Status:     status,
					IsActive:   isActive,
					VirtualAccount: []*models.VirtualAccountKWallet{
						{
							VirtualAccount: d.Novac,
							Bank:           "",
							BankCode:       "",
						},
					},
				}

				_, err = s.service.repository.InsertKWalletRepositoryTx(s.service.ctx, tx, kWallet)
				if err != nil {
					log.Error().Err(err).Msg("insert k-wallet repository error")
					return err
				}

				err = s.service.repository.UpdateMasDionVaCustPayBalRepository(s.service.ctx, tx, func(db *gorm.DB) *gorm.DB {
					db = db.Where("id = ?", d.Id)
					db = db.Update("migration", true)

					return db
				})
				if err != nil {
					log.Error().Err(err).Msg("update mas dion va cust pay bal repository error")
					return err
				}
				totalExecute += len(data)
			}
			return nil
		})
		if err != nil {
			log.Error().Err(err).Msg("error query transaction")
			log.Printf("migarasi k-wallet mas dion gagal pada: %v", err)
			return err
		}
		lastId = data[len(data)-1].Id
		fmt.Printf("%d k-wallet mas dion selesai, last-id: %v", totalExecute, lastId)
		//offSet += batchSize
	}

	return nil
}

func (s *MigrationService) MigrasiKWalletMasAmmar() error {
	var err error
	// batch size
	lastId := 0
	batchSize := 10
	offSet := 0
	totalExecute := 0

	timeNow := time.Now()

	for {
		var data []*models.MasAmmarKWalletMember

		data, err = s.service.repository.FindMasAmmarKWalletMemberRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
			db = db.Limit(batchSize)
			if lastId > 0 {
				db = db.Where("rec_id > ?", lastId)
			}
			db = db.Where("id_member LIKE 'ID%' OR id_member LIKE 'EI%' OR id_member LIKE 'ET%'")
			db = db.Where("status = '1'")
			db = db.Where("nama NOT LIKE '%JHON DOE%'")
			db = db.Where("migration = ?", false)
			db = db.Order("rec_id asc")

			return db
		})

		fmt.Printf("%d k-wallet mas amar selesai, last-id: %v", totalExecute, lastId)

		if len(data) == 0 {
			fmt.Println("migrasi k-wallet mas dion selesai")
			break
		}

		err = s.service.repository.WithTransaction(func(tx *gorm.DB) error {
			for _, d := range data {
				log.Info().Interface("memberId", d.IdMember).Msg("migrasi k-wallet mas ammar")

				var saldoMember *models.MasAmmarKWalletMemberSaldo

				saldoMember, err = s.service.repository.GetMasAmmarKWalletMemberSaldoRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
					db = db.Where("id_member = ?", d.IdMember)

					return db
				})
				if err != nil {
					log.Error().Err(err).Msg("error query get mas ammar k-wallet member saldo")
					return err
				}

				var kWallet *models.KWallet
				kWallet, err = s.service.repository.GetKWalletRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
					db = db.Where("member_id = ? and currency = ?", d.IdMember, enums.CURRENCY_IDR)
					return db
				})
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					log.Error().Err(err).Msg("error get k-wallet repository")
					return err
				}

				if kWallet == nil {
					var genVa *models.MasAmmarKWalletGenVa

					genVa, err = s.service.repository.GetMasAmmarKWalletGenVaRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
						db = db.Where("member_id = ?", d.IdMember)

						return db
					})
					if err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							continue
						}
						log.Error().Err(err).Msg("error query get mas ammar k-wallet gen va")

						return err
					}

					var genVaDetail []*models.MasAmmarKWalletGenVaDetail

					genVaDetail, err = s.service.repository.FindMasAmmarKwalletGenVaDetailRepository(s.service.ctx, func(db *gorm.DB) *gorm.DB {
						db = db.Where("id_gen_va = ?", genVa.RecId)
						db = db.Where("va_number != '' or va_number != NULL")

						return db
					})
					if err != nil {
						log.Error().Err(err).Msg("error query get mas ammar k-wallet gen va detail")
						return err
					}

					virtualAccounts := make([]*models.VirtualAccountKWallet, 0)

					for _, vaDetail := range genVaDetail {
						virtualAccounts = append(virtualAccounts, &models.VirtualAccountKWallet{
							VirtualAccount: vaDetail.VaNumber,
							Bank:           "",
							BankCode:       vaDetail.BankCode,
						})
					}

					idGenVa := fmt.Sprintf("%08s", fmt.Sprint(genVa.RecId))

					orderId := time.Now().Unix()
					orderString := fmt.Sprint(orderId)
					noRekening := fmt.Sprintf("%s%s", orderString, idGenVa)

					status := enums.KWalletStatusActive

					kWallet = &models.KWallet{
						MemberID:       d.IdMember,
						FullName:       d.Nama,
						NoRekening:     noRekening,
						GenVA:          idGenVa,
						Balance:        saldoMember.LastSaldo,
						Currency:       enums.CURRENCY_IDR,
						Symbol:         enums.SYMBOL_CURRENCY_IDR,
						Status:         status,
						IsActive:       true,
						VirtualAccount: virtualAccounts,
					}

					_, err = s.service.repository.InsertKWalletRepositoryTx(s.service.ctx, tx, kWallet)
					if err != nil {
						log.Error().Err(err).Msg("insert k-wallet repository error")
						return err
					}
				} else {
					kWallet.AddBalance(saldoMember.LastSaldo)

					err = s.service.repository.UpdateKWalletRepositoryTx(s.service.ctx, tx, func(db *gorm.DB) *gorm.DB {
						db = db.Where("id = ?", kWallet.ID)
						db = db.UpdateColumns(map[string]interface{}{
							"balance":    kWallet.Balance,
							"updated_at": time.Now(),
						})
						return db
					})
					if err != nil {
						log.Error().Err(err).Msg("update k-wallet repository error")
						return err
					}

					month := timeNow.Month()
					year := time.Now().Year()

					kWalletTransaction := &models.KWalletTransaction{
						KWalletID:                kWallet.ID,
						KWalletTypeTransactionID: 1,
						PaymentID:                fmt.Sprint(d.RecId),
						Title:                    "migrasi k-wallet mas ammar",
						PaymentCode:              fmt.Sprint(time.Now().Unix()),
						TransactionCode:          fmt.Sprint(time.Now().Unix()),
						TransactionType:          "migrasi saldo",
						Direction:                enums.K_WALLET_DIRECTION_IN,
						CounterpartyName:         "",
						CounterpartyBank:         "",
						PaymentChannel:           "K-WALLET-MAS-AMMAR",
						Description:              "penambahan saldo dari data migrasi mas ammar",
						Balance:                  kWallet.Balance,
						Debit:                    decimal.NewFromInt(0),
						Credit:                   saldoMember.LastSaldo,
						Amount:                   saldoMember.LastSaldo,
						Currency:                 enums.CURRENCY_IDR,
						Symbol:                   enums.SYMBOL_CURRENCY_IDR,
						Status:                   "success",
						Month:                    int64(month),
						Year:                     int64(year),
						Date:                     timeNow,
						Time:                     timeNow.Format(time.TimeOnly),
						DateTime:                 timeNow,
					}

					err = s.service.repository.InsertKWalletTransactionRepositoryTx(s.service.ctx, tx, kWalletTransaction)
					if err != nil {
						log.Error().Err(err).Msg("insert k-wallet transaction repository error")
						return err
					}
				}

				err = s.service.repository.UpdateMasAmmarKWalletMemberRepository(s.service.ctx, tx, func(db *gorm.DB) *gorm.DB {
					db = db.Where("rec_id = ?", d.RecId)
					db = db.Update("migration", true)

					return db
				})
				if err != nil {
					log.Error().Err(err).Msg("update mas amar k-wallet mas ammar repository error")
					return err
				}

				totalExecute += len(data)
			}
			return nil
		})
		if err != nil {
			log.Error().Err(err).Msg("error query transaction")
			log.Printf("migarasi k-wallet mas amar gagal pada: %v", err)
			return err
		}

		lastId = data[len(data)-1].RecId
		fmt.Printf("%d k-wallet mas amar selesai, last-id: %v", totalExecute, lastId)
		fmt.Printf("%d k-wallet mas amar selesai", totalExecute)
		offSet += batchSize
	}

	return nil
}
