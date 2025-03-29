package middlewares

import (
	"go-json/internal/security"
	"net/http"
	"slices"
)

func ProtectedHandler(next http.Handler, token security.TokenService, roles []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):]
		claims, err := token.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Pastikan user memiliki salah satu role yang diperbolehkan
		authorized := false
		for _, userRole := range claims.Role {
			if slices.Contains(roles, userRole) {
				authorized = true
				break
			}
		}

		if !authorized {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Simpan email ke header untuk digunakan di handler berikutnya
		r.Header.Set("email", claims.Email)
		next.ServeHTTP(w, r)
	})
}
