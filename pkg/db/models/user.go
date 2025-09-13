package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `gorm:"type:uuid,primaryKey;default;gen_random_uuid()" json:"id"`

	Email      string     `gorm:"unique" json:"email"`
	Password   string     `gorm:"not null, size:255" json:"-"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	BirthDate  *time.Time `json:"birth_date"`
	IsActive   bool       `json:"is_active"`
	IsVerified bool       `json:"is_verified"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	return nil
}
