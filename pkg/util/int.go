package util

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseBoolFallback(i any, fallback bool) bool {
	result, ok := parseBool(i)
	if !ok {
		result = false
	}

	return result
}

func parseBool(i any) (result bool, ok bool) {
	ok = true

	switch v := i.(type) {
	case bool:
		result = ok
	case string:
		result = strings.ToLower(v) == "true"
	default:
		ok = false
	}

	return result, ok
}

func ParseIntFallback(i any, fallback int) int {
	result, ok := parseInt(i)
	if !ok {
		return fallback
	}

	return result
}

func parseInt(i any) (result int, ok bool) {

	ok = true

	switch v := i.(type) {
	case int:
		result = v
	case string:
		var err error
		parse, err := strconv.Atoi(v)
		if err != nil {
			ok = false
		}
		result = parse
	default:
		ok = false
	}

	return
}

func ParseStringFallback(i any, fallback string) string {
	result, ok := parseString(i)
	if !ok {
		return fallback
	}

	return result
}

func parseString(i any) (result string, ok bool) {
	ok = true

	switch v := i.(type) {
	case string:
		result = v
	case nil:
		ok = false
	default:
		str := fmt.Sprintf("%v", i)
		result = str
	}

	return
}
