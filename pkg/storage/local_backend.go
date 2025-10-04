package storage

import (
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"share-docs/pkg/logger"

	"github.com/gabriel-vasile/mimetype"
)

type LocalStorage struct {
	StorageBackend
}

func NewLocalStorage(uploadPath string, logger logger.Logger) *LocalStorage {
	return &LocalStorage{
		StorageBackend: StorageBackend{
			UploadPath: uploadPath,
			Logger:     logger,
		},
	}
}

func (s *LocalStorage) Upload(file multipart.File, object string) (*StorageObject, error) {
	fileName := fmt.Sprintf("%s/%s", s.UploadPath, s.normaliseFilename(object))
	fmt.Printf("Will create file: %s\n", fileName)

	f, err := os.Create(fileName)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	mimeType, err := mimetype.DetectReader(f)

	if err != nil {
		return nil, err
	}

	filebytes, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	bytesWritten, err := f.Write(filebytes)

	if err != nil {
		return nil, err
	}

	if bytesWritten == 0 {
		return nil, ErrNoBytesWritten
	}

	hash := md5.Sum(filebytes)

	so := &StorageObject{
		Name:          fileName,
		Path:          object,
		MimeType:      mimeType.String(),
		FileSizeBytes: int64(len(filebytes)),
		FileHash:      fmt.Sprintf("%x", hash),
	}

	s.Logger.WithField("bytes_written", bytesWritten).Info("Bytes written")

	return so, nil
}
