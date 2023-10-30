package cart

type Result struct {
	Valid      bool
	TotalPrice float64
	LineItems  []LineItem
}

type Calculator struct {
}

func NewCalculator() *Calculator {
	return &Calculator{}
}

func (c *Calculator) Calculate(items []LineItem) (*Result, error) {
	return &Result{
		Valid:      true,
		TotalPrice: 0,
		LineItems:  items,
	}, nil
}

type LineItem struct {
	Description string
	Quantity    int
	TaxRate     float64
	Price       float64 // not how you'd like to represent this, but it's a toy example
}
