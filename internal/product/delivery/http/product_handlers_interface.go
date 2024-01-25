package http

import "github.com/gin-gonic/gin"

type ProductController interface {
	GetProducts(c *gin.Context)
	GetProductsByType(c *gin.Context)
}
