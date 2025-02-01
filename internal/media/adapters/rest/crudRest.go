package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"styl-monolith/internal/media/adapters/rest/models"
	"styl-monolith/internal/media/core/services"
)

type CrudPostHandler interface {
	PostCreatePost(ech echo.Context) error
	PostCreateComment(ech echo.Context) error
	PostCreateReport(ech echo.Context) error
	PostUploadPostImages(ech echo.Context) error
	GetPostById(ech echo.Context, id string) error
	GetPostImagesById(ech echo.Context, id string) error
	DeletePost(ech echo.Context, id string) error
	DeleteComment(ech echo.Context, id string) error
	DeleteReport(ech echo.Context, id string) error
	DeleteImage(ech echo.Context, id string) error
	UpdatePost(ech echo.Context, id string) error
	UpdateComment(ech echo.Context, id string) error
	GetListReports(ech echo.Context) error
	GetListImages(ech echo.Context) error
}

type CrudPostStruct struct {
	logger       *logrus.Logger
	mediaService services.MediaService
}

func NewCrudPostHandler(logger *logrus.Logger, mediaService services.MediaService) CrudPostHandler {
	return &CrudPostStruct{
		logger:       logger,
		mediaService: mediaService,
	}
}

func (c CrudPostStruct) PostCreatePost(ech echo.Context) error {
	var payload models.PostPayload
	if err := ech.Bind(&payload); err != nil {
		c.logger.Error("Failed to bind PostPayload: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Bad request payload"})
	}

	if err := payload.Validate(); err != nil {
		c.logger.Error("Validation error: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Validation failed"})
	}

	post, err := c.mediaService.CreatePost(payload.ToDomainPost())
	if err != nil {
		c.logger.Error("Failed to create post: ", err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create post"})
	}
	return ech.JSON(http.StatusCreated, models.NewPostResponseFromDomain(post))
}

func (c CrudPostStruct) PostCreateComment(ech echo.Context) error {
	var payload models.CommentPayload
	if err := ech.Bind(&payload); err != nil {
		c.logger.Error("Failed to bind CommentPayload: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Bad request payload"})
	}

	if err := payload.Validate(); err != nil {
		c.logger.Error("Validation error: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Validation failed"})
	}

	comment, err := c.mediaService.CreateComment(payload.ToDomainComment())
	if err != nil {
		c.logger.Error("Failed to create comment: ", err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create comment"})
	}
	return ech.JSON(http.StatusCreated, models.NewCommentResponseFromDomain(comment))
}

func (c CrudPostStruct) PostCreateReport(ech echo.Context) error {
	var payload models.ReportPayload
	if err := ech.Bind(&payload); err != nil {
		c.logger.Error("Failed to bind ReportPayload: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Bad request payload"})
	}

	if err := payload.Validate(); err != nil {
		c.logger.Error("Validation error: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Validation failed"})
	}

	report, err := c.mediaService.CreateReport(payload.ToDomainReport())
	if err != nil {
		c.logger.Error("Failed to create report: ", err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create report"})
	}
	return ech.JSON(http.StatusCreated, models.NewReportResponseFromDomainReport(report))
}

func (c CrudPostStruct) PostUploadPostImages(ech echo.Context) error {
	var payload models.ImagePayload
	if err := ech.Bind(&payload); err != nil {
		c.logger.Error("Failed to bind ImagePayload: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Bad request payload"})
	}

	if err := payload.Validate(); err != nil {
		c.logger.Error("Validation error: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Validation failed"})
	}

	image, err := c.mediaService.UploadImage(payload.ToDomainImage())
	if err != nil {
		c.logger.Error("Failed to upload image: ", err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to upload image"})
	}
	return ech.JSON(http.StatusCreated, models.NewImageResponseFromDomain(image))
}

func (c CrudPostStruct) GetPostById(ech echo.Context, id string) error {
	postId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid post ID"})
	}

	post, err := c.mediaService.GetPost(uint(postId))
	if err != nil {
		c.logger.Error("Failed to retrieve post: ", err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get post"})
	}
	return ech.JSON(http.StatusOK, models.NewPostResponseFromDomain(post))
}

func (c CrudPostStruct) DeletePost(ech echo.Context, id string) error {
	postId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid post ID"})
	}

	if err := c.mediaService.DeletePost(uint(postId)); err != nil {
		c.logger.Error("Failed to delete post: ", err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete post"})
	}
	return ech.JSON(http.StatusNoContent, nil)
}

func (c CrudPostStruct) GetPostImagesById(ech echo.Context, id string) error {
	postId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.logger.Error("Invalid post ID: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid post ID"})
	}

	images, _, err := c.mediaService.ListImages(uint(postId), 1, 100)
	if err != nil {
		c.logger.Error("Failed to retrieve post images: ", err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve images"})
	}

	response := models.NewPaginatedImagesResponse(images, 1, 100)
	return ech.JSON(http.StatusOK, response)
}

func (c CrudPostStruct) DeleteComment(ech echo.Context, id string) error {
	commentId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.logger.Error("Invalid comment ID: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid comment ID"})
	}

	if err := c.mediaService.DeleteComment(uint(commentId)); err != nil {
		c.logger.Error("Failed to delete comment: ", err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete comment"})
	}
	return ech.JSON(http.StatusNoContent, nil)
}

func (c CrudPostStruct) DeleteReport(ech echo.Context, id string) error {
	reportId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.logger.Error("Invalid report ID: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid report ID"})
	}

	if err := c.mediaService.DeleteReport(uint(reportId)); err != nil {
		c.logger.Error("Failed to delete report: ", err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete report"})
	}
	return ech.JSON(http.StatusNoContent, nil)
}

func (c CrudPostStruct) DeleteImage(ech echo.Context, id string) error {
	imageKey := id // Assuming imageKey is passed as the ID
	if err := c.mediaService.DeleteImage(imageKey); err != nil {
		c.logger.Error("Failed to delete image: ", err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete image"})
	}
	return ech.JSON(http.StatusNoContent, nil)
}

func (c CrudPostStruct) UpdatePost(ech echo.Context, id string) error {
	var payload models.PostPayload
	if err := ech.Bind(&payload); err != nil {
		c.logger.Error("Failed to bind PostPayload: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Bad request payload"})
	}

	if err := payload.Validate(); err != nil {
		c.logger.Error("Validation error: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Validation failed"})
	}

	postId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.logger.Error("Invalid post ID: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid post ID"})
	}

	updatedPost, err := c.mediaService.UpdatePost(uint(postId), payload.ToDomainPost())
	if err != nil {
		c.logger.Error("Failed to update post: ", err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update post"})
	}
	return ech.JSON(http.StatusOK, models.NewPostResponseFromDomain(updatedPost))
}

func (c CrudPostStruct) UpdateComment(ech echo.Context, id string) error {
	var payload models.CommentPayload
	if err := ech.Bind(&payload); err != nil {
		c.logger.Error("Failed to bind CommentPayload: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Bad request payload"})
	}

	if err := payload.Validate(); err != nil {
		c.logger.Error("Validation error: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Validation failed"})
	}

	commentId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.logger.Error("Invalid comment ID: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid comment ID"})
	}

	updatedComment, err := c.mediaService.UpdateComment(uint(commentId), payload.ToDomainComment())
	if err != nil {
		c.logger.Error("Failed to update comment: ", err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update comment"})
	}
	return ech.JSON(http.StatusOK, models.NewCommentResponseFromDomain(updatedComment))
}

func (c CrudPostStruct) GetListReports(ech echo.Context) error {
	pageParam := ech.QueryParam("page")
	sizeParam := ech.QueryParam("size")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeParam)
	if err != nil || size < 1 {
		size = 10
	}

	reports, _, err := c.mediaService.ListReports(page, size)
	if err != nil {
		c.logger.Error("Failed to list reports: ", err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to list reports"})
	}

	response := models.NewPaginatedReportsResponse(reports, page, size)
	return ech.JSON(http.StatusOK, response)
}

func (c CrudPostStruct) GetListImages(ech echo.Context) error {
	pageParam := ech.QueryParam("page")
	sizeParam := ech.QueryParam("size")
	idParam := ech.QueryParam("id")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeParam)
	if err != nil || size < 1 {
		size = 10
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || size <= 1 {
		c.logger.Error("Validation error: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Validation failed"})
	}

	// Assuming all images are paginated
	images, _, err := c.mediaService.ListImages(uint(id), page, size)
	if err != nil {
		c.logger.Error("Failed to list images: ", err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to list images"})
	}

	response := models.NewPaginatedImagesResponse(images, page, size)
	return ech.JSON(http.StatusOK, response)
}
