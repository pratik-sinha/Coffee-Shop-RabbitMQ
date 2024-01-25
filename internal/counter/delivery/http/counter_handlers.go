package http

import (
	"coffee-shop/config"
	"coffee-shop/internal/counter/models"
	"coffee-shop/internal/counter/service"

	errors "coffee-shop/pkg/custom_errors"
	"coffee-shop/pkg/utils"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("")

type counterController struct {
	cs  service.CounterService
	cfg *config.Config
}

func NewCounterController(cfg *config.Config, counterService service.CounterService) CounterController {
	return &counterController{cfg: cfg, cs: counterService}
}

func (cc *counterController) PlaceOrder(c *gin.Context) {
	ctx, span := tracer.Start(utils.GetRequestCtx(c), "CounterHandlerHttp.PlaceOrder")
	defer span.End()
	reqInfo, _ := httputil.DumpRequest(c.Request, true)
	var req models.PlaceOrderReq
	err := c.BindJSON(&req)
	if err != nil {
		ctxError := errors.BadRequest.Wrap(span, true, errors.AddRequestContextToError(string(reqInfo), err), "Error while parsing request body")
		errors.HandleHttpError(c, ctxError)
		return
	}
	profileId, _ := c.Get("profile_id")
	req.ProfileID = profileId.(string)
	err = cc.cs.PlaceOrder(ctx, req)
	if err != nil {
		ctxError := errors.AddRequestContextToError(string(reqInfo), err)
		errors.HandleHttpError(c, ctxError)
		return
	}

	errors.HandleSuccess(c, nil)
}

func (cc *counterController) GetOrders(c *gin.Context) {
	ctx, span := tracer.Start(utils.GetRequestCtx(c), "CounterHandlerHttp.PlaceOrder")
	defer span.End()
	reqInfo, _ := httputil.DumpRequest(c.Request, true)
	var req models.GetOrdersReq
	profileId, _ := c.Get("profile_id")
	req.ProfileID = profileId.(string)
	res, err := cc.cs.GetOrders(ctx, req)
	if err != nil {
		ctxError := errors.AddRequestContextToError(string(reqInfo), err)
		errors.HandleHttpError(c, ctxError)
		return
	}

	errors.HandleSuccess(c, res)
}
