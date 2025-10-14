package handlers

import (
	"fmt"
	"mime/multipart"
	"share-docs/pkg/services"

	"github.com/gin-gonic/gin"
)

type DocHandler struct {
	BaseHandler
	documentService services.DocumentService
	storageService  services.StorageService
}

func NewDocHandler(ds services.DocumentService, ss services.StorageService, bs BaseHandler) *DocHandler {
	return &DocHandler{
		BaseHandler:     bs,
		documentService: ds,
		storageService:  ss,
	}
}

type DocHandlerInterface interface {
	CreateDocument(c *gin.Context)
	GetDocument(c *gin.Context)
	GetFile(c *gin.Context)
}

// as a user, I should be able to upload a document
// on request
// 1. grab the document from the request
// 2. create a document reference with upload_success field
// 3. use the storage service to store the document under /docs/{user_id}/{document_id}
// 4. on success, update document reference
// 4. return document reference

type CreateDocumentRequest struct {
	File     *multipart.FileHeader `form:"file" binding:"required"`
	IsPublic bool                  `form:"is_public"`
}

func (h *DocHandler) CreateDocument(c *gin.Context) {
	log := h.GetLogger(c)
	var req CreateDocumentRequest
	if err := h.BindFormAndValidate(c, &req); err != nil {
		h.BadRequest(c, fmt.Sprintf("Invalid request! %s", err.Error()))
		return
	}

	if req.File.Size <= 0 {
		h.BadRequest(c, fmt.Sprintf("Empty file! Size: %d", req.File.Size))
		return
	}

	userID, err := h.GetUserIDFromContext(c)

	if err != nil {
		log.WithError(err).Error("Failed getting UserID")
		h.Unauthorized(c, fmt.Sprintf("Failed getting UserID!"))
		return
	}

	f, err := req.File.Open()
	defer f.Close()

	if err != nil {
		log.WithError(err).Error("Failed opening file")
		h.BadRequest(c, "Failed opening file!")
		return
	}

	filepath := fmt.Sprintf("%s/", userID)

	so, err := h.storageService.UploadDocument(f, filepath, req.File.Filename)

	if err != nil {
		log.WithError(err).Error("Failed uploading document")
		h.InternalError(c, fmt.Sprintf("Failed uploading document!"))
		return
	}

	(*so).IsPublic = req.IsPublic

	log.WithField("storage_object", so).Info("Storage object debug")
	doc, err := h.documentService.CreateDocument(userID, *so)

	if err != nil {
		log.WithError(err).Error("Failed creating document reference")
		h.InternalError(c, fmt.Sprintf("Failed creating document reference"))
		return
	}

	h.Created(c, doc, "Successfully created a document!")
}

func (h *DocHandler) GetDocument(c *gin.Context) {
	documentId := c.Param("id")

	document, err := h.documentService.GetDocument(documentId)

	if err != nil {
		h.handlerRetrieveDocumentError(c, err)
		return
	}

	h.Success(c, document, "document found")
}

func (h *DocHandler) GetFile(c *gin.Context) {
	documentId := c.Param("id")

	document, err := h.documentService.GetDocument(documentId)

	if err != nil {
		h.handlerRetrieveDocumentError(c, err)
		return
	}

	if !document.IsPublic {
		h.Unauthorized(c, "document is not public")
		return
	}

	c.File(document.OriginalFilename)
}

func (h *DocHandler) handlerRetrieveDocumentError(c *gin.Context, err error) {
	log := h.GetLogger(c)
	log.WithError(err).Error("Failed to retrieve document")

	switch err {
	case services.ErrDocumentNotFound:
		h.NotFound(c, "Document not found")
		return
	case services.ErrInvalidId:
		h.BadRequest(c, "Invalid document ID")
		return
	default:
		h.InternalError(c, "Internal server error")
		return
	}
}
