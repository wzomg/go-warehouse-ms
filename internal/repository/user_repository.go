package repository

import (
	"go-warehouse-ms/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(userID string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("userid = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Exists(userID string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.User{}).Where("userid = ?", userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}
