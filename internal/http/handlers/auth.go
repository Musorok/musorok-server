package handlers

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/musorok/server/internal/core/auth"
	"github.com/musorok/server/internal/domain"
	"github.com/musorok/server/internal/services"
)

type AuthHandler struct{
	Users *services.UserService
	JWTSecret string
	JWTRefreshSecret string
	AccessTTL time.Duration
	RefreshTTL time.Duration
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct{ Phone, Email, Name, Password string }
	if err := c.BindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"bad json"}); return }
	u, err := h.Users.Create(c, req.Phone, req.Email, req.Name, req.Password, domain.RoleUser)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusCreated, gin.H{"id": u.ID})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct{ Login, Password string }
	if err := c.BindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"bad json"}); return }
	u, err := h.Users.Authenticate(c, req.Login, req.Password)
	if err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"}); return }
	acc, _ := auth.NewAccessToken(h.JWTSecret, u.ID.String(), string(u.Role), h.AccessTTL)
	ref, _ := auth.NewAccessToken(h.JWTRefreshSecret, u.ID.String(), string(u.Role), h.RefreshTTL)
	c.JSON(http.StatusOK, gin.H{"access": acc, "refresh": ref})
}

func (h *AuthHandler) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{"id": c.GetString("uid")},
		"hasActiveSubscription": false,
		"remainingBags": 0,
		"hasSavedAddresses": false,
	})
}
