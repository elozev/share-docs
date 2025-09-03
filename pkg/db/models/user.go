package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string
	Name     string
	Birthday *time.Time
}
