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

// ----- Builder -----
type calculatorBuilder struct {
	taxRates  cart.TaxRates
	discounts []cart.Discounter
}

func (c calculatorBuilder) WithTaxRates(tx cart.TaxRates) calculatorBuilder {
	c.taxRates = tx
	return c
}

func (c calculatorBuilder) WithDiscounts(ds []cart.Discounter) calculatorBuilder {
	c.discounts = ds
	return c
}

func (c calculatorBuilder) Build() *cart.Calculator {
	return cart.NewCalculator(c.taxRates, c.discounts)
}

func calculator() calculatorBuilder {
	return calculatorBuilder{
		taxRates:  swedishTaxRates(),
		discounts: nil,
	}
}

// ----- /Builder -----

// ----- Option -----
type calculatorOptionBuilder struct {
	taxRates  cart.TaxRates
	discounts []cart.Discounter
}

type calculatorOption func(c *calculatorOptionBuilder)

func withTaxRates(tx cart.TaxRates) calculatorOption {
	return func(c *calculatorOptionBuilder) {
		c.taxRates = tx
	}
}

func withDiscounts(ds []cart.Discounter) calculatorOption {
	return func(c *calculatorOptionBuilder) {
		c.discounts = ds
	}
}

func defaultCalculator(options ...calculatorOption) *cart.Calculator {
	calc := calculatorOptionBuilder{
		taxRates:  swedishTaxRates(),
		discounts: nil,
	}

	for _, opt := range options {
		opt(&calc)
	}

	return cart.NewCalculator(calc.taxRates, calc.discounts)
}

// ----- /Option -----

func TestCalculator_Calculate(t *testing.T) {
	for _, tc := range []struct {
		name         string
		calculator   *cart.Calculator
		items        []cart.LineItem
		expectError  func(t *testing.T, err error)
		expectResult func(t *testing.T, result *cart.Result)
	}{
		{
			name:       "Sums to 0 with an empty cart",
			calculator: defaultCalculator(),
			items: []cart.LineItem{
				{},
			},
			expectError:  noError,
			expectResult: totals(0, 0),
		},
		{
			name:       "Calculate an item with a tax rate",
			calculator: defaultCalculator(),
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
			name:       "Calculate an item where quantity is not 1",
			calculator: defaultCalculator(),
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
			// Strategy (options) pattern
			calculator: defaultCalculator(
				withTaxRates(cart.NewStaticTaxRates()), // no tax rate is ever valid
			),
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
			name: "Calculates with discounts applied",
			calculator: defaultCalculator(
				withDiscounts([]cart.Discounter{
					cart.NewDiscountForItem(
						"Ripe Banana",
						cart.Discount{"Expiring soon", 0.2},
					),
				}),
			),
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
			result, err := tc.calculator.Calculate(tc.items)

			tc.expectError(t, err)
			tc.expectResult(t, result)
		})
	}
}
