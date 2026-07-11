package api

import (
	"net/http"

	"github.com/ArlindoSilva/gobid/internal/jsonutils"
	"github.com/gorilla/csrf"
)

func (api *Api) HandleGetCSRFToken(w http.ResponseWriter, r *http.Request) {
	// Obtém o token CSRF gerado pelo middleware CSRF
	csrfToken := csrf.Token(r)
	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"csrf_token": csrfToken,
	})
}

// AuthMiddleware is a middleware that checks if the user is authenticated by checking if the session contains the "AuthenticatedUserId" key. If the key is not present, it returns a 401 Unauthorized response. Otherwise, it calls the next handler in the chain.
func (api *Api) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.Sessions.Exists(r.Context(), "AuthenticatedUserId") {
			jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{
				"message": "must be logged in",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}