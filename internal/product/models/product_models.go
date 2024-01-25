package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id          *primitive.ObjectID `bson:"_id,omitempty"`
	Type        int32               `bson:"type"`
	KitchenType int32               `bson:"kitchen_type"`
	Name        string              `bson:"name"`
	Price       float32             `bson:"price"`
	Image       string              `bson:"image"`
}

type Products []Product

func (p *Products) GetProductDto() []ProductDto {
	res := []ProductDto{}
	for _, v := range *p {
		res = append(res, ProductDto{
			Type:        v.Type,
			Name:        v.Name,
			Image:       v.Image,
			Price:       v.Price,
			KitchenType: v.KitchenType,
		})
	}
	return res
}
