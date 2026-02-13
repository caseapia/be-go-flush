package invite

import (
	"crypto/rand"
	"math/big"
)

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateCode() (string, error) {
	codeLength := 6
	code := make([]byte, codeLength)

	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}

		code[i] = letters[n.Int64()]
	}

	return string(code), nil
}
