package user

import (
	"VyacheslavKuchumov/test-backend/config"
	"VyacheslavKuchumov/test-backend/service/auth"
	"VyacheslavKuchumov/test-backend/types"
	"VyacheslavKuchumov/test-backend/utils"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	token, err := h.createSessionToken(payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	auth.SetAuthCookie(w, token)
	utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err := h.registerUser(payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) HandleWebLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/login?error="+url.QueryEscape("invalid form payload"), http.StatusSeeOther)
		return
	}

	token, err := h.createSessionToken(types.LoginUserPayload{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	})
	if err != nil {
		http.Redirect(w, r, "/login?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	auth.SetAuthCookie(w, token)
	http.Redirect(w, r, "/goals", http.StatusSeeOther)
}

func (h *Handler) HandleWebRegister(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/register?error="+url.QueryEscape("invalid form payload"), http.StatusSeeOther)
		return
	}

	err := h.registerUser(types.RegisterUserPayload{
		FirstName: r.FormValue("firstName"),
		LastName:  r.FormValue("lastName"),
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
	})
	if err != nil {
		http.Redirect(w, r, "/register?error="+url.QueryEscape(err.Error()), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/login?ok="+url.QueryEscape("registration successful"), http.StatusSeeOther)
}

func (h *Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	auth.ClearAuthCookie(w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *Handler) registerUser(payload types.RegisterUserPayload) error {
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return fmt.Errorf("Invalid payload %v", errors)
	}

	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		return fmt.Errorf("User with email %s already exists", payload.Email)
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return err
	}

	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) createSessionToken(payload types.LoginUserPayload) (string, error) {
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		return "", fmt.Errorf("Invalid payload %v", errors)
	}

	u, err := h.store.GetUserByEmail(payload.Email)
	if err != nil {
		return "", fmt.Errorf("User not found, invalid email or password")
	}

	if !auth.ComparePasswords(u.Password, payload.Password) {
		return "", fmt.Errorf("User not found, invalid email or password")
	}
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
