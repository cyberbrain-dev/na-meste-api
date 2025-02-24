package endpoints

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/cyberbrain-dev/na-meste-api/internal/models"
	"github.com/cyberbrain-dev/na-meste-api/internal/models/abstractions"
	"github.com/cyberbrain-dev/na-meste-api/pkg/errfmt"
	"github.com/cyberbrain-dev/na-meste-api/pkg/hashing"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

// Returns a handler for user registration
func RegisterUser(logger *slog.Logger, repo abstractions.UsersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// name of the endpoint
		ep := "endpoinds.Register"

		type response struct {
			Status string `json:"status"`
			Error  string `json:"error,omitempty"`
		}

		decoder := json.NewDecoder(r.Body)
		encoder := json.NewEncoder(w)

		w.Header().Set("Content-Type", "application/json")

		// editing the logger
		logger := logger.With(
			slog.String("ep", ep),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Request with all the info needed for registration
		var req struct {
			Username string `json:"username" validate:"required"`
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required"`
			Role     string `json:"role" validate:"required"`

			CollegeID uint `json:"college_id" validate:"required"`
		}

		err := decoder.Decode(&req)
		if errors.Is(err, io.EOF) {
			logger.Error("request body is empty")

			w.WriteHeader(http.StatusBadRequest)

			encoder.Encode(response{
				Status: "Error",
				Error:  "Request body is empty",
			})

			return
		}
		if err != nil {
			logger.Error("cannot decode the request body")

			w.WriteHeader(http.StatusBadRequest)

			encoder.Encode(response{
				Status: "Error",
				Error:  "Cannot decode the request body",
			})

			return
		}

		logger.Info(
			"request body decoded",
			slog.Any("request", req),
		)

		if err := validator.New().Struct(&req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			logger.Error("invalid request", slog.Any("err", err))

			w.WriteHeader(http.StatusBadRequest)

			encoder.Encode(response{
				Status: "Error",
				Error:  errfmt.ValidationErrorsToString(validateErr),
			})

			return
		}

		user := models.User{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: hashing.HashSHA256(req.Password),
			Role:         req.Role,

			CollegeID: req.CollegeID,
		}

		if err := repo.Create(&user); err != nil {
			logger.Error("cannot add user to db", slog.Any("err", err))

			w.WriteHeader(http.StatusBadRequest)

			encoder.Encode(response{
				Status: "Error",
				Error:  "Cannot add user to db",
			})

			return
		}

		logger.Info(
			"user has been successfully added",
			slog.String("email", user.Email),
		)

		w.WriteHeader(http.StatusCreated)

		encoder.Encode(response{
			Status: "OK",
		})
	}
}
