package xordistance

import (
	"fmt"
	"math/big"
)

func XORBytes(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length of byte slices is not equivalent: %d != %d", len(a), len(b))
	}
	buf := make([]byte, len(a))
	for i := range a {
		buf[i] = a[i] ^ b[i]
	}
	return buf, nil
}

func BytesToInt(input_bytes []byte) *big.Int {
	z := new(big.Int)
	z.SetBytes(input_bytes)
	return z
}
