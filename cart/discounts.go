package cart

type Discount struct {
	Description   string
	PercentageOff float64
}

type ItemDiscount struct {
	ItemDescription string
	Discount        Discount
}

func (i *ItemDiscount) Apply(li *LineItem) {
	if li.Description != i.ItemDescription {
		return
	}

	li.Discount = i.Discount
}

func NewDiscountForItem(itemDescription string, discount Discount) *ItemDiscount {
	return &ItemDiscount{
		ItemDescription: itemDescription,
		Discount:        discount,
	}
}
