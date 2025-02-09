package rest

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"styl-monolith/internal/media/adapters/rest/models"
	"styl-monolith/internal/media/core/domain"
	"styl-monolith/internal/media/core/services"
	"styl-monolith/pkg/logger"
)

type CrudPostHandler interface {
	CreatePost(ech echo.Context) error
	CreateComment(ech echo.Context) error
	CreateReport(ech echo.Context) error
	CreateLike(ech echo.Context) error
	GetPostById(ech echo.Context, id string) error
	PostImagesById(ech echo.Context, id string) error
	DeletePost(ech echo.Context, id string) error
	DeleteComment(ech echo.Context, id string) error
	DeleteReport(ech echo.Context, id string) error
	DeleteImage(ech echo.Context, id string) error
	UpdatePost(ech echo.Context, id string) error
	UpdateComment(ech echo.Context, id string) error
	ListReports(ech echo.Context) error
	ListImages(ech echo.Context) error
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

func (c CrudPostStruct) CreatePost(ech echo.Context) error {
	var payload models.PostPayload
	form, err := ech.MultipartForm()
	if err != nil {
		c.logger.Error(logger.FailedToParseMultipartForm, err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to read form data"})
	}

	body := form.Value["body"]
	if len(body) == 0 {
		c.logger.Error(logger.MissingBodyField)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Missing 'body' field in the form data"})
	}

	// Parse the "body" value and bind it to PostPayload
	if err = json.Unmarshal([]byte(body[0]), &payload); err != nil {
		c.logger.Error(logger.FailedToBindPostPayload, err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid JSON structure in 'body' field"})
	}

	if err = payload.Validate(); err != nil {
		c.logger.Error(logger.ValidationError, err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Validation failed"})
	}
	var imageToUpload []domain.PostImage
	files := form.File["images"]

	if len(files) > 8 || len(files) < 1 {
		c.logger.Error(logger.TooManyImagesUploaded)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "You can upload up to 8 images only"})
	}
	post, err := c.mediaService.CreatePost(payload.ToDomainPost())

	if err != nil {
		c.logger.Error(logger.FailedToCreatePost, err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create post"})
	}

	for index, file := range files {
		img, err := domain.ToPostImageDomainFromMultipartForm(file, index)
		if err != nil {
			c.logger.Error(logger.FailedToParseMultipartForm, err)
			err = c.mediaService.HardDeletePost(post.Id)
			c.logger.Error(logger.FailedToDeletePost, err)
			return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to parse multipart form"})
		}
		userEmail := ech.Get("email").(string)
		postImage, err := c.mediaService.UploadImage(img, file, userEmail, post.Id)
		if err != nil {
			c.logger.Error(logger.FailedToUploadImage, err)
			err = c.mediaService.HardDeletePost(post.Id)
			if err != nil {
				c.logger.Error(logger.FailedToDeletePost, err)
			}
			return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to parse multipart form"})
		}
		imageToUpload = append(imageToUpload, postImage)
	}

	post.PostImages = imageToUpload
	post, err = c.mediaService.UpdatePost(post.Id, post)
	c.logger.Info(logger.SuccessfullyCreatedPost)
	return ech.JSON(http.StatusCreated, models.NewPostResponseFromDomain(post))
}

func (c CrudPostStruct) CreateComment(ech echo.Context) error {
	var payload models.CommentPayload
	if err := ech.Bind(&payload); err != nil {
		c.logger.Error(logger.FailedToBindCommentPayload, err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Bad request payload"})
	}

	if err := payload.Validate(); err != nil {
		c.logger.Error(logger.ValidationError, err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Validation failed"})
	}

	comment, err := c.mediaService.CreateComment(payload.ToDomainComment())
	if err != nil {
		c.logger.Error(logger.FailedToCreateComment, err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create comment"})
	}

	c.logger.Info(logger.SuccessfullyCreatedComment)
	return ech.JSON(http.StatusCreated, models.NewCommentResponseFromDomain(comment))
}

func (c CrudPostStruct) CreateReport(ech echo.Context) error {
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

func (c CrudPostStruct) CreateLike(ech echo.Context) error {
	var payload models.LikePayload
	if err := ech.Bind(&payload); err != nil {
		c.logger.Error("Failed to bind LikePayload: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Bad request payload"})
	}

	if err := payload.Validate(); err != nil {
		c.logger.Error("Validation error: ", err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Validation failed"})
	}

	like, err := c.mediaService.CreateLike(payload.ToDomainLike())
	if err != nil {
		c.logger.Error("Failed to create like: ", err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create like"})
	}
	return ech.JSON(http.StatusCreated, models.NewResponseLikeFromDomain(like))
}

func (c CrudPostStruct) GetPostById(ech echo.Context, id string) error {
	postId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.logger.Error(logger.InvalidPostID, err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid post ID"})
	}

	post, err := c.mediaService.GetPost(uint(postId))
	if err != nil {
		c.logger.Error(logger.FailedToRetrievePost, err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get post"})
	}

	c.logger.Info(logger.SuccessfullyRetrievedPost)
	return ech.JSON(http.StatusOK, models.NewPostResponseFromDomain(post))
}

func (c CrudPostStruct) DeletePost(ech echo.Context, id string) error {
	postId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.logger.Error(logger.InvalidPostID, err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid post ID"})
	}

	if err := c.mediaService.DeletePost(uint(postId)); err != nil {
		c.logger.Error(logger.FailedToDeletePost, err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete post"})
	}

	c.logger.Info(logger.SuccessfullyDeletedPost)
	return ech.JSON(http.StatusNoContent, nil)
}

func (c CrudPostStruct) PostImagesById(ech echo.Context, id string) error {
	postId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.logger.Error(logger.InvalidPostID, err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid post ID"})
	}

	images, _, err := c.mediaService.ListImages(uint(postId), 1, 100)
	if err != nil {
		c.logger.Error(logger.FailedToRetrieveImages, err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve images"})
	}

	c.logger.Info(logger.SuccessfullyRetrievedImages)
	response := models.NewPaginatedImagesResponse(images, 1, 100)
	return ech.JSON(http.StatusOK, response)
}

func (c CrudPostStruct) DeleteComment(ech echo.Context, id string) error {
	commentId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.logger.Error(logger.InvalidCommentID, err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid comment ID"})
	}

	if err := c.mediaService.DeleteComment(uint(commentId)); err != nil {
		c.logger.Error(logger.FailedToDeleteComment, err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete comment"})
	}

	c.logger.Info(logger.SuccessfullyDeletedComment)
	return ech.JSON(http.StatusNoContent, nil)
}

func (c CrudPostStruct) DeleteReport(ech echo.Context, id string) error {
	reportId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.logger.Error(logger.InvalidReportID, err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid report ID"})
	}

	if err := c.mediaService.DeleteReport(uint(reportId)); err != nil {
		c.logger.Error(logger.FailedToDeleteReport, err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete report"})
	}

	c.logger.Info(logger.SuccessfullyDeletedReport)
	return ech.JSON(http.StatusNoContent, nil)
}

func (c CrudPostStruct) DeleteImage(ech echo.Context, id string) error {
	imageKey := id // Assuming imageKey is passed as the ID
	if err := c.mediaService.DeleteImage(imageKey); err != nil {
		c.logger.Error(logger.FailedToDeleteImage, err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete image"})
	}

	c.logger.Info(logger.SuccessfullyDeletedImage)
	return ech.JSON(http.StatusNoContent, nil)
}

func (c CrudPostStruct) UpdatePost(ech echo.Context, id string) error {
	var payload models.PostPayload
	if err := ech.Bind(&payload); err != nil {
		c.logger.Error(logger.FailedToBindPostPayload, err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Bad request payload"})
	}

	if err := payload.Validate(); err != nil {
		c.logger.Error(logger.ValidationError, err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Validation failed"})
	}

	postId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.logger.Error(logger.InvalidPostID, err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid post ID"})
	}

	updatedPost, err := c.mediaService.UpdatePost(uint(postId), payload.ToDomainPost())
	if err != nil {
		c.logger.Error(logger.FailedToUpdatePost, err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update post"})
	}

	c.logger.Info(logger.SuccessfullyUpdatedPost)
	return ech.JSON(http.StatusOK, models.NewPostResponseFromDomain(updatedPost))
}

func (c CrudPostStruct) UpdateComment(ech echo.Context, id string) error {
	var payload models.CommentPayload
	if err := ech.Bind(&payload); err != nil {
		c.logger.Error(logger.FailedToBindCommentPayload, err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Bad request payload"})
	}

	if err := payload.Validate(); err != nil {
		c.logger.Error(logger.ValidationError, err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Validation failed"})
	}

	commentId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.logger.Error(logger.InvalidCommentID, err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid comment ID"})
	}

	updatedComment, err := c.mediaService.UpdateComment(uint(commentId), payload.ToDomainComment())
	if err != nil {
		c.logger.Error(logger.FailedToUpdateComment, err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update comment"})
	}

	c.logger.Info(logger.SuccessfullyUpdatedComment)
	return ech.JSON(http.StatusOK, models.NewCommentResponseFromDomain(updatedComment))
}

func (c CrudPostStruct) ListReports(ech echo.Context) error {
	pageParam := ech.QueryParam("page")
	sizeParam := ech.QueryParam("size")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		c.logger.Warn(logger.InvalidPageParam, err)
		page = 1
	}

	size, err := strconv.Atoi(sizeParam)
	if err != nil || size < 1 {
		c.logger.Warn(logger.InvalidSizeParam, err)
		size = 10
	}

	reports, _, err := c.mediaService.ListReports(page, size)
	if err != nil {
		c.logger.Error(logger.FailedToListReports, err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to list reports"})
	}

	c.logger.Info(logger.SuccessfullyListedReports)
	response := models.NewPaginatedReportsResponse(reports, page, size)
	return ech.JSON(http.StatusOK, response)
}

func (c CrudPostStruct) ListImages(ech echo.Context) error {
	pageParam := ech.QueryParam("page")
	sizeParam := ech.QueryParam("size")
	idParam := ech.QueryParam("id")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		c.logger.Warn(logger.InvalidPageParam, err)
		page = 1
	}

	size, err := strconv.Atoi(sizeParam)
	if err != nil || size < 1 {
		c.logger.Warn(logger.InvalidSizeParam, err)
		size = 10
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.logger.Error(logger.ValidationError, err)
		return ech.JSON(http.StatusBadRequest, map[string]string{"message": "Validation failed"})
	}

	images, _, err := c.mediaService.ListImages(uint(id), page, size)
	if err != nil {
		c.logger.Error(logger.FailedToListImages, err)
		return ech.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to list images"})
	}

	c.logger.Info(logger.SuccessfullyListedImages)
	response := models.NewPaginatedImagesResponse(images, page, size)
	return ech.JSON(http.StatusOK, response)
}
