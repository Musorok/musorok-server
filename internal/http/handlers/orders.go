package handlers

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/musorok/server/internal/domain"
	"github.com/musorok/server/internal/core/payments/paynetworks"
	"gorm.io/gorm"
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
