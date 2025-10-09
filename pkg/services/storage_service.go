package services

import (
	"fmt"
	"mime/multipart"
	"os"
	"share-docs/pkg/logger"
	"share-docs/pkg/storage"
	"share-docs/pkg/util"
)

type StorageService struct {
	// TODO: expand
	sb storage.StorageBackendInterface
}

type StorageServiceInterface interface {
	UploadDocument(file multipart.File, object string) error
}

func NewStorageService(storageType string, logger *logger.Logger) *StorageService {
	supportedBackends := map[string]bool{
		"local": true,
		"gcp":   false,
	}

	var sb storage.StorageBackendInterface

	if !supportedBackends[storageType] {
		panic(fmt.Sprintf("storage backend %s not supported", storageType))
	} else if storageType == "local" {
		path := util.MustGetEnv("STORAGE_LOCAL_PATH")

		_, err := os.Stat(path)

		if os.IsNotExist(err) {
			panic(err)
		}

		sb = storage.NewLocalStorage(path, *logger)
	}

	return &StorageService{
		sb: sb,
	}
}

func (s *StorageService) UploadDocument(file multipart.File, path string, filename string) (*storage.StorageObject, error) {
	so, err := s.sb.Upload(file, path, filename)

	if err != nil {
		return nil, err
	}

	return so, nil
}
