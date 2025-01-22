package repository

// import (
// 	"coffee-shop/internal/product/models"
// 	"coffee-shop/pkg/custom_errors"
// 	"context"
// 	"encoding/json"
// 	"time"

// 	"github.com/go-redis/redis/v8"
// )

// //tracer = otel.Tracer("")

// const (
// 	key = "products"
// )

// type productRedisRepo struct {
// 	redisClient *redis.Client
// }

// func NewProductRedisRepository(redisClient *redis.Client) ProductRedisRepository {
// 	return &productRedisRepo{redisClient: redisClient}
// }

// func (p *productRedisRepo) SetProducts(ctx context.Context, products []models.ProductDto) error {

// 	ctx, span := tracer.Start(ctx, "ProductRedisRepository.SetProducts")
// 	defer span.End()

// 	sessBytes, err := json.Marshal(products)
// 	if err != nil {
// 		return custom_errors.InternalError.Wrap(span, true, err, "Error while marshaling products object")
// 	}
// 	if err = p.redisClient.Set(ctx, key, sessBytes, time.Second*time.Duration(2400)).Err(); err != nil {
// 		return custom_errors.InternalError.Wrap(span, true, err, "Error while inserting products object in redis")
// 	}
// 	return nil
// }

// func (p *productRedisRepo) GetProducts(ctx context.Context) ([]models.ProductDto, error) {
// 	ctx, span := tracer.Start(ctx, "ProductRedisRepository.GetProducts")
// 	defer span.End()

// 	productRes := p.redisClient.Get(ctx, key)
// 	if productRes.Err() == redis.Nil {
// 		return nil, nil
// 	}
// 	productBytes, err := productRes.Bytes()
// 	if err != nil {
// 		return nil, custom_errors.InternalError.Wrap(span, true, err, "Error while retrieving session object from redis")
// 	}

// 	res := []models.ProductDto{}
// 	if err = json.Unmarshal(productBytes, &res); err != nil {
// 		return nil, custom_errors.InternalError.Wrap(span, true, err, "Error while unmarshling products object")
// 	}
// 	return res, nil
// }
