package models

// create User model
type User struct {
	// auto increment id
	ID       int    `gorm:"primary_key,autoIncrement" json:"id"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// UserRepository ..
type UserRepository interface {
	Create(user *User) error
	GetAllUsers() ([]User, error)
	GetUserByID(id string) (User, error)
	UpdateUser(user *User) error
	DeleteUser(id string) error
	GetUserByUsernameAndPassword(username string, password string) (User, error)
}
