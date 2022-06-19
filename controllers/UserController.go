package controllers

// import User model
import (
	"encoding/json"
	"net/http"

	"restAPI/models"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

// BaseHandler will hold everything that controller needs
type BaseHandler struct {
	userRepository models.UserRepository
}

// NewHandler returns a new BaseHandler
func NewHandler(userRepository models.UserRepository) *BaseHandler {
	return &BaseHandler{
		userRepository: userRepository,
	}
}

func (h *BaseHandler) Create(w http.ResponseWriter, r *http.Request) {

	// decode the request body into a new user struct
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// save the user to the database
	err = h.userRepository.Create(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Get
func (h *BaseHandler) Get(w http.ResponseWriter, r *http.Request) {
	// get the user from the database
	var vars = mux.Vars(r)
	var id = vars["id"]
	user, err := h.userRepository.GetUserByID(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// GetAll
func (h *BaseHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// get all users from the database
	users, err := h.userRepository.GetAllUsers()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the users
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// Update
func (h *BaseHandler) Update(w http.ResponseWriter, r *http.Request) {

	var id = mux.Vars(r)["id"]

	// decode the request body into a new user struct
	user := models.User{}
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if id != user.ID {
		http.Error(w, "ID in URL must match ID in request body", http.StatusBadRequest)
		return
	}

	// update the user in the database
	err = h.userRepository.UpdateUser(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the updated user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Delete
func (h *BaseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// delete the user from the database
	var id = mux.Vars(r)["id"]
	err := h.userRepository.DeleteUser(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return a status notifying the client the user was deleted
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted"})
}
