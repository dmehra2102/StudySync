package service

import (
	"github.com/dmehra2102/StudySync/internal/domain"
	"github.com/dmehra2102/StudySync/pkg/auth"
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
