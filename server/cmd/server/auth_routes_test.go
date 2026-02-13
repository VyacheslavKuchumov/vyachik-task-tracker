package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProtectedEndpointsRequireAuthorization(t *testing.T) {
	srv := NewServer(":0", nil)
	handler := srv.router()

	protectedCases := []struct {
		name   string
		method string
		path   string
		body   []byte
	}{
		{name: "swagger ui", method: http.MethodGet, path: "/swagger/index.html"},
		{name: "list goals", method: http.MethodGet, path: "/api/v1/goals"},
		{name: "create goal", method: http.MethodPost, path: "/api/v1/goals", body: []byte(`{}`)},
		{name: "update goal", method: http.MethodPut, path: "/api/v1/goals/1", body: []byte(`{}`)},
		{name: "delete goal", method: http.MethodDelete, path: "/api/v1/goals/1"},
		{name: "create task", method: http.MethodPost, path: "/api/v1/goals/1/tasks", body: []byte(`{}`)},
		{name: "assigned tasks", method: http.MethodGet, path: "/api/v1/tasks/assigned"},
		{name: "update task", method: http.MethodPut, path: "/api/v1/tasks/1", body: []byte(`{}`)},
		{name: "delete task", method: http.MethodDelete, path: "/api/v1/tasks/1"},
		{name: "assign task", method: http.MethodPut, path: "/api/v1/tasks/1/assign", body: []byte(`{}`)},
	}

	for _, tc := range protectedCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, bytes.NewBuffer(tc.body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusForbidden {
				t.Fatalf("expected %d, got %d", http.StatusForbidden, rr.Code)
			}
		})
	}
}

func TestLoginAndRegisterArePublic(t *testing.T) {
	srv := NewServer(":0", nil)
	handler := srv.router()

	cases := []struct {
		name string
		path string
	}{
		{name: "login", path: "/api/v1/login"},
		{name: "register", path: "/api/v1/register"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tc.path, bytes.NewBufferString(`{}`))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code == http.StatusForbidden {
				t.Fatalf("expected non-%d for public endpoint, got %d", http.StatusForbidden, rr.Code)
			}
		})
	}
}
