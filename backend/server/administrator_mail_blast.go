package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"conf/mailer"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5/middleware"
)

type AdministratorMailBlastRequest struct {
	Subject       string                                   `json:"subject"`
	PlaintextBody string                                   `json:"plaintextBody"`
	HtmlBody      string                                   `json:"htmlBody"`
	Recipients    []AdministratorMailBlastRecipientRequest `json:"recipients"`
}

type AdministratorMailBlastRecipientRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (s *ServerDependency) AdministratorMailBlast(w http.ResponseWriter, r *http.Request) {
	requestId := middleware.GetReqID(r.Context())
	sentry.GetHubFromContext(r.Context()).Scope().SetTag("request-id", requestId)

	if !s.featureFlag.EnableAdministratorMode {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var requestBody AdministratorMailBlastRequest
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

	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	_, ok, err := s.administratorDomain.Validate(r.Context(), token)
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
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message":    "Invalid authentication",
			"request_id": requestId,
		})
		return
	}

	var unsuccessfulDestinations []string
	for _, recipient := range requestBody.Recipients {
		mail := &mailer.Mail{
			RecipientName:  recipient.Name,
			RecipientEmail: recipient.Email,
			Subject:        requestBody.Subject,
			PlainTextBody:  strings.ReplaceAll(requestBody.PlaintextBody, "___REPLACE_WITH_NAME___", recipient.Name),
			HtmlBody:       strings.ReplaceAll(requestBody.HtmlBody, "___REPLACE_WITH_NAME___", recipient.Name),
		}

		err := s.mailSender.Send(r.Context(), mail)
		if err != nil {
			unsuccessfulDestinations = append(unsuccessfulDestinations, recipient.Email)
			sentry.GetHubFromContext(r.Context()).CaptureException(err)
			continue
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message":                   "Done",
		"unsuccessful_destinations": unsuccessfulDestinations,
		"request_id":                requestId,
	})
	return
}
