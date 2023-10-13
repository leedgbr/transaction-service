package forex_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"transaction-service/internal/forex"
)

func TestConverter(t *testing.T) {
	converter := &forex.Converter{}

	tcs := []struct {
		name         string
		amount       int
		exchangeRate float64
		wantAmount   int
	}{
		{
			name:         "zero",
			amount:       0,
			exchangeRate: 0,
			wantAmount:   0,
		},
		{
			name:         "small amount",
			amount:       1,
			exchangeRate: 0,
			wantAmount:   0,
		},
		{
			name:         "small exchange rate",
			amount:       0,
			exchangeRate: 1,
			wantAmount:   0,
		},
		{
			name:         "negative amount",
			amount:       -1,
			exchangeRate: 1,
			wantAmount:   -1,
		},
		{
			name:         "big - convert a billion dollars with an exchange rate way higher than that of iranian rial",
			amount:       1000000000,
			exchangeRate: 1000000,
			wantAmount:   1000000000000000,
		},
		{
			name:         "big negative - convert a billion dollars with an exchange rate way higher than that of iranian rial",
			amount:       -1000000000,
			exchangeRate: 1000000,
			wantAmount:   -1000000000000000,
		},
		{
			name:         "rounding - up",
			amount:       25,
			exchangeRate: 0.75,
			wantAmount:   19,
		},
		{
			name:         "rounding - up",
			amount:       25,
			exchangeRate: 0.74,
			wantAmount:   19,
		},
		{
			name:         "rounding - up",
			amount:       25,
			exchangeRate: 0.745,
			wantAmount:   19,
		},
		{
			name:         "rounding - down",
			amount:       25,
			exchangeRate: 0.73,
			wantAmount:   18,
		},
		{
			name:         "rounding - down",
			amount:       25,
			exchangeRate: 0.725,
			wantAmount:   18,
		},
		{
			name:         "rounding - lots of decimal places",
			amount:       25,
			exchangeRate: 0.77777777777777777777777,
			wantAmount:   19,
		},
	}
	for _, tc := range tcs {
		t.Run(fmt.Sprintf("name: %s, amount: %d, exchangeRate: %f", tc.name, tc.amount, tc.exchangeRate),
			func(t *testing.T) {
				result := converter.Convert(tc.amount, tc.exchangeRate)
				assert.Equal(t, tc.wantAmount, result)
			},
		)
	}
}
