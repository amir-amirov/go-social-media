package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func Generate6DigitCode() (string, error) {
	max := big.NewInt(1000000) // 0 to 999999
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
