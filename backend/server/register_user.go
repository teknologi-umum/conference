package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"conf/user"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5/middleware"
)

type RegisterUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *ServerDependency) RegisterUser(w http.ResponseWriter, r *http.Request) {
	requestId := middleware.GetReqID(r.Context())
	sentry.GetHubFromContext(r.Context()).Scope().SetTag("request-id", requestId)

	if !s.featureFlag.EnableRegistration {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotAcceptable)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message":    "Registration is closed",
			"request_id": requestId,
		})
		return
	}

	requestBody := RegisterUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message":    "Invalid request body",
			"errors":     err.Error(),
			"request_id": requestId,
		})
		return
	}

	err := s.userDomain.CreateParticipant(
		r.Context(),
		user.CreateParticipantRequest{
			Name:  requestBody.Name,
			Email: requestBody.Email,
		},
	)
	if err != nil {
		var validationError *user.ValidationError
		if errors.As(err, &validationError) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"message":    "Validation error",
				"errors":     validationError.Errors,
				"request_id": requestId,
			})
			return
		}

		sentry.GetHubFromContext(r.Context()).CaptureException(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message":    "Internal server error",
			"errors":     "Internal server error",
			"request_id": requestId,
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	return
}
