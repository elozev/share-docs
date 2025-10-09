package storage

import (
	"errors"
	"fmt"
	"mime/multipart"
	"share-docs/pkg/logger"
	"strings"
	"time"
)

type StorageBackend struct {
	UploadPath string
	Logger     logger.Logger
}

type StorageObject struct {
	Name          string
	Path          string
	MimeType      string
	FileSizeBytes int64
	FileHash      string
}

var (
	ErrNoBytesWritten = errors.New("Nothing was written to destination")
)

type StorageBackendInterface interface {
	Upload(file multipart.File, path string, filename string) (*StorageObject, error)
	// Get(object string) (*StorageObject, error)
	// Delete(object string) error
}

func (s *StorageBackend) normaliseFilename(name string) string {
	lc := strings.ToLower(name)

	splitFile := strings.Split(lc, ".")

	timestamp := time.Now().Unix()
	noWhiteSpace := strings.ReplaceAll(fmt.Sprintf("%s-%d", splitFile[0], timestamp), " ", "-")

	res := strings.Join([]string{
		noWhiteSpace,
		splitFile[1],
	}, ".")

	return res
}
