package services

import (
	"share-docs/pkg/db/models"

	"gorm.io/gorm"
)

type DocumentServiceInterface interface {
	// TODO: add parameters
	CreateDocument() (*models.Document, error)

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

func (s *DocumentService) CreateDocument() (*models.Document, error) {
	var d models.Document
	return &d, nil
}
