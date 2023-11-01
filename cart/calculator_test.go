package cart_test

import (
	"testing"

	"github.com/gaqzi/presentation-test-suite/cart"

	"github.com/stretchr/testify/require"
)

func noError(t *testing.T, err error) {
	t.Helper()
	require.NoError(t, err, "uh-oh, we got an unexpected error!")
}

func totals(totalAmount, totalTaxAmount float64) func(*testing.T, *cart.Result) {
	return func(t *testing.T, result *cart.Result) {
		t.Helper()
		require.True(t, result.Valid)
		require.Equal(t, totalAmount, result.TotalAmount)
		require.Equal(t, totalTaxAmount, result.TotalTaxAmount)
	}
}

func swedishTaxRates() cart.TaxRates {
	return cart.NewStaticTaxRates(
		cart.TaxRate(0.25, 0.20),
		cart.TaxRate(0.12, 0.1071),
		cart.TaxRate(0.06, 0.566),
		cart.TaxRate(0, 0),
	)
}

func noDiscounts() []cart.Discounter {
	return nil
}

func TestCalculator_Calculate(t *testing.T) {
	for _, tc := range []struct {
		name         string
		taxRates     func() cart.TaxRates
		discounts    func() []cart.Discounter
		expectError  func(t *testing.T, err error)
		expectResult func(t *testing.T, result *cart.Result)
		items        []cart.LineItem
	}{
		{
			name:      "Sums to 0 with an empty cart",
			taxRates:  swedishTaxRates,
			discounts: noDiscounts,
			items: []cart.LineItem{
				{},
			},
			expectError:  noError,
			expectResult: totals(0, 0),
		},
		{
			name:      "Calculate an item with a tax rate",
			taxRates:  swedishTaxRates,
			discounts: noDiscounts,
			items: []cart.LineItem{
				{
					Description: "Overpriced Banana",
					Quantity:    1,
					TaxRate:     0.12,
					Price:       1,
				},
			},
			expectError:  noError,
			expectResult: totals(1, 0.1071),
		},
		{
			name:      "Calculate an item where quantity is not 1",
			taxRates:  swedishTaxRates,
			discounts: noDiscounts,
			items: []cart.LineItem{
				{
					Description: "Overpriced Banana",
					Quantity:    2,
					TaxRate:     0.12,
					Price:       1,
				},
			},
			expectError:  noError,
			expectResult: totals(2, 0.2142),
		},
		{
			name: "Stops calculating when there's an invalid tax rate",
			taxRates: func() cart.TaxRates {
				return cart.NewStaticTaxRates() // no tax rate is ever valid
			},
			discounts: noDiscounts,
			items: []cart.LineItem{
				{
					Description: "Invalid Banana",
					Quantity:    1,
					TaxRate:     0.66,
					Price:       1,
				},
			},
			expectError: func(t *testing.T, err error) {
				require.ErrorIs(t, err, &cart.UnknownTaxRate{0.66}, "expected an unknown tax rate as the error")
			},
			expectResult: func(t *testing.T, result *cart.Result) {
				require.Empty(t, result)
			},
		},
		{
			name:     "Calculates with discounts applied",
			taxRates: swedishTaxRates,
			discounts: func() []cart.Discounter {
				return []cart.Discounter{
					cart.NewDiscountForItem(
						"Ripe Banana",
						cart.Discount{"Expiring soon", 0.2},
					),
				}
			},
			items: []cart.LineItem{
				{
					Description: "Ripe Banana",
					Quantity:    1,
					TaxRate:     0.12,
					Price:       1,
				},
			},
			expectError:  noError,
			expectResult: totals(0.8, 0.08568),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calc := cart.NewCalculator(tc.taxRates(), tc.discounts())

			result, err := calc.Calculate(tc.items)

			tc.expectError(t, err)
			tc.expectResult(t, result)
		})
	}
}
