package service

import (
	"github.com/jixlox0/studoto-backend/internal/models"
	"github.com/jixlox0/studoto-backend/internal/repository"
	"github.com/jixlox0/studoto-backend/pkg/auth"
)

type UserService interface {
	GetUserByID(id uint) (*models.UserResponse, error)
	GetUserByEmail(email string) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	jwtAuth  *auth.JWTAuth
}

func NewUserService(userRepo repository.UserRepository, jwtAuth *auth.JWTAuth) UserService {
	return &userService{
		userRepo: userRepo,
		jwtAuth:  jwtAuth,
	}
}

func (s *userService) GetUserByID(id uint) (*models.UserResponse, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.FindByEmail(email)
}
