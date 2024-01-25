package http

import (
	"coffee-shop/config"
	"coffee-shop/internal/user/models"
	"coffee-shop/internal/user/service"

	errors "coffee-shop/pkg/custom_errors"
	"coffee-shop/pkg/utils"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("")

type userController struct {
	us  service.UserService
	cfg *config.Config
}

func NewUserController(cfg *config.Config, userService service.UserService) UserController {
	return &userController{cfg: cfg, us: userService}
}

func (u *userController) Register(c *gin.Context) {
	ctx, span := tracer.Start(utils.GetRequestCtx(c), "UserHandlerHttp.Register")
	defer span.End()
	reqInfo, _ := httputil.DumpRequest(c.Request, true)
	var req models.RegisterReq
	err := c.BindJSON(&req)
	if err != nil {
		ctxError := errors.BadRequest.Wrap(span, true, errors.AddRequestContextToError(string(reqInfo), err), "Error while parsing request body")
		errors.HandleHttpError(c, ctxError)
		return
	}

	err = u.us.Register(ctx, req)
	if err != nil {
		ctxError := errors.AddRequestContextToError(string(reqInfo), err)
		errors.HandleHttpError(c, ctxError)
		return
	}

	errors.HandleSuccess(c, nil)
}

func (u *userController) Login(c *gin.Context) {
	ctx, span := tracer.Start(utils.GetRequestCtx(c), "UserHandlerHttp.Login")
	defer span.End()
	reqInfo, _ := httputil.DumpRequest(c.Request, true)
	var req models.LoginReq
	err := c.BindJSON(&req)
	if err != nil {
		ctxError := errors.BadRequest.Wrap(span, true, errors.AddRequestContextToError(string(reqInfo), err), "Error while parsing request body")
		errors.HandleHttpError(c, ctxError)
		return
	}

	res, err := u.us.Login(ctx, req)
	if err != nil {
		ctxError := errors.AddRequestContextToError(string(reqInfo), err)
		errors.HandleHttpError(c, ctxError)
		return
	}

	utils.CreateSessionCookie(c, u.cfg, res.Token)

	errors.HandleSuccess(c, nil)
}

func (u *userController) Logout(c *gin.Context) {
	ctx, span := tracer.Start(utils.GetRequestCtx(c), "UserHandlerHttp.Logout")
	defer span.End()
	reqInfo, _ := httputil.DumpRequest(c.Request, true)
	var req models.LogoutReq
	sessionId, _ := c.Get("session_id")
	req.SessionId = sessionId.(string)
	err := u.us.Logout(ctx, req)
	if err != nil {
		ctxError := errors.AddRequestContextToError(string(reqInfo), err)
		errors.HandleHttpError(c, ctxError)
		return
	}
	utils.DeleteSessionCookie(c, u.cfg)

	errors.HandleSuccess(c, nil)
}

func (u *userController) IsCookieValid(c *gin.Context) {
	errors.HandleSuccess(c, nil)
}
