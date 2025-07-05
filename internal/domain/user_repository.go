package domain

type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
}
