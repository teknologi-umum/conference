package server

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"conf/ticketing"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/crypto/bcrypt"
)

type DayTicketScanRequest struct {
	Code string `json:"code"`
	Key  string `json:"key"`
}

func (s *ServerDependency) DayTicketScan(w http.ResponseWriter, r *http.Request) {
	requestId := middleware.GetReqID(r.Context())
	sentry.GetHubFromContext(r.Context()).Scope().SetTag("request-id", requestId)

	var requestBody DayTicketScanRequest
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

	// Validate key
	decodedPassphrase, err := hex.DecodeString(s.validateTicketKey)
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

	if err := bcrypt.CompareHashAndPassword(decodedPassphrase, []byte(requestBody.Key)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"message":    "Wrong passphrase",
				"errors":     "",
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

	verifiedTicket, err := s.ticketDomain.VerifyTicket(r.Context(), []byte(requestBody.Code))
	if err != nil {
		var validationError *ticketing.ValidationError
		if errors.As(err, &validationError) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"message":    "Validation error",
				"errors":     validationError.Error(),
				"request_id": requestId,
			})
			return
		}

		if errors.Is(err, ticketing.ErrInvalidTicket) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotAcceptable)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"message":    "Invalid ticket",
				"errors":     err.Error(),
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

	userEntry, err := s.userDomain.GetUserByEmail(r.Context(), verifiedTicket.Email)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "Ticket confirmed",
		"student": verifiedTicket.Student,
		"name":    userEntry.Name,
		"type":    userEntry.Type,
		"email":   verifiedTicket.Email,
	})
	return
}
