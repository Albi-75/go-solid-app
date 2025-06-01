package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"go-solid-app/backend/internal/db"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine, dbConn *sql.DB) {
	q := db.New(dbConn)

	r.GET("/health", func(c *gin.Context) {
		if err := dbConn.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "db error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// List all todos
	r.GET("/todos", func(c *gin.Context) {
		todos, err := q.ListTodos(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, todos)
	})

	// Create a todo
	r.POST("/todos", func(c *gin.Context) {
		var req struct {
			Description string `json:"description"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.Description == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid description"})
			return
		}
		todo, err := q.CreateTodo(c, req.Description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, todo)
	})

	// Mark as done
	r.PUT("/todos/:id/done", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		err = q.MarkTodoDone(c, int32(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})
}
