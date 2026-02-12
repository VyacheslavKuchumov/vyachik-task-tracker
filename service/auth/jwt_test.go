package auth

import (
	"net/http"
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
	t.Run("reads bearer token from authorization header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer abc123")

		token := getTokenFromRequest(req)
		if token != "abc123" {
			t.Fatalf("expected abc123, got %s", token)
		}
	})

	t.Run("falls back to auth cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{
			Name:  AuthCookieName,
			Value: "cookie-token",
		})

		token := getTokenFromRequest(req)
		if token != "cookie-token" {
			t.Fatalf("expected cookie-token, got %s", token)
		}
	})
}
