package main

import (
	"coffee-shop/config"
	_ "coffee-shop/internal/product/migrations"
	"coffee-shop/internal/product/server"
	"coffee-shop/pkg/db/mongo"
	"coffee-shop/pkg/db/redis"
	"coffee-shop/pkg/logger"
	"coffee-shop/pkg/metrics"
	"coffee-shop/pkg/tracing"
	"coffee-shop/pkg/utils"
	"log"
	"os"

	migrate "github.com/xakep666/mongo-migrate"
)

func main() {
	log.Println("Starting api server")
	configPath := utils.GetConfigPath(os.Getenv("config"))
	log.Println(configPath)

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}
	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof("LogLevel: %s, Mode: %s", cfg.Logger.Level, cfg.Environment.Env)
	mongoDb, mongoClient, disconnectMongo, err := mongo.ConnectMongoDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	}
	defer disconnectMongo()

	migrate.SetDatabase(mongoDb)
	if err := migrate.Up(migrate.AllAvailable); err != nil {
		appLogger.Fatal(err)
	}
	redisClient, disconnectRedis := redis.NewRedisClient(cfg)
	defer disconnectRedis()

	shutdownTracing, err := tracing.InitTracerProvider(cfg, cfg.User_service.ServiceName)
	if err != nil {
		appLogger.Fatalf("Tracing init: %s", err)
	}
	defer shutdownTracing()

	shutdownMetrics, err := metrics.InitMetricsProvider(cfg, cfg.User_service.ServiceName)
	if err != nil {
		appLogger.Fatalf("Metrics init: %s", err)
	}
	defer shutdownMetrics()

	s := server.NewServer(cfg, mongoDb, mongoClient, redisClient, appLogger)

	if err := s.Run(); err != nil {
		appLogger.Fatal(err)
	}

}
