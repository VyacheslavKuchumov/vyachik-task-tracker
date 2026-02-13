package user

import (
	"VyacheslavKuchumov/test-backend/config"
	"VyacheslavKuchumov/test-backend/service/auth"
	"VyacheslavKuchumov/test-backend/types"
	"VyacheslavKuchumov/test-backend/utils"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

// HandleLogin godoc
// @Summary Login
// @Description Authenticate a user and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body types.LoginUserPayload true "Login payload"
// @Success 200 {object} types.LoginResponse
// @Failure 400 {object} types.ErrorResponse
// @Router /login [post]
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
	utils.WriteJSON(w, http.StatusOK, types.LoginResponse{Token: token})
}

// HandleRegister godoc
// @Summary Register
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param payload body types.RegisterUserPayload true "Registration payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} types.ErrorResponse
// @Router /register [post]
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

// HandleGetProfile godoc
// @Summary Get profile
// @Description Get authenticated user's profile
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} types.UserProfile
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /profile [get]
func (h *Handler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID <= 0 {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, toUserProfile(user))
}

// HandleUpdateProfile godoc
// @Summary Update profile
// @Description Update authenticated user's first and last name
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body types.UpdateProfilePayload true "Profile payload"
// @Success 200 {object} types.UserProfile
// @Failure 400 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /profile [put]
func (h *Handler) HandleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID <= 0 {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	var payload types.UpdateProfilePayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	user, err := h.store.UpdateUserProfile(userID, payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, toUserProfile(user))
}

// HandleUpdatePassword godoc
// @Summary Update password
// @Description Change authenticated user's password
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body types.UpdatePasswordPayload true "Password payload"
// @Success 204 {object} nil
// @Failure 400 {object} types.ErrorResponse
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /profile/password [put]
func (h *Handler) HandleUpdatePassword(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID <= 0 {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	var payload types.UpdatePasswordPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if !auth.ComparePasswords(user.Password, payload.CurrentPassword) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("current password is invalid"))
		return
	}

	hashedPassword, err := auth.HashPassword(payload.NewPassword)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := h.store.UpdateUserPassword(userID, hashedPassword); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleListUsers godoc
// @Summary List users lookup
// @Description List users for assignment lookups
// @Tags users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} types.UserLookup
// @Failure 403 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /users/lookup [get]
func (h *Handler) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromContext(r.Context())
	if userID <= 0 {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		return
	}

	users, err := h.store.ListUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, users)
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

func toUserProfile(user *types.User) types.UserProfile {
	return types.UserProfile{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
