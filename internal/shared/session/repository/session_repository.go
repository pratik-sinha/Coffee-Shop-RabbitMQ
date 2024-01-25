package repository

import (
	"coffee-shop/internal/shared/session/models"
	"coffee-shop/pkg/custom_errors"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("")

const (
	basePrefix = "sessions"
)

// Session repository
type sessionRepo struct {
	redisClient *redis.Client
}

// Session repository constructor
func NewSessionRepository(redisClient *redis.Client) SessionRepository {
	return &sessionRepo{redisClient: redisClient}
}

// Create session in redis
func (s *sessionRepo) CreateSession(ctx context.Context, sess *models.Session, expire int) (*string, error) {

	ctx, span := tracer.Start(ctx, "SessionRepository.CreateSession")
	defer span.End()

	sessionKey := s.createKey(uuid.New().String())

	sessBytes, err := json.Marshal(&sess)
	if err != nil {
		return nil, custom_errors.InternalError.Wrap(span, true, err, "Error while marshaling session object")
	}
	if err = s.redisClient.Set(ctx, sessionKey, sessBytes, time.Second*time.Duration(expire)).Err(); err != nil {
		return nil, custom_errors.InternalError.Wrap(span, true, err, "Error while inserting session object in redis")
	}
	return &sessionKey, nil
}

// Get session by id
func (s *sessionRepo) GetSessionByID(ctx context.Context, sessionKey string) (*models.Session, error) {
	ctx, span := tracer.Start(ctx, "SessionRepository.GetSessionByID")
	defer span.End()

	sessBytes, err := s.redisClient.Get(ctx, sessionKey).Bytes()
	if err != nil {
		return nil, custom_errors.InternalError.Wrap(span, true, err, "Error while retrieving session object from redis")
	}

	sess := &models.Session{}
	if err = json.Unmarshal(sessBytes, &sess); err != nil {
		return nil, custom_errors.InternalError.Wrap(span, true, err, "Error while unmarshling session object")
	}
	return sess, nil
}

// Delete session by id
func (s *sessionRepo) DeleteByID(ctx context.Context, sessionID string) error {
	ctx, span := tracer.Start(ctx, "SessionRepository.DeleteByID")
	defer span.End()

	if err := s.redisClient.Del(ctx, sessionID).Err(); err != nil {
		return custom_errors.InternalError.Wrap(span, true, err, "Error while removing session object")
	}
	return nil
}

func (s *sessionRepo) createKey(sessionID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, sessionID)
}
