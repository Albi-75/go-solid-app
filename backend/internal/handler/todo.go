package handler

import (
	"database/sql"
	"go-solid-app/backend/internal/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	q      *db.Queries
	dbConn *sql.DB
}

func NewHandler(dbConn *sql.DB) *Handler {
	return &Handler{
		q:      db.New(dbConn),
		dbConn: dbConn,
	}
}

type CreateTodoRequest struct {
	Description string `json:"description"`
}

type TodoResponse struct {
	ID          int32  `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
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
// @Success 200 {array} TodoResponse
// @Failure 500 {object} map[string]string
// @Router /todos [get]
func (h *Handler) ListTodos(c *gin.Context) {
	todos, err := h.q.ListTodos(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]TodoResponse, 0, len(todos))
	for _, t := range todos {
		resp = append(resp, TodoResponse{
			ID:          t.ID,
			Description: t.Description,
			Completed:   t.Completed,
		})
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Add a new todo
// @Tags Todos
// @Accept json
// @Produce json
// @Param todo body handler.CreateTodoRequest true "Todo to add"
// @Success 201 {object} TodoResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos [post]
func (h *Handler) CreateTodo(c *gin.Context) {
	var req CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid description"})
		return
	}
	todo, err := h.q.CreateTodo(c, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := TodoResponse{
		ID:          todo.ID,
		Description: todo.Description,
		Completed:   todo.Completed,
	}
	c.JSON(http.StatusCreated, resp)
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

// @Summary Delete a todo
// @Tags Todos
// @Param id path int true "Todo ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /todos/{id} [delete]
func (h *Handler) DeleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err = h.q.DeleteTodo(c, int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
