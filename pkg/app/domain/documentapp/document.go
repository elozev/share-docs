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
	IsPublic    bool    `json:"is_public"`

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
		IsPublic:    md.IsPublic,

		User: userapp.ToAppUser(md.User),
	}
}

type UpdateDocument struct {
	Title       *string `json:"title" validate:"omitempty"`
	Description *string `json:"description" validate:"omitempty"`
	Tags        *string `json:"tags" validate:"omitempty"`
	IsPublic    *bool   `json:"is_public" validate:"omitempty"`
}

func (ud *UpdateDocument) HasAtLeastOneField() bool {
	return ud.Title != nil || ud.Description != nil || ud.Tags != nil || ud.IsPublic != nil
}

func (ud *UpdateDocument) ToModelDocument() models.Document {
	var d models.Document
	if ud.Title != nil {
		d.Title = ud.Title
	}

	if ud.Description != nil {
		d.Description = ud.Description
	}

	if ud.Tags != nil {
		d.Tags = ud.Tags
	}

	if ud.IsPublic != nil {
		d.IsPublic = *ud.IsPublic
	}

	return d
}
