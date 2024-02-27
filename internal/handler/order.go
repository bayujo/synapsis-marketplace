package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bayujo/synapsis-marketplace/internal/order"
	"github.com/gorilla/mux"
)

type OrderHandler struct {
	OrderUsecase order.OrderUseCase
}

func NewOrderHandler(orderUsecase order.OrderUseCase) *OrderHandler {
	return &OrderHandler{OrderUsecase: orderUsecase}
}

func (h *OrderHandler) CheckoutCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		UserID int `json:"user_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.OrderUsecase.CheckoutCart(ctx, req.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetOrderDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	orderIDStr := vars["orderID"]

	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := h.OrderUsecase.GetOrderDetails(ctx, orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) PayOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

    vars := mux.Vars(r)
    orderIDStr := vars["orderID"]

    orderID, err := strconv.Atoi(orderIDStr)
    if err != nil {
        http.Error(w, "Invalid order ID", http.StatusBadRequest)
        return
    }

    err = h.OrderUsecase.PayOrder(ctx, orderID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

    vars := mux.Vars(r)
    orderIDStr := vars["orderID"]

    orderID, err := strconv.Atoi(orderIDStr)
    if err != nil {
        http.Error(w, "Invalid order ID", http.StatusBadRequest)
        return
    }

    err = h.OrderUsecase.CancelOrder(ctx, orderID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}