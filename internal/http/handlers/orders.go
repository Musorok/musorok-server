package handlers

import (
    "net/http"
    "time"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "gorm.io/gorm"
    
    "github.com/musorok/server/internal/domain"
    "github.com/musorok/server/internal/core/payments/paynetworks"
)

type OrdersHandler struct{
	DB *gorm.DB
	Pay *paynetworks.Client
}

func (h *OrdersHandler) Quote(c *gin.Context) {
	var req struct{ AddressID string `json:"address_id"`; BagsCount int `json:"bags_count"` }
	if err := c.BindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"bad json"}); return }
	if req.BagsCount <= 0 { c.JSON(http.StatusBadRequest, gin.H{"error":"bags_count must be > 0"}); return }
	var addr domain.Address
	if err := h.DB.First(&addr, "id = ?", req.AddressID).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error":"address not found"}); return }
	canServe := addr.PolygonID != nil
	price := 249 * req.BagsCount
	c.JSON(http.StatusOK, gin.H{"price_kzt": price, "can_serve": canServe, "polygon_name": addr.PolygonName})
}

func (h *OrdersHandler) Create(c *gin.Context) {
	uid := c.GetString("uid")
	var req struct{
		AddressID string `json:"address_id"`
		BagsCount int `json:"bags_count"`
		TimeOption domain.TimeOption `json:"time_option"`
		ScheduledAt *time.Time `json:"scheduled_at"`
		Comment string `json:"comment"`
		Promocode string `json:"promocode"`
	}
	if err := c.BindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"bad json"}); return }
	if req.BagsCount <= 0 { c.JSON(http.StatusBadRequest, gin.H{"error":"bags_count must be > 0"}); return }
	var addr domain.Address
	if err := h.DB.First(&addr, "id = ?", req.AddressID).Error; err != nil { c.JSON(http.StatusNotFound, gin.H{"error":"address not found"}); return }
	if addr.PolygonID == nil { c.JSON(http.StatusUnprocessableEntity, gin.H{"error":"этот район пока не обслуживается"}); return }

	price := 249 * req.BagsCount
	uuidv, _ := uuid.Parse(uid)
	order := domain.Order{
		UserID: uuidv, AddressID: addr.ID, PolygonID: *addr.PolygonID,
		Type: domain.OrderOneTime, BagsCount: req.BagsCount, PriceKZT: price,
		Comment: req.Comment, TimeOption: req.TimeOption, ScheduledAt: req.ScheduledAt,
		Status: domain.StatusNew,
	}
	if err := h.DB.Create(&order).Error; err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }

	intent, _ := h.Pay.CreatePaymentIntent(c, price, map[string]string{"order_id": order.ID.String()})
	c.JSON(http.StatusCreated, gin.H{"order": order, "payment": gin.H{"id": intent.ID, "paymentUrl": intent.PaymentURL}})
}

// History returns a paginated list of the authenticated user's past orders.
// Query parameters: page (default 1), limit (default 20), sort (created_at desc|asc).
func (h *OrdersHandler) History(c *gin.Context) {
    uid := c.GetString("uid")
    var uID uuid.UUID
    var err error
    if uID, err = uuid.Parse(uid); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
        return
    }
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    if page < 1 { page = 1 }
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    if limit <= 0 { limit = 20 }
    sort := c.DefaultQuery("sort", "desc")
    orderStr := "created_at desc"
    if sort == "asc" { orderStr = "created_at asc" }
    offset := (page - 1) * limit
    var orders []domain.Order
    var total int64
    h.DB.Model(&domain.Order{}).Where("user_id = ?", uID).Count(&total)
    h.DB.Where("user_id = ?", uID).Order(orderStr).Offset(offset).Limit(limit).Find(&orders)
    c.JSON(http.StatusOK, gin.H{
        "orders": orders,
        "total_count": total,
        "page": page,
        "limit": limit,
    })
}
