package handlers

import (
    "net/http"
    "time"
    
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// PromocodesHandler exposes an endpoint for promocode validation. In the
// current stub implementation it always returns that the code is invalid.
// A real implementation would look up the promocode in the database,
// check validity dates and usage limits, and compute discounts.
// PromocodesHandler provides endpoints for validating promocodes and generating
// promo codes via admin. It holds a DB reference for queries.
type PromocodesHandler struct{
    DB *gorm.DB
}

// Validate checks whether a promocode is valid and returns discount information.
// For now it always returns valid=false.
func (h *PromocodesHandler) Validate(c *gin.Context) {
    var req struct { Code string `json:"code"` }
    // allow code from query param as fallback
    code := c.Query("code")
    if err := c.BindJSON(&req); err == nil && req.Code != "" {
        code = req.Code
    }
    if code == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "code required"})
        return
    }
    // query promocode
    var promo struct{
        ID string
        DiscountType string
        Value int
    }
    now := time.Now()
    row := h.DB.Raw(
        `SELECT id, discount_type, value FROM promocodes WHERE code = ? AND is_active = true AND active_from <= ? AND active_to >= ? AND (usage_limit = 0 OR used_count < usage_limit)`,
        code, now, now,
    ).Row()
    if row.Err() != nil {
        c.JSON(http.StatusOK, gin.H{"valid": false})
        return
    }
    var id, dtype string
    var val int
    if err := row.Scan(&id, &dtype, &val); err != nil {
        c.JSON(http.StatusOK, gin.H{"valid": false})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "valid": true,
        "discount_type": dtype,
        "value": val,
    })
}