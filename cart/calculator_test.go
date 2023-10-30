package cart

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculator_Calculate(t *testing.T) {
	for _, tc := range []struct {
		name        string
		expectError bool
		items       []LineItem
		totalPrice  float64
	}{
		{
			name:        "Sums to 0 with an empty cart",
			expectError: false,
			items: []LineItem{
				{},
			},
			totalPrice: 0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calc := NewCalculator()

			result, err := calc.Calculate(tc.items)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.True(t, result.Valid)
				require.Equal(t, tc.totalPrice, result.TotalPrice)
			}
		})
	}
}
