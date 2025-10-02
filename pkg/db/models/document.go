package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `gorm:"type:uuid,primaryKey;default;gen_random_uuid()"`

	// File information
	OriginalFilename string `json:"original_filename"`
	FilePath         string `json:"file_path"`
	FileSize         int64  `json:"file_size"`
	MimeType         string `json:"mime_type"`
	FileHash         string `json:"file_hash"`

	// Metadata
	Title       *string `gorm:"size:255" json:"title"`
	Description *string `gorm:"size:1000" json:"description"`
	Tags        string  `gorm:"size:500" json:"tags"`

	// Relationships
	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"-"`
	User   User      `gorm:"foreignKey:UserID" json:"user"`
}

func (d *Document) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}
