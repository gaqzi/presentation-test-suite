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
	// Amount calculates the tax amount of an inclusive tax price or returns UnknownTaxRate.
	Amount(rate float64, price float64) (float64, error)
}

type Calculator struct {
	taxRates TaxRates
}

func NewCalculator(taxRates TaxRates) *Calculator {
	return &Calculator{
		taxRates: taxRates,
	}
}

func (c *Calculator) Calculate(items []LineItem) (*Result, error) {
	var totalTaxAmount float64
	var totalAmount float64

	for _, li := range items {
		priceAmount := li.Price * float64(li.Quantity)
		amount, err := c.taxRates.Amount(li.TaxRate, priceAmount)
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
}
