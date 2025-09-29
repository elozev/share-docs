package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"share-docs/pkg/logger"
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

func (s *LocalStorage) Upload(file multipart.File, object string) error {
	fileName := fmt.Sprintf("%s/%s", s.UploadPath, object)
	fmt.Printf("Will create file: %s\n", fileName)

	f, err := os.Create(fileName)
	defer f.Close()

	if err != nil {
		return err
	}

	filebytes, err := io.ReadAll(file)

	if err != nil {
		return err
	}

	bytesWritten, err := f.Write(filebytes)

	if err != nil {
		return err
	}

	if bytesWritten == 0 {
		return ErrNoBytesWritten
	}

	s.Logger.WithField("bytes_written", bytesWritten).Info("Bytes written")

	return nil
}
