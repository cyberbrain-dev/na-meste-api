package endpoints

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/cyberbrain-dev/na-meste-api/internal/models/abstractions"
	"github.com/cyberbrain-dev/na-meste-api/pkg/authentication"
	"github.com/cyberbrain-dev/na-meste-api/pkg/errfmt"
	"github.com/cyberbrain-dev/na-meste-api/pkg/hashing"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

// Provides an endpoint for logging in the application and getting the JWT
func Login(logger *slog.Logger, repo abstractions.UsersRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// name of the endpoint
		ep := "endpoinds.Login"

		// a struct for server's response
		type response struct {
			Status string `json:"status"`
			Token  string `json:"jwt,omitempty"`
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

		// client's request for logging in
		var req struct {
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required"`
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

		user, err := repo.Get(req.Email)
		// if the user ain't exist
		if err != nil {
			logger.Error("user with this email does not exist")

			w.WriteHeader(http.StatusNotFound)

			encoder.Encode(response{
				Status: "Error",
				Error:  "User with this email does not exist",
			})

			return
		}

		// checking the password
		reqPasswordHash := hashing.HashSHA256(req.Password)
		// if the password is incorrect
		if reqPasswordHash != user.PasswordHash {
			logger.Error("password is incorrect")

			w.WriteHeader(http.StatusUnauthorized)

			encoder.Encode(response{
				Status: "Error",
				Error:  "Password is incorrect",
			})

			return
		}

		// if everything is fine, generating a JWT for this user
		token, err := authentication.GenerateJWT(user.ID, user.Role)
		// if smth goes wrong
		if err != nil {
			logger.Error("failed to generate the JWT", slog.Any("err", err))

			w.WriteHeader(http.StatusInternalServerError)

			encoder.Encode(response{
				Status: "Error",
				Error:  "Failed to log in, try later again",
			})

			return
		}

		logger.Info("successfully logged in", slog.Any("user_id", user.ID))

		w.WriteHeader(http.StatusOK)

		encoder.Encode(response{
			Status: "OK",
			Token:  token,
		})

		return
	}
}
