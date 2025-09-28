package service

import (
	"context"
	"errors"

	"github.com/Jrozo97/reminderapp-backend/internal/domain"
	"github.com/Jrozo97/reminderapp-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) Register(ctx context.Context, name, email, password string) error {
	// verificar si ya existe
	_, err := s.repo.FindByEmail(ctx, email)
	if err == nil {
		return errors.New("el usuario ya existe")
	}

	// encriptar password
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
	}

	return s.repo.CreateUser(ctx, user)
}
