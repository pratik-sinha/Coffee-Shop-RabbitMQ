package interceptors

import (
	"coffee-shop/config"
	"coffee-shop/pkg/custom_errors"
	"coffee-shop/pkg/logger"
	"coffee-shop/pkg/metrics"
	"context"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

type InterceptorManager struct {
	logger logger.Logger
	cfg    *config.Config
	metr   metrics.Metrics
}

func NewInterceptorManager(logger logger.Logger, cfg *config.Config, metr metrics.Metrics) *InterceptorManager {
	return &InterceptorManager{logger: logger, cfg: cfg, metr: metr}
}

func (i *InterceptorManager) GrpcLogger(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	startTime := time.Now()
	res, err := handler(ctx, req)
	timeDur := time.Since(startTime).Microseconds()
	status := http.StatusOK
	if err != nil {
		status = custom_errors.GetHttpStatusCode(custom_errors.GetErrorType(err))
		i.logger.Errorf("Grpc Method:%s status:%d duration:%d err:%v  ", info.FullMethod, status, timeDur, err)
	} else {
		i.logger.Infof("Grpc Method:%s status:%d duration:%d  ", info.FullMethod, status, timeDur)

	}

	return res, err
}

func (i *InterceptorManager) RecordMetrics(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	startTime := time.Now()
	res, err := handler(ctx, req)
	status := http.StatusOK
	if err != nil {
		status = custom_errors.GetHttpStatusCode(custom_errors.GetErrorType(err))
	}
	timeDur := time.Since(startTime).Milliseconds()
	i.metr.IncHits(ctx, status, info.FullMethod, info.FullMethod)
	i.metr.ObserveResponseTime(ctx, status, info.FullMethod, info.FullMethod, float64(timeDur))
	return res, err
}
