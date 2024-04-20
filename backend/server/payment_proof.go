package server

import (
	"encoding/json"
	"errors"
	"mime"
	"net/http"
	"path"
	"slices"

	"conf/ticketing"
	"conf/user"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func (s *ServerDependency) UploadPaymentProof(w http.ResponseWriter, r *http.Request) {
	requestId := middleware.GetReqID(r.Context())
	sentry.GetHubFromContext(r.Context()).Scope().SetTag("request-id", requestId)

	if !s.featureFlag.EnablePaymentProofUpload {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	
	if err := r.ParseMultipartForm(32 << 10); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message":    "Parsing error",
			"errors":     err.Error(),
			"request_id": requestId,
		})
		return
	}

	email := r.FormValue("email")
	if email == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message":    "Validation error",
			"errors":     "Email field is required",
			"request_id": requestId,
		})
		return
	}

	photoFile, photoFormHeader, err := r.FormFile("photo")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"message":    "Validation error",
				"errors":     "Photo field is required",
				"request_id": requestId,
			})
			return
		}

		sentry.GetHubFromContext(r.Context()).CaptureException(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message":    "Reading form file",
			"errors":     err.Error(),
			"request_id": requestId,
		})
		return
	}
	defer func() {
		err := photoFile.Close()
		if err != nil {
			log.Error().Err(err).Str("request_id", requestId).Msg("Closing photo file")
		}
	}()

	photoExtension := path.Ext(photoFormHeader.Filename)
	// Guard the content type, the only content type allowed is images.
	if !slices.Contains([]string{".gif", ".jpeg", ".jpg", ".png", ".webp"}, photoExtension) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"message":    "Unknown photo file type",
			"errors":     "Unknown photo file type",
			"request_id": requestId,
		})
		return
	}

	photoContentType := mime.TypeByExtension(photoExtension)

	userEntry, err := s.userDomain.GetUserByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, user.ErrUserEmailNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusPreconditionFailed)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"message":    "User not found",
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

	err = s.ticketDomain.StorePaymentReceipt(r.Context(), userEntry, photoFile, photoContentType)
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
