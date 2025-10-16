package handlers

import (
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "gorm.io/gorm"
    
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
    DB *gorm.DB
    // in-memory storage for OTP codes by phone. In production this would
    // integrate with WhatsApp API and a persistent cache (Redis). For the
    // purpose of this demo the codes are stored in memory and reset on
    // process restart.
    otpStore map[string]string
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

// SendOTP sends a one-time password to the user's phone via WhatsApp. In
// this demo implementation it simply records a static code in memory.
func (h *AuthHandler) SendOTP(c *gin.Context) {
    var req struct { Phone string `json:"phone"` }
    if err := c.BindJSON(&req); err != nil || req.Phone == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "phone required"})
        return
    }
    if h.otpStore == nil { h.otpStore = make(map[string]string) }
    // generate or assign a static OTP. For testing we use 0000
    code := "0000"
    h.otpStore[req.Phone] = code
    // Here you would integrate with 360Dialog/Twilio to send the code via WhatsApp
    c.JSON(http.StatusOK, gin.H{"sent": true})
}

// VerifyOTP verifies the provided code for the given phone number. If the
// code matches, it logs in or creates the user and returns JWT tokens.
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
    var req struct { Phone string `json:"phone"`; Code string `json:"code"` }
    if err := c.BindJSON(&req); err != nil || req.Phone == "" || req.Code == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "phone and code required"})
        return
    }
    if h.otpStore == nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "otp not sent"})
        return
    }
    exp, ok := h.otpStore[req.Phone]
    if !ok || exp != req.Code {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid code"})
        return
    }
    // OTP is valid; find or create user
    var u domain.User
    if err := h.DB.Where("phone = ?", req.Phone).First(&u).Error; err != nil {
        // create new user with empty name and default role
        newUser, err := h.Users.Create(c, req.Phone, "", "", "", domain.RoleUser)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        u = *newUser
    }
    // generate tokens
    acc, _ := auth.NewAccessToken(h.JWTSecret, u.ID.String(), string(u.Role), h.AccessTTL)
    ref, _ := auth.NewAccessToken(h.JWTRefreshSecret, u.ID.String(), string(u.Role), h.RefreshTTL)
    // delete OTP from store
    delete(h.otpStore, req.Phone)
    c.JSON(http.StatusOK, gin.H{"access": acc, "refresh": ref})
}

// DeleteAccount removes the authenticated user's account and associated data
// (addresses, sessions, device tokens and subscriptions) from the database.
func (h *AuthHandler) DeleteAccount(c *gin.Context) {
    uidStr := c.GetString("uid")
    id, err := uuid.Parse(uidStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }
    // delete associated rows using raw table names for entities without Go structs
    // addresses
    h.DB.Where("user_id = ?", id).Delete(&domain.Address{})
    // sessions (table: sessions)
    h.DB.Exec("DELETE FROM sessions WHERE user_id = ?", id)
    // device tokens (table: device_tokens)
    h.DB.Exec("DELETE FROM device_tokens WHERE user_id = ?", id)
    // subscriptions
    h.DB.Where("user_id = ?", id).Delete(&domain.Subscription{})
    // orders could remain for history; optionally mark is_deleted
    // finally delete user
    h.DB.Where("id = ?", id).Delete(&domain.User{})
    c.JSON(http.StatusOK, gin.H{"deleted": true})
}
