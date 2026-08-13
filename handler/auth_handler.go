package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"product-data-management-platform/model"
	"product-data-management-platform/service"
)

type AuthHandler struct{ Service *service.AuthService }

func NewAuthHandler(s *service.AuthService) *AuthHandler { return &AuthHandler{Service: s} }
func (h *AuthHandler) Register(c *gin.Context) {
	var input model.RegisterRequest
	if c.ShouldBindJSON(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "a valid email and password of at least 8 characters are required"})
		return
	}
	user, err := h.Service.Register(input)
	if errors.Is(err, service.ErrEmailTaken) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "unable to register user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"id": user.ID, "email": user.Email}})
}
func (h *AuthHandler) Login(c *gin.Context) {
	var input model.LoginRequest
	if c.ShouldBindJSON(&input) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}
	token, user, err := h.Service.Login(input)
	if errors.Is(err, service.ErrInvalidCredentials) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{"error": "unable to log in"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token, "user": gin.H{"id": user.ID, "email": user.Email}})
}
