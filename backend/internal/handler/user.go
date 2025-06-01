package handler

import (
	"net/http"

	"go-solid-app/backend/internal/db"
	"go-solid-app/backend/internal/util"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body handler.RegisterRequest true "User to register"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	hash, err := util.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}
	user, err := h.q.CreateUser(c, db.CreateUserParams{
		Username:     req.Username,
		PasswordHash: hash,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "username": user.Username})
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// @Summary Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body handler.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	user, err := h.q.GetUserByUsername(c, req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if !util.CheckPassword(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := util.GenerateJWT(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
