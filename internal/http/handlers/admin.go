package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// AdminHandler provides stub implementations for administrative operations.
// Administrative endpoints allow management of polygons, couriers and
// retrieval of system metrics. These handlers currently return
// Not Implemented responses.
// AdminHandler implements administrative operations such as managing
// polygons, couriers, promocodes and metrics. It holds a DB reference for
// performing queries and writes. For unimplemented methods it returns
// HTTP 501.
type AdminHandler struct{
    DB *gorm.DB
}

// ListPolygons returns a list of polygons configured in the system. In a
// real implementation this would query the database and return geojson and
// status flags.
func (h *AdminHandler) ListPolygons(c *gin.Context) {
    c.JSON(http.StatusOK, []gin.H{})
}

// CreatePolygon creates a new polygon with the supplied name, city and
// geojson. Currently returns 501.
func (h *AdminHandler) CreatePolygon(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, gin.H{"error": "create polygon not implemented"})
}

// UpdatePolygon updates an existing polygon identified by id. Currently
// returns 501.
func (h *AdminHandler) UpdatePolygon(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, gin.H{"error": "update polygon not implemented"})
}

// CreateCourier creates a new courier user or assigns an existing user as a
// courier. Currently returns 501.
func (h *AdminHandler) CreateCourier(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, gin.H{"error": "create courier not implemented"})
}

// UpdateCourier updates a courier's status, polygon assignment or work hours.
func (h *AdminHandler) UpdateCourier(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, gin.H{"error": "update courier not implemented"})
}

// Metrics returns summary statistics such as orders per polygon and SLA.
func (h *AdminHandler) Metrics(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, gin.H{"error": "metrics not implemented"})
}

// CreatePromocode allows an admin to create a new promocode. The request
// should include code, discount_type (FIXED or PERCENT), value, active_from,
// active_to and optional usage_limit (0 for unlimited). The handler
// inserts a new row into the promocodes table.
func (h *AdminHandler) CreatePromocode(c *gin.Context) {
    var req struct{
        Code string `json:"code"`
        DiscountType string `json:"discount_type"`
        Value int `json:"value"`
        ActiveFrom time.Time `json:"active_from"`
        ActiveTo time.Time `json:"active_to"`
        UsageLimit int `json:"usage_limit"`
    }
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }
    if req.Code == "" || req.Value <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "code and positive value required"})
        return
    }
    // default usage limit to zero (unlimited) if not provided
    promo := struct{
        Code string
        DiscountType string
        Value int
        ActiveFrom time.Time
        ActiveTo time.Time
        UsageLimit int
        IsActive bool
    }{
        Code: req.Code,
        DiscountType: req.DiscountType,
        Value: req.Value,
        ActiveFrom: req.ActiveFrom,
        ActiveTo: req.ActiveTo,
        UsageLimit: req.UsageLimit,
        IsActive: true,
    }
    if promo.UsageLimit < 0 { promo.UsageLimit = 0 }
    if err := h.DB.Exec(
        `INSERT INTO promocodes (code, discount_type, value, active_from, active_to, usage_limit, used_count, is_active) VALUES (?, ?, ?, ?, ?, ?, 0, true)`,
        promo.Code, promo.DiscountType, promo.Value, promo.ActiveFrom, promo.ActiveTo, promo.UsageLimit,
    ).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"code": promo.Code})
}