package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type CrudPostHandler interface {
	PostCreatePost(ech echo.Context)
	PostCreateComment(ech echo.Context)
	PostCreateReport(ech echo.Context)
	PostUploadPostImages(ech echo.Context)
	GetPostById(ech echo.Context, id string)
	GetPostImagesById(ech echo.Context, id string)
	DeletePost(ech echo.Context, id string)
	DeleteComment(ech echo.Context, id string)
	DeleteReport(ech echo.Context, id string)
	DeleteImage(ech echo.Context, id string)
	UpdatePost(ech echo.Context, id string)
	UpdateComment(ech echo.Context, id string)
	GetListReports(ech echo.Context)
	GetListImages(ech echo.Context)
}

type CrudPostStruct struct {
	logger *logrus.Logger
}

func NewCrudPostHandler(logger *logrus.Logger) CrudPostHandler {
	return &CrudPostStruct{
		logger: logger,
	}
}

func (c CrudPostStruct) PostCreatePost(ech echo.Context) {
	//TODO implement me
	panic("implement me")
}

func (c CrudPostStruct) PostCreateComment(ech echo.Context) {
	//TODO implement me
	panic("implement me")
}

func (c CrudPostStruct) PostCreateReport(ech echo.Context) {
	//TODO implement me
	panic("implement me")
}

func (c CrudPostStruct) PostUploadPostImages(ech echo.Context) {
	//TODO implement me
	panic("implement me")
}

func (c CrudPostStruct) GetPostById(ech echo.Context, id string) {
	//TODO implement me
	panic("implement me")
}

func (c CrudPostStruct) GetPostImagesById(ech echo.Context, id string) {
	//TODO implement me
	panic("implement me")
}

func (c CrudPostStruct) DeletePost(ech echo.Context, id string) {
	//TODO implement me
	panic("implement me")
}

func (c CrudPostStruct) DeleteComment(ech echo.Context, id string) {
	//TODO implement me
	panic("implement me")
}

func (c CrudPostStruct) DeleteReport(ech echo.Context, id string) {
	//TODO implement me
	panic("implement me")
}

func (c CrudPostStruct) DeleteImage(ech echo.Context, id string) {
	//TODO implement me
	panic("implement me")
}

func (c CrudPostStruct) UpdatePost(ech echo.Context, id string) {
	//TODO implement me
	panic("implement me")
}

func (c CrudPostStruct) UpdateComment(ech echo.Context, id string) {
	//TODO implement me
	panic("implement me")
}

func (c CrudPostStruct) GetListReports(ech echo.Context) {
	//TODO implement me
	panic("implement me")
}

func (c CrudPostStruct) GetListImages(ech echo.Context) {
	//TODO implement me
	panic("implement me")
}
