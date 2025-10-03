package auth_test

import (
	"github.com/roxensox/chirpy/internal/auth"
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	// Slice of test cases
	test_cases := []struct {
		password1 string
		password2 string
		expected  bool
	}{
		{
			password1: "testpass",
			password2: "testpass",
			expected:  true,
		},
		{
			password1: "testpass",
			password2: "Testpass",
			expected:  false,
		},
		{
			password1: "",
			password2: "",
			expected:  true,
		},
		{
			password1: "",
			password2: " ",
			expected:  false,
		},
		{
			password1: "password",
			password2: "passwordandsomemore",
			expected:  false,
		},
	}

	// Loops through the test cases
	for _, tc := range test_cases {
		hash, err := auth.HashPassword(tc.password2)
		if err != nil {
			t.Errorf("Failed to hash password %s", tc.password2)
		}
		res, err := auth.CheckPasswordHash(tc.password1, hash)
		if err != nil {
			t.Errorf("Failed to compare password %s and hash", tc.password1)
		}
		if res != tc.expected {
			t.Errorf("Comparing %s and %s -> expected %v got %v", tc.password1, tc.password2, tc.expected, res)
		}
	}
}
