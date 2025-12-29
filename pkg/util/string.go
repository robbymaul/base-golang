package util

import "strings"

func StringWithoutSpace(s string) string {
	temp := strings.Split(strings.Trim(s, " "), " ")

	return strings.Join(temp, "")
}

func StringToBoolean(s string) bool {
	return s == "true"
}
