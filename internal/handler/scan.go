package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"scanfetcher/internal/domain"
	"scanfetcher/internal/service"
	"scanfetcher/internal/usecase"

	"github.com/gorilla/mux"
)

type ScanHandler struct {
	usecase *usecase.ScanUseCase
	svc     *service.ScanService
}

func NewScanHandler(uc *usecase.ScanUseCase, svc *service.ScanService) *ScanHandler {
	return &ScanHandler{usecase: uc, svc: svc}
}

func (h *ScanHandler) CreateScan(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Name         string `json:"name"`
		Url          string `json:"url"`
		LastScanRead string `json:"last_scan_read"`
		WebsiteID    uint   `json:"website_id"`
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

	scans, err := h.usecase.GetScansForWebsite(uint(websiteID))
	if err != nil {
		http.Error(w, "Failed to fetch scans", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scans)
}

func (h *ScanHandler) RunAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute)
	defer cancel()

	scans, err := h.svc.RunAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(struct {
		Count  int         `json:"count"`
		Result interface{} `json:"result"`
	}{Count: len(scans), Result: scans})
}

func (h *ScanHandler) RunForWebsite(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("website_id")
	id, err := strconv.ParseInt(q, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid website_id", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Minute)
	defer cancel()

	scans, err := h.svc.RunForWebsiteID(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(struct {
		Count  int         `json:"count"`
		Result interface{} `json:"result"`
	}{Count: len(scans), Result: scans})
}
