package middleware

import (
	"coffee-shop/config"
	session_service "coffee-shop/internal/shared/session/service"
	"coffee-shop/pkg/logger"
	"coffee-shop/pkg/token"
)

type MiddlewareManager struct {
	sessUC session_service.SessionService
	t      token.Maker
	cfg    *config.Config
	logger logger.Logger
}

func NewMiddlewareManager(sessUC session_service.SessionService, t token.Maker, cfg *config.Config, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{sessUC: sessUC, cfg: cfg, logger: logger, t: t}
}
