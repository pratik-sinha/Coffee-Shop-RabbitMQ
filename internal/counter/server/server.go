package server

import (
	"coffee-shop/config"
	product "coffee-shop/internal/counter/external/product"
	external "coffee-shop/internal/counter/external/publisher"
	counter_repository "coffee-shop/internal/counter/repository"
	counter_service "coffee-shop/internal/counter/service"
	session_repository "coffee-shop/internal/shared/session/repository"
	session_service "coffee-shop/internal/shared/session/service"

	"coffee-shop/pkg/logger"
	"coffee-shop/pkg/metrics"
	"coffee-shop/pkg/rabbitmq/consumer"
	"coffee-shop/pkg/token"

	mongoDB "coffee-shop/pkg/db/mongo"
	"coffee-shop/pkg/middleware"
	"coffee-shop/pkg/validator"

	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	counter_http "coffee-shop/internal/counter/delivery/http"
	"coffee-shop/internal/counter/delivery/rabbitmq"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

type Server struct {
	gin *gin.Engine
	//	grpc        *grpc.Server
	cfg         *config.Config
	mw          *middleware.MiddlewareManager
	amqConn     *amqp091.Connection
	mongoConn   *mongo.Database
	mongoClient *mongo.Client
	redisConn   *redis.Client
	logger      logger.Logger
}

func NewServer(cfg *config.Config, mongoConn *mongo.Database, mongoClient *mongo.Client, redisConn *redis.Client, amqConn *amqp091.Connection, logger logger.Logger) *Server {
	return &Server{gin: gin.Default(), cfg: cfg, mongoConn: mongoConn, mongoClient: mongoClient, redisConn: redisConn, amqConn: amqConn, logger: logger}
}

func (s *Server) Run() error {
	m, err := metrics.CreateMetrics(s.cfg.User_service.ServiceName)
	if err != nil {
		return err
	}
	validator := validator.NewValidator()
	mongoTx := mongoDB.NewTxInterface(s.mongoClient)

	//im := interceptors.NewInterceptorManager(s.logger, s.cfg, m)
	tokenMaker, err := token.NewPasetoMaker(s.cfg.Http_global.TokenSymmetricKey)
	if err != nil {
		return err
	}

	eventConsumer := consumer.NewConsumer(s.amqConn)
	if err != nil {
		return err
	}

	counterConsumer := rabbitmq.NewCounterConsumer(eventConsumer)
	counterConsumer.Configure()

	baristaEventPublisher, err := external.NewBaristaEventPublisher(s.amqConn)
	if err != nil {
		return err
	}
	baristaEventPublisher.Configure()

	kitchenEventPublisher, err := external.NewKitchenEventPublisher(s.amqConn)
	if err != nil {
		return err
	}
	kitchenEventPublisher.Configure()

	productClient, err := product.NewGRPCProductClient(s.cfg)
	if err != nil {
		return err
	}

	sessRepo := session_repository.NewSessionRepository(s.redisConn)
	sessService := session_service.NewSessionService(sessRepo)

	counterRepo := counter_repository.NewCounterRepository(s.mongoConn)
	counterService := counter_service.NewCounterService(counterRepo, validator, mongoTx, productClient, kitchenEventPublisher, baristaEventPublisher)
	userController := counter_http.NewCounterController(s.cfg, counterService)

	s.mw = middleware.NewMiddlewareManager(sessService, tokenMaker, s.cfg, s.logger)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go counterConsumer.StartConsumer(counterService)

	s.RunHttpServer(m, userController, quit)

	return nil
}

func (s *Server) RunHttpServer(m metrics.Metrics, counterController counter_http.CounterController, quit chan os.Signal) {
	server := &http.Server{
		Addr:           s.cfg.Counter_service.HttpPort,
		Handler:        s.gin,
		ReadTimeout:    time.Second * s.cfg.Http_global.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Http_global.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	s.gin.Use(middleware.CORS())
	s.gin.Use(middleware.MaxBodyLimit())
	s.gin.Use(otelgin.Middleware(s.cfg.Counter_service.ServiceName))
	s.gin.Use(middleware.RequestID())

	s.gin.Use(middleware.RecordMetrics(m))

	counter_http.CounterRoutes(s.gin.Group("/counter"), s.mw, counterController)

	s.gin.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Up and Running..."})
	})

	// Graceful server shutdown - https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Failed to initialize http server: %v\n", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()

	s.logger.Infof("Http Server Listening on port %v\n", s.cfg.Counter_service.HttpPort)

	<-quit
	log.Println("Shutting down http server...")
	server.Shutdown(ctx)
}
