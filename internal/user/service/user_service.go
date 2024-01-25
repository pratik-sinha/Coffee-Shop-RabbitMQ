package service

import (
	session_models "coffee-shop/internal/shared/session/models"
	session_service "coffee-shop/internal/shared/session/service"
	"time"

	"coffee-shop/internal/user/models"
	"coffee-shop/internal/user/repository"

	errors "coffee-shop/pkg/custom_errors"
	"coffee-shop/pkg/db/mongo"
	"coffee-shop/pkg/helpers"
	"coffee-shop/pkg/token"
	"coffee-shop/pkg/validator"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("")

type userService struct {
	ur repository.UserRepository
	s  session_service.SessionService
	t  token.Maker
	v  validator.ValidatorInterface
	tx mongo.MongoTxInterface
}

func NewUserService(userRepo repository.UserRepository, v validator.ValidatorInterface, tx mongo.MongoTxInterface, s session_service.SessionService, t token.Maker) UserService {
	return &userService{ur: userRepo, v: v, tx: tx, s: s, t: t}
}

func (u *userService) Register(ctx context.Context, req models.RegisterReq) error {
	ctx, span := tracer.Start(ctx, "UserService.Register")
	defer span.End()
	err := u.v.Struct(req)
	if err != nil {
		return errors.BadRequest.Wrap(span, true, err, "Invalid request body")
	}

	user, err := u.ur.GetUser(ctx, bson.M{"$or": []bson.M{{"user_name": req.UserName}, {"email": req.UserName}}})
	if err != nil {
		return err
	}

	if user != nil {
		return errors.UIError.New(span, false, "Username/Email already in use!")
	}

	passwordHash, err := helpers.HashAndSalt(req.Password)

	if err != nil {
		return errors.InternalError.Wrapf(span, true, err, "Error while hashing password!")
	}

	userModel := models.User{
		Email:        req.Email,
		UserName:     req.UserName,
		PasswordHash: passwordHash,
		CreatedAt:    helpers.GetUTCTimeStamp(),
		UpdatedAt:    helpers.GetUTCTimeStamp(),
	}

	err = u.ur.CreateUser(ctx, userModel)
	if err != nil {
		return err
	}
	return err
}

func (u *userService) Login(ctx context.Context, req models.LoginReq) (*models.LoginRes, error) {
	ctx, span := tracer.Start(ctx, "UserService.Login")
	defer span.End()
	err := u.v.Struct(req)
	if err != nil {
		return nil, errors.BadRequest.Wrap(span, true, err, "Invalid request body")
	}

	user, err := u.ur.GetUser(ctx, bson.M{"$or": []bson.M{{"user_name": req.UserName}, {"email": req.UserName}}})
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.UIError.Newf(span, false, "User not found with user_name: %s", req.UserName)
	}

	_, err = helpers.ComparePasswords(user.PasswordHash, req.Password)
	if err != nil {
		return nil, err
	}

	sessionKey, err := u.s.CreateSession(ctx, &session_models.Session{ProfileID: (*user.Id).Hex()}, 3600)
	if err != nil {
		return nil, err
	}

	token, payload, err := u.t.CreateToken(*sessionKey, 3600*time.Second)
	if err != nil {
		return nil, errors.InternalError.Wrapf(span, true, err, "Error while generating token for payload: %#v", payload)
	}

	return &models.LoginRes{Token: token}, nil
}

func (u *userService) Logout(ctx context.Context, req models.LogoutReq) error {
	ctx, span := tracer.Start(ctx, "UserService.Logout")
	defer span.End()

	err := u.s.DeleteByID(ctx, req.SessionId)
	if err != nil {
		return errors.InternalError.Wrapf(span, true, err, "Error while deleting session for sessionId %s", req.SessionId)
	}
	return nil
}
