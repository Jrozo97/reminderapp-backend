package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/Jrozo97/reminderapp-backend/internal/domain"
	"github.com/Jrozo97/reminderapp-backend/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo}
}

func (s *UserService) Register(ctx context.Context, name, email, password string) error {
	// verificar si ya existe
	_, err := s.repo.FindByEmail(ctx, email)
	if err == nil {
		return errors.New("el usuario ya existe")
	}

	// encriptar password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error al encriptar la contraseña")
	}

	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("usuario o contraseña incorrectos")
	}

	// Comparar password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New("usuario o contraseña incorrectos")
	}

	// Generar token JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET no configurado")
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID.Hex(),
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // expira en 24h
		"issuedAt": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
