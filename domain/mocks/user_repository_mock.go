package mocks

import (
	"cafe/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Create(c context.Context, user *domain.User) error {
	args := m.Called(c, user)
	return args.Error(0)
}

func (m *UserRepository) Fetch(c context.Context) ([]domain.User, error) {
	args := m.Called(c)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *UserRepository) GetByEmail(c context.Context, email string) (domain.User, error) {
	args := m.Called(c, email)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserRepository) GetByID(c context.Context, id string) (domain.User, error) {
	args := m.Called(c, id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserRepository) AssignRole(c context.Context, user *domain.User, role *domain.Role) error {
	args := m.Called(c, user, role)
	return args.Error(0)
}
