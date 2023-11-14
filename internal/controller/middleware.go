package controller

type ctxKey int8

const (
	ctxKeyUser ctxKey = iota
)

func (h *Handler) authenticateUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			user models.User
			err	 error
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

		if user.ExpiresAt.Before(time.Now()) {
			if err := h.services.DeleteSessionToken(cookie.Value); err != nil {
				h.errorPage(w, http.StatusInternalServerError, err.Error())
				return
			}
			next.ServeHttp(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, models.User{})))
			return
		}

		next.ServeHttp(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, user)))
	}
}