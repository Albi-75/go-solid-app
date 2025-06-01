package main

import (
	"database/sql"
	"go-solid-app/backend/internal/handler"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-solid-app/backend/cmd/app/docs"

	swaggerFiles "github.com/swaggo/files"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	r := gin.Default()
	handler.Register(r, db)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
