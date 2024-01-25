package migrations

import (
	"coffee-shop/internal/product/models"
	"context"

	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	migrate.Register(func(db *mongo.Database) error {
		products := []interface{}{models.Product{
			Name:        "CAPPUCCINO",
			Type:        0,
			Price:       4.5,
			KitchenType: 0,
			Image:       "img/CAPPUCCINO.png",
		}, models.Product{
			Name:        "COFFEE_BLACK",
			Type:        1,
			Price:       3,
			KitchenType: 0,
			Image:       "img/COFFEE_BLACK.png",
		}, models.Product{
			Name:        "COFFEE_WITH_ROOM",
			Type:        2,
			Price:       3,
			KitchenType: 0,
			Image:       "img/COFFEE_WITH_ROOM.png",
		}, models.Product{
			Name:        "ESPRESSO",
			Type:        3,
			Price:       3.5,
			KitchenType: 0,
			Image:       "img/ESPRESSO.png",
		}, models.Product{
			Name:        "ESPRESSO_DOUBLE",
			Type:        4,
			Price:       4.5,
			KitchenType: 0,
			Image:       "img/ESPRESSO_DOUBLE.png",
		}, models.Product{
			Name:        "LATTE",
			Type:        5,
			Price:       4.5,
			KitchenType: 0,
			Image:       "img/LATTE.png",
		}, models.Product{
			Name:        "CAKEPOP",
			Type:        6,
			Price:       2.5,
			KitchenType: 1,
			Image:       "img/CAKEPOP.png",
		}, models.Product{
			Name:        "CROISSANT",
			Type:        7,
			Price:       3.25,
			KitchenType: 1,
			Image:       "img/CROISSANT.png",
		}, models.Product{
			Name:        "MUFFIN",
			Type:        8,
			Price:       3,
			KitchenType: 1,
			Image:       "img/MUFFIN.png",
		}, models.Product{
			Name:        "CROISSANT_CHOCOLATE",
			Type:        9,
			Price:       3.5,
			KitchenType: 1,
			Image:       "img/CROISSANT_CHOCOLATE.png",
		}}

		err := db.CreateCollection(context.Background(), "products")
		if err != nil {
			return err
		}
		_, err = db.Collection("products").InsertMany(context.Background(), products)
		if err != nil {
			return err
		}

		return nil
	}, func(db *mongo.Database) error {
		return nil
	})
}
