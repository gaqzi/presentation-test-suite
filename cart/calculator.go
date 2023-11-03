package cart

import (
	"fmt"
)

type Result struct {
	Valid          bool
	TotalAmount    float64
	TotalTaxAmount float64
	LineItems      []LineItem
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
	discounts Discounts
}

func NewCalculator(discounts []Discounter) *Calculator {
	return &Calculator{
		discounts: discounts,
	}
}

type Discounts []Discounter

func (d Discounts) Apply(items []LineItem) []LineItem {
	var returns []LineItem

	for _, item := range items {
		for _, d := range d {
			d.Apply(&item)
		}
		returns = append(returns, item)
	}

	return returns
}

func (c *Calculator) Calculate(items LineItems) *Result {
	items = c.discounts.Apply(items)
	totals := items.Totals()

	return &Result{
		Valid:          true,
		TotalAmount:    totals.TotalAmount,
		TotalTaxAmount: totals.TotalTaxAmount,
		LineItems:      items,
	}
}

type LineItem struct {
	Description string
	Quantity    int
	Price       float64 // inclusive of tax, and probably shouldn't do float64 for money in the real world
	TaxRate     TaxRate
	Discount    Discount
}

func (i *LineItem) TotalPrice() float64 {
	total := i.Price * float64(i.Quantity)
	if i.Discount.PercentageOff > 0 {
		total -= total * i.Discount.PercentageOff
	}

	return total
}

// TaxableAmount calculates the inclusive amount of tax in price using the current tax rate.
func (i *LineItem) TaxableAmount() float64 {
	return i.TotalPrice() * i.TaxRate.Remove
}

type LineItems []LineItem

func (li LineItems) Totals() Result {
	var result Result

	for _, i := range li {
		result.TotalAmount += i.TotalPrice()
		result.TotalTaxAmount += i.TaxableAmount()
	}

	return result
}
