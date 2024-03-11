package nocodbmock

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"

	"conf/nocodb"
	"github.com/go-chi/chi/v5"
)

type errorResponse struct {
	Message string `json:"msg"`
}

type listTableResponse struct {
	List     []map[string]any `json:"list"`
	PageInfo nocodb.PageInfo  `json:"pageInfo"`
}

type creationSuccessfulResponse struct {
	ID int64 `json:"Id"`
}

func NewNocoDBMockServer() (*httptest.Server, error) {
	documentStorage := newInMemoryStorage()

	r := chi.NewRouter()

	r.Get("/api/v2/tables/{tableId}/records", func(w http.ResponseWriter, r *http.Request) {
		tableId := chi.URLParam(r, "tableId")
		if tableId == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(errorResponse{Message: "BadRequest [ERROR]: tableId is empty"})
			return
		}

		records, err := documentStorage.GetByTableId(tableId)
		if err != nil && !errors.Is(err, errNotFound) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(errorResponse{Message: "BadRequest [ERROR]: " + err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(listTableResponse{
			List: records,
			PageInfo: nocodb.PageInfo{
				TotalRows:   int64(len(records)),
				Page:        1,
				PageSize:    1,
				IsFirstPage: true,
				IsLastPage:  true,
			},
		})
		return
	})

	r.Post("/api/v2/tables/{tableId}/records", func(w http.ResponseWriter, r *http.Request) {
		tableId := chi.URLParam(r, "tableId")
		if tableId == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(errorResponse{Message: "BadRequest [ERROR]: tableId is empty"})
			return
		}

		var records []map[string]any
		err := json.NewDecoder(r.Body).Decode(&records)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(errorResponse{Message: "BadRequest [ERROR]: Invalid request body"})
			return
		}

		recordIds, err := documentStorage.Insert(tableId, records)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errorResponse{Message: err.Error()})
			return
		}

		var response []creationSuccessfulResponse
		for _, id := range recordIds {
			response = append(response, creationSuccessfulResponse{ID: id})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
		return
	})

	r.Patch("/api/v2/tables/{tableId}/records", func(w http.ResponseWriter, r *http.Request) {
		tableId := chi.URLParam(r, "tableId")
		if tableId == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(errorResponse{Message: "BadRequest [ERROR]: tableId is empty"})
			return
		}

		var records []map[string]any
		err := json.NewDecoder(r.Body).Decode(&records)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(errorResponse{Message: "BadRequest [ERROR]: Invalid request body"})
			return
		}

		recordIds, err := documentStorage.Update(tableId, records)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(errorResponse{Message: err.Error()})
			return
		}

		var response []creationSuccessfulResponse
		for _, id := range recordIds {
			response = append(response, creationSuccessfulResponse{ID: id})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
		return
	})

	r.Get("/api/v2/tables/{tableId}/records/{recordId}", func(w http.ResponseWriter, r *http.Request) {
		tableId := chi.URLParam(r, "tableId")
		if tableId == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(errorResponse{Message: "BadRequest [ERROR]: tableId is empty"})
			return
		}

		recordId := chi.URLParam(r, "recordId")
		if recordId == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(errorResponse{Message: "BadRequest [ERROR]: recordId is empty"})
			return
		}

		recordIdAsInt64, err := strconv.ParseInt(recordId, 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(errorResponse{Message: "BadRequest [ERROR]: invalid recordId value"})
			return
		}

		record, err := documentStorage.GetByRecordId(tableId, recordIdAsInt64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(errorResponse{Message: "BadRequest [ERROR]: " + err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(record)
		return
	})

	server := httptest.NewServer(r)

	return server, nil
}
