package server

import (
	"coffee-shop/config"
	product_handler "coffee-shop/internal/product/delivery/grpc"
	product_repository "coffee-shop/internal/product/repository"
	product_service "coffee-shop/internal/product/service"
	"context"
	"strings"

	"net"

	"coffee-shop/pkg/interceptors"
	"coffee-shop/pkg/logger"
	"coffee-shop/pkg/metrics"
	"coffee-shop/pkg/pb"

	"coffee-shop/pkg/validator"

	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

type Server struct {
	//gin *gin.Engine
	grpc        *grpc.Server
	cfg         *config.Config
	im          *interceptors.InterceptorManager
	mongoConn   *mongo.Database
	mongoClient *mongo.Client
	redisConn   *redis.Client
	logger      logger.Logger
}

func NewServer(cfg *config.Config, mongoConn *mongo.Database, mongoClient *mongo.Client, redisConn *redis.Client, logger logger.Logger) *Server {
	return &Server{cfg: cfg, mongoConn: mongoConn, mongoClient: mongoClient, redisConn: redisConn, logger: logger}
}

func (s *Server) Run() error {
	m, err := metrics.CreateMetrics(s.cfg.Product_service.ServiceName)
	if err != nil {
		return err
	}
	validator := validator.NewValidator()

	productRepo := product_repository.NewProductRepository(s.mongoConn)
	productService := product_service.NewProductService(productRepo, validator)

	s.im = interceptors.NewInterceptorManager(s.logger, s.cfg, m)

	go s.RunGrpcServer(productService, s.im)

	s.RunGrpcGatewayServer(productService)

	//s.RunHttpServer(m, userController, quit)

	return nil
}

func (s *Server) RunGrpcServer(productService product_service.ProductService, i *interceptors.InterceptorManager) {
	l, err := net.Listen("tcp", s.cfg.Product_service.GrpcPort)

	if err != nil {
		s.logger.Fatalf("Failed to initialize grpc server: %v\n", err)
	}

	defer l.Close()

	interceptors := []grpc.UnaryServerInterceptor{i.RecordMetrics}

	if s.cfg.Environment.Env != "prod" {
		interceptors = append(interceptors, i.GrpcLogger)
	}

	s.grpc = grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: s.cfg.Grpc_global.MaxConnectionIdle * time.Minute,
		Timeout:           s.cfg.Grpc_global.Timeout * time.Second,
		MaxConnectionAge:  s.cfg.Grpc_global.MaxConnectionAge * time.Minute,
		Time:              s.cfg.Grpc_global.Timeout * time.Minute,
	}), grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()), grpc.ChainUnaryInterceptor(interceptors...))

	productHandler := product_handler.NewProductHandler(productService)
	pb.RegisterProductServiceServer(s.grpc, productHandler)

	if s.cfg.Environment.Env != "prod" {
		reflection.Register(s.grpc)
	}

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.grpc.Serve(l); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Failed to initialize grpc server: %v\n", err)
		}
	}()

	s.logger.Infof("Grpc Server Listening on port %v\n", s.cfg.Product_service.GrpcPort)

	<-quit

	log.Println("Shutting down grpc server...")
	s.grpc.GracefulStop()

}

func (s *Server) RunGrpcGatewayServer(productService product_service.ProductService) {
	// conn, err := grpc.DialContext(
	// 	context.Background(),
	// 	s.cfg.Product_service.GrpcPort,
	// 	grpc.WithBlock(),
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// )
	// if err != nil {
	// 	log.Fatalln("Failed to dial server:", err)
	// }

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{

		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: false,
		},
	})

	preflightHandler := func(w http.ResponseWriter, r *http.Request) {
		headers := []string{"Content-Type", "Accept", "Authorization"}
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
		methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	}

	allowCORS := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if origin := r.Header.Get("Origin"); origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
					preflightHandler(w, r)
					return
				}
			}
			h.ServeHTTP(w, r)
		})
	}

	grpcMux := runtime.NewServeMux(jsonOption)
	mux := http.NewServeMux()
	mux.Handle("/", allowCORS(grpcMux))

	// mux := http.NewServeMux()
	// mux.Handle("/", grpcMux)

	// mux.
	productHandler := product_handler.NewProductHandler(productService)
	//err = pb.RegisterProductServiceHandler(context.Background(), grpcMux, conn)
	err := pb.RegisterProductServiceHandlerServer(context.Background(), grpcMux, productHandler)
	if err != nil {
		s.logger.Fatal("Error while starting grpc gateway server ", err)
	}

	gwServer := &http.Server{
		Addr:           s.cfg.Product_service.GrpcGatewayPort,
		Handler:        grpcMux,
		MaxHeaderBytes: maxHeaderBytes,
		ReadTimeout:    s.cfg.Http_global.ReadTimeout,
		WriteTimeout:   s.cfg.Http_global.WriteTimeout,
	}

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	// log.Fatalln(gwServer.ListenAndServe())

	go func() {
		//err = gwServer.ListenAndServe()
		err = http.ListenAndServe(s.cfg.Product_service.GrpcGatewayPort, mux)

		if err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Failed to initialize http server: %v\n", err)
		}
	}()
	s.logger.Infof("Grpc Gateway Server Listening on port %v\n", s.cfg.Product_service.GrpcGatewayPort)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer cancel()
	log.Println("Shutting down grpc gateway server...")
	gwServer.Shutdown(ctx)

}

// func (s *Server) RunHttpServer(m metrics.Metrics, userController user_http.UserController, quit chan os.Signal) {
// 	server := &http.Server{
// 		Addr:           s.cfg.User_service.HttpPort,
// 		Handler:        s.gin,
// 		ReadTimeout:    time.Second * s.cfg.Http_global.ReadTimeout,
// 		WriteTimeout:   time.Second * s.cfg.Http_global.WriteTimeout,
// 		MaxHeaderBytes: maxHeaderBytes,
// 	}

// 	s.gin.Use(middleware.CORS())
// 	s.gin.Use(middleware.MaxBodyLimit())
// 	s.gin.Use(otelgin.Middleware(s.cfg.User_service.ServiceName))
// 	s.gin.Use(middleware.RequestID())

// 	s.gin.Use(middleware.RecordMetrics(m))

// 	user_http.UserRoutes(s.gin.Group("/user"), s.mw, userController)

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

// 	s.logger.Infof("Http Server Listening on port %v\n", s.cfg.User_service.HttpPort)

// 	<-quit
// 	log.Println("Shutting down http server...")
// 	server.Shutdown(ctx)
// }
