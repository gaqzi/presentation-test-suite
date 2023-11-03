package cart

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

type LineItemApplicator interface {
	Apply(items []LineItem) []LineItem
}

type LineItemTotalser interface {
	Totals(items []LineItem) Result
}

type Calculator struct {
	discounts LineItemApplicator
	totaler   LineItemTotalser
}

func NewCalculator(discounts LineItemApplicator, totaler LineItemTotalser) *Calculator {
	return &Calculator{
		discounts: discounts,
		totaler:   totaler,
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

func (c *Calculator) Calculate(items []LineItem) *Result {
	items = c.discounts.Apply(items)
	totals := c.totaler.Totals(items)

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

type LineItemsTotaler struct{}

func (lit *LineItemsTotaler) Totals(items []LineItem) Result {
	var result Result

	for _, i := range items {
		result.TotalAmount += i.TotalPrice()
		result.TotalTaxAmount += i.TaxableAmount()
	}

	return result
}
