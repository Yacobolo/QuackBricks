package catalog

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) RegisterCatalogEntry(w http.ResponseWriter, r *http.Request) {
	var req *CatalogEntryInput

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.Service.Register(r.Context(), req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating catalog entry: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Catalog entry created successfully"))
}

func (h *Handler) ListCatalogEntries(w http.ResponseWriter, r *http.Request) {
	catalogEntries, err := h.Service.ListCatalogEntries(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error listing catalog entries: %v", err), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(catalogEntries)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling catalog entries: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func (h *Handler) DeleteCatalogEntry(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Catalog ID is required", http.StatusBadRequest)
		return
	}

	err := h.Service.DeleteCatalogEntry(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting catalog entry: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetCatalogEntry(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Catalog ID is required", http.StatusBadRequest)
		return
	}

	catalogEntry, err := h.Service.GetCatalogEntry(r.Context(), id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting catalog entry: %v", err), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(catalogEntry)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling catalog entry: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
