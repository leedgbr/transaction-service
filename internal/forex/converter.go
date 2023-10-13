package forex

import "math/big"

// Converter is responsible for currency conversion given an amount and an exchange rate
type Converter struct{}

// Convert performs the exchange rate calculation as accurately as possible and rounds to the nearest cent.
func (c *Converter) Convert(amount int, exchangeRate float64) int {
	rate := big.NewFloat(exchangeRate)
	originalAmount := big.NewFloat(float64(amount))

	targetCurrencyAmount := &big.Float{}
	targetCurrencyAmount.Mul(originalAmount, rate)

	rounded := roundToNearestBigInt(targetCurrencyAmount)
	return int(rounded.Int64())
}

// roundToNearestBigInt rounds the supplied big.Float to the nearest big.Int
func roundToNearestBigInt(value *big.Float) *big.Int {
	newValue := &big.Float{}
	newValue.Copy(value)
	delta := 0.5
	if newValue.Sign() < 0 {
		delta = -0.5
	}
	newValue.Add(newValue, big.NewFloat(delta))
	rounded, _ := newValue.Int(nil)
	return rounded
}
