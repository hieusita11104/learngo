package domain

import "fmt"

type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) Validate() error {
	if len(u.Email) == 0 {
		return fmt.Errorf("email is required")
	}

	if len(u.Password) < 6 {
		return fmt.Errorf("password must be at least 6 characters long")
	}

	return nil
}

func (u *User) ValidateLogin() error {
	if len(u.Email) == 0 {
		return fmt.Errorf("email is required")
	}
	if len(u.Password) == 0 {
		return fmt.Errorf("password is required")
	}
	return nil
}
