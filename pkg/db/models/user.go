package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid,primary_key;default;gen_random_uuid()"`

	Email    string
	Name     string
	Birthday *time.Time
}
