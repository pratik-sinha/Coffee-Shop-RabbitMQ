package server

import (
	"coffee-shop/config"
	session_repository "coffee-shop/internal/shared/session/repository"
	session_service "coffee-shop/internal/shared/session/service"
	user_repository "coffee-shop/internal/user/repository"
	user_service "coffee-shop/internal/user/service"
	"fmt"
	"strings"

	"coffee-shop/pkg/logger"
	"coffee-shop/pkg/metrics"
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

	user_http "coffee-shop/internal/user/delivery/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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
	mongoConn   *mongo.Database
	mongoClient *mongo.Client
	redisConn   *redis.Client
	logger      logger.Logger
}

func NewServer(cfg *config.Config, mongoConn *mongo.Database, mongoClient *mongo.Client, redisConn *redis.Client, logger logger.Logger) *Server {
	return &Server{gin: gin.Default(), cfg: cfg, mongoConn: mongoConn, mongoClient: mongoClient, redisConn: redisConn, logger: logger}
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

	sessRepo := session_repository.NewSessionRepository(s.redisConn)
	sessService := session_service.NewSessionService(sessRepo)

	userRepo := user_repository.NewUserRepository(s.mongoConn)
	userService := user_service.NewUserService(userRepo, validator, mongoTx, sessService, tokenMaker)
	userController := user_http.NewUserController(s.cfg, userService)

	s.mw = middleware.NewMiddlewareManager(sessService, tokenMaker, s.cfg, s.logger)

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	//go s.RunGrpcServer(accService, im, quit)

	//go s.RunGrpcGatewayServer(accService, quit)

	s.RunHttpServer(m, userController, quit)

	return nil
}

func (s *Server) RunHttpServer(m metrics.Metrics, userController user_http.UserController, quit chan os.Signal) {
	fmt.Print(strings.Split(":", s.cfg.User_service.HttpPort))
	server := &http.Server{
		Addr:           s.cfg.User_service.HttpPort,
		Handler:        s.gin,
		ReadTimeout:    time.Second * s.cfg.Http_global.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Http_global.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	s.gin.Use(middleware.CORS())
	s.gin.Use(middleware.MaxBodyLimit())
	s.gin.Use(otelgin.Middleware(s.cfg.User_service.ServiceName))
	s.gin.Use(middleware.RequestID())

	s.gin.Use(middleware.RecordMetrics(m))

	user_http.UserRoutes(s.gin.Group("/user"), s.mw, userController)

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

	s.logger.Infof("Http Server Listening on port %v\n", s.cfg.User_service.HttpPort)

	<-quit
	log.Println("Shutting down http server...")
	server.Shutdown(ctx)
}

// func (s *Server) RunGrpcServer(accService service.AccountService, i *interceptors.InterceptorManager, quit chan os.Signal) {
// 	l, err := net.Listen("tcp", s.cfg.Grpc.Port)

// 	if err != nil {
// 		s.logger.Fatalf("Failed to initialize grpc server: %v\n", err)
// 	}

// 	defer l.Close()

// 	interceptors := []grpc.UnaryServerInterceptor{i.RecordMetrics}

// 	if s.cfg.Environment.Env != "prod" {
// 		interceptors = append(interceptors, i.GrpcLogger)
// 	}

// 	s.grpc = grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
// 		MaxConnectionIdle: s.cfg.Grpc.MaxConnectionIdle * time.Minute,
// 		Timeout:           s.cfg.Grpc.Timeout * time.Second,
// 		MaxConnectionAge:  s.cfg.Grpc.MaxConnectionAge * time.Minute,
// 		Time:              s.cfg.Grpc.Timeout * time.Minute,
// 	}), grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()), grpc.ChainUnaryInterceptor(interceptors...))

// 	accControllerGrpc := account_grpc.NewAccountHandler(accService)
// 	pb.RegisterBankServer(s.grpc, accControllerGrpc)

// 	if s.cfg.Environment.Env != "prod" {
// 		reflection.Register(s.grpc)
// 	}

// 	go func() {
// 		if err := s.grpc.Serve(l); err != nil && err != http.ErrServerClosed {
// 			s.logger.Fatalf("Failed to initialize grpc server: %v\n", err)
// 		}
// 	}()

// 	s.logger.Infof("Grpc Server Listening on port %v\n", s.cfg.Grpc.Port)

// 	<-quit

// 	log.Println("Shutting down grpc server...")
// 	s.grpc.GracefulStop()

// }

// func (s *Server) RunGrpcGatewayServer(accService service.AccountService, quit chan os.Signal) {
// 	// l, err := net.Listen("tcp", s.cfg.Grpc.GatewayPort)

// 	// if err != nil {
// 	// 	s.logger.Fatalf("Failed to initialize grpc gatewat server: %v\n", err)
// 	// }

// 	// defer l.Close()

// 	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
// 		MarshalOptions: protojson.MarshalOptions{
// 			UseProtoNames: true,
// 		},
// 		UnmarshalOptions: protojson.UnmarshalOptions{
// 			DiscardUnknown: true,
// 		},
// 	})

// 	grpcMux := runtime.NewServeMux(jsonOption)
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	mux := http.NewServeMux()
// 	mux.Handle("/", grpcMux)

// 	accControllerGrpc := account_grpc.NewAccountHandler(accService)
// 	err := pb.RegisterBankHandlerServer(ctx, grpcMux, accControllerGrpc)

// 	if err != nil {
// 		s.logger.Fatal("Error while starting grpc gateway server ", err)
// 	}

// 	go func() {
// 		err = http.ListenAndServe(s.cfg.Grpc.GatewayPort, mux)
// 		if err != nil {
// 			s.logger.Fatal("Error while starting grpc gateway server ", err)
// 		}
// 	}()
// 	s.logger.Infof("Grpc Gateway Server Listening on port %v\n", s.cfg.Grpc.GatewayPort)
// 	<-quit

// 	log.Println("Shutting down grpc gateway server...")
// }
