package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AuthHandler struct {
	JWKS keyfunc.Keyfunc
}

func NewAuthHandler(jwksURL string) (*AuthHandler, error) {
	ctx := context.Background()

	jwks, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("failed to load JWKS: %w", err)
	}

	return &AuthHandler{JWKS: jwks}, nil
}

func (a *AuthHandler) validateToken(tokenStr string) (*User, error) {

	token, err := jwt.Parse(tokenStr, a.JWKS.Keyfunc, jwt.WithValidMethods([]string{"RS256"}))

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := &User{
			ID:    claims["oid"].(string),
			Name:  claims["name"].(string),
			Email: claims["upn"].(string),
		}
		return user, nil
	}

	return nil, errors.New("invalid token")
}

func (a *AuthHandler) getUser(w http.ResponseWriter, r *http.Request) (*User, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("missing or invalid authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	user, err := a.validateToken(token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func AuthMiddleware(authHandler *AuthHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := authHandler.getUser(w, r)
			if err != nil {
				http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Store user in context
			ctx := context.WithValue(r.Context(), "user", user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserFromContext(r *http.Request) *User {
	user, ok := r.Context().Value("user").(*User)
	if !ok {
		return nil
	}
	return user
}
