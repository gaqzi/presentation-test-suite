package cart

import (
	"errors"
	"fmt"
)

type UnknownTaxRate struct {
	Rate float64
}

func (u *UnknownTaxRate) Error() string {
	return fmt.Sprintf("tax rate is unknown: %f", u.Rate)
}

func (u *UnknownTaxRate) Is(target error) bool {
	var t *UnknownTaxRate
	if !errors.As(target, &t) {
		return false
	}

	return t.Rate == u.Rate
}

type StaticTaxRates struct {
	lookup map[float64]float64
}

// TaxableAmount calculates the inclusive amount of tax in price from rate.
func (s *StaticTaxRates) TaxableAmount(rate float64, price float64) (float64, error) {
	backwardRate, ok := s.lookup[rate]
	if !ok {
		return 0, &UnknownTaxRate{Rate: rate}
	}

	return price * backwardRate, nil
}

type TaxRateAdder func(*StaticTaxRates)

// TaxRate configures the calculation of a tax rate.
// add is the tax rate described when it's *added* to a price.
// remove is the tax rate when it's *removed* from an inclusive price.
func TaxRate(add, remove float64) TaxRateAdder {
	return func(rates *StaticTaxRates) {
		rates.lookup[add] = remove
	}
}

func NewStaticTaxRates(rates ...TaxRateAdder) *StaticTaxRates {
	tr := &StaticTaxRates{
		lookup: map[float64]float64{},
	}

	for _, r := range rates {
		r(tr)
	}

	return tr
}
