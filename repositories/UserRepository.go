package repositories

import (
	"restAPI/models"

	"github.com/jinzhu/gorm"
)

// UserRepo implements models.UserRepository
type BaseRepository struct {
	db *gorm.DB
}

// newUserRepository
func NewUserRepository(db *gorm.DB) *BaseRepository {

	// create the user table in the database
	db.AutoMigrate(&models.User{})
	return &BaseRepository{
		db: db,
	}
}

// Create a new user
func (r *BaseRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetAllUsers returns all users
func (r *BaseRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

// getUserByID returns a user by id
func (r *BaseRepository) GetUserByID(id string) (models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return user, err
}

// UpdateUser updates a user
func (r *BaseRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

// DeleteUser deletes a user
func (r *BaseRepository) DeleteUser(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.User{}).Error
}

func (r *BaseRepository) GetUserByUsernameAndPassword(username string, password string) (models.User, error) {
	var user models.User
	err := r.db.Where("username = ? AND password = ?", username, password).First(&user).Error
	return user, err
}
