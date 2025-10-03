package auth_test

import (
	"github.com/roxensox/chirpy/internal/auth"
	"testing"
)

func TestMakeRefreshToken(t *testing.T) {
	_, err := auth.MakeRefreshToken()
	if err != nil {
		t.Errorf("Failed to make refresh token")
	}
}
