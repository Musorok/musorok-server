package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/musorok/server/internal/core/geospatial"
	"github.com/musorok/server/internal/domain"
	"github.com/musorok/server/internal/services"
	"gorm.io/gorm"
)

type AddressHandler struct{
	DB *gorm.DB
	Addrs *services.AddressService
}

func (h *AddressHandler) List(c *gin.Context) {
	var items []domain.Address
	if err := h.Addrs.List(c, c.GetString("uid"), &items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusOK, items)
}

func (h *AddressHandler) Create(c *gin.Context) {
	uid := c.GetString("uid")
	uuidv, _ := uuid.Parse(uid)
	var req domain.Address
	if err := c.BindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"bad json"}); return }
	req.UserID = uuidv
	var polys []domain.Polygon
	if err := h.DB.Where("is_active = true").Find(&polys).Error; err == nil {
		for _, p := range polys {
			if geospatial.ContainsPoint(p.GeoJSON, req.Lng, req.Lat) {
				req.PolygonID = &p.ID; req.PolygonName = &p.Name
				break
			}
		}
	}
	if req.PolygonID == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error":"этот район пока не обслуживается"}); return
	}
	if err := h.Addrs.Create(c, &req, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return
	}
	c.JSON(http.StatusCreated, req)
}
