package repository

import (
	"coffee-shop/internal/user/models"
	"coffee-shop/pkg/custom_errors"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("")

type userRepository struct {
	conn *mongo.Database
}

func NewUserRepository(m *mongo.Database) UserRepository {
	return &userRepository{conn: m}
}

func (u *userRepository) CreateUser(ctx context.Context, user models.User) error {
	ctx, span := tracer.Start(ctx, "UserRepository.CreateUser")
	defer span.End()
	_, err := u.conn.Collection("user_profiles").InsertOne(ctx, user)
	if err != nil {
		err := custom_errors.InternalError.Wrapf(span, true, err, "Error while creating new user: %#v", user)
		return err
	}
	return nil
}

func (u *userRepository) GetUser(ctx context.Context, filter bson.M) (*models.User, error) {
	ctx, span := tracer.Start(ctx, "UserRepository.GetUser")
	defer span.End()

	var res models.User

	err := u.conn.Collection("user_profiles").FindOne(ctx, filter).Decode(&res)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, custom_errors.InternalError.Wrapf(span, true, err, "Error while retrieving user for Filter: %#v", filter)
	}

	return &res, nil
}
