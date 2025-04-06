package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"scanfetcher/internal/usecase"

	"github.com/gorilla/mux"
)

type WebsiteHandler struct {
	usecase *usecase.WebsiteUseCase
}

func NewWebsiteHandler(uc *usecase.WebsiteUseCase) *WebsiteHandler {
	return &WebsiteHandler{usecase: uc}
}

func (h *WebsiteHandler) CreateWebsite(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Name string `json:"name"`
	}

	var body input
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if err := h.usecase.CreateWebsite(body.Name); err != nil {
		log.Printf("Error creating website: %v", err)
		http.Error(w, "Failed to create website", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *WebsiteHandler) GetWebsite(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid website ID", http.StatusBadRequest)
		return
	}

	website, err := h.usecase.GetWebsite(id)
	if err != nil {
		http.Error(w, "Website not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(website)
}

func (h *WebsiteHandler) ListWebsites(w http.ResponseWriter, r *http.Request) {
	websites, err := h.usecase.ListWebsites()
	if err != nil {
		http.Error(w, "Failed to fetch websites", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(websites)
}
