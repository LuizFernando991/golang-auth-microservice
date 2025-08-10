package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/LuizFernando991/golang-auth-microservice/internal/model"
	"github.com/LuizFernando991/golang-auth-microservice/internal/repository"
	"github.com/LuizFernando991/golang-auth-microservice/internal/service"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(ctx context.Context, u *model.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}
func (m *MockRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	u := args.Get(0)
	if u == nil {
		return nil, args.Error(1)
	}
	return u.(*model.User), args.Error(1)
}
func (m *MockRepo) FindByID(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	u := args.Get(0)
	if u == nil {
		return nil, args.Error(1)
	}
	return u.(*model.User), args.Error(1)
}
func (m *MockRepo) SaveRefreshToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, token, expiresAt)
	return args.Error(0)
}
func (m *MockRepo) FindRefreshToken(ctx context.Context, token string) (*repository.RefreshTokenRow, error) {
	return nil, errors.New("not implemented")
}
func (m *MockRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}
func (m *MockRepo) DeleteAllRefreshTokensForUser(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestRegisterSuccess(t *testing.T) {
	repo := &MockRepo{}
	repo.On("FindByEmail", mock.Anything, "new@example.com").Return(nil, errors.New("not found"))
	repo.On("Create", mock.Anything, mock.Anything).Return(nil)

	svc := service.NewAuthService(repo, "secret", time.Minute*15, time.Hour*24*7, 4)
	u, err := svc.Register(context.Background(), model.CreateUserDTO{Email: "new@example.com", Password: "12345678"})
	assert.NoError(t, err)
	assert.Equal(t, "new@example.com", u.Email)
	repo.AssertExpectations(t)
}
