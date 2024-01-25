//go:generate mockgen -source usecase.go -destination mock/usecase.go -package mock
package service

import (
	"coffee-shop/internal/shared/session/models"
	"context"
)

// Session Service
type SessionService interface {
	CreateSession(ctx context.Context, session *models.Session, expire int) (*string, error)
	GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error)
	DeleteByID(ctx context.Context, sessionID string) error
}
