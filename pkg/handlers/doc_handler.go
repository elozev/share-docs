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
}

// as a user, I should be able to upload a document
// on request
// 1. grab the document from the request
// 2. create a document reference with upload_success field
// 3. use the storage service to store the document under /docs/{user_id}/{document_id}
// 4. on success, update document reference
// 4. return document reference

type CreateDocumentRequest struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
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

	log.WithField("file_size", req.File.Size).Info("Uploaded file size")

	f, err := req.File.Open()
	defer f.Close()

	if err != nil {
		log.WithError(err).Error("Failed opening file")
		h.BadRequest(c, "Failed opening file!")
		return
	}

	so, err := h.storageService.UploadDocument(f, req.File.Filename)

	if err != nil {
		log.WithError(err).Error("Failed uploading document")
		h.InternalError(c, fmt.Sprintf("Failed uploading document!"))
		return
	}

	h.Created(c, so, "Successfully created a document!")
}
