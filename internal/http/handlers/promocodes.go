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
    var req struct {
        Code string `json:"code" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "code required"})
        return
    }

    type Promo struct {
        ID           string    `json:"id"`
        Code         string    `json:"code"`
        DiscountType string    `json:"discount_type"`
        Value        int       `json:"value"`
        ActiveFrom   time.Time `json:"active_from"`
        ActiveTo     time.Time `json:"active_to"`
    }

    var p Promo
    row := h.DB.
        Raw(`
            SELECT id, code, discount_type, value, active_from, active_to
            FROM promocodes
            WHERE code = ?
              AND is_active = true
              AND now() BETWEEN active_from AND active_to
              AND (usage_limit IS NULL OR used_count < usage_limit)
        `, req.Code).Row()

    if err := row.Scan(&p.ID, &p.Code, &p.DiscountType, &p.Value, &p.ActiveFrom, &p.ActiveTo); err != nil {
        // не найдено / невалидно
        c.JSON(http.StatusOK, gin.H{"valid": false})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "valid":      true,
        "promocode":  p,
    })
}
