package http

import "github.com/gin-gonic/gin"

type CounterController interface {
	PlaceOrder(c *gin.Context)
	GetOrders(c *gin.Context)
}
