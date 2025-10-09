package services

import (
	"errors"
	"fmt"
	"share-docs/pkg/db/models"
	"share-docs/pkg/storage"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TODO: do not json serialise documents on the model level, but rather here

type DocumentServiceInterface interface {
	// TODO: add parameters
	CreateDocument(userID uuid.UUID, o storage.StorageObject) (*models.Document, error)
	GetDocument(documentID uuid.UUID) (*models.Document, error)
}

var (
	ErrDocumentNotFound = errors.New("document not found")
)

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
		FileSize:         o.FileSizeBytes,
		MimeType:         o.MimeType,
		FileHash:         o.FileHash,

		UserID: userID,
	}

	if result := s.db.Preload("User").Create(document); result.Error != nil {
		// TODO: use logger
		return nil, fmt.Errorf("failed to create a document")
	}

	return document, nil
}

func (s *DocumentService) GetDocument(documentStringID string) (*models.Document, error) {
	documentID, err := uuid.Parse(documentStringID)

	if err != nil {
		return nil, ErrInvalidId
	}

	var document *models.Document

	result := s.db.Preload("User").First(&document, documentID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrDocumentNotFound
		}

		return nil, err
	}

	return document, nil
}
