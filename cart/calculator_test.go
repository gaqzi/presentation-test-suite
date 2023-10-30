package cart

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculator_Calculate(t *testing.T) {
	for _, tc := range []struct {
		name           string
		expectError    bool
		items          []LineItem
		totalAmount    float64
		totalTaxAmount float64
	}{
		{
			name: "Sums to 0 with an empty cart",
			items: []LineItem{
				{},
			},
			totalAmount: 0,
		},
		{
			name: "Calculate an item with a tax rate",
			items: []LineItem{
				{
					Description: "Overpriced Banana",
					Quantity:    1,
					TaxRate:     0.12,
					Price:       1,
				},
			},
			totalAmount:    1,
			totalTaxAmount: 0.1071,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			taxRates := NewStaticTaxRates(
				TaxRate(0.25, 0.20),
				TaxRate(0.12, 0.1071),
				TaxRate(0.6, 5.66),
				TaxRate(0, 0),
			)
			calc := NewCalculator(taxRates)

			result, err := calc.Calculate(tc.items)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.True(t, result.Valid)
				require.Equal(t, tc.totalAmount, result.TotalAmount)
				require.Equal(t, tc.totalTaxAmount, result.TotalTaxAmount)
			}
		})
	}
}
