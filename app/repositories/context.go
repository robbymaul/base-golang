package repositories

import (
	"context"
	"errors"
	"fmt"
	"paymentserviceklink/app/client/espay"
	"paymentserviceklink/app/client/midtrans"
	clientsenangpay "paymentserviceklink/app/client/senangpay"
	"paymentserviceklink/app/models"
	"paymentserviceklink/app/strategy"
	"paymentserviceklink/config"
	"paymentserviceklink/pkg/pagination"
	"strings"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type RepositoryContext struct {
	ctx context.Context
	db  *gorm.DB
	cfg *config.Config
	*Adapter
}

func (rc *RepositoryContext) ReleaseTx(tx *gorm.DB, err error) error {
	if err != nil {
		if errRollback := tx.Rollback().Error; errRollback != nil {
			//panic(fmt.Errorf("failed to rollback database transaction \n > [%v]", errRollback))
			return errRollback
		}
	} else {
		if errCommit := tx.Commit().Error; errCommit != nil {
			//panic(fmt.Errorf("failed to commit database transaction \n > [%v]", errCommit))
			return errCommit
		}
	}

	return nil
}

type transactionFn func(tx *gorm.DB) error

func (rc *RepositoryContext) WithTransaction(callback transactionFn) error {
	var err error
	tx := rc.db.Begin()
	defer func() {
		if err2 := rc.ReleaseTx(tx, err); err2 != nil {
			log.Error().Msg(fmt.Sprintf("failed transaction database with error = [%v]", err2))
		}
	}()

	err = callback(tx)
	if err != nil {
		return err
	}

	return nil
}

func (rc *RepositoryContext) SetConfigurationPayment(configurations *models.Configuration) {
	log.Debug().Interface("configurations", configurations).Msg("set configuration payment")
	rc.Senangpay = clientsenangpay.NewSenangpay(rc.HttpClient, configurations)
	log.Debug().Interface("senangpay", rc.Senangpay).Msg("new senangpay")
	rc.Midtrans = midtrans.NewMidtrans(rc.HttpClient, configurations)
	log.Debug().Interface("midtrans", rc.Midtrans).Msg("new midtrans")
	rc.Espay = espay.NewEspay(rc.HttpClient, configurations)
	log.Debug().Interface("espay", rc.Espay).Msg("new espay")
	rc.Strategy = strategy.NewStrategy(rc.Senangpay, rc.Midtrans, rc.Espay)
}

func (rc *RepositoryContext) SearchQuery(filters []*pagination.Filter, joinOperator string) (string, []interface{}) {
	var queryParts []string
	var args []interface{}

	if len(filters) < 1 {
		return "", nil
	}

	for _, filter := range filters {
		var part string
		//operator := ""
		//
		//if idx > 0 {
		//	operator = joinOperator
		//}

		//searchQuery += rc.SearchQuery(filter, pages.JoinOperator)
		switch filter.Operator {
		case "eq":
			//searchQuery += fmt.Sprintf(" %v %v = '%v' ", operator, filter.ID, filter.Value)
			part = fmt.Sprintf("%s = ?", filter.ID)
		case "ne":
			//searchQuery += fmt.Sprintf(" %v %v != '%v' ", operator, filter.ID, filter.Value)
			part = fmt.Sprintf("%s != ?", filter.ID)
		case "like":
			//searchQuery += fmt.Sprintf(" %v %v ILIKE '%v' ", operator, filter.ID, fmt.Sprintf("%%%v%%", filter.Value))
			part = fmt.Sprintf("%s ILIKE ?", filter.ID)
		case "notLike":
			//searchQuery += fmt.Sprintf(" %v %v NOT ILIKE '%v' ", operator, filter.ID, fmt.Sprintf("%%%v%%", filter.Value))
			part = fmt.Sprintf("%s NOT ILIKE ?", filter.ID)
		case "gt":
			//searchQuery += fmt.Sprintf(" %v %v > '%v' ", operator, filter.ID, filter.Value)
			part = fmt.Sprintf("%s > ?", filter.ID)
		case "gte":
			//searchQuery += fmt.Sprintf(" %v %v >= '%v' ", operator, filter.ID, filter.Value)
			part = fmt.Sprintf("%s >= ?", filter.ID)
		case "lt":
			//searchQuery += fmt.Sprintf(" %v %v < '%v' ", operator, filter.ID, filter.Value)
			part = fmt.Sprintf("%s < ?", filter.ID)
		case "lte":
			//searchQuery += fmt.Sprintf(" %v %v <= '%v' ", operator, filter.ID, filter.Value)
			part = fmt.Sprintf("%s <= ?", filter.ID)
		default:
			continue
		}

		queryParts = append(queryParts, filter.Table+"."+part)
		args = append(args, filter.Value)
	}

	separator := " " + joinOperator + " "
	query := strings.Join(queryParts, separator)

	log.Debug().Interface("query search", query).Interface("args", args).Msg("search query")
	return query, args
}

func (rc *RepositoryContext) IsRetryableError(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	if err != nil && (rc.contains(err.Error(), "deadlock detected")) {
		return true
	}

	if err != nil && (rc.contains(err.Error(), "duplicate key value")) {
		return true
	}

	return false
}

func (rc *RepositoryContext) contains(errMsg string, substring string) bool {
	return len(errMsg) >= len(substring) && (func() bool {
		for i := 0; i+len(substring) <= len(errMsg); i++ {
			match := true
			for j := range substring {
				if errMsg[i+j] != substring[j] {
					match = false
					break
				}

				if match {
					return true
				}
			}
		}

		return false
	}())
}
