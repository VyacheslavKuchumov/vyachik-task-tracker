package user

import (
	"VyacheslavKuchumov/test-backend/service/auth"
	"VyacheslavKuchumov/test-backend/types"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should fail if the user payload is invalid", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "invalid_user",
			LastName:  "fsdsd",
			Email:     "fdgdfg",
			Password:  "asdfd",
		}
		marshalled, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.HandleFunc("/register", handler.HandleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should create a user correctly", func(t *testing.T) {
		payload := types.RegisterUserPayload{
			FirstName: "correct_user",
			LastName:  "fsdsd",
			Email:     "test123@gmail.com",
			Password:  "asdfd",
		}
		marshaled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshaled))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.HandleFunc("/register", handler.HandleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("should fail login when user does not exist", func(t *testing.T) {
		payload := types.LoginUserPayload{
			Email:    "missing@example.com",
			Password: "password",
		}
		marshaled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshaled))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.HandleFunc("/login", handler.HandleLogin)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should login correctly", func(t *testing.T) {
		hash, err := auth.HashPassword("asdfd")
		if err != nil {
			t.Fatal(err)
		}
		userStore.userByEmail["test123@gmail.com"] = &types.User{
			ID:       1,
			Email:    "test123@gmail.com",
			Password: hash,
		}

		payload := types.LoginUserPayload{
			Email:    "test123@gmail.com",
			Password: "asdfd",
		}
		marshaled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshaled))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.HandleFunc("/login", handler.HandleLogin)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

type mockUserStore struct {
	userByEmail map[string]*types.User
}

func (m *mockUserStore) ensure() {
	if m.userByEmail == nil {
		m.userByEmail = map[string]*types.User{}
	}
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	m.ensure()
	u, ok := m.userByEmail[email]
	if !ok {
		return nil, fmt.Errorf("User doesn't exist")
	}
	return u, nil
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}

func (m *mockUserStore) CreateUser(types.User) error {
	return nil
}
