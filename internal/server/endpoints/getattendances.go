package endpoints

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/cyberbrain-dev/na-meste-api/internal/models"
	"github.com/cyberbrain-dev/na-meste-api/internal/models/abstractions"
	"github.com/cyberbrain-dev/na-meste-api/pkg/errfmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

// An endpoint for getting the endpoints
func GetAttendances(logger *slog.Logger, repo abstractions.AttendancesRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// name of the endpoint
		ep := "endpoinds.GetAttendances"

		// a struct for server's response
		type response struct {
			Status      string               `json:"status"`
			Error       string               `json:"error,omitempty"`
			Attendances []*models.Attendance `json:"attendances,omitempty"`
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

		// client's request for getting the attendances
		var req struct {
			StudentID uint      `json:"student_id" validate:"required,number"`
			StartDate time.Time `json:"start_date" validate:"required"`
			EndDate   time.Time `json:"end_date" validate:"required"`
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

		// checking if the dates are correct
		if req.StartDate.Unix() > req.EndDate.Unix() {
			logger.Error("start date must be less than end date")

			w.WriteHeader(http.StatusBadRequest)

			encoder.Encode(response{
				Status: "Error",
				Error:  "start date must be less than end date",
			})

			return
		}

		// getting the attendances
		atts, err := repo.GetByStudentAndDatespan(
			req.StudentID,
			req.StartDate,
			req.EndDate,
		)
		// if an error occurs
		if err != nil {
			logger.Error("cannot get the attendances")

			w.WriteHeader(http.StatusInternalServerError)

			encoder.Encode(response{
				Status: "Error",
				Error:  "cannot get the attendances",
			})

			return
		}

		logger.Info("successfully got the attendances", slog.Any("student_id", req.StudentID))

		w.WriteHeader(http.StatusOK)

		encoder.Encode(response{
			Status:      "OK",
			Attendances: atts,
		})

		return
	}
}
