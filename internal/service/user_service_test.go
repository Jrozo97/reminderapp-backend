package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Jrozo97/reminderapp-backend/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Mock del UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestRegister_NewUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// Caso: usuario no existe → se crea
	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").
		Return(nil, errors.New("not found"))
	mockRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*domain.User")).
		Return(nil)

	err := service.Register(context.Background(), "Test User", "test@example.com", "password123")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRegister_ExistingUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// Caso: usuario ya existe
	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").
		Return(&domain.User{Email: "test@example.com"}, nil)

	err := service.Register(context.Background(), "Test User", "test@example.com", "password123")

	assert.Error(t, err)
	assert.Equal(t, "el usuario ya existe", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestLogin_CorrectCredentials(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// Crear hash del password correcto
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	existingUser := &domain.User{
		Email:    "test@example.com",
		Password: string(hashed),
	}

	// Caso: credenciales correctas
	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").
		Return(existingUser, nil)

	// Set secret para JWT
	t.Setenv("JWT_SECRET", "testsecret")

	token, err := service.Login(context.Background(), "test@example.com", "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	// Usuario con otro password
	hashed, _ := bcrypt.GenerateFromPassword([]byte("otherpass"), bcrypt.DefaultCost)
	existingUser := &domain.User{
		Email:    "test@example.com",
		Password: string(hashed),
	}

	mockRepo.On("FindByEmail", mock.Anything, "test@example.com").
		Return(existingUser, nil)

	t.Setenv("JWT_SECRET", "testsecret")

	token, err := service.Login(context.Background(), "test@example.com", "wrongpass")

	assert.Error(t, err)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}
