package helpers

import "math/big"

func NewFloat(x float64, precision uint) *big.Float {
	return big.NewFloat(x).SetPrec(precision)
}
