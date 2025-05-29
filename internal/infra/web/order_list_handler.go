package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	
	"github.com/devfullcycle/20-CleanArch/internal/usecase"
)

type OrderListHandler struct {
	listOrdersUseCase usecase.ListOrdersUseCase
}

func NewOrderListHandler(listOrdersUseCase usecase.ListOrdersUseCase) *OrderListHandler {
	return &OrderListHandler{
		listOrdersUseCase: listOrdersUseCase,
	}
}

func handleError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (h *OrderListHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters from query string
	page := 1
	perPage := 10

	if r.URL.Query().Get("page") != "" {
		var err error
		page, err = strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			handleError(w, http.StatusBadRequest, "Invalid page parameter")
			return
		}
	}

	if r.URL.Query().Get("per_page") != "" {
		var err error
		perPage, err = strconv.Atoi(r.URL.Query().Get("per_page"))
		if err != nil {
			handleError(w, http.StatusBadRequest, "Invalid per_page parameter")
			return
		}
	}

	// Prepare input for use case
	input := usecase.ListOrdersInput{
		Page:    page,
		PerPage: perPage,
	}

	// Execute use case
	output, err := h.listOrdersUseCase.Execute(r.Context(), input)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(output); err != nil {
		handleError(w, http.StatusInternalServerError, err.Error())
	}
}
