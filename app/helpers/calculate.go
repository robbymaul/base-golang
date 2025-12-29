package helpers

import "github.com/shopspring/decimal"

func CalculatePercentage(amount, feeBasis10000 decimal.Decimal) (total, fee decimal.Decimal) {
	// product = amount * feeBasis10000
	product := amount.Mul(feeBasis10000)

	// divisor = 10000
	divisor := decimal.NewFromInt(10000)

	// fee = product / 10000 (pembagian biasa)
	fee = product.Div(divisor)

	// Cek apakah ada sisa (untuk ceiling / pembulatan ke atas)
	remainder := product.Mod(divisor)
	if remainder.GreaterThan(decimal.Zero) {
		// fee = fee + 1
		fee = fee.Add(decimal.NewFromInt(1))
	}

	total = amount.Add(fee)
	return total, fee
}
