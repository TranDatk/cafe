package mocks

import (
	"cafe/domain"
	"context"

	"github.com/stretchr/testify/mock"
)

type RoleRepository struct {
	mock.Mock
}

func (m *RoleRepository) Create(c context.Context, role *domain.Role) error {
	args := m.Called(c, role)
	return args.Error(0)
}

func (m *RoleRepository) Fetch(c context.Context) ([]domain.Role, error) {
	args := m.Called(c)
	return args.Get(0).([]domain.Role), args.Error(1)
}

func (m *RoleRepository) GetByID(c context.Context, id string) (domain.Role, error) {
	args := m.Called(c, id)
	return args.Get(0).(domain.Role), args.Error(1)
}

func (m *RoleRepository) GetByName(c context.Context, name string) (domain.Role, error) {
	args := m.Called(c, name)
	return args.Get(0).(domain.Role), args.Error(1)
}
