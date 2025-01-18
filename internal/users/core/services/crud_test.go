package services_test

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"styl-monolith/internal/users/core/domain"
	"styl-monolith/internal/users/core/services"
	"testing"
)

// Mock implementation for ports.UserPort
type MockUserPort struct {
	mock.Mock
}

func (m *MockUserPort) CreateUser(user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserPort) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserPort) GetUserByID(id uint) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserPort) GetUserByEmail(email string) (domain.User, error) {
	args := m.Called(email)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserPort) GetUserByUsername(username string) (domain.User, error) {
	args := m.Called(username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserPort) GetUserByPhone(phone string) (domain.User, error) {
	args := m.Called(phone)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserPort) ListUsers(offset, limit int) ([]domain.User, int, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]domain.User), args.Int(1), args.Error(2)
}

func (m *MockUserPort) UpdateUser(user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserPort) PartialUserUpdate(id uint, user domain.User) (domain.User, error) {
	args := m.Called(id, user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *MockUserPort) UpdateUserRole(user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

// Mock implementation for ports.RolePort
type MockRolePort struct {
	mock.Mock
}

func (m *MockRolePort) CreateRole(role domain.Role) (domain.Role, error) {
	args := m.Called(role)
	return args.Get(0).(domain.Role), args.Error(1)
}

func (m *MockRolePort) DeleteRole(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRolePort) GetRoleByID(id uint) (domain.Role, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Role), args.Error(1)
}

func (m *MockRolePort) GetRoleByName(name string) (domain.Role, error) {
	args := m.Called(name)
	return args.Get(0).(domain.Role), args.Error(1)
}

func (m *MockRolePort) ListRoles(offset, limit int) ([]domain.Role, int, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]domain.Role), args.Int(1), args.Error(2)
}

func (m *MockRolePort) UpdateRole(role domain.Role) (domain.Role, error) {
	args := m.Called(role)
	return args.Get(0).(domain.Role), args.Error(1)
}

func (m *MockRolePort) PartialRoleUpdate(id uint, role domain.Role) (domain.Role, error) {
	args := m.Called(id, role)
	return args.Get(0).(domain.Role), args.Error(1)
}

func TestCrudService_CreateUser(t *testing.T) {
	mockUserPort := new(MockUserPort)
	mockRolePort := new(MockRolePort)
	logger := logrus.New()
	service := services.NewCrudService(logger, mockUserPort, mockRolePort)

	requestUser := domain.User{ID: 1, Name: "John Doe", Email: "john@example.com"}
	expectedUser := domain.User{ID: 1, Name: "John Doe", Email: "john@example.com"}

	mockUserPort.On("CreateUser", requestUser).Return(expectedUser, nil)

	result, err := service.CreateUser(requestUser)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	mockUserPort.AssertExpectations(t)
}

func TestCrudService_CreateRole(t *testing.T) {
	mockUserPort := new(MockUserPort)
	mockRolePort := new(MockRolePort)
	logger := logrus.New()
	service := services.NewCrudService(logger, mockUserPort, mockRolePort)

	requestRole := domain.Role{ID: 1, Name: "Admin"}
	expectedRole := domain.Role{ID: 1, Name: "Admin"}

	mockRolePort.On("CreateRole", requestRole).Return(expectedRole, nil)

	result, err := service.CreateRole(requestRole)

	assert.NoError(t, err)
	assert.Equal(t, expectedRole, result)
	mockRolePort.AssertExpectations(t)
}

func TestCrudService_DeleteUserById(t *testing.T) {
	mockUserPort := new(MockUserPort)
	mockRolePort := new(MockRolePort)
	logger := logrus.New()
	service := services.NewCrudService(logger, mockUserPort, mockRolePort)

	userId := uint(1)

	mockUserPort.On("DeleteUser", userId).Return(nil)

	err := service.DeleteUserById(userId)

	assert.NoError(t, err)
	mockUserPort.AssertExpectations(t)
}

func TestCrudService_DeleteRoleById(t *testing.T) {
	mockUserPort := new(MockUserPort)
	mockRolePort := new(MockRolePort)
	logger := logrus.New()
	service := services.NewCrudService(logger, mockUserPort, mockRolePort)

	roleId := uint(1)

	mockRolePort.On("DeleteRole", roleId).Return(nil)

	err := service.DeleteRoleById(roleId)

	assert.NoError(t, err)
	mockRolePort.AssertExpectations(t)
}

func TestCrudService_GetUserById(t *testing.T) {
	mockUserPort := new(MockUserPort)
	mockRolePort := new(MockRolePort)
	logger := logrus.New()
	service := services.NewCrudService(logger, mockUserPort, mockRolePort)

	userId := uint(1)
	expectedUser := domain.User{ID: userId, Name: "John Doe", Email: "john@example.com"}

	mockUserPort.On("GetUserByID", userId).Return(expectedUser, nil)

	result, err := service.GetUserById(userId)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	mockUserPort.AssertExpectations(t)
}

func TestCrudService_ListUsers(t *testing.T) {
	mockUserPort := new(MockUserPort)
	mockRolePort := new(MockRolePort)
	logger := logrus.New()
	service := services.NewCrudService(logger, mockUserPort, mockRolePort)

	offset, limit := 0, 10
	users := []domain.User{{ID: 1, Name: "John Doe"}, {ID: 2, Name: "Jane Doe"}}

	mockUserPort.On("ListUsers", offset, limit).Return(users, 2, nil)

	result, count, err := service.ListUsers(offset, limit)

	assert.NoError(t, err)
	assert.Equal(t, 2, count)
	assert.Equal(t, users, result)
	mockUserPort.AssertExpectations(t)
}

func TestCrudService_UpdateUserById(t *testing.T) {
	mockUserPort := new(MockUserPort)
	mockRolePort := new(MockRolePort)
	logger := logrus.New()
	service := services.NewCrudService(logger, mockUserPort, mockRolePort)

	userId := uint(1)
	oldUser := domain.User{ID: userId, Name: "John Doe"}
	updatedUser := domain.User{ID: userId, Name: "John Smith"}

	mockUserPort.On("GetUserByID", userId).Return(oldUser, nil)
	mockUserPort.On("UpdateUser", updatedUser).Return(updatedUser, nil)

	result, err := service.UpdateUserById(userId, updatedUser)

	assert.NoError(t, err)
	assert.Equal(t, updatedUser, result)
	mockUserPort.AssertExpectations(t)
}
