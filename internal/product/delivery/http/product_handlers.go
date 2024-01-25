package http

import (
	"coffee-shop/config"
	"coffee-shop/internal/product/models"
	"coffee-shop/internal/product/service"

	errors "coffee-shop/pkg/custom_errors"
	"coffee-shop/pkg/utils"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("")

type productController struct {
	ps  service.ProductService
	cfg *config.Config
}

func NewProductController(cfg *config.Config, productService service.ProductService) ProductController {
	return &productController{cfg: cfg, ps: productService}
}

func (p *productController) GetProducts(c *gin.Context) {
	ctx, span := tracer.Start(utils.GetRequestCtx(c), "ProductHandlerHttp.GetProducts")
	defer span.End()
	reqInfo, _ := httputil.DumpRequest(c.Request, true)
	res, err := p.ps.GetProducts(ctx)
	if err != nil {
		ctxError := errors.AddRequestContextToError(string(reqInfo), err)
		errors.HandleHttpError(c, ctxError)
		return
	}

	errors.HandleSuccess(c, res)
}

func (p *productController) GetProductsByType(c *gin.Context) {
	ctx, span := tracer.Start(utils.GetRequestCtx(c), "ProductHandlerHttp.GetProductsByType")
	defer span.End()
	reqInfo, _ := httputil.DumpRequest(c.Request, true)
	var req models.GetProductsByTypeReq
	err := c.BindJSON(&req)
	if err != nil {
		ctxError := errors.BadRequest.Wrap(span, true, errors.AddRequestContextToError(string(reqInfo), err), "Error while parsing request body")
		errors.HandleHttpError(c, ctxError)
		return
	}

	res, err := p.ps.GetProductsByType(ctx, req)
	if err != nil {
		ctxError := errors.AddRequestContextToError(string(reqInfo), err)
		errors.HandleHttpError(c, ctxError)
		return
	}

	errors.HandleSuccess(c, res)
}
