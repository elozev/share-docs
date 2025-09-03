package main

import (
	"context"
	"fmt"
	database "share-docs/pkg/db"
	"share-docs/pkg/db/models"

	"gorm.io/gorm"
)

func main() {
	// r := routes.SetupRouter()

	// s := &http.Server{
	// 	Addr: ":8080",
	// 	Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	// 		if req.Method == "HEAD" {
	// 			req.Method = "GET"
	// 		}
	// 		r.ServeHTTP(w, req)
	// 	}),
	// 	ReadTimeout:    10 * time.Second,
	// 	WriteTimeout:   10 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }

	// s.ListenAndServe()

	db := database.Connect()

	user := models.User{
		Email:    "test@example.com",
		Name:     "Test User",
		Birthday: nil,
	}

	ctx := context.Background()
	err := gorm.G[models.User](db).Create(ctx, &user)

	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully created user!")
}
