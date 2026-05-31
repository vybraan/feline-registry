package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestServer() *http.ServeMux {
	store := &Store{}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /names", handleRegister(store))
	mux.HandleFunc("GET /names", handleList(store))
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
	})
	return mux
}

func TestRegisterName(t *testing.T) {
	mux := setupTestServer()
	body := bytes.NewBufferString(`{"name":"lion"}`)
	req := httptest.NewRequest(http.MethodPost, "/names", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, w.Code)
	}
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["name"] != "lion" {
		t.Fatalf("expected lion, got %q", resp["name"])
	}
}

func TestListNames(t *testing.T) {
	mux := setupTestServer()
	postBody := bytes.NewBufferString(`{"name":"tiger"}`)
	postReq := httptest.NewRequest(http.MethodPost, "/names", postBody)
	postReq.Header.Set("Content-Type", "application/json")
	postRes := httptest.NewRecorder()
	mux.ServeHTTP(postRes, postReq)
	if postRes.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, postRes.Code)
	}
	getReq := httptest.NewRequest(http.MethodGet, "/names", nil)
	getRes := httptest.NewRecorder()
	mux.ServeHTTP(getRes, getReq)
	if getRes.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, getRes.Code)
	}
	var resp struct {
		Names []string `json:"names"`
	}
	json.Unmarshal(getRes.Body.Bytes(), &resp)
	if len(resp.Names) != 1 || resp.Names[0] != "tiger" {
		t.Fatalf("expected [tiger], got %v", resp.Names)
	}
}

func TestRegisterNameRejectsEmptyName(t *testing.T) {
	mux := setupTestServer()
	body := bytes.NewBufferString(`{"name":"   "}`)
	req := httptest.NewRequest(http.MethodPost, "/names", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHealthz(t *testing.T) {
	mux := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestReadyz(t *testing.T) {
	mux := setupTestServer()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
}
