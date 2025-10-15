package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// AdminHandler provides stub implementations for administrative operations.
// Administrative endpoints allow management of polygons, couriers and
// retrieval of system metrics. These handlers currently return
// Not Implemented responses.
type AdminHandler struct{}

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