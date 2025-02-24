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
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

// Returns a handler for college creation
func CreateCollege(logger *slog.Logger, repo abstractions.CollegesRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ep := "endpoints.CreateCollege"

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

		var req struct {
			Name string `json:"name" validate:"required"`
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

		college := models.College{
			Name: req.Name,
		}

		if err := repo.Create(&college); err != nil {
			logger.Error("cannot add college to db", slog.Any("err", err))

			w.WriteHeader(http.StatusBadRequest)

			encoder.Encode(response{
				Status: "Error",
				Error:  "Cannot add college to db",
			})

			return
		}

		logger.Info(
			"college has been successfully added",
			slog.String("name", college.Name),
		)

		w.WriteHeader(http.StatusCreated)

		encoder.Encode(response{
			Status: "OK",
		})
	}
}
