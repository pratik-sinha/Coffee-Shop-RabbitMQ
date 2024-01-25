package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment     EnvironmentConfig
	Http_global     Http_GlobalConfig
	Grpc_global     Grpc_GlobalConfig
	User_service    UserServiceConfig
	Product_service ProductServiceConfig
	Counter_service CounterServiceConfig
	Kitchen_service KitchenServiceConfig
	Barista_service BaristaServiceConfig
	Mongo           MongoConfig
	Redis           RedisConfig
	RabbitMQ        RabbitMQConfig
	Cookie          Cookie
	Session         Session
	Logger          Logger
	Otel            OtelConfig
}

type OtelConfig struct {
	Host string
}

type EnvironmentConfig struct {
	Env string
}

type Http_GlobalConfig struct {
	TokenSymmetricKey string
	CookieName        string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	SSL               bool
	CtxDefaultTimeout time.Duration
	CSRF              bool
	Debug             bool
}

type Cookie struct {
	Name     string
	Domain   string
	MaxAge   int
	Secure   bool
	HTTPOnly bool
}

type Session struct {
	Prefix string
	Name   string
	Expire int
}

type Grpc_GlobalConfig struct {
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	CtxDefaultTimeout time.Duration
	MaxConnectionIdle time.Duration
	Timeout           time.Duration
	MaxConnectionAge  time.Duration
	Time              time.Duration
}

type UserServiceConfig struct {
	ServiceName string
	HttpPort    string
}

type ProductServiceConfig struct {
	ServiceName     string
	GrpcPort        string
	GrpcGatewayPort string
}

type CounterServiceConfig struct {
	ServiceName string
	HttpPort    string
}

type KitchenServiceConfig struct {
	ServiceName string
}

type BaristaServiceConfig struct {
	ServiceName string
}

type MongoConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

type RedisConfig struct {
	Host     string
	Password string
}

type RabbitMQConfig struct {
	URL string
}

type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

func LoadConfig(fileName string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(fileName)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}
	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
