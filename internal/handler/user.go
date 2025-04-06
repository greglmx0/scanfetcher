package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"scanfetcher/internal/domain"
	"scanfetcher/internal/usecase"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: uc}
}

func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.userUseCase.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type SignupInput struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input SignupInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Vérifie si l'utilisateur existe déjà
	if _, err := h.userUseCase.GetUserByEmail(input.Email); err == nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Crée un user domaine
	user := domain.User{
		Name:  input.Name,
		Email: input.Email,
	}

	// Crée l'utilisateur
	if err := h.userUseCase.Create(user, input.Password); err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Réponse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user) // ou nil si tu veux pas renvoyer les données
}
