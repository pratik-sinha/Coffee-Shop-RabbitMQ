package main

import (
	"coffee-shop/config"
	"coffee-shop/internal/user/server"
	"coffee-shop/pkg/db/mongo"
	"coffee-shop/pkg/db/redis"
	"coffee-shop/pkg/logger"
	"coffee-shop/pkg/utils"
	"log"
	"os"
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

	redisClient, disconnectRedis := redis.NewRedisClient(cfg)
	defer disconnectRedis()

	// shutdownTracing, err := tracing.InitTracerProvider(cfg, cfg.User_service.ServiceName)
	// if err != nil {
	// 	appLogger.Fatalf("Tracing init: %s", err)
	// }
	// defer shutdownTracing()

	// shutdownMetrics, err := metrics.InitMetricsProvider(cfg, cfg.User_service.ServiceName)
	// if err != nil {
	// 	appLogger.Fatalf("Metrics init: %s", err)
	// }
	// defer shutdownMetrics()

	s := server.NewServer(cfg, mongoDb, mongoClient, redisClient, appLogger)

	if err := s.Run(); err != nil {
		appLogger.Fatal(err)
	}

}
