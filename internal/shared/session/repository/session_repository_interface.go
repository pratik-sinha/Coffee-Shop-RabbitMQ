//go:generate mockgen -source redis_repository.go -destination mock/redis_repository.go -package mock
package repository

import (
	"coffee-shop/internal/shared/session/models"
	"context"
)

// Session repository
type SessionRepository interface {
	CreateSession(ctx context.Context, session *models.Session, expire int) (*string, error)
	GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error)
	DeleteByID(ctx context.Context, sessionID string) error
}
