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
