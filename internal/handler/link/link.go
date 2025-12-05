package link

import (
	"encoding/json"
	"net/http"

	linkModels "github.com/SenechkaP/link-checker/internal/model/link"
	linkService "github.com/SenechkaP/link-checker/internal/service/link"
)

type LinkHandler struct {
	service *linkService.LinkService
}

func NewLinkHandler(service *linkService.LinkService) *LinkHandler {
	return &LinkHandler{
		service: service,
	}
}

func (h *LinkHandler) GetStatusesHandler(w http.ResponseWriter, r *http.Request) {
	var req linkModels.LinksRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetStatuses(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *LinkHandler) GetStatusesByNumsHandler(w http.ResponseWriter, r *http.Request) {
	var req linkModels.LinksListRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	pdfBytes, err := h.service.GetStatusesByNums(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"links.pdf\"")
	w.Write(pdfBytes)
}
