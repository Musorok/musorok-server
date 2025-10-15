package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// PromocodesHandler exposes an endpoint for promocode validation. In the
// current stub implementation it always returns that the code is invalid.
// A real implementation would look up the promocode in the database,
// check validity dates and usage limits, and compute discounts.
type PromocodesHandler struct{}

// Validate checks whether a promocode is valid and returns discount information.
// For now it always returns valid=false.
func (h *PromocodesHandler) Validate(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "valid": false,
        "discount_type": nil,
        "value": nil,
    })
}