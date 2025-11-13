package service

import (
	"errors"

	"github.com/dmehra2102/StudySync/internal/domain"
	"github.com/dmehra2102/StudySync/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo domain.UserRepository
	jwtAuth  *auth.JWTAuth
}

func NewUserService(userRepo domain.UserRepository, jwtAuth *auth.JWTAuth) *UserService {
	return &UserService{
		userRepo: userRepo,
		jwtAuth:  jwtAuth,
	}
}

type RegisterRequest struct {
	Email     string `json:"email" bidning:"required,email"`
	Password  string `json:"password" bidning:"required,min=6"`
	FirstName string `json:"first_name" bidning:"required"`
	LastName  string `json:"last_name" bidning:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User  *domain.User `json:"user"`
	Token string       `json:"token"`
}

func (s *UserService) Register(req RegisterRequest) (*AuthResponse, error) {
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists with this email")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "student",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := s.jwtAuth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *UserService) Login(req LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwtAuth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		User:  user,
		Token: token,
	}, nil
}

func (s *UserService) GetProfile(userID uint) (*domain.User, error) {
	return s.userRepo.FindByID(userID)
}
