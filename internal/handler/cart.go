package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bayujo/synapsis-marketplace/internal/cart"
	"github.com/gorilla/mux"
)

type CartHandler struct {
	CartUsecase cart.ShoppingCartUseCase
}

func NewCartHandler(cartUsecase cart.ShoppingCartUseCase) *CartHandler {
	return &CartHandler{CartUsecase: cartUsecase}
}

func (h *CartHandler) AddCartItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		UserID    int `json:"user_id"`
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.CartUsecase.AddProductToCart(ctx, req.UserID, req.ProductID, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *CartHandler) GetCartItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	userIDStr := vars["userID"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	items, err := h.CartUsecase.GetCartItems(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(items)
}

func (h *CartHandler) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	productIDStr := vars["productID"]

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = h.CartUsecase.DeleteCartItem(ctx, userID, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}