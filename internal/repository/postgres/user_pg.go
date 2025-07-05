package postgres

import (
	"learngo/internal/domain"
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex"`
	Password  string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	CreatedBy uint
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	UpdatedBy uint
}

func toUserEntity(m *UserModel) *domain.User {
	return &domain.User{
		ID:       m.ID,
		Email:    m.Email,
		Password: m.Password,
	}
}

func toUserModel(e *domain.User) *UserModel {
	return &UserModel{
		ID:       uint(e.ID),
		Email:    e.Email,
		Password: e.Password,
	}
}

type UserPostgresRepo struct {
	DB *gorm.DB
}

func NewUserPostgresRepo(db *gorm.DB) domain.UserRepository {
	if err := db.AutoMigrate(&UserModel{}); err != nil {
		panic(err)
	}
	return &UserPostgresRepo{DB: db}
}

// GetAll retrieves all users from the database.
func (r *UserPostgresRepo) GetAll() ([]*domain.User, error) {
	var models []UserModel
	if err := r.DB.Find(&models).Error; err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(models))
	for i, model := range models {
		users[i] = toUserEntity(&model)
	}
	return users, nil
}

func (r *UserPostgresRepo) Create(u *domain.User) error {
	model := toUserModel(u)
	if err := r.DB.Create(model).Error; err != nil {
		return err
	}
	u.ID = model.ID
	return nil
}

func (r *UserPostgresRepo) GetByEmail(email string) (*domain.User, error) {
	var model UserModel
	if err := r.DB.Where("email = ?", email).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No user found
		}
		return nil, err // Other error
	}
	return toUserEntity(&model), nil
}
