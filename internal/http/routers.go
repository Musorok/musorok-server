package httpapi

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"github.com/musorok/server/internal/http/middleware"
	"github.com/musorok/server/internal/http/handlers"
	"github.com/musorok/server/internal/services"
	"github.com/musorok/server/internal/core/payments/paynetworks"
)

var upgrader = websocket.Upgrader{ CheckOrigin: func(r *http.Request) bool { return true } }

func NewRouter(db *gorm.DB, secret, refresh string, accessTTLSeconds, refreshTTLSeconds int64, pay *paynetworks.Client) *gin.Engine {
	r := gin.Default()
	r.GET("/v1/health", func(c *gin.Context){ c.JSON(http.StatusOK, gin.H{"status":"ok"}) })
	r.StaticFS("/docs", http.Dir("docs"))

	users := services.NewUserService(db)
	authH := &handlers.AuthHandler{Users: users, JWTSecret: secret, JWTRefreshSecret: refresh, AccessTTL: time.Duration(accessTTLSeconds)*time.Second, RefreshTTL: time.Duration(refreshTTLSeconds)*time.Second}

	r.POST("/v1/auth/register", authH.Register)
	r.POST("/v1/auth/login", authH.Login)
	r.POST("/v1/auth/refresh", func(c *gin.Context){ c.JSON(http.StatusOK, gin.H{}) }) // TODO

	api := r.Group("/v1", middleware.JWT(secret))
	api.GET("/me", authH.Me)

	addrH := &handlers.AddressHandler{DB: db, Addrs: services.NewAddressService(db)}
	api.GET("/addresses", addrH.List)
	api.POST("/addresses", addrH.Create)

	ordersH := &handlers.OrdersHandler{DB: db, Pay: pay}
	api.POST("/orders/quote", ordersH.Quote)
	api.POST("/orders", ordersH.Create)

	payH := &handlers.PaymentsHandler{}
	r.POST("/v1/payments/webhook", payH.Webhook)

	return r
}
