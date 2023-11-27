package controller

import (
	"context"
	"forum/internal/models"
	"net/http"
	"time"
)

// ctxKey is a custom type used as a key for context values in middleware functions.
type ctxKey int

const (
	// ctxKeyUser is a context key for storing user information in middleware functions.
	ctxKeyUser ctxKey = iota
)

func (h *Handler) authenticateUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			user models.User
			err  error
		)

		cookie, err := r.Cookie("sessionID")
		if err != nil {
			h.errorPage(w, http.StatusUnauthorized, err.Error())
			return
		}

		user, err = h.services.GetSessionToken(cookie.Value)
		if err != nil {
			h.errorPage(w, http.StatusUnauthorized, err.Error())
			return
		}

		if err != nil || user.ExpiresAt.Before(time.Now()) {
			// Clear the invalid or expired session cookie
			http.SetCookie(w, &http.Cookie{
				Name:    "sessionID",
				Value:   "",
				Expires: time.Now().Add(-1 * time.Hour),
			})

			// Redirect to login or show error page
			h.errorPage(w, http.StatusUnauthorized, "Session expired. Please log in again.")
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, user)))
	}
}
