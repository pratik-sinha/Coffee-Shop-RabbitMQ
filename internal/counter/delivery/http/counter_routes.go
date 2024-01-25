package http

import (
	"coffee-shop/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func CounterRoutes(route *gin.RouterGroup, mw *middleware.MiddlewareManager, counterController CounterController) {
	route.POST("/placeorder", mw.AuthenticateSession(), counterController.PlaceOrder)
	route.POST("/getOrders", mw.AuthenticateSession(), counterController.GetOrders)
}
