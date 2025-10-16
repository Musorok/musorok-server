package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "gorm.io/gorm"

    "github.com/musorok/server/internal/domain"
)

// CourierHandler defines stub endpoints for courier operations. In a full
// implementation, couriers could authenticate, view orders in their polygon,
// accept orders and update statuses. Currently, these return 501 Not
// Implemented responses.
// CourierHandler implements courier-related endpoints such as viewing the
// courier profile, listing orders, accepting orders, updating statuses and
// checking balances/withdrawing earnings. It requires a gorm.DB to
// perform database operations.
type CourierHandler struct{
    DB *gorm.DB
}

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
    // identify courier
    courierID := c.GetString("uid")
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "order id required"})
        return
    }
    // fetch order
    var order domain.Order
    if err := h.DB.First(&order, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
        return
    }
    // ensure order is in a state that can be accepted
    if order.Status != domain.StatusNew && order.Status != domain.StatusPaid {
        c.JSON(http.StatusBadRequest, gin.H{"error": "order cannot be accepted in its current status"})
        return
    }
    // check that order has no courier yet
    if order.CourierID != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "order already assigned"})
        return
    }
    // ensure courier belongs to same polygon as the order
    // fetch courier from couriers table
    var courier struct{
        ID uuid.UUID
        PolygonID uuid.UUID
        IsActive bool
    }
    if err := h.DB.Raw("SELECT id, polygon_id, is_active FROM couriers WHERE user_id = ?", courierID).Scan(&courier).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "courier not found"})
        return
    }
    if !courier.IsActive {
        c.JSON(http.StatusBadRequest, gin.H{"error": "courier inactive"})
        return
    }
    if courier.PolygonID != order.PolygonID {
        c.JSON(http.StatusBadRequest, gin.H{"error": "order outside of courier polygon"})
        return
    }
    // assign order
    cid := uuid.MustParse(courierID)
    order.CourierID = &cid
    order.Status = domain.StatusAssigned
    if err := h.DB.Save(&order).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, order)
}

// UpdateOrderStatus allows the courier to update the status of an order. The
// request would specify the next status (e.g. PICKING_UP, DONE) and optional
// metadata. Currently not implemented.
func (h *CourierHandler) UpdateOrderStatus(c *gin.Context) {
    courierID := c.GetString("uid")
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "order id required"})
        return
    }
    var req struct{
        ToStatus string `json:"to_status"`
        Meta map[string]interface{} `json:"meta"`
    }
    if err := c.BindJSON(&req); err != nil || req.ToStatus == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "to_status required"})
        return
    }
    // fetch order
    var order domain.Order
    if err := h.DB.First(&order, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
        return
    }
    // verify courier is assigned to this order
    if order.CourierID == nil || order.CourierID.String() != courierID {
        c.JSON(http.StatusForbidden, gin.H{"error": "not assigned to courier"})
        return
    }
    // allowed transitions
    var newStatus domain.OrderStatus
    switch req.ToStatus {
    case string(domain.StatusPickingUp):
        // only from ASSIGNED
        if order.Status != domain.StatusAssigned {
            c.JSON(http.StatusBadRequest, gin.H{"error": "cannot pick up order in current status"})
            return
        }
        newStatus = domain.StatusPickingUp
    case string(domain.StatusDone):
        // only from PICKING_UP
        if order.Status != domain.StatusPickingUp {
            c.JSON(http.StatusBadRequest, gin.H{"error": "cannot complete order in current status"})
            return
        }
        newStatus = domain.StatusDone
    case string(domain.StatusCanceled):
        // allow cancel from ASSIGNED or PICKING_UP
        if order.Status != domain.StatusAssigned && order.Status != domain.StatusPickingUp {
            c.JSON(http.StatusBadRequest, gin.H{"error": "cannot cancel order in current status"})
            return
        }
        newStatus = domain.StatusCanceled
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status transition"})
        return
    }
    // update status
    order.Status = newStatus
    if err := h.DB.Save(&order).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    // if completed, create settlement and update courier balance
    if newStatus == domain.StatusDone {
        // compute settlement amount: per bag fixed rate (example 200 KZT per bag)
        amount := order.BagsCount * 200
        // create settlement record if not exists
        settlement := domain.OrderSettlement{
            OrderID: order.ID,
            CourierID: *order.CourierID,
            BagsCount: order.BagsCount,
            AmountKZT: amount,
        }
        if err := h.DB.Create(&settlement).Error; err == nil {
            // update courier balance
            var bal domain.CourierBalance
            if err := h.DB.First(&bal, "courier_id = ?", order.CourierID).Error; err == nil {
                bal.TotalEarnedKZT += int64(amount)
                bal.UpdatedAt = time.Now()
                h.DB.Save(&bal)
            }
        }
    }
    c.JSON(http.StatusOK, order)
}

// Balance returns the courier's current available balance and aggregated
// earnings/withdrawals. If the courier does not yet have a balance row it
// will be created with zeros. The returned JSON includes totalEarned,
// totalWithdrawn and available (earned minus withdrawn).
func (h *CourierHandler) Balance(c *gin.Context) {
    courierID := c.GetString("uid")
    // ensure courier exists
    var bal domain.CourierBalance
    // create row if missing
    tx := h.DB.First(&bal, "courier_id = ?", courierID)
    if tx.Error != nil {
        // no existing balance, create with zero values
        cid := uuid.MustParse(courierID)
        bal = domain.CourierBalance{CourierID: cid, TotalEarnedKZT: 0, TotalWithdrawnKZT: 0}
        if err := h.DB.Create(&bal).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    }
    available := bal.TotalEarnedKZT - bal.TotalWithdrawnKZT
    c.JSON(http.StatusOK, gin.H{
        "totalEarned": bal.TotalEarnedKZT,
        "totalWithdrawn": bal.TotalWithdrawnKZT,
        "available": available,
    })
}

// Withdraw allows a courier to request withdrawal of a portion of their
// available earnings. The request body should contain {"amount": number}.
// A payout request row is created with status REQUESTED. The balance is not
// immediately decreased; funds are deducted when the payout is processed.
func (h *CourierHandler) Withdraw(c *gin.Context) {
    courierID := c.GetString("uid")
    var req struct{ Amount int `json:"amount"` }
    if err := c.BindJSON(&req); err != nil || req.Amount <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be > 0"})
        return
    }
    // fetch balance
    var bal domain.CourierBalance
    if err := h.DB.First(&bal, "courier_id = ?", courierID).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "balance not found"})
        return
    }
    available := bal.TotalEarnedKZT - bal.TotalWithdrawnKZT
    if int64(req.Amount) > available {
        c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance"})
        return
    }
    // create payout request
    pr := domain.PayoutRequest{
        CourierID: uuid.MustParse(courierID),
        AmountKZT: req.Amount,
        Status: domain.PayoutRequested,
    }
    if err := h.DB.Create(&pr).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"id": pr.ID, "status": pr.Status})
}