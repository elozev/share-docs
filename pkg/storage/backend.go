package storage

import (
	"errors"
	"mime/multipart"
	"share-docs/pkg/logger"
)

type StorageBackend struct {
	UploadPath string
	Logger     logger.Logger
}

type StorageObject struct {
	Name     string
	Path     string
	MimeType string
}

var (
	ErrNoBytesWritten = errors.New("Nothing was written to destination")
)

type StorageBackendInterface interface {
	Upload(file multipart.File, object string) error
	// Get(object string) (*StorageObject, error)
	// Delete(object string) error
}
