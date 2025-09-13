package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseHandler struct {
	db *gorm.DB
	// TODO: add logger
	//
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

func (h *BaseHandler) Success(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func (h *BaseHandler) SuccessWithMeta(c *gin.Context, data interface{}, message string, meta *Meta) {
	c.JSON(http.StatusOK, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func (h *BaseHandler) Created(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func (h *BaseHandler) failedRequest(c *gin.Context, message string, code int) {
	if code < 400 || code >= 600 {
		code = http.StatusBadRequest
	}
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
	userID, exists := c.Get("userID")

	if !exists {
		return uuid.Nil, fmt.Errorf("user ID not found in contexts")
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
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
