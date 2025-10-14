package documentapp

import (
	"share-docs/pkg/app/domain/userapp"
	"share-docs/pkg/db/models"
)

type Document struct {
	ID string `json:"id"`

	OriginalFilename string `json:"original_filename"`
	FileSize         int64  `json:"file_size"`
	MimeType         string `json:"mime_type"`

	Title       *string `json:"title"`
	Description *string `json:"description"`
	Tags        *string `json:"tags"`

	User userapp.User `json:"user"`
}

func ToAppDocument(md models.Document) Document {
	return Document{
		ID:               md.ID.String(),
		OriginalFilename: md.OriginalFilename,
		FileSize:         md.FileSize,
		MimeType:         md.MimeType,

		Title:       md.Title,
		Description: md.Description,
		Tags:        md.Tags,

		User: userapp.ToAppUser(md.User),
	}
}
