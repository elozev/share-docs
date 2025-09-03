package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	connectionString := "host=localhost user=postgres password=postgres dbname=share_docs port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database!")

	return db
}
