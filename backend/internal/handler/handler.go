package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"go-solid-app/backend/internal/db"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	q *db.Queries
	dbConn *sql.DB
}

func NewHandler(dbConn *sql.DB) *Handler {
	return &Handler{
		q: db.New(dbConn),
		dbConn: dbConn,
	}
}

// @Summary Health check
// @Tags Health
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /health [get]
func (h *Handler) Health(c *gin.Context) {
	if err := h.dbConn.Ping(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// @Summary List all todos
// @Tags Todos
// @Produce json
// @Success 200 {array} db.Todo
// @Failure 500 {object} map[string]string
// @Router /todos [get]
func (h *Handler) ListTodos(c *gin.Context) {
	todos, err := h.q.ListTodos(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, todos)
}

// @Summary Add a new todo
// @Tags Todos
// @Accept json
// @Produce json
// @Param todo body struct{Description string} true "Todo to add"
// @Success 201 {object} db.Todo
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos [post]
func (h *Handler) CreateTodo(c *gin.Context) {
	var req struct {
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid description"})
		return
	}
	todo, err := h.q.CreateTodo(c, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, todo)
}

// @Summary Mark todo as done
// @Tags Todos
// @Param id path int true "Todo ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos/{id}/done [put]
func (h *Handler) MarkDone(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err = h.q.MarkTodoDone(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func Register(r *gin.Engine, dbConn *sql.DB) {
	h := NewHandler(dbConn)
	r.GET("/health", h.Health)
	r.GET("/todos", h.ListTodos)
	r.POST("/todos", h.CreateTodo)
	r.PUT("/todos/:id/done", h.MarkDone)
}
