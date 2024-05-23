package services

import (
	"context"
	pb "product-service/genproto/product_service"
	"product-service/internal/entity"
)

func (d *productRPC) CreateOrder(ctx context.Context, in *pb.Order) (*pb.GetWithID, error) {
	order, err := d.productUsecase.CreateOrder(ctx, &entity.Order{
		Id:        in.Id,
		ProductID: in.ProductId,
		UserID:    in.UserId,
		Count:     in.Count,
		Status:    in.Status,
	})

	if err != nil {
		return nil, err
	}

	return &pb.GetWithID{Id: order.Id}, nil
}

func (d *productRPC) GetOrder(ctx context.Context, in *pb.GetWithID) (*pb.Order, error) {
	order, err := d.productUsecase.GetOrder(ctx, map[string]string{"id": in.Id})

	if err != nil {
		return nil, err
	}

	return &pb.Order{
		Id:        order.Id,
		ProductId: order.ProductID,
		UserId:    order.UserID,
		Count:     order.Count,
		Status:    order.Status,
	}, nil
}

func (p *productRPC) DeleteOrder(ctx context.Context, in *pb.GetWithID) (*pb.MoveResponse, error) {
	err := p.productUsecase.DeleteOrder(ctx, map[string]string{
		"order_id": in.Id,
	})

	if err != nil {
		return nil, err
	}

	return &pb.MoveResponse{Status: true}, nil
}

func (p *productRPC) UserOrderHistory(ctx context.Context, in *pb.SearchRequest) (*pb.ListProduct, error) {
	products, err := p.productUsecase.UserOrderHistory(ctx, &entity.SearchRequest{
		Page:   in.Page,
		Limit:  in.Limit,
		Params: in.Params,
	})

	if err != nil {
		return nil, err
	}

	var response pb.ListProduct
	for _, product := range products.Products {
		response.Products = append(response.Products, &pb.Product{
			Id:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Category:    product.Category,
			MadeIn:      product.MadeIn,
			Color:       product.Color,
			Count:       product.Count,
			Cost:        product.Cost,
			Discount:    product.Discount,
			AgeMin:      product.AgeMin,
			AgeMax:      product.AgeMax,
			ForGender:   product.ForGender,
			ProductSize: product.Size,
		})
	}
	response.TotalCount = products.TotalCount

	return &response, nil
}
