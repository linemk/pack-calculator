package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/linemk/pack-calculator/internal/calculator"
	"github.com/linemk/pack-calculator/internal/store"
)

type Handler struct {
	calc  *calculator.Calculator
	store *store.Store
}

func New(c *calculator.Calculator, s *store.Store) *Handler {
	return &Handler{calc: c, store: s}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/packs", h.getPacks).Methods(http.MethodGet)
	r.HandleFunc("/api/packs", h.setPacks).Methods(http.MethodPut)
	r.HandleFunc("/api/calculate", h.calculate).Methods(http.MethodPost)
	r.HandleFunc("/health", h.health).Methods(http.MethodGet)
}

func (h *Handler) getPacks(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"sizes": h.store.Get()})
}

func (h *Handler) setPacks(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Sizes []int `json:"sizes"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, errResponse(err.Error()))
		return
	}
	if err := h.store.Set(body.Sizes); err != nil {
		writeJSON(w, http.StatusBadRequest, errResponse(err.Error()))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"sizes": h.store.Get()})
}

func (h *Handler) calculate(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Order int `json:"order"`
	}
	if err := readJSON(r, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, errResponse(err.Error()))
		return
	}
	packs, err := h.calc.Calculate(body.Order)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errResponse(err.Error()))
		return
	}
	total := 0
	for size, qty := range packs {
		total += size * qty
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"packs":       packs,
		"total_items": total,
	})
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func readJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return fmt.Errorf("decode request body: %w", err)
	}
	return nil
}

func errResponse(msg string) map[string]string {
	return map[string]string{"error": msg}
}
