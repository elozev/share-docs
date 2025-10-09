package userapp

import (
	"share-docs/pkg/db/models"
	"time"
)

type User struct {
	ID        string     `json:"id"`
	Email     string     `json:"email"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	BirthDate *time.Time `json:"birth_date"`
}

func ToAppUser(mu models.User) User {
	return User{
		ID:        mu.ID.String(),
		Email:     mu.Email,
		FirstName: mu.FirstName,
		LastName:  mu.LastName,
		BirthDate: mu.BirthDate,
	}
}
