package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	nums := make([]byte, 32)
	_, err := rand.Read(nums)
	if err != nil {
		fmt.Println("Failed to generate number")
		return "", err
	}

	outStr := hex.EncodeToString(nums)
	return outStr, nil
}
