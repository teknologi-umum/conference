package server

import (
	"encoding/json"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5/middleware"
)

type AdministratorLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	OtpCode  string `json:"otp"`
}

func (s *ServerDependency) AdministratorLogin(w http.ResponseWriter, r *http.Request) {
	requestId := middleware.GetReqID(r.Context())
	sentry.GetHubFromContext(r.Context()).Scope().SetTag("request-id", requestId)

	if !s.featureFlag.EnableAdministratorMode {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var requestBody AdministratorLoginRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message":    "Invalid request body",
			"errors":     err.Error(),
			"request_id": requestId,
		})
		return
	}

	token, ok, err := s.administratorDomain.Authenticate(r.Context(), requestBody.Username, requestBody.Password, requestBody.OtpCode)
	if err != nil {
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

	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message":    "Invalid authentication",
			"request_id": requestId,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"token":      token,
		"request_id": requestId,
	})
	return
}
