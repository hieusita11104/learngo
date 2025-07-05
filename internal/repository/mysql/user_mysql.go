package mysql

import (
	"learngo/internal/domain"

	"gorm.io/gorm"
)

type UserModel struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Email    string `gorm:"uniqueIndex;type:varchar(255);not null"`
	Password string `gorm:"type:varchar(255);not null"`
}

func (UserModel) TableName() string {
	return "users"
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
		ID:       e.ID,
		Email:    e.Email,
		Password: e.Password,
	}
}

type UserMySQLRepo struct {
	DB *gorm.DB
}

func NewUserMySQLRepo(db *gorm.DB) domain.UserRepository {
	if err := db.AutoMigrate(&UserModel{}); err != nil {
		return nil
	}
	return &UserMySQLRepo{DB: db}
}

func (r *UserMySQLRepo) Create(u *domain.User) error {
	model := toUserModel(u)
	if err := r.DB.Create(model).Error; err != nil {
		return err
	}
	u.ID = model.ID
	return nil
}

func (r *UserMySQLRepo) GetByEmail(email string) (*domain.User, error) {
	var model UserModel
	if err := r.DB.Where("email = ?", email).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No user found
		}
		return nil, err // Other error
	}
	return toUserEntity(&model), nil
}
