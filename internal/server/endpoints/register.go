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

var vld = validator.New()

// Returns a handler for user registration
func Register(logger *slog.Logger, repo abstractions.UsersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// name of the endpoint
		ep := "endpoinds.Register"

		// an anonymous struct for server's response
		type response struct {
			Status string `json:"status"`
			Error  string `json:"error,omitempty"`
		}

		// decoder of the body's json
		decoder := json.NewDecoder(r.Body)
		// encodes the response to the response body
		encoder := json.NewEncoder(w)

		// setting the type of response
		w.Header().Set("Content-Type", "application/json")

		// editing the logger
		logger := logger.With(
			slog.String("ep", ep),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// request with all the info needed for registration
		var req struct {
			Username string `json:"username" validate:"required"`
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required"`
			Role     string `json:"role" validate:"required"`

			CollegeID uint `json:"college_id" validate:"required"`
		}

		// decoding the request's body
		err := decoder.Decode(&req)
		// if the body's empty
		if errors.Is(err, io.EOF) {
			logger.Error("request body is empty")

			w.WriteHeader(http.StatusBadRequest)

			encoder.Encode(response{
				Status: "Error",
				Error:  "Request body is empty",
			})

			return
		}
		// if another error occurs
		if err != nil {
			logger.Error("cannot decode the request body")

			w.WriteHeader(http.StatusBadRequest)

			encoder.Encode(response{
				Status: "Error",
				Error:  "Cannot decode the request body",
			})

			return
		}

		// logging...
		logger.Info(
			"request body decoded",
			slog.Any("request", req),
		)

		// validating the request
		if err := vld.Struct(&req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			logger.Error("invalid request", slog.Any("err", err))

			w.WriteHeader(http.StatusBadRequest)

			encoder.Encode(response{
				Status: "Error",
				Error:  errfmt.ValidationErrorsToString(validateErr),
			})

			return
		}

		// creating a user model
		user := models.User{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: hashing.HashSHA256(req.Password),
			Role:         req.Role,

			CollegeID: req.CollegeID,
		}

		// trying to write user to a db
		// and handling an error if the one occurs
		if err := repo.Create(&user); err != nil {
			logger.Error("cannot add user to db", slog.Any("err", err))

			w.WriteHeader(http.StatusBadRequest)

			encoder.Encode(response{
				Status: "Error",
				Error:  "Cannot add user to db",
			})

			return
		}

		// logging...
		logger.Info(
			"user has been successfully added",
			slog.String("email", user.Email),
		)

		// OK response
		w.WriteHeader(http.StatusCreated)

		encoder.Encode(response{
			Status: "OK",
		})
	}
}
