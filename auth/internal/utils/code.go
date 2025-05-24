package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateVerificationCode() (string, error) {
	const codeLength = 6
	const digits = "0123456789"
	code := make([]byte, codeLength)

	for i := 0; i < codeLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[num.Int64()]
	}

	return string(code), nil
}
