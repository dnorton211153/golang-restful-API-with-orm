package controllers

// import User model
import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"restAPI/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// BaseHandler will hold everything that controller needs
type BaseHandler struct {
	userRepository models.UserRepository
	sessions       map[string]Session
}

type Session struct {
	username string
	expiry   time.Time
}

func Expired(s Session) bool {
	return s.expiry.Before(time.Now())
}

// NewHandler returns a new BaseHandler
func NewHandler(userRepository models.UserRepository, sessions map[string]Session) *BaseHandler {
	return &BaseHandler{
		userRepository: userRepository,
		sessions:       sessions,
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

	idInt, err := strconv.Atoi(id)
	if idInt != user.ID {
		http.Error(w, "id in URL (e.g. /user/2) must match id in request body (e.g. {\"id\":2,...})", http.StatusBadRequest)
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

// Create a struct that models the structure of a user in the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func (h *BaseHandler) Login(w http.ResponseWriter, r *http.Request) {

	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.userRepository.GetUserByUsernameAndPassword(creds.Username, creds.Password)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	h.sessions[sessionToken] = Session{
		username: user.Username,
		expiry:   expiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
