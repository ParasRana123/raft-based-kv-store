package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"raft-kv-store/kv"
)

// Server handles HTTP requests for the key-value store.
type Server struct {
	store kv.KVStore
	mux   *http.ServeMux
}

// SetRequest represents the JSON payload for setting a value.
type SetRequest struct {
	Value string `json:"value"`
}

// KVResponse represents a standard key-value JSON response.
type KVResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse represents a generic message response.
type MessageResponse struct {
	Message string `json:"message"`
}

// HealthResponse represents the health check response.
type HealthResponse struct {
	Status string `json:"status"`
}

// NewServer creates and initializes a new HTTP Server.
func NewServer(store kv.KVStore) *Server {
	s := &Server{
		store: store,
		mux:   http.NewServeMux(),
	}
	s.registerRoutes()
	return s
}

// ServeHTTP implements the http.Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("GET /health", s.handleHealth)
	s.mux.HandleFunc("PUT /kv/{key}", s.handlePutKV)
	s.mux.HandleFunc("GET /kv/{key}", s.handleGetKV)
	s.mux.HandleFunc("DELETE /kv/{key}", s.handleDeleteKV)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, HealthResponse{Status: "ok"})
}

func (s *Server) handlePutKV(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "key is required"})
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "unable to read request body"})
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "request body cannot be empty"})
		return
	}

	value := ""
	var req SetRequest
	if err := json.Unmarshal(body, &req); err == nil && req.Value != "" {
		value = req.Value
	} else {
		// Fallback to raw text body
		value = strings.TrimSpace(string(body))
	}

	s.store.Set(key, value)
	writeJSON(w, http.StatusOK, KVResponse{Key: key, Value: value})
}

func (s *Server) handleGetKV(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "key is required"})
		return
	}

	value, exists := s.store.Get(key)
	if !exists {
		writeJSON(w, http.StatusNotFound, ErrorResponse{Error: "key not found"})
		return
	}

	writeJSON(w, http.StatusOK, KVResponse{Key: key, Value: value})
}

func (s *Server) handleDeleteKV(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "key is required"})
		return
	}

	s.store.Delete(key)
	writeJSON(w, http.StatusOK, MessageResponse{Message: "key deleted"})
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}
