package http

import "github.com/gin-gonic/gin"

type UserController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	IsCookieValid(c *gin.Context)
}
