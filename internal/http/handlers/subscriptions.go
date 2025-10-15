package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// SubscriptionsHandler provides stub implementations for subscription‑related endpoints.
// In a full implementation this handler would coordinate with services to manage
// subscription plans, create and cancel subscriptions, and create orders from
// subscriptions. The current implementation returns 501 Not Implemented
// responses to indicate that the functionality is not yet available.
type SubscriptionsHandler struct{}

// ListPlans returns available subscription plans. In a real implementation this
// would fetch plan data from configuration or database. Here it returns an empty
// array to satisfy client expectations.
func (h *SubscriptionsHandler) ListPlans(c *gin.Context) {
    c.JSON(http.StatusOK, []gin.H{})
}

// Create handles creation of a new subscription. The request would normally
// include a plan identifier and optional promocode. A payment intent would be
// created and the subscription stored. Currently returns 501.
func (h *SubscriptionsHandler) Create(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, gin.H{"error": "subscriptions creation not implemented"})
}

// Current returns the current active subscription for the authenticated user.
func (h *SubscriptionsHandler) Current(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, gin.H{"error": "fetch current subscription not implemented"})
}

// Cancel cancels an existing subscription identified by path parameter id.
func (h *SubscriptionsHandler) Cancel(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, gin.H{"error": "cancel subscription not implemented"})
}

// CreateOrderFromSubscription creates an order drawing from the remaining bags
// of an active subscription. The request would include bags_count, address_id,
// time_option and optional scheduled_at/comment. Currently returns 501.
func (h *SubscriptionsHandler) CreateOrderFromSubscription(c *gin.Context) {
    c.JSON(http.StatusNotImplemented, gin.H{"error": "create order from subscription not implemented"})
}