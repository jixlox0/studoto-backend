package service

import (
	"errors"
	"math/rand"
	"time"

	"github.com/jixlox0/studoto-backend/internal/models"
	"github.com/jixlox0/studoto-backend/internal/repository"
	"github.com/jixlox0/studoto-backend/pkg/auth"
	"github.com/jixlox0/studoto-backend/pkg/oauth"
	"github.com/jixlox0/studoto-backend/pkg/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req *models.CreateUserRequest) (*models.AuthResponse, error)
	Login(req *models.LoginRequest) (*models.AuthResponse, error)
	OAuthLogin(provider, code string) (*models.AuthResponse, error)
	GetOAuthURL(provider string) (string, error)
}

type authService struct {
	userRepo     repository.UserRepository
	jwtAuth      *auth.JWTAuth
	oauthService oauth.OAuthService
}

func NewAuthService(userRepo repository.UserRepository, jwtAuth *auth.JWTAuth, oauthService oauth.OAuthService) AuthService {
	return &authService{
		userRepo:     userRepo,
		jwtAuth:      jwtAuth,
		oauthService: oauthService,
	}
}

func (s *authService) Register(req *models.CreateUserRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		UUID:         uuid.Generate(uuid.PrefixUser),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// Generate token
	token, err := s.jwtAuth.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *authService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, err := s.jwtAuth.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *authService) GetOAuthURL(provider string) (string, error) {
	state := generateState()

	switch provider {
	case "google":
		return s.oauthService.GetGoogleAuthURL(state), nil
	case "github":
		return s.oauthService.GetGitHubAuthURL(state), nil
	default:
		return "", errors.New("unsupported OAuth provider")
	}
}

func (s *authService) OAuthLogin(provider, code string) (*models.AuthResponse, error) {
	var oauthUser *oauth.OAuthUser
	var err error

	switch provider {
	case "google":
		oauthUser, err = s.oauthService.ExchangeGoogleCode(code)
	case "github":
		oauthUser, err = s.oauthService.ExchangeGitHubCode(code)
	default:
		return nil, errors.New("unsupported OAuth provider")
	}

	if err != nil {
		return nil, err
	}

	// Check if user exists by provider
	user, err := s.userRepo.FindByProvider(provider, oauthUser.ID)
	if err != nil {
		// User doesn't exist, create new user
		user = &models.User{
			UUID:       uuid.Generate(uuid.PrefixUser),
			Email:      oauthUser.Email,
			Name:       oauthUser.Name,
			AvatarURL:  oauthUser.AvatarURL,
			Provider:   provider,
			ProviderID: oauthUser.ID,
		}

		if err := s.userRepo.Create(user); err != nil {
			return nil, err
		}
	} else {
		// Update user info if needed
		if user.AvatarURL != oauthUser.AvatarURL || user.Name != oauthUser.Name {
			user.AvatarURL = oauthUser.AvatarURL
			user.Name = oauthUser.Name
			s.userRepo.Update(user)
		}
	}

	// Generate token
	token, err := s.jwtAuth.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func generateState() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}
