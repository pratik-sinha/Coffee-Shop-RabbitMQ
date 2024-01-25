package http

import (
	"coffee-shop/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(route *gin.RouterGroup, mw *middleware.MiddlewareManager, userController UserController) {
	route.POST("/register", userController.Register)
	route.POST("/login", userController.Login)
	route.POST("/logout", mw.AuthenticateSession(), userController.Logout)
	route.POST("/isCookieValid", mw.AuthenticateSession(), userController.IsCookieValid)

}
