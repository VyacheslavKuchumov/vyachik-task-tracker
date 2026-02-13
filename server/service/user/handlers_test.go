package user

import (
	"VyacheslavKuchumov/test-backend/service/auth"
	"VyacheslavKuchumov/test-backend/types"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{userByEmail: map[string]*types.User{}}
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

		cookies := rr.Result().Cookies()
		if len(cookies) == 0 || cookies[0].Name != auth.AuthCookieName {
			t.Errorf("Expected auth cookie %q to be set", auth.AuthCookieName)
		}
	})

	t.Run("should return profile for authenticated user", func(t *testing.T) {
		userStore.userByID[1] = &types.User{
			ID:        1,
			FirstName: "Test",
			LastName:  "User",
			Email:     "test123@gmail.com",
		}

		req := httptest.NewRequest(http.MethodGet, "/profile", nil)
		req = req.WithContext(context.WithValue(req.Context(), auth.UserKey, 1))
		rr := httptest.NewRecorder()

		handler.HandleGetProfile(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should update profile", func(t *testing.T) {
		userStore.userByID[1] = &types.User{
			ID:        1,
			FirstName: "Old",
			LastName:  "Name",
			Email:     "test123@gmail.com",
		}

		payload := types.UpdateProfilePayload{
			FirstName: "New",
			LastName:  "Name",
		}
		marshaled, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(marshaled))
		req = req.WithContext(context.WithValue(req.Context(), auth.UserKey, 1))
		rr := httptest.NewRecorder()

		handler.HandleUpdateProfile(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should update password", func(t *testing.T) {
		hash, err := auth.HashPassword("oldpass")
		if err != nil {
			t.Fatal(err)
		}
		userStore.userByID[1] = &types.User{
			ID:        1,
			FirstName: "Test",
			LastName:  "User",
			Email:     "test123@gmail.com",
			Password:  hash,
		}

		payload := types.UpdatePasswordPayload{
			CurrentPassword: "oldpass",
			NewPassword:     "newpass",
		}
		marshaled, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/profile/password", bytes.NewBuffer(marshaled))
		req = req.WithContext(context.WithValue(req.Context(), auth.UserKey, 1))
		rr := httptest.NewRecorder()

		handler.HandleUpdatePassword(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rr.Code)
		}
	})

}

type mockUserStore struct {
	userByEmail map[string]*types.User
	userByID    map[int]*types.User
}

func (m *mockUserStore) ensure() {
	if m.userByEmail == nil {
		m.userByEmail = map[string]*types.User{}
	}
	if m.userByID == nil {
		m.userByID = map[int]*types.User{}
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
	m.ensure()
	u, ok := m.userByID[id]
	if !ok {
		return nil, fmt.Errorf("User doesn't exist")
	}
	return u, nil
}

func (m *mockUserStore) CreateUser(user types.User) error {
	m.ensure()
	if user.ID == 0 {
		user.ID = len(m.userByID) + 1
	}
	copyUser := user
	m.userByID[copyUser.ID] = &copyUser
	m.userByEmail[copyUser.Email] = &copyUser
	return nil
}

func (m *mockUserStore) UpdateUserProfile(userID int, payload types.UpdateProfilePayload) (*types.User, error) {
	m.ensure()
	u, ok := m.userByID[userID]
	if !ok {
		return nil, fmt.Errorf("User doesn't exist")
	}

	u.FirstName = payload.FirstName
	u.LastName = payload.LastName
	return u, nil
}

func (m *mockUserStore) UpdateUserPassword(userID int, hashedPassword string) error {
	m.ensure()
	u, ok := m.userByID[userID]
	if !ok {
		return fmt.Errorf("User doesn't exist")
	}
	u.Password = hashedPassword
	return nil
}

func (m *mockUserStore) ListUsers() ([]*types.UserLookup, error) {
	m.ensure()
	users := make([]*types.UserLookup, 0, len(m.userByID))
	for _, u := range m.userByID {
		users = append(users, &types.UserLookup{
			ID:   u.ID,
			Name: u.FirstName + " " + u.LastName,
		})
	}
	return users, nil
}
