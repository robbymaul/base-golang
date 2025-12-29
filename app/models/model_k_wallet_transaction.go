package models

import (
	"errors"
	"fmt"
	"paymentserviceklink/app/enums"
	"paymentserviceklink/pkg/pagination"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

type KWalletTransaction struct {
	ID                       int64                  `gorm:"column:id;primaryKey;notNull;autoIncrement"`
	KWalletID                int64                  `gorm:"column:k_wallet_id;foreignKey"`
	KWalletTypeTransactionID int64                  `gorm:"k_wallet_type_transaction_id;foreignKey"`
	PaymentID                string                 `gorm:"column:payment_id"`
	Title                    string                 `gorm:"column:title"`
	PaymentCode              string                 `gorm:"column:payment_code"`
	TransactionCode          string                 `gorm:"column:transaction_code"`
	TransactionType          string                 `gorm:"column:transaction_type"`
	Direction                enums.KWalletDirection `gorm:"column:direction"`
	CounterpartyName         string                 `gorm:"column:counterparty_name"`
	CounterpartyBank         enums.Channel          `gorm:"column:counterparty_bank"`
	PaymentChannel           string                 `gorm:"column:payment_channel"`
	Description              string                 `gorm:"column:description"`
	Balance                  decimal.Decimal        `gorm:"column:balance"`
	Debit                    decimal.Decimal        `gorm:"column:debit"`
	Credit                   decimal.Decimal        `gorm:"column:credit"`
	Amount                   decimal.Decimal        `gorm:"column:amount"`
	Currency                 enums.Currency         `gorm:"column:currency"`
	Symbol                   enums.SymbolCurrency   `gorm:"column:symbol"`
	Status                   enums.PaymentStatus    `gorm:"column:status"`
	Month                    int64                  `gorm:"column:month"`
	Year                     int64                  `gorm:"column:year"`
	Date                     time.Time              `gorm:"column:date"`
	Time                     string                 `gorm:"column:time"`
	DateTime                 time.Time              `gorm:"column:datetime"`
	BaseField
}

func (*KWalletTransaction) TableName() string {
	return "k_wallet_transaction"
}

func (k *KWalletTransaction) TopupTransaction(balance decimal.Decimal, amount decimal.Decimal) {
	k.Title = "Topup"
	timeNow := time.Now()
	k.Month = int64(timeNow.Month())
	k.Year = int64(timeNow.Year())
	k.Date = timeNow
	k.Time = timeNow.Format(time.TimeOnly)
	k.DateTime = timeNow
	k.Direction = enums.K_WALLET_DIRECTION_IN
	k.Debit = decimal.NewFromInt(0)
	k.Credit = amount
	k.Amount = amount
	k.Balance = balance
}

func (k *KWalletTransaction) PaymentTransaction(balance decimal.Decimal, amount decimal.Decimal) {
	log.Debug().Interface("balance", balance).Interface("amount", amount).Interface("transaction balance", k.Balance).Msg("before = payment transaction")
	k.Title = "Pembayaran"
	timeNow := time.Now()
	k.Month = int64(timeNow.Month())
	k.Year = int64(timeNow.Year())
	k.Date = timeNow
	k.Time = timeNow.Format(time.TimeOnly)
	k.DateTime = timeNow
	k.Direction = enums.K_WALLET_DIRECTION_OUT
	k.Debit = amount
	k.Credit = decimal.NewFromInt(0)
	k.Amount = amount
	k.Balance = balance
	log.Debug().Interface("balance", balance).Interface("amount", amount).Interface("transaction balance", k.Balance).Msg("after = payment transaction")
}

func AllowedFilterColumnKWalletTransaction() map[string]FilterColumn {
	return map[string]FilterColumn{
		"no_rekening": {
			Operator: []string{"eq"},
			Variant:  "string",
			Table:    "k_wallet",
		},
		"member_id": {
			Operator: []string{"eq"},
			Variant:  "string",
			Table:    "k_wallet",
		},
		"payment_id": {
			Operator: []string{"eq", "like"},
			Variant:  "string",
			Table:    "k_wallet_transaction",
		},
		"year": {
			Operator: []string{"eq"},
			Variant:  "number",
			Table:    "k_wallet_transaction",
		},
		"month": {
			Operator: []string{"eq"},
			Variant:  "number",
			Table:    "k_wallet_transaction",
		},
		"date": {
			Operator: []string{"gt", "gte", "lt", "lte"},
			Variant:  "date",
			Table:    "k_wallet_transaction",
		},
		"datetime": {
			Operator: []string{"gt", "gte", "lt", "lte"},
			Variant:  "datetime",
			Table:    "k_wallet_transaction",
		},
	}
}

func TransformFilterColumnKWalletTransaction(filters []*pagination.Filter, filterColumn map[string]FilterColumn) ([]*pagination.Filter, error) {
	if len(filters) == 0 {
		return filters, nil
	}

	var err error
	filterData := make([]*pagination.Filter, 0)

	var fromDate time.Time
	var toDate time.Time

	for idx, filter := range filters {
		if filter.ID == "from_date" {

			if filter.Value == "" {
				return nil, errors.New("filter column from date is required value")
			} else {
				fromDate, err = time.Parse(time.DateOnly, filter.Value.(string))
				if err != nil {
					return nil, fmt.Errorf("filter column from date format %v", time.DateOnly)
				}
			}

			filter.ID = "date"
			filter.Value = fromDate.Format(time.DateOnly)
		}

		if filter.ID == "to_date" {
			if filter.Value == "" {
				return nil, errors.New("filter column to date is required value")
			} else {
				toDate, err = time.Parse(time.DateOnly, filter.Value.(string))
				if err != nil {
					return nil, fmt.Errorf("filter column to date format %v", time.DateOnly)
				}
			}

			filter.ID = "date"
			filter.Value = toDate.Format(time.DateOnly)
		}

		log.Debug().Interface("index", idx).Interface("filter.id", filter.ID).Interface("value", filter.Value).Msg("loop filter data transform")
		filterData = append(filters, filter)
	}

	if !fromDate.Before(toDate) {
		return nil, errors.New("filter column to date must be greater than from date")
	}

	maxDuration := 30 * 24 * time.Hour

	validationDuration := ValidationDuration(fromDate, toDate, maxDuration)
	log.Debug().Interface("validationDuration", validationDuration).Msg("validation duration")
	if !validationDuration {
		return nil, errors.New("filter column to date can only be up to 30 days")
	}

	log.Debug().Interface("filterData", filterData).Msg("filter data")

	return filterData, err
}

func ValidationDuration(fromDate time.Time, toDate time.Time, duration time.Duration) bool {
	log.Debug().Interface("fromDate", fromDate).Interface("toDate", toDate).Interface("toDate.Sub(fromDate)", toDate.Sub(fromDate)).Interface("duration", duration).Msg("validation duration")

	return toDate.Sub(fromDate) <= duration
}
