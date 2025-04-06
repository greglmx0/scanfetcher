package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"scanfetcher/internal/domain"
	"scanfetcher/internal/usecase"

	"github.com/gorilla/mux"
)

type ScanHandler struct {
	usecase *usecase.ScanUseCase
}

func NewScanHandler(uc *usecase.ScanUseCase) *ScanHandler {
	return &ScanHandler{usecase: uc}
}

func (h *ScanHandler) CreateScan(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Name         string `json:"name"`
		Url          string `json:"url"`
		LastScanRead string `json:"last_scan_read"`
		WebsiteID    uint    `json:"website_id"`
	}

	var body input
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	scan := domain.Scan{
		Name:         body.Name,
		Url:          body.Url,
		LastScanRead: body.LastScanRead,
		WebsiteID:    body.WebsiteID,
	}

	if err := h.usecase.CreateScan(scan); err != nil {
		log.Printf("Error creating scan: %v", err)
		http.Error(w, "Failed to create scan", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ScanHandler) GetScansForWebsite(w http.ResponseWriter, r *http.Request) {
	websiteIDStr := mux.Vars(r)["websiteID"]
	websiteID, err := strconv.Atoi(websiteIDStr)
	if err != nil {
		http.Error(w, "Invalid website ID", http.StatusBadRequest)
		return
	}

	scans, err := h.usecase.GetScansForWebsite(websiteID)
	if err != nil {
		http.Error(w, "Failed to fetch scans", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scans)
}
