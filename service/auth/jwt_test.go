package auth

import (
	"net/http/httptest"
	"testing"
)

func TestCreateJWT(t *testing.T) {
	secret := []byte("secret")

	token, err := CreateJWT(secret, 1)
	if err != nil {
		t.Errorf("error creating JWT: %v", err)
	}

	if token == "" {
		t.Error("expected token to be not empty")
	}
}

func TestGetTokenFromRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer abc123")

	token := getTokenFromRequest(req)
	if token != "abc123" {
		t.Fatalf("expected abc123, got %s", token)
	}
}
