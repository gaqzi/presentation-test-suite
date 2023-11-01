package cart_test

import (
	"testing"

	"github.com/gaqzi/presentation-test-suite/cart"

	"github.com/stretchr/testify/require"
)

func TestCalculator_Calculate(t *testing.T) {
	for _, tc := range []struct {
		name           string
		expectError    bool
		items          []cart.LineItem
		totalAmount    float64
		totalTaxAmount float64
	}{
		{
			name: "Sums to 0 with an empty cart",
			items: []cart.LineItem{
				{},
			},
			totalAmount: 0,
		},
		{
			name: "Calculate an item with a tax rate",
			items: []cart.LineItem{
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
		{
			name: "Calculate an item where quantity is not 1",
			items: []cart.LineItem{
				{
					Description: "Overpriced Banana",
					Quantity:    2,
					TaxRate:     0.12,
					Price:       1,
				},
			},
			totalAmount:    2,
			totalTaxAmount: 0.2142,
		},
		{
			name: "Stops calculating when there's an invalid tax rate",
			items: []cart.LineItem{
				{
					Description: "Invalid Banana",
					Quantity:    1,
					TaxRate:     0.66,
					Price:       1,
				},
			},
			expectError: true,
		},
		{
			name: "Calculates with discounts applied",
			items: []cart.LineItem{
				{
					Description: "Ripe Banana",
					Quantity:    1,
					TaxRate:     0.12,
					Price:       1,
				},
			},
			totalAmount:    0.8,
			totalTaxAmount: 0.08568,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			taxRates := cart.NewStaticTaxRates(
				cart.TaxRate(0.25, 0.20),
				cart.TaxRate(0.12, 0.1071),
				cart.TaxRate(0.06, 0.566),
				cart.TaxRate(0, 0),
			)
			discountRules := cart.NewDiscountForItem(
				"Ripe Banana",
				cart.Discount{"Expiring soon", 0.2},
			)
			calc := cart.NewCalculator(taxRates, []cart.Discounter{discountRules})

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
