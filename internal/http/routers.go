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
    // redirect root to the swagger documentation.  This makes it easy to open the API docs without
    // needing to remember the /docs path.
    r.GET("/", func(c *gin.Context) {
        c.Redirect(http.StatusFound, "/docs")
    })
    // health check
    r.GET("/v1/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
    // serve swagger documentation at /docs
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

    // subscriptions and promocodes routes
    subH := &handlers.SubscriptionsHandler{}
    promoH := &handlers.PromocodesHandler{}
    api.GET("/subscriptions/plans", subH.ListPlans)
    api.POST("/subscriptions", subH.Create)
    api.GET("/subscriptions/current", subH.Current)
    api.POST("/subscriptions/:id/cancel", subH.Cancel)
    api.POST("/subscription-orders", subH.CreateOrderFromSubscription)
    api.POST("/promocodes/validate", promoH.Validate)

    // courier routes (login and protected actions)
    courierH := &handlers.CourierHandler{}
    // unauthenticated login
    r.POST("/v1/courier/auth/login", courierH.Login)
    courierGroup := r.Group("/v1/courier", middleware.JWT(secret))
    courierGroup.GET("/me", courierH.Me)
    courierGroup.GET("/orders", courierH.ListOrders)
    courierGroup.POST("/orders/:id/accept", courierH.AcceptOrder)
    courierGroup.POST("/orders/:id/status", courierH.UpdateOrderStatus)

    // admin routes, protected by admin role (checked in handlers or middleware)
    adminH := &handlers.AdminHandler{}
    adminGroup := r.Group("/v1/admin", middleware.JWT(secret))
    adminGroup.GET("/polygons", adminH.ListPolygons)
    adminGroup.POST("/polygons", adminH.CreatePolygon)
    adminGroup.PUT("/polygons/:id", adminH.UpdatePolygon)
    adminGroup.POST("/couriers", adminH.CreateCourier)
    adminGroup.PUT("/couriers/:id", adminH.UpdateCourier)
    adminGroup.GET("/metrics", adminH.Metrics)

	return r
}
