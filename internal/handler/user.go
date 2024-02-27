package handler

import (
	"encoding/json"
	"github.com/bayujo/synapsis-marketplace/internal/user"
	"net/http"
)

type UserHandler struct {
	UserUsecase user.UserUsecase
}

func NewUserHandler(userUsecase user.UserUsecase) *UserHandler {
	return &UserHandler{
		UserUsecase: userUsecase,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req user.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.UserUsecase.RegisterUser(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req user.User
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.UserUsecase.AuthenticateUser(ctx, req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
