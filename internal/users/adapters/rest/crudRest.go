package rest

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"styl-monolith/internal/users/core/services"
)

type CrudRest interface {
	CreateUser(name, email, country, lastIp string)
	CreateRole()
	DeleteUser()
	DeleteRole()
	UpdateUser()
	UpdateRole()
	GetUser(ech echo.Context, id uint) error
	GetRole(ech echo.Context, id uint) error
	ListUsers(ech echo.Context) error
	ListRoles()
}

type crudRestStruct struct {
	crudService *services.CrudService
	logger      *logrus.Logger
}

func NewCrudRestUser(crudService *services.CrudService, logger *logrus.Logger) CrudRest {
	return &crudRestStruct{crudService: crudService, logger: logger}
}

func (c *crudRestStruct) CreateUser(name, email, country, lastIp string) {
	c.logger.Info("CreateUser method called")
	// Implementation goes here
}

func (c *crudRestStruct) CreateRole() {
	c.logger.Info("CreateRole method called")
	// Implementation goes here
}

func (c *crudRestStruct) DeleteUser() {
	c.logger.Info("DeleteUser method called")
	// Implementation goes here
}

func (c *crudRestStruct) DeleteRole() {
	c.logger.Info("DeleteRole method called")
	// Implementation goes here
}

func (c *crudRestStruct) UpdateUser() {
	c.logger.Info("UpdateUser method called")
	// Implementation goes here
}

func (c *crudRestStruct) UpdateRole() {
	c.logger.Info("UpdateRole method called")
	// Implementation goes here
}

func (c *crudRestStruct) GetUser(ech echo.Context, id uint) error {
	c.logger.Info("GetUser method called")
	return ech.String(http.StatusInternalServerError, fmt.Sprintf("Hello, role %d", id))
}

func (c *crudRestStruct) GetRole(ech echo.Context, id uint) error {
	c.logger.Info("GetRole method called")
	return ech.String(http.StatusOK, fmt.Sprintf("Hello, role %d", id))
	// Implementation goes here
}

func (c *crudRestStruct) ListUsers(ech echo.Context) error {
	c.logger.Info("ListUsers method called")
	return ech.String(http.StatusInternalServerError, "Hello, world")
}

func (c *crudRestStruct) ListRoles() {
	c.logger.Info("ListRoles method called")
	// Implementation goes here
}
