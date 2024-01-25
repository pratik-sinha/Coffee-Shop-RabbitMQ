package http

import (
	"coffee-shop/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(route *gin.RouterGroup, mw *middleware.MiddlewareManager, productController ProductController) {
	route.POST("/getproducts", productController.GetProducts)
	route.POST("/getproductsbytype", productController.GetProductsByType)
}
