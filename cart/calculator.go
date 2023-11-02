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
	discounts []Discounter
}

func NewCalculator(discounts []Discounter) *Calculator {
	return &Calculator{
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
		totalTaxAmount += li.TaxableAmount()
		totalAmount += li.TotalPrice()
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
	Price       float64 // not how you'd like to represent this, but it's a toy example
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
