package auth_test

import (
	"github.com/google/uuid"
	"github.com/roxensox/chirpy/internal/auth"
	"net/http"
	"testing"
	"time"
)

func TestJWT(t *testing.T) {
	test_UUID, err := uuid.NewUUID()
	if err != nil {
		t.Errorf("Failed to generate UUID")
		return
	}

	duration := time.Second * 30

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
			expected:    true,
		},
		{
			UID:         test_UUID,
			tokenSecret: "Luna",
			compSecret:  "Luna",
			expiresIn:   time.Millisecond * 0,
			expected:    false,
		},
		{
			UID:         test_UUID,
			tokenSecret: "Pippin",
			compSecret:  "Luna",
			expiresIn:   duration,
			expected:    false,
		},
	}

	for i, tc := range test_cases {
		jwt_str, err := auth.MakeJWT(tc.UID, tc.tokenSecret, tc.expiresIn)
		if err != nil {
			t.Errorf("Failed at MakeJWT:\n\tError: %v", err)
			return
		}
		UID, err := auth.ValidateJWT(jwt_str, tc.compSecret)
		if err != nil {
			if tc.tokenSecret == tc.compSecret && tc.expiresIn > 0 {
				t.Errorf("Case #%d failed at ValidateJWT:\n\tError: %v", i, err)
				return
			}
		}
		if (UID == tc.UID) != tc.expected {
			t.Errorf("UUIDs don't match")
			return
		}
	}
}

func TestGetBearerToken(t *testing.T) {
	testUID, _ := uuid.NewUUID()
	testjwt, _ := auth.MakeJWT(testUID, "Pippin", 5*time.Second)
	test_cases := []struct {
		token    string
		header   http.Header
		expected bool
	}{
		// Test with matching key
		{
			token: testjwt,
			header: http.Header{
				"Authorization": []string{testjwt},
			},
			expected: true,
		},
		// Test with multiple keys
		{
			token: "test",
			header: http.Header{
				"Authorization": []string{"test", "tset"},
			},
			expected: true,
		},
		// Test with multiple keys where match isn't first
		{
			token: "test",
			header: http.Header{
				"Authorization": []string{"tset", "test"},
			},
			expected: false,
		},
		// Test with mismatching key
		{
			token: "test",
			header: http.Header{
				"Authorization": []string{"schucks"},
			},
			expected: false,
		},
	}

	for _, h := range test_cases {
		tkn, err := auth.GetBearerToken(h.header)
		if err != nil {
			t.Errorf("Unable to get token")
			return
		}
		if (tkn == h.token) != h.expected {
			if h.expected {
				t.Errorf("Token mismatch: %s != %s", tkn, h.token)
			} else {
				t.Errorf("Tokens shouldn't match: %s == %s", tkn, h.token)
			}
		}
	}
}
