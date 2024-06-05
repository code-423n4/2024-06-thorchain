package common

import (
	"crypto/rand"
	"math/big"
)

const hexBytes = "abcdefABCDEF0123456789"

// RandHexString generates random hex string used for test purpose
func RandHexString(n int) string {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(hexBytes))))
		if err != nil {
			return ""
		}
		ret[i] = hexBytes[num.Int64()]
	}

	return string(ret)
}
