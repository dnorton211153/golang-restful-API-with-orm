package controllers

// import User model
import (
	"encoding/json"
	"net/http"

	"restAPI/models"

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
