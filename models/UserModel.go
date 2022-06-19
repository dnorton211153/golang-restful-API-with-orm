package models

// create User model
type User struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
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
}
