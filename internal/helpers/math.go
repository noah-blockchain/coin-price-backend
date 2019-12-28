package helpers

import (
	"math/big"
)

const precision = 100

// default amount of qNoahs in 1 Noah
var qNoahInNoah = big.NewFloat(1000000000000000000)
var feeDefaultMultiplier = big.NewInt(1000000000000000)

// default amount of unit in one noah
const unitInNoah = 1000

func QNoahStr2Noah(value string) string {
	if value == "" {
		return "0"
	}

	floatValue, err := new(big.Float).SetPrec(500).SetString(value)
	CheckErrBool(err)

	return new(big.Float).SetPrec(500).Quo(floatValue, qNoahInNoah).Text('f', 18)
}
