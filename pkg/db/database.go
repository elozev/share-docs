package db

import (
	"fmt"
	"share-docs/pkg/util"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/joho/godotenv/autoload"
)

var (
	host     = util.GetEnv("DB_HOST", "localhost")
	port     = util.GetEnv("DB_PORT", "5432")
	user     = util.GetEnv("DB_USER", "postgres")
	password = util.GetEnv("DB_PASSWORD", "postgres")
	dbname   = util.GetEnv("DB_NAME", "share_docs")
)

func Connect() *gorm.DB {
	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database!")

	return db
}
