package usecase

import (
	"errors"
	"time"

	"learngo/internal/domain"
	"learngo/pkg/token"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	Repo domain.UserRepository
}

func NewUserUsecase(r domain.UserRepository) *UserUsecase {
	return &UserUsecase{Repo: r}
}

func (uc *UserUsecase) CreateUser(u *domain.User) error {
	// Check if email exists
	existing, err := uc.Repo.GetByEmail(u.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("email already exists")
	}
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}
	// Save to database
	u.Password = string(hashedPassword)
	return uc.Repo.Create(u)
}

func (uc *UserUsecase) LoginUser(u *domain.User) (string, string, error) {
	// Validate user
	existing, err := uc.Repo.GetByEmail(u.Email)
	if err != nil {
		return "", "", err
	}
	if existing == nil {
		return "", "", errors.New("invalid email or password")
	}
	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(existing.Password), []byte(u.Password)); err != nil {
		return "", "", errors.New("invalid email or password")
	}

	accessToken, _ := token.GenerateToken(existing.ID, 15*time.Minute)
	refreshToken, _ := token.GenerateToken(existing.ID, 7*24*time.Hour)

	return accessToken, refreshToken, nil
}

func (uc *UserUsecase) RefreshToken(userID uint) (string, error) {
	accessToken, err := token.GenerateToken(userID, 15*time.Minute)
	if err != nil {
		return "", errors.New("failed to generate access token")
	}

	return accessToken, nil
}
