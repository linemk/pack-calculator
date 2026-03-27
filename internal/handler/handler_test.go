package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/linemk/pack-calculator/internal/calculator"
	"github.com/linemk/pack-calculator/internal/handler"
	"github.com/linemk/pack-calculator/internal/store"
)

func setup(t *testing.T) (*mux.Router, *store.Store) {
	t.Helper()
	s := store.New()
	c := calculator.New(s)
	h := handler.New(c, s)
	r := mux.NewRouter()
	h.RegisterRoutes(r)
	return r, s
}

func TestGetPacks(t *testing.T) {
	t.Parallel()
	r, _ := setup(t)
	req := httptest.NewRequest(http.MethodGet, "/api/packs", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string][]int
	json.NewDecoder(w.Body).Decode(&resp)
	if len(resp["sizes"]) == 0 {
		t.Fatal("expected non-empty sizes")
	}
}

func TestSetPacks(t *testing.T) {
	t.Parallel()
	r, _ := setup(t)
	body, _ := json.Marshal(map[string][]int{"sizes": {23, 31, 53}})
	req := httptest.NewRequest(http.MethodPut, "/api/packs", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCalculate(t *testing.T) {
	t.Parallel()
	r, s := setup(t)
	s.Set([]int{23, 31, 53})

	body, _ := json.Marshal(map[string]int{"order": 500000})
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Packs      map[string]int `json:"packs"`
		TotalItems int            `json:"total_items"`
	}
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.TotalItems != 500000 {
		t.Errorf("expected total_items=500000, got %d", resp.TotalItems)
	}
}

func TestCalculateBadRequest(t *testing.T) {
	t.Parallel()
	r, _ := setup(t)
	body, _ := json.Marshal(map[string]int{"order": -1})
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHealth(t *testing.T) {
	t.Parallel()
	r, _ := setup(t)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
