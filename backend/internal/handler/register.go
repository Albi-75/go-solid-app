package handler

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// Register sets up all routes for the application.
func Register(r *gin.Engine, dbConn *sql.DB) {
	h := NewHandler(dbConn)

	// Healthcheck
	r.GET("/health", h.Health)

	// Todos
	r.GET("/todos", h.ListTodos)
	r.POST("/todos", h.CreateTodo)
	r.PUT("/todos/:id/done", h.MarkDone)
	r.DELETE("/todos/:id", h.DeleteTodo) // If you have DeleteTodo implemented

	// Authentication
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
}
