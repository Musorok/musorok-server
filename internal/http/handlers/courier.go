package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// CourierHandler defines stub endpoints for courier operations. In a full
// implementation, couriers could authenticate, view orders in their polygon,
// accept orders and update statuses. Currently, these return 501 Not
// Implemented responses.
type CourierHandler struct{}

// Login authenticates a courier using phone/email and password. A JWT is
// returned on success. Currently not implemented.
func (h *CourierHandler) Login(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, gin.H{"error": "courier login not implemented"})
}

// Me returns information about the logged in courier, such as assigned polygon
// and daily statistics.
func (h *CourierHandler) Me(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, gin.H{"error": "courier profile not implemented"})
}

// ListOrders lists orders for the courier, optionally filtered by status
// query parameter. Currently returns empty list.
func (h *CourierHandler) ListOrders(c *gin.Context) {
    c.JSON(http.StatusOK, []gin.H{})
}

// AcceptOrder allows the courier to accept an order by id. In a real
// implementation the order would be reserved for the courier.
func (h *CourierHandler) AcceptOrder(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, gin.H{"error": "accept order not implemented"})
}

// UpdateOrderStatus allows the courier to update the status of an order. The
// request would specify the next status (e.g. PICKING_UP, DONE) and optional
// metadata. Currently not implemented.
func (h *CourierHandler) UpdateOrderStatus(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, gin.H{"error": "update order status not implemented"})
}