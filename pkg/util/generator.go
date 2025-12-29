package util

import gonanoid "github.com/matoous/go-nanoid/v2"

const (
	ALPHA_NUMERIC_LOW_UPPER_CHAR_SET = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	ALPHA_NUMERIC_SET                = "0123456789"
)

func GenerateRandomString(i int) string {
	alpha := gonanoid.MustGenerate(ALPHA_NUMERIC_LOW_UPPER_CHAR_SET, i)
	return alpha
}

func GenerateRandomAlphaNumericString(i int) string {
	alpha := gonanoid.MustGenerate(ALPHA_NUMERIC_SET, i)
	return alpha
}
