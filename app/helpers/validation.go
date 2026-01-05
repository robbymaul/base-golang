package helpers

import (
	"slices"
)

func IsInList[T comparable](item T, list ...T) (bool, []T) {
	if len(list) == 0 {
		return true, list
	}
	return slices.Contains(list, item), list
}
