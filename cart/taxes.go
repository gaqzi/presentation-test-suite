package cart

type TaxRate struct {
	// Add is the percentage to add to an untaxed amount to make it inclusive of tax.
	Add float64
	// Remove is the percentage to remove from an inclusive tax amount to get the amount that is tax.
	Remove float64
}
