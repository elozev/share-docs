package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `gorm:"type:uuid,primaryKey;default;gen_random_uuid()"`

	// File information
	OriginalFilename string
	FilePath         string
	FileSize         int64
	MimeType         string
	FileHash         string

	// Metadata
	Title       *string `gorm:"size:255"`
	Description *string `gorm:"size:1000"`
	Tags        *string `gorm:"size:500"`
	IsPublic    bool    `gorm:"type:bool"`

	// Relationships
	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	User   User      `gorm:"foreignKey:UserID"`
}

func (d *Document) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

func (d *Document) AfterCreate(tx *gorm.DB) error {
	// TODO: send message to pubsub, clear cache, etc
	return nil
}
