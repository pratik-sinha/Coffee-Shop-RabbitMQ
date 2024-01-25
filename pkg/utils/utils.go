package utils

import (
	"coffee-shop/config"
	"context"

	"github.com/gin-gonic/gin"
)

// Get config path for local or docker
func GetConfigPath(env string) string {
	if env == "docker" {
		return "./config/config.docker"
	}
	return "./config/config.local"
}

func CreateSessionCookie(c *gin.Context, cfg *config.Config, token string) {
	c.SetCookie(cfg.Cookie.Name, token, cfg.Cookie.MaxAge, "/", cfg.Cookie.Domain, cfg.Cookie.Secure, cfg.Cookie.HTTPOnly)
}

// Delete session
func DeleteSessionCookie(c *gin.Context, cfg *config.Config) {
	c.SetCookie(cfg.Cookie.Name, "", -1, "/", cfg.Cookie.Domain, false, true)
}

type ReqIDCtxKey struct{}

func GetRequestCtx(c *gin.Context) context.Context {
	return context.WithValue(c.Request.Context(), ReqIDCtxKey{}, GetRequestID(c))
}

func GetRequestID(c *gin.Context) string {
	return c.Request.Header.Get("X-Request-ID")
}
