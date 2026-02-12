package service

import (
	"errors"

	"go-warehouse-ms/internal/model"
	"go-warehouse-ms/internal/repository"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrPasswordWrong  = errors.New("password wrong")
	ErrUserExists     = errors.New("user exists")
	ErrInvalidPayload = errors.New("invalid payload")
)

type AuthService struct {
	users *repository.UserRepository
}

func NewAuthService(users *repository.UserRepository) *AuthService {
	return &AuthService{users: users}
}

func (s *AuthService) Login(userID, password string) (*model.User, error) {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	if user.UserPwd != password {
		return nil, ErrPasswordWrong
	}
	return user, nil
}

func (s *AuthService) Register(user *model.User) error {
	if user == nil || user.UserID == "" || user.UserPwd == "" {
		return ErrInvalidPayload
	}
	exists, err := s.users.Exists(user.UserID)
	if err != nil {
		return err
	}
	if exists {
		return ErrUserExists
	}
	return s.users.Create(user)
}
