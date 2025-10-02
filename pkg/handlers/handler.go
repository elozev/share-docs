package handlers

import (
	"fmt"
	"net/http"
	"share-docs/pkg/logger"
	"share-docs/pkg/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseHandlerInterface interface {
	Success(c *gin.Context, data interface{}, message string)
	SuccessWithMeta(c *gin.Context, data interface{}, message string, meta *Meta)
	Created(c *gin.Context, data interface{}, message string)
	BadRequest(c *gin.Context, message string)
	Unauthorized(c *gin.Context, message string)
	Forbidden(c *gin.Context, message string)
	NotFound(c *gin.Context, message string)
	InternalError(c *gin.Context, message string)
	ValidationError(c *gin.Context, errors map[string]string)
	GetUserIDFromContext(c *gin.Context) (uuid.UUID, error)
	GetUUIDParam(c *gin.Context, param string) (uuid.UUID, error)
	GetPaginationParams(c *gin.Context) (page, limit int)
	BindAndValidate(c *gin.Context, obj interface{}) error
	GetLogger(c *gin.Context) *logger.Logger
}

type BaseHandler struct {
	db     *gorm.DB
	logger *logger.Logger
}

type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page       int   `json:"page,omitempty"`
	Limit      int   `json:"limit,omitempty"`
	Total      int64 `json:"total,omitempty"`
	TotalPages int   `json:"total_pages,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code,omitemtpy"`
}

func NewBaseHandler(db *gorm.DB, log *logger.Logger) *BaseHandler {
	return &BaseHandler{
		db:     db,
		logger: log,
	}
}

func (h *BaseHandler) Success(c *gin.Context, data interface{}, message string) {
	log := h.GetLogger(c)
	log.WithFields(map[string]interface{}{
		"status":  http.StatusOK,
		"message": message,
	}).Info("Successful response")

	c.JSON(http.StatusOK, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func (h *BaseHandler) SuccessWithMeta(c *gin.Context, data interface{}, message string, meta *Meta) {
	log := h.GetLogger(c)

	log.WithFields(map[string]interface{}{
		"status":  http.StatusOK,
		"message": message,
	}).Info("Succesful response")

	c.JSON(http.StatusOK, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func (h *BaseHandler) Created(c *gin.Context, data interface{}, message string) {
	log := h.GetLogger(c)
	log.WithFields(map[string]interface{}{
		"status":  http.StatusCreated,
		"message": message,
	}).Info("Successfully created record!")

	c.JSON(http.StatusCreated, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func (h *BaseHandler) failedRequest(c *gin.Context, message string, code int) {
	log := h.GetLogger(c)

	if code < 400 || code >= 600 {
		code = http.StatusBadRequest
	}

	log.WithFields(map[string]interface{}{
		"status":  code,
		"message": message,
	}).Error("Request failed")

	c.JSON(code, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
	})
}

func (h *BaseHandler) BadRequest(c *gin.Context, message string) {
	h.failedRequest(c, message, http.StatusBadRequest)
}

func (h *BaseHandler) Unauthorized(c *gin.Context, message string) {
	h.failedRequest(c, message, http.StatusUnauthorized)
}

func (h *BaseHandler) Forbidden(c *gin.Context, message string) {
	h.failedRequest(c, message, http.StatusForbidden)
}

func (h *BaseHandler) NotFound(c *gin.Context, message string) {
	h.failedRequest(c, message, http.StatusNotFound)
}

func (h *BaseHandler) InternalError(c *gin.Context, message string) {
	h.failedRequest(c, message, http.StatusInternalServerError)
}

func (h *BaseHandler) ValidationError(c *gin.Context, errors map[string]string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"error":   "Validation failed",
		"details": errors,
	})
}

func (h *BaseHandler) GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userID := c.GetString("UserID")

	// if !exists {
	// 	return uuid.Nil, fmt.Errorf("user ID not found in contexts")
	// }

	fmt.Println("GetUserIDFromContext: userID", userID)

	id, err := uuid.Parse(userID)

	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID format")
	}

	return id, nil
}

func (h *BaseHandler) GetUUIDParam(c *gin.Context, param string) (uuid.UUID, error) {
	paramStr := c.Param(param)
	if paramStr == "" {
		return uuid.Nil, fmt.Errorf("%s parameter is required", param)
	}

	id, err := uuid.Parse(paramStr)

	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID format for %s", param)
	}

	return id, nil
}

func (h *BaseHandler) GetPaginationParams(c *gin.Context) (page, limit int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 1000 {
		limit = 10
	}
	return page, limit
}

func (h *BaseHandler) BindAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return err
	}

	return nil
}

func (h *BaseHandler) BindFormAndValidate(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBind(obj); err != nil {
		return err
	}

	return nil
}

func (h *BaseHandler) GetLogger(c *gin.Context) *logger.Logger {
	return util.GetLoggerFromContext(c)
}
