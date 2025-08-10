package handler

import (
	"net/http"

	"github.com/LuizFernando991/golang-auth-microservice/internal/model"
	"github.com/LuizFernando991/golang-auth-microservice/internal/service"
	"github.com/LuizFernando991/golang-auth-microservice/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	svc service.AuthService
	v   *validator.Validate
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc, v: validator.New()}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var dto model.CreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.v.Struct(dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := h.svc.Register(c.Request.Context(), dto)
	if err != nil {
		if err == util.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": u})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var dto model.LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.v.Struct(dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	access, refresh, err := h.svc.Login(c.Request.Context(), dto)
	if err != nil {
		if err == util.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": refresh})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken := c.GetHeader("Refresh_Token")

	if refreshToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing refresh_token header"})
		return
	}
	access, refresh, err := h.svc.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		if err == util.ErrRefreshTokenNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": refresh})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken := c.GetHeader("Refresh_Token")

	if refreshToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing refresh_token header"})
		return
	}

	if err := h.svc.Logout(c.Request.Context(), refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

func (h *AuthHandler) Me(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var userID int64
	switch v := uid.(type) {
	case int64:
		userID = v
	case float64:
		userID = int64(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id type"})
		return
	}
	u, err := h.svc.GetUserById(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": u})
}
