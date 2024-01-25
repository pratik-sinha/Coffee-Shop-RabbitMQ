package service

import (
	"coffee-shop/internal/shared/session/models"
	"coffee-shop/internal/shared/session/repository"
	"context"

	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("")

// Session use case
type sessionUC struct {
	sessionRepo repository.SessionRepository
}

// New session use case constructor
func NewSessionService(sessionRepo repository.SessionRepository) SessionService {
	return &sessionUC{sessionRepo: sessionRepo}
}

// Create new session
func (u *sessionUC) CreateSession(ctx context.Context, session *models.Session, expire int) (*string, error) {
	ctx, span := tracer.Start(ctx, "UserRepository.CreateUser")
	defer span.End()

	return u.sessionRepo.CreateSession(ctx, session, expire)
}

// Delete session by id
func (u *sessionUC) DeleteByID(ctx context.Context, sessionID string) error {
	ctx, span := tracer.Start(ctx, "UserRepository.CreateUser")
	defer span.End()
	return u.sessionRepo.DeleteByID(ctx, sessionID)
}

// get session by id
func (u *sessionUC) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	ctx, span := tracer.Start(ctx, "UserRepository.CreateUser")
	defer span.End()
	return u.sessionRepo.GetSessionByID(ctx, sessionID)
}
