package auth

import (
	"VyacheslavKuchumov/test-backend/config"
	"VyacheslavKuchumov/test-backend/types"
	"VyacheslavKuchumov/test-backend/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "userID"
const AuthCookieName = "task_tracker_token"

func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func SetAuthCookie(w http.ResponseWriter, token string) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)
	expiresAt := time.Now().Add(expiration)
	http.SetCookie(w, &http.Cookie{
		Name:     AuthCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int(config.Envs.JWTExpirationInSeconds),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromRequest(r, store)
		if err != nil {
			log.Printf("Failed to authorize request: %v", err)
			permissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, userID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func getUserIDFromRequest(r *http.Request, store types.UserStore) (int, error) {
	tokenString := getTokenFromRequest(r)
	token, err := validateToken(tokenString)
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)
	str, ok := claims["userID"].(string)
	if !ok {
		return 0, fmt.Errorf("missing userID claim")
	}
	userID, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	u, err := store.GetUserByID(userID)
	if err != nil {
		return 0, err
	}

	return u.ID, nil
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenAuth = strings.TrimSpace(tokenAuth)
	if tokenAuth == "" {
		cookie, err := r.Cookie(AuthCookieName)
		if err != nil {
			return ""
		}
		return strings.TrimSpace(cookie.Value)
	}

	if strings.HasPrefix(strings.ToLower(tokenAuth), "bearer ") {
		return strings.TrimSpace(tokenAuth[7:])
	}

	return tokenAuth
}

func validateToken(t string) (*jwt.Token, error) {
	return jwt.Parse(t, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("Permission denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)

	if !ok {
		return -1
	}

	return userID
}
