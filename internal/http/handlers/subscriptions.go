package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "gorm.io/gorm"

    "github.com/musorok/server/internal/domain"
    "github.com/musorok/server/internal/core/payments/paynetworks"
)

// SubscriptionsHandler provides stub implementations for subscription‑related endpoints.
// In a full implementation this handler would coordinate with services to manage
// subscription plans, create and cancel subscriptions, and create orders from
// subscriptions. The current implementation returns 501 Not Implemented
// responses to indicate that the functionality is not yet available.
// SubscriptionsHandler manages subscription plans and orders derived from
// subscriptions. It holds a DB and payment client for persistence and
// payment intents.
type SubscriptionsHandler struct{
    DB *gorm.DB
    Pay *paynetworks.Client
}

// ListPlans returns available subscription plans. In a real implementation this
// would fetch plan data from configuration or database. Here it returns an empty
// array to satisfy client expectations.
func (h *SubscriptionsHandler) ListPlans(c *gin.Context) {
    // return predefined plans with prices and total bags
    plans := []gin.H{
        {"plan": domain.PlanP7, "price": 1569, "total_bags": 7},
        {"plan": domain.PlanP15, "price": 3175, "total_bags": 15},
        {"plan": domain.PlanP30, "price": 5976, "total_bags": 30},
    }
    c.JSON(http.StatusOK, plans)
}

// Create handles creation of a new subscription. The request would normally
// include a plan identifier and optional promocode. A payment intent would be
// created and the subscription stored. Currently returns 501.
func (h *SubscriptionsHandler) Create(c *gin.Context) {
    uid := c.GetString("uid")
    var req struct{
        Plan domain.SubscriptionPlan `json:"plan"`
        Promocode string `json:"promocode"`
    }
    if err := c.BindJSON(&req); err != nil || req.Plan == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "plan required"})
        return
    }
    // determine price and total bags
    var price int
    var total int
    switch req.Plan {
    case domain.PlanP7:
        price, total = 1569, 7
    case domain.PlanP15:
        price, total = 3175, 15
    case domain.PlanP30:
        price, total = 5976, 30
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan"})
        return
    }
    // TODO: apply promocode discount if provided via h.DB
    // create subscription with status ACTIVE but pending payment
    userID, _ := uuid.Parse(uid)
    sub := domain.Subscription{
        UserID: userID,
        Plan: req.Plan,
        TotalBags: total,
        RemainingBags: total,
        PriceKZT: price,
        Status: domain.SubActive,
        StartedAt: time.Now(),
    }
    if err := h.DB.Create(&sub).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    // create payment intent
    meta := map[string]string{"subscription_id": sub.ID.String()}
    intent, _ := h.Pay.CreatePaymentIntent(c, price, meta)
    c.JSON(http.StatusCreated, gin.H{"subscription": sub, "payment": gin.H{"id": intent.ID, "paymentUrl": intent.PaymentURL}})
}

// Current returns the current active subscription for the authenticated user.
func (h *SubscriptionsHandler) Current(c *gin.Context) {
    uid := c.GetString("uid")
    userID, _ := uuid.Parse(uid)
    var sub domain.Subscription
    if err := h.DB.Where("user_id = ? AND status = ?", userID, domain.SubActive).Order("started_at desc").First(&sub).Error; err != nil {
        c.JSON(http.StatusOK, gin.H{"subscription": nil})
        return
    }
    c.JSON(http.StatusOK, gin.H{"subscription": sub})
}

// Cancel cancels an existing subscription identified by path parameter id.
func (h *SubscriptionsHandler) Cancel(c *gin.Context) {
    subID := c.Param("id")
    if subID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
        return
    }
    var sub domain.Subscription
    if err := h.DB.First(&sub, "id = ?", subID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
        return
    }
    sub.Status = domain.SubCanceled
    if err := h.DB.Save(&sub).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"canceled": true})
}

// CreateOrderFromSubscription creates an order drawing from the remaining bags
// of an active subscription. The request would include bags_count, address_id,
// time_option and optional scheduled_at/comment. Currently returns 501.
func (h *SubscriptionsHandler) CreateOrderFromSubscription(c *gin.Context) {
    uid := c.GetString("uid")
    userID, _ := uuid.Parse(uid)
    // parse request
    var req struct{
        BagsCount int `json:"bags_count"`
        AddressID string `json:"address_id"`
        TimeOption domain.TimeOption `json:"time_option"`
        ScheduledAt *time.Time `json:"scheduled_at"`
        Comment string `json:"comment"`
    }
    if err := c.BindJSON(&req); err != nil || req.BagsCount <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }
    // find active subscription
    var sub domain.Subscription
    if err := h.DB.Where("user_id = ? AND status = ?", userID, domain.SubActive).Order("started_at desc").First(&sub).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "no active subscription"})
        return
    }
    if sub.RemainingBags < req.BagsCount {
        c.JSON(http.StatusBadRequest, gin.H{"error": "not enough remaining bags"})
        return
    }
    // find address
    var addr domain.Address
    if err := h.DB.First(&addr, "id = ?", req.AddressID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "address not found"})
        return
    }
    if addr.PolygonID == nil {
        c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "этот район пока не обслуживается"})
        return
    }
    // create order linked to subscription, price zero (paid via subscription)
    order := domain.Order{
        UserID: userID,
        AddressID: addr.ID,
        PolygonID: *addr.PolygonID,
        Type: domain.OrderSubscription,
        BagsCount: req.BagsCount,
        PriceKZT: 0,
        Comment: req.Comment,
        TimeOption: req.TimeOption,
        ScheduledAt: req.ScheduledAt,
        Status: domain.StatusNew,
    }
    if err := h.DB.Create(&order).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    // update subscription remaining bags
    sub.RemainingBags -= req.BagsCount
    if err := h.DB.Save(&sub).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"order": order})
}