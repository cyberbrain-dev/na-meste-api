package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/cyberbrain-dev/na-meste-api/pkg/authentication"
	"github.com/go-chi/chi/v5/middleware"
)

// Returns a middleware function that checks the role of the user
func CheckRole(logger *slog.Logger, requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mw := "middleware.CheckRole"

		// editing the logger
		logger := logger.With(
			slog.String("mw", mw),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// getting header that contains the JWT
		authHeader := r.Header.Get("Authorization")
		// if smth goes wrong
		if authHeader == "" {
			logger.Error("no token provided")

			http.Error(w, "no token provided", http.StatusUnauthorized)
			return
		}

		// getting the token itself without "Bearer " prefix
		tokenString := strings.TrimPrefix("Bearer ", authHeader)
		// if there's no Bearer prefix
		if tokenString == authHeader {
			logger.Error("invalid Authorization format")

			http.Error(w, "invalid Authorization format", http.StatusUnauthorized)
			return
		}

		claims, err := authentication.ParseJWT(tokenString)
		if err != nil {
			logger.Error(
				"failed to parse the token",
				slog.Any("err", err),
			)
		}

		// checking the role
		userRole := claims.Role
		if userRole != requiredRole {
			logger.Error(
				"access is forbidden",
				slog.Int("user_id", int(claims.UserID)),
			)

			http.Error(
				w,
				"forbidden: insufficient permissions",
				http.StatusForbidden,
			)
			return
		}

		// if everything gors fine,
		// moving to next endpoint or middleware
		next(w, r)
	}
}
