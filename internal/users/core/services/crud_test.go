package services_test

import (
	"errors"
	"styl-monolith/internal/users/core/services"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"styl-monolith/internal/users/core/domain"
)

// Mock implementations of the ports
type mockUserPort struct {
	mock.Mock
}

func (m *mockUserPort) CreateUser(user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *mockUserPort) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockUserPort) UpdateUser(user domain.User) (domain.User, error) {
	args := m.Called(user)
	return args.Get(0).(domain.User), args.Error(1)
}

// ... other UserPort methods mocked as needed

type mockRolePort struct {
	mock.Mock
}

func (m *mockRolePort) CreateRole(role domain.Role) (domain.Role, error) {
	args := m.Called(role)
	return args.Get(0).(domain.Role), args.Error(1)
}

func (m *mockRolePort) GetRoleByName(name string) (domain.Role, error) {
	args := m.Called(name)
	return args.Get(0).(domain.Role), args.Error(1)
}

// ... other RolePort methods mocked as needed

type mockUserReportRepository struct {
	mock.Mock
}

// ... mock methods for UserReportRepository as needed

func TestCreateUser(t *testing.T) {
	logger := logrus.New()
	userRepoMock := new(mockUserPort)
	roleRepoMock := new(mockRolePort)
	reportRepoMock := new(mockUserReportRepository)

	service := services.NewCrudService(logger, userRepoMock, roleRepoMock, reportRepoMock)

	requestUser := domain.User{
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "johndoe@example.com",
	}

	defaultRole := domain.Role{
		ID:   1,
		Name: "default",
	}

	expectedUser := domain.User{
		ID:       1,
		Name:     "John Doe",
		Username: "johndoe",
		Email:    "johndoe@example.com",
		Role:     defaultRole,
	}

	// Mocking expectations
	roleRepoMock.On("GetRoleByName", "default").Return(defaultRole, nil)
	userRepoMock.On("CreateUser", mock.Anything).Return(expectedUser, nil)

	// Test CreateUser method
	createdUser, err := service.CreateUser(requestUser)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, createdUser)
	roleRepoMock.AssertExpectations(t)
	userRepoMock.AssertExpectations(t)
}

func TestDeleteUserById(t *testing.T) {
	logger := logrus.New()
	userRepoMock := new(mockUserPort)
	roleRepoMock := new(mockRolePort)
	reportRepoMock := new(mockUserReportRepository)

	service := services.NewCrudService(logger, userRepoMock, roleRepoMock, reportRepoMock)

	userId := uint(1)

	// Mocking expectations
	userRepoMock.On("DeleteUser", userId).Return(nil)

	// Test DeleteUserById method
	err := service.DeleteUserById(userId)

	// Assertions
	assert.NoError(t, err)
	userRepoMock.AssertExpectations(t)
}

func TestGetUserById(t *testing.T) {
	logger := logrus.New()
	userRepoMock := new(mockUserPort)
	roleRepoMock := new(mockRolePort)
	reportRepoMock := new(mockUserReportRepository)

	service := services.NewCrudService(logger, userRepoMock, roleRepoMock, reportRepoMock)

	userId := uint(1)

	role := domain.Role{ID: 1, Name: "admin"}
	user := domain.User{ID: userId, Name: "John", Role: role}

	// Mocking expectations
	userRepoMock.On("GetUserByID", userId).Return(user, nil)
	roleRepoMock.On("GetRoleByID", role.ID).Return(role, nil)

	// Test GetUserById method
	resultUser, err := service.GetUserById(userId)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, user, resultUser)
	userRepoMock.AssertExpectations(t)
	roleRepoMock.AssertExpectations(t)
}

func TestCreateRole(t *testing.T) {
	logger := logrus.New()
	userRepoMock := new(mockUserPort)
	roleRepoMock := new(mockRolePort)
	reportRepoMock := new(mockUserReportRepository)

	service := services.NewCrudService(logger, userRepoMock, roleRepoMock, reportRepoMock)

	inputRole := domain.Role{
		Name:        "admin",
		Description: "Administrator role",
	}

	expectedRole := domain.Role{
		ID:          1,
		Name:        "admin",
		Description: "Administrator role",
	}

	// Mocking expectations
	roleRepoMock.On("CreateRole", inputRole).Return(expectedRole, nil)

	// Test CreateRole method
	createdRole, err := service.CreateRole(inputRole)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedRole, createdRole)
	roleRepoMock.AssertExpectations(t)
}

func TestUpdateRoleById_ErrorOnGetRole(t *testing.T) {
	logger := logrus.New()
	userRepoMock := new(mockUserPort)
	roleRepoMock := new(mockRolePort)
	reportRepoMock := new(mockUserReportRepository)

	service := services.NewCrudService(logger, userRepoMock, roleRepoMock, reportRepoMock)

	roleId := uint(1)
	inputRole := domain.Role{Name: "admin", Description: "Updated description"}

	// Mocking expectations
	roleRepoMock.On("GetRoleByID", roleId).Return(domain.Role{}, errors.New("not found"))

	// Test UpdateRoleById method
	_, err := service.UpdateRoleById(roleId, inputRole)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "not found", err.Error())
	roleRepoMock.AssertExpectations(t)
}
