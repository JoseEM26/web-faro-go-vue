package middleware

import (
	"context"
	"net/http"
	"strings"

	"proyecto-go/pkg/response"
	"proyecto-go/pkg/token"
)

type contextKey string

// UserIDKey es la clave para obtener el ID del usuario autenticado desde el contexto.
// Uso en un handler: userID := r.Context().Value(middleware.UserIDKey).(uint)
const UserIDKey contextKey = "user_id"

// Auth devuelve un middleware que valida el JWT del header Authorization.
// Rutas sin este middleware son publicas. Rutas con el son protegidas.
func Auth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				response.Error(w, http.StatusUnauthorized, "se requiere autenticacion")
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := token.Validate(tokenStr, secret)
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "token invalido o expirado")
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
