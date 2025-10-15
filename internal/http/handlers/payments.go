package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type PaymentsHandler struct{}

func (h *PaymentsHandler) Webhook(c *gin.Context) {
	c.Status(http.StatusOK) // TODO: verify signature + update payments/orders
}
