package helpers

import (
	"fmt"
	"paymentserviceklink/app/models"
	"paymentserviceklink/pkg/pagination"
	"slices"
	"time"

	"github.com/rs/zerolog/log"
)

func IsInList[T comparable](item T, list ...T) (bool, []T) {
	if len(list) == 0 {
		return true, list
	}
	return slices.Contains(list, item), list
}

func FilterColumnValidation(filter []*pagination.Filter, filterColumn map[string]models.FilterColumn) ([]*pagination.Filter, error) {
	for _, f := range filter {
		value, ok := filterColumn[f.ID]
		if !ok {
			return nil, fmt.Errorf(fmt.Sprintf("filter column id %s does not exist or not allowed", f.ID))
		}

		if !slices.Contains(filterColumn[f.ID].Operator, f.Operator) {
			return nil, fmt.Errorf("filter column id %s cannot use operator %v", f.ID, f.Operator)
		}

		if filterColumn[f.ID].Variant != f.Variant {
			return nil, fmt.Errorf("filter column id %s variant is %v not suitable filter variant %v", f.ID, filterColumn[f.ID].Variant, f.Variant)
		}

		if filterColumn[f.ID].IsRequired && f.Value == "" {
			return nil, fmt.Errorf("filter column id %s is required", f.ID)
		}

		f.Table = value.Table
	}

	return filter, nil
}

func ValidationDuration(fromDate, toDate time.Time, maxDuration time.Duration) (string, bool) {
	diff := toDate.Sub(fromDate)
	daysSent := int(diff.Hours() / 24)
	maxDays := int(maxDuration.Hours() / 24)

	message := fmt.Sprintf(
		"The date range is only allowed up to %d days, but the duration you provided is %d days.",
		maxDays, daysSent,
	)

	if daysSent < 1 {
		return message, false
	}

	isValid := daysSent <= maxDays

	log.Debug().
		Time("fromDate", fromDate).
		Time("toDate", toDate).
		Int("hari_dikirim", daysSent).
		Int("batas_maksimal_hari", maxDays).
		Bool("valid", isValid).
		Msg("validation duration")

	return message, isValid
}
