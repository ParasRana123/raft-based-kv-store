package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"raft-kv-store/kv"
)

func setupTestServer() *Server {
	store := kv.NewStore()
	return NewServer(store)
}

func TestHealthEndpoint(t *testing.T) {
	server := setupTestServer()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "ok" {
		t.Fatalf("expected status 'ok', got '%s'", resp.Status)
	}
}

func TestPutAndGetEndpointJSON(t *testing.T) {
	server := setupTestServer()

	// 1. PUT /kv/name with JSON body
	payload := `{"value": "Paras"}`
	putReq := httptest.NewRequest(http.MethodPut, "/kv/name", bytes.NewBufferString(payload))
	putReq.Header.Set("Content-Type", "application/json")
	putRR := httptest.NewRecorder()

	server.ServeHTTP(putRR, putReq)

	if putRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusOK, putRR.Code, putRR.Body.String())
	}

	var putResp KVResponse
	if err := json.Unmarshal(putRR.Body.Bytes(), &putResp); err != nil {
		t.Fatalf("failed to parse PUT response: %v", err)
	}
	if putResp.Key != "name" || putResp.Value != "Paras" {
		t.Fatalf("unexpected PUT response: %+v", putResp)
	}

	// 2. GET /kv/name
	getReq := httptest.NewRequest(http.MethodGet, "/kv/name", nil)
	getRR := httptest.NewRecorder()

	server.ServeHTTP(getRR, getReq)

	if getRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getRR.Code)
	}

	var getResp KVResponse
	if err := json.Unmarshal(getRR.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("failed to parse GET response: %v", err)
	}
	if getResp.Key != "name" || getResp.Value != "Paras" {
		t.Fatalf("expected name=Paras, got %+v", getResp)
	}
}

func TestPutAndGetEndpointRawBody(t *testing.T) {
	server := setupTestServer()

	// PUT with raw text
	putReq := httptest.NewRequest(http.MethodPut, "/kv/city", bytes.NewBufferString("New Delhi"))
	putRR := httptest.NewRecorder()

	server.ServeHTTP(putRR, putReq)

	if putRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, putRR.Code)
	}

	// GET /kv/city
	getReq := httptest.NewRequest(http.MethodGet, "/kv/city", nil)
	getRR := httptest.NewRecorder()

	server.ServeHTTP(getRR, getReq)

	var getResp KVResponse
	if err := json.Unmarshal(getRR.Body.Bytes(), &getResp); err != nil {
		t.Fatalf("failed to parse GET response: %v", err)
	}
	if getResp.Value != "New Delhi" {
		t.Fatalf("expected 'New Delhi', got '%s'", getResp.Value)
	}
}

func TestGetNotFound(t *testing.T) {
	server := setupTestServer()

	getReq := httptest.NewRequest(http.MethodGet, "/kv/missing_key", nil)
	getRR := httptest.NewRecorder()

	server.ServeHTTP(getRR, getReq)

	if getRR.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, getRR.Code)
	}
}

func TestDeleteEndpoint(t *testing.T) {
	server := setupTestServer()

	// Put key first
	putReq := httptest.NewRequest(http.MethodPut, "/kv/temp", bytes.NewBufferString("val"))
	putRR := httptest.NewRecorder()
	server.ServeHTTP(putRR, putReq)

	// DELETE /kv/temp
	delReq := httptest.NewRequest(http.MethodDelete, "/kv/temp", nil)
	delRR := httptest.NewRecorder()
	server.ServeHTTP(delRR, delReq)

	if delRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, delRR.Code)
	}

	// Verify it's gone
	getReq := httptest.NewRequest(http.MethodGet, "/kv/temp", nil)
	getRR := httptest.NewRecorder()
	server.ServeHTTP(getRR, getReq)

	if getRR.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got %d", getRR.Code)
	}
}

func TestPutEmptyBody(t *testing.T) {
	server := setupTestServer()

	putReq := httptest.NewRequest(http.MethodPut, "/kv/empty", bytes.NewBufferString(""))
	putRR := httptest.NewRecorder()
	server.ServeHTTP(putRR, putReq)

	if putRR.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for empty body, got %d", putRR.Code)
	}
}
