package cart_test

import (
	"testing"

	"github.com/gaqzi/presentation-test-suite/cart"

	"github.com/stretchr/testify/require"
)

func totals(totalAmount, totalTaxAmount float64) func(*testing.T, *cart.Result) {
	return func(t *testing.T, result *cart.Result) {
		t.Helper()
		require.True(t, result.Valid)
		require.Equal(t, totalAmount, result.TotalAmount)
		require.Equal(t, totalTaxAmount, result.TotalTaxAmount)
	}
}

func noDiscounts() []cart.Discounter {
	return nil
}

func overpricedBanana(changes ...func(i *cart.LineItem)) cart.LineItem {
	item := cart.LineItem{
		Description: "Overpriced Banana",
		Quantity:    1,
		Price:       1,
		TaxRate:     taxRateFood,
		Discount:    cart.Discount{},
	}

	for _, change := range changes {
		change(&item)
	}

	return item
}

// ----- Builder -----
type calculatorBuilder struct {
	discounts []cart.Discounter
}

func (c calculatorBuilder) WithDiscounts(ds []cart.Discounter) calculatorBuilder {
	c.discounts = ds
	return c
}

func calculator() calculatorBuilder {
	return calculatorBuilder{
		discounts: nil,
	}
}

// ----- /Builder -----

// ----- Option -----
type calculatorOptionBuilder struct {
	discounts cart.Discounts
	totaler   cart.LineItemTotalser
}

type calculatorOption func(c *calculatorOptionBuilder)

func withDiscounts(ds []cart.Discounter) calculatorOption {
	return func(c *calculatorOptionBuilder) {
		c.discounts = ds
	}
}

func defaultCalculator(options ...calculatorOption) *cart.Calculator {
	calc := calculatorOptionBuilder{
		discounts: nil,
		totaler:   &cart.LineItemsTotaler{},
	}

	for _, opt := range options {
		opt(&calc)
	}

	return cart.NewCalculator(calc.discounts, calc.totaler)
}

// ----- /Option -----

var (
	taxRateFood = cart.TaxRate{Add: 0.12, Remove: 0.1071}
)

// func TestCalculator_Calculate(t *testing.T) {
// 	for _, tc := range []struct {
// 		name         string
// 		calculator   *cart.Calculator
// 		items        []cart.LineItem
// 		expectResult func(t *testing.T, result *cart.Result)
// 	}{
// 		{
// 			name:       "Sums to 0 with an empty cart",
// 			calculator: defaultCalculator(),
// 			items: []cart.LineItem{
// 				{},
// 			},
// 			expectResult: totals(0, 0),
// 		},
// 		{
// 			name:       "Calculate an item with a tax rate",
// 			calculator: defaultCalculator(),
// 			items: []cart.LineItem{
// 				{
// 					Description: "Overpriced Banana",
// 					Quantity:    1,
// 					TaxRate:     taxRateFood,
// 					Price:       1,
// 				},
// 			},
// 			expectResult: totals(1, 0.1071),
// 		},
// 		{
// 			name:       "Calculate an item where quantity is not 1",
// 			calculator: defaultCalculator(),
// 			items: []cart.LineItem{
// 				{
// 					Description: "Overpriced Banana",
// 					Quantity:    2,
// 					TaxRate:     taxRateFood,
// 					Price:       1,
// 				},
// 			},
// 			expectResult: totals(2, 0.2142),
// 		},
// 		{
// 			name: "Calculates with discounts applied",
// 			calculator: defaultCalculator(
// 				withDiscounts([]cart.Discounter{
// 					cart.NewDiscountForItem(
// 						"Ripe Banana",
// 						cart.Discount{"Expiring soon", 0.2},
// 					),
// 				}),
// 			),
// 			items: []cart.LineItem{
// 				{
// 					Description: "Ripe Banana",
// 					Quantity:    1,
// 					TaxRate:     taxRateFood,
// 					Price:       1,
// 				},
// 			},
// 			expectResult: totals(0.8, 0.08568),
// 		},
// 	} {
// 		t.Run(tc.name, func(t *testing.T) {
// 			result := tc.calculator.Calculate(tc.items)
//
// 			tc.expectResult(t, result)
// 		})
// 	}
// }

func TestCalculator_Calculate(t *testing.T) {
	mockDiscounts := cartmock.NewMockLineItemApplicator()
}

func TestLineItem_TotalPrice(t *testing.T) {
	for _, tc := range []struct {
		description string
		item        cart.LineItem
		expected    float64
	}{
		{
			description: "An item with price costs nothing",
			item: cart.LineItem{
				Description: "Air",
				Quantity:    1,
				Price:       0,
			},
			expected: 0,
		},
		{
			description: "A single item with a price sums up to that price",
			item: cart.LineItem{
				Description: "Overpriced Banana",
				Quantity:    1,
				Price:       1,
			},
			expected: 1,
		},
		{
			description: "An item with quantity of 2 doubles the single price",
			item: cart.LineItem{
				Description: "Overpriced Banana",
				Quantity:    2,
				Price:       1,
			},
			expected: 2,
		},
		{
			description: "An item will apply it",
			item: cart.LineItem{
				Description: "Ripe Banana",
				Quantity:    1,
				Price:       1,
				Discount: cart.Discount{
					Description:   "Expiring soon",
					PercentageOff: 0.2,
				},
			},
			expected: 0.8,
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			require.Equal(t, tc.expected, tc.item.TotalPrice())
		})
	}
}

func TestLineItem_TaxableAmount(t *testing.T) {
	for _, tc := range []struct {
		description string
		item        cart.LineItem
		expected    float64
	}{
		{
			description: "Tax of 0 on remove returns 0",
			item: cart.LineItem{
				Description: "Overpriced Banana",
				TaxRate:     cart.TaxRate{0, 0},
				Price:       1,
			},
			expected: 0,
		},
		{
			description: "Tax of 10% on remove returns that amount",
			item: cart.LineItem{
				Description: "Overpriced Banana",
				TaxRate:     cart.TaxRate{0, 0.1},
				Quantity:    1,
				Price:       1,
			},
			expected: 0.1,
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			require.Equal(t, tc.expected, tc.item.TaxableAmount())
		})
	}
}

func TestDiscounts_Apply(t *testing.T) {
	for _, tc := range []struct {
		description string
		discounts   cart.Discounts
		items       []cart.LineItem
		expected    []cart.LineItem
	}{
		{
			description: "When no discounts apply it doesn't add one to the items",
			discounts:   nil,
			items:       []cart.LineItem{overpricedBanana()},
			expected:    []cart.LineItem{overpricedBanana()},
		},
		{
			description: "Only apply a discount to the matching item",
			discounts: cart.Discounts{
				cart.NewDiscountForItem("Ripe Banana", cart.Discount{
					Description:   "Expiring soon",
					PercentageOff: 0.2,
				}),
			},
			items: []cart.LineItem{
				overpricedBanana(),
				overpricedBanana(func(i *cart.LineItem) {
					i.Description = "Ripe Banana"
				}),
			},
			expected: []cart.LineItem{
				overpricedBanana(),
				overpricedBanana(func(i *cart.LineItem) {
					i.Description = "Ripe Banana"
					// The discount wasn't part of the input items but were added
					i.Discount = cart.Discount{
						Description:   "Expiring soon",
						PercentageOff: 0.2,
					}
				}),
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			actual := tc.discounts.Apply(tc.items)

			require.Equal(t, tc.expected, actual)
		})
	}
}

func TestLineItems_Totals(t *testing.T) {
	for _, tc := range []struct {
		description string
		items       []cart.LineItem
		expected    cart.Result
	}{
		{
			description: "No line items sum to nothing",
			items:       nil,
			expected: cart.Result{
				TotalAmount:    0,
				TotalTaxAmount: 0,
			},
		},
		{
			description: "A single item returns its price and taxable amount",
			items: []cart.LineItem{
				overpricedBanana(),
			},
			expected: cart.Result{
				TotalAmount:    1,
				TotalTaxAmount: 0.1071,
			},
		},
		{
			description: "Multiple items are summed",
			items: []cart.LineItem{
				overpricedBanana(),
				overpricedBanana(func(i *cart.LineItem) {
					i.Description = "Green Banana"
					i.Price = 0.5
				}),
			},
			expected: cart.Result{
				TotalAmount: 1.5,
				// why you shouldn't use floats for money without a proper strategy
				TotalTaxAmount: 0.16065000000000002,
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			totaler := cart.LineItemsTotaler{}

			actual := totaler.Totals(tc.items)

			require.Equal(t, tc.expected, actual)
		})
	}
}
