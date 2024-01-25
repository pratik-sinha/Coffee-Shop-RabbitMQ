package server

import (
	"coffee-shop/config"
	external "coffee-shop/internal/kitchen/external/publisher"
	kitchen_repository "coffee-shop/internal/kitchen/repository"
	kitchen_service "coffee-shop/internal/kitchen/service"
	"coffee-shop/pkg/logger"
	"coffee-shop/pkg/rabbitmq/consumer"
	"coffee-shop/pkg/rabbitmq/publisher"

	mongoDB "coffee-shop/pkg/db/mongo"
	"coffee-shop/pkg/middleware"
	"coffee-shop/pkg/validator"

	"os"
	"os/signal"
	"syscall"

	"coffee-shop/internal/kitchen/delivery/rabbitmq"

	"github.com/go-redis/redis/v8"
	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

type Server struct {
	//gin *gin.Engine
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
	return &Server{cfg: cfg, mongoConn: mongoConn, mongoClient: mongoClient, redisConn: redisConn, amqConn: amqConn, logger: logger}
}

func (s *Server) Run() error {
	// m, err := metrics.CreateMetrics(s.cfg.User_service.ServiceName)
	// if err != nil {
	// 	return err
	// }
	validator := validator.NewValidator()
	mongoTx := mongoDB.NewTxInterface(s.mongoClient)

	eventConsumer := consumer.NewConsumer(s.amqConn)

	kitchenConsumer := rabbitmq.NewKitchenConsumer(eventConsumer)
	kitchenConsumer.Configure()

	eventPublisher, err := publisher.NewPublisher(s.amqConn)
	if err != nil {
		return err
	}

	counterEventPublisher := external.NewCounterEventPublisher(eventPublisher)
	counterEventPublisher.Configure()

	kitchenRepo := kitchen_repository.NewKitchenRepository(s.mongoConn)
	counterService := kitchen_service.NewKitchenService(kitchenRepo, validator, mongoTx, counterEventPublisher)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go kitchenConsumer.StartConsumer(counterService)
	s.logger.Info("Kitchen consumer running!")
	<-quit
	s.logger.Info("Kitchen consumer exiting!")

	return nil
}

// func (s *Server) RunHttpServer(m metrics.Metrics, counterController counter_http.CounterController, quit chan os.Signal) {
// 	server := &http.Server{
// 		Addr:           s.cfg.Counter_service.HttpPort,
// 		Handler:        s.gin,
// 		ReadTimeout:    time.Second * s.cfg.Http_global.ReadTimeout,
// 		WriteTimeout:   time.Second * s.cfg.Http_global.WriteTimeout,
// 		MaxHeaderBytes: maxHeaderBytes,
// 	}

// 	s.gin.Use(middleware.CORS())
// 	s.gin.Use(middleware.MaxBodyLimit())
// 	s.gin.Use(otelgin.Middleware(s.cfg.Counter_service.ServiceName))
// 	s.gin.Use(middleware.RequestID())

// 	s.gin.Use(middleware.RecordMetrics(m))

// 	counter_http.CounterRoutes(s.gin.Group("/kitchen"), s.mw, counterController)

// 	s.gin.GET("/", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{"data": "Up and Running..."})
// 	})

// 	// Graceful server shutdown - https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
// 	go func() {
// 		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			s.logger.Fatalf("Failed to initialize http server: %v\n", err)
// 		}
// 	}()

// 	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
// 	defer cancel()

// 	s.logger.Infof("Http Server Listening on port %v\n", s.cfg.Counter_service.HttpPort)

// 	<-quit
// 	log.Println("Shutting down http server...")
// 	server.Shutdown(ctx)
// }
