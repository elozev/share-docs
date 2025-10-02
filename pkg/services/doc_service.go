package services

import (
	"fmt"
	"share-docs/pkg/db/models"
	"share-docs/pkg/storage"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DocumentServiceInterface interface {
	// TODO: add parameters
	CreateDocument(userID uuid.UUID, o storage.StorageObject) (*models.Document, error)

	GetDocument() (*models.Document, error)
}

type DocumentService struct {
	db *gorm.DB
}

func NewDocumentService(db *gorm.DB) *DocumentService {
	return &DocumentService{
		db: db,
	}
}

func (s *DocumentService) CreateDocument(userID uuid.UUID, o storage.StorageObject) (*models.Document, error) {
	document := &models.Document{

		OriginalFilename: o.Name,
		FilePath:         o.Path,
		// TODO: get from storageobject
		FileSize: 0,
		MimeType: o.MimeType,
		// TODO: generate and store in storageobject
		FileHash: "",

		UserID: userID,
	}

	if result := s.db.Create(document); result.Error != nil {
		// TODO: use logger
		return nil, fmt.Errorf("failed to create a document")
	}

	return document, nil
}
