package auth_test

import (
	"github.com/google/uuid"
	"github.com/roxensox/chirpy/internal/auth"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	test_UUID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to generate UUID")
		return
	}
	duration := time.Millisecond * 200
	if err != nil {
		t.Errorf("Failed to generate duration")
		return
	}
	test_cases := []struct {
		UID         uuid.UUID
		tokenSecret string
		compSecret  string
		expiresIn   time.Duration
		sleepFor    time.Duration
		expected    bool
	}{
		{
			UID:         test_UUID,
			tokenSecret: "Pippin",
			compSecret:  "Pippin",
			expiresIn:   duration,
			sleepFor:    time.Millisecond * 0,
			expected:    true,
		},
		{
			UID:         test_UUID,
			tokenSecret: "Luna",
			compSecret:  "Luna",
			expiresIn:   duration,
			sleepFor:    duration + (time.Millisecond * 100),
			expected:    false,
		},
		{
			UID:         test_UUID,
			tokenSecret: "Pippin",
			compSecret:  "Luna",
			expiresIn:   duration,
			sleepFor:    time.Millisecond * 0,
			expected:    false,
		},
	}

	for _, tc := range test_cases {
		jwt_str, err := auth.MakeJWT(tc.UID, tc.tokenSecret, tc.expiresIn)
		if err != nil {
			t.Errorf("Failed at MakeJWT:\n\tError: %v", err)
			return
		}
		time.Sleep(tc.sleepFor)
		UID, err := auth.ValidateJWT(jwt_str, tc.compSecret)
		if err != nil {
			if tc.tokenSecret == tc.compSecret && tc.sleepFor < duration {
				t.Errorf("Failed at ValidateJWT:\n\tError: %v", err)
				return
			}
		}
		if (UID == tc.UID) != tc.expected {
			t.Errorf("UUIDs don't match")
			return
		}
	}
}
