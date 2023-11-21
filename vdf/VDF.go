package vdf

import "math"
import "math/big"


func VDF(version uint64, prevSuccessfulWrites uint64) *big.Int {
	vNextBig := big.NewInt(int64(version))

	p := new(big.Int)
	p.SetString(P, 16)

	Ns := int64(N_BASE) * int64(math.Sqrt(float64(prevSuccessfulWrites)))
	NsBig := big.NewInt(Ns)

	exponent := new(big.Int).Exp(big.NewInt(2), NsBig, nil)
	exponentiatedVNext := new(big.Int).Exp(vNextBig, exponent, nil)
	return new(big.Int).Exp(big.NewInt(G), exponentiatedVNext, p)
}