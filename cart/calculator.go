package cart

import "fmt"

type Result struct {
	Valid          bool
	TotalAmount    float64
	TotalTaxAmount float64
	LineItems      []LineItem
}

func (u *UnknownTaxRate) Error() string {
	return fmt.Sprintf("tax rate is unknown: %f", u.Rate)
}

type TaxRates interface {
	// TaxableAmount calculates the tax amount of an inclusive tax price or returns UnknownTaxRate.
	TaxableAmount(rate float64, price float64) (float64, error)
}

type Discounter interface {
	// Apply decides whether to apply a discount to a LineItem.
	Apply(*LineItem)
}

type Calculator struct {
	taxRates  TaxRates
	discounts []Discounter
}

func NewCalculator(taxRates TaxRates, discounts []Discounter) *Calculator {
	return &Calculator{
		taxRates:  taxRates,
		discounts: discounts,
	}
}

func (c *Calculator) Calculate(items []LineItem) (*Result, error) {
	for i := 0; i < len(items); i++ {
		for _, d := range c.discounts {
			d.Apply(&items[i])
		}
	}

	var totalTaxAmount float64
	var totalAmount float64

	for _, li := range items {
		priceAmount := li.Price * float64(li.Quantity)
		if li.Discount.PercentageOff > 0 {
			priceAmount -= priceAmount * li.Discount.PercentageOff
		}

		amount, err := c.taxRates.TaxableAmount(li.TaxRate, priceAmount)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate tax amount for %q: %w", li.Description, err)
		}

		totalTaxAmount += amount
		totalAmount += priceAmount
	}

	return &Result{
		Valid:          true,
		TotalAmount:    totalAmount,
		TotalTaxAmount: totalTaxAmount,
		LineItems:      items,
	}, nil
}

type LineItem struct {
	Description string
	Quantity    int
	TaxRate     float64
	Price       float64 // not how you'd like to represent this, but it's a toy example
	Discount    Discount
}
