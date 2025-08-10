package service

import (
	"context"
	"time"

	"github.com/LuizFernando991/golang-auth-microservice/internal/model"
	"github.com/LuizFernando991/golang-auth-microservice/internal/repository"
	"github.com/LuizFernando991/golang-auth-microservice/internal/util"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, dto model.CreateUserDTO) (*model.User, error)
	Login(ctx context.Context, dto model.LoginDTO) (accessToken string, refreshToken string, err error)
	Refresh(ctx context.Context, refreshToken string) (newAccess string, newRefresh string, err error)
	Logout(ctx context.Context, refreshToken string) error
	GetUserById(ctx context.Context, userID int64) (*model.User, error)
}

type authService struct {
	repo       repository.UserRepository
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
	bcryptCost int
}

func NewAuthService(repo repository.UserRepository, jwtSecret string, accessTTL, refreshTTL time.Duration, bcryptCost int) AuthService {
	return &authService{
		repo:       repo,
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		bcryptCost: bcryptCost,
	}
}

func (s *authService) Register(ctx context.Context, dto model.CreateUserDTO) (*model.User, error) {
	if existing, _ := s.repo.FindByEmail(ctx, dto.Email); existing != nil && existing.ID != 0 {
		return nil, util.ErrUserExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), s.bcryptCost)
	if err != nil {
		return nil, err
	}
	u := &model.User{
		Email:        dto.Email,
		PasswordHash: string(hash),
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return u, nil
}

func (s *authService) Login(ctx context.Context, dto model.LoginDTO) (string, string, error) {
	u, err := s.repo.FindByEmail(ctx, dto.Email)
	if err != nil || u == nil {
		return "", "", util.ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(dto.Password)); err != nil {
		return "", "", util.ErrInvalidCredentials
	}
	access, err := util.NewAccessToken(s.jwtSecret, s.accessTTL, u.ID)
	if err != nil {
		return "", "", err
	}
	refresh, err := util.GenerateSecureToken(32)
	if err != nil {
		return "", "", err
	}
	expires := time.Now().Add(s.refreshTTL)
	if err := s.repo.SaveRefreshToken(ctx, u.ID, refresh, expires); err != nil {
		return "", "", err
	}
	return access, refresh, nil
}

func (s *authService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	rt, err := s.repo.FindRefreshToken(ctx, refreshToken)
	if err != nil || rt == nil {
		return "", "", util.ErrRefreshTokenNotFound
	}
	if time.Now().After(rt.ExpiresAt) {
		_ = s.repo.DeleteRefreshToken(ctx, refreshToken)
		return "", "", util.ErrRefreshTokenNotFound
	}

	// Rotation: delete old and issue new
	if err := s.repo.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return "", "", err
	}

	access, err := util.NewAccessToken(s.jwtSecret, s.accessTTL, rt.UserID)
	if err != nil {
		return "", "", err
	}
	newRefresh, err := util.GenerateSecureToken(32)
	if err != nil {
		return "", "", err
	}
	expires := time.Now().Add(s.refreshTTL)
	if err := s.repo.SaveRefreshToken(ctx, rt.UserID, newRefresh, expires); err != nil {
		return "", "", err
	}
	return access, newRefresh, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	return s.repo.DeleteRefreshToken(ctx, refreshToken)
}

func (s *authService) GetUserById(ctx context.Context, userID int64) (*model.User, error) {
	u, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	u.PasswordHash = ""
	return u, nil
}
