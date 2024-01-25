package models

import "coffee-shop/pkg/pb"

type GetProductsByTypeReq struct {
	ProductTypes []int32 `json:"product_types" validate:"required"`
}

func (p *GetProductsByTypeReq) ConvertFromProto(req *pb.GetProductsByTypeReq) {

	p.ProductTypes = req.ProductTypes
}

type ProductDto struct {
	Type        int32   `json:"type"`
	Name        string  `json:"name"`
	KitchenType int32   `json:"kitchen_type"`
	Image       string  `json:"image"`
	Price       float32 `json:"price"`
}

type GetProductsRes struct {
	Products []ProductDto `json:"products"`
}

func (p *GetProductsRes) ConvertToProto() *pb.GetProductsRes {
	var res []*pb.Product
	for _, product := range p.Products {
		res = append(res, &pb.Product{
			Type:        product.Type,
			Name:        product.Name,
			Image:       product.Image,
			Price:       product.Price,
			KitchenType: product.KitchenType,
		})
	}
	return &pb.GetProductsRes{
		Products: res,
	}
}
