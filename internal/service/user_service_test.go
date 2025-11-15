package service_test

import (
	"testing"

	"github.com/dmehra2102/StudySync/internal/domain"
	"github.com/dmehra2102/StudySync/internal/service"
	"github.com/dmehra2102/StudySync/pkg/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserService_Register(t *testing.T) {
	jwtAuth := auth.NewJWTAuth("test-secret", 24*60*60)

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := service.NewUserService(mockRepo, jwtAuth)

		req := service.RegisterRequest{
			Email:     "test@example.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
		}

		mockRepo.On("FindByEmail", req.Email).Return((*domain.User)(nil), assert.AnError)
		mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

		response, err := userService.Register(req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.Token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserAlreadyExists", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := service.NewUserService(mockRepo, jwtAuth)

		req := service.RegisterRequest{
			Email:     "existing@example.com",
			Password:  "password123",
			FirstName: "John",
			LastName:  "Doe",
		}

		existingUser := &domain.User{Email: req.Email}
		mockRepo.On("FindByEmail", req.Email).Return(existingUser, nil)

		response, err := userService.Register(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "user already exists with this email", err.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Login(t *testing.T) {
	jwtAuth := auth.NewJWTAuth("test-secret", 24*60*60)

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := service.NewUserService(mockRepo, jwtAuth)

		req := service.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		user := &domain.User{
			Email:    req.Email,
			Password: string(hashedPassword),
		}

		mockRepo.On("FindByEmail", req.Email).Return(user, nil)

		response, err := userService.Login(req)

		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.Token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		userService := service.NewUserService(mockRepo, jwtAuth)

		req := service.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

		user := &domain.User{
			Email:    req.Email,
			Password: string(hashedPassword),
		}

		mockRepo.On("FindByEmail", req.Email).Return(user, nil)

		response, err := userService.Login(req)

		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "invalid credentials", err.Error())
		mockRepo.AssertExpectations(t)
	})
}
