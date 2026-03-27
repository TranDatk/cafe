package usecase

import (
	"cafe/domain"
	"cafe/domain/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate_Success(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockRoleRepo := new(mocks.RoleRepository)
	mockRowTx := new(mocks.Transaction)

	user := &domain.User{
		ID:    "1",
		Name:  "Test User",
		Email: "test@example.com",
	}

	role := domain.Role{
		ID:   "role-1",
		Name: domain.UserRole,
	}

	// Expect WithinTransaction to be called and succeed
	mockRowTx.On("WithinTransaction", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(context.Context) error)
		_ = fn(context.Background())
	})

	mockUserRepo.On("Create", mock.Anything, user).Return(nil)
	mockRoleRepo.On("GetByName", mock.Anything, domain.UserRole).Return(role, nil)
	mockUserRepo.On("AssignRole", mock.Anything, user, &role).Return(nil)

	usecase := NewSignupUsecase(mockUserRepo, mockRoleRepo, mockRowTx, time.Second*2)

	err := usecase.Create(context.Background(), user)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockRoleRepo.AssertExpectations(t)
	mockRowTx.AssertExpectations(t)
}

func TestCreate_UserCreateError(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockRoleRepo := new(mocks.RoleRepository)
	mockRowTx := new(mocks.Transaction)

	user := &domain.User{
		ID:    "1",
		Name:  "Test User",
		Email: "test@example.com",
	}

	expectedErr := errors.New("db error")

	// Expect WithinTransaction to return the error from the function
	mockRowTx.On("WithinTransaction", mock.Anything, mock.Anything).Return(expectedErr)

	usecase := NewSignupUsecase(mockUserRepo, mockRoleRepo, mockRowTx, time.Second*2)

	err := usecase.Create(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestCreate_RoleNotFoundError(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	mockRoleRepo := new(mocks.RoleRepository)
	mockRowTx := new(mocks.Transaction)

	user := &domain.User{
		ID:    "1",
		Name:  "Test User",
		Email: "test@example.com",
	}

	expectedErr := errors.New("role not found")

	mockRowTx.On("WithinTransaction", mock.Anything, mock.Anything).Return(expectedErr)

	usecase := NewSignupUsecase(mockUserRepo, mockRoleRepo, mockRowTx, time.Second*2)

	err := usecase.Create(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
}
