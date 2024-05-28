package services

import (
	"context"
	pb "product-service/genproto/product_service"
	"product-service/internal/entity"
)

func (d *productRPC) CreateProduct(ctx context.Context, in *pb.Product) (*pb.GetWithID, error) {
	respProduct, err := d.productUsecase.CreateProduct(ctx, &entity.Product{
		Id:          in.Id,
		Name:        in.Name,
		Description: in.Description,
		Category:    in.Category,
		MadeIn:      in.MadeIn,
		Color:       in.Color,
		Count:       in.Count,
		Cost:        in.Cost,
		Discount:    in.Discount,
		AgeMin:      in.AgeMin,
		AgeMax:      in.AgeMax,
		ForGender:   in.ForGender,
		Size:        in.ProductSize,
	})

	if err != nil {
		return nil, err
	}

	return &pb.GetWithID{Id: respProduct.Id}, nil
}

func (d *productRPC) UpdateProduct(ctx context.Context, in *pb.Product) (*pb.Product, error) {
	err := d.productUsecase.UpdateProduct(ctx, &entity.Product{
		Id:          in.Id,
		Name:        in.Name,
		Description: in.Description,
		Category:    in.Category,
		MadeIn:      in.MadeIn,
		Color:       in.Color,
		Count:       in.Count,
		Cost:        in.Cost,
		Discount:    in.Discount,
		AgeMin:      in.AgeMin,
		AgeMax:      in.AgeMax,
		ForGender:   in.ForGender,
		Size:        in.ProductSize,
	})

	if err != nil {
		return nil, err
	}

	return in, nil
}

func (d *productRPC) DeleteProduct(ctx context.Context, in *pb.GetWithID) (*pb.MoveResponse, error) {
	err := d.productUsecase.DeleteProduct(ctx, in.Id)

	if err != nil {
		return nil, err
	}

	return &pb.MoveResponse{Status: true}, nil
}

func (d *productRPC) GetProduct(ctx context.Context, in *pb.GetWithID) (*pb.Product, error) {
	product, err := d.productUsecase.GetProduct(ctx, map[string]string{"id": in.Id})

	if err != nil {
		return nil, err
	}

	return &pb.Product{
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
	}, nil
}

func (d *productRPC) ListProducts(ctx context.Context, in *pb.ListRequest) (*pb.ListProduct, error) {
	filter := &entity.ListRequest{
		Limit: in.Limit,
		Page:  in.Page,
	}

	products, err := d.productUsecase.ListProducts(ctx, filter)
	if err != nil {
		return nil, err
	}

	pbProducts := &pb.ListProduct{}
	for _, product := range products.Products {
		pbProducts.Products = append(pbProducts.Products, &pb.Product{
			Id:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Category:    product.Category,
			MadeIn:      product.MadeIn,
			Color:       product.Color,
			Count:       int64(product.Count),
			Cost:        float32(product.Cost),
			Discount:    float32(product.Discount),
			AgeMin:      product.AgeMin,
			AgeMax:      product.AgeMax,
			ForGender:   product.ForGender,
			ProductSize: product.Size,
		})
	}

	return &pb.ListProduct{
		Products:   pbProducts.Products,
		TotalCount: products.TotalCount,
	}, nil
}

func (d *productRPC) SearchProduct(ctx context.Context, in *pb.SearchRequest) (*pb.ListProduct, error) {
	products, err := d.productUsecase.SearchProduct(ctx, &entity.SearchRequest{
		Page:  in.Page,
		Limit: in.Limit,
		Params: map[string]string{
			"name": in.Params["name"],
		},
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

func (d *productRPC) GetDisableProducts(ctx context.Context, in *pb.ListRequest) (*pb.ListOrder, error) {
	orders, err := d.productUsecase.GetDisableProducts(ctx, &entity.ListRequest{
		Page:  in.Page,
		Limit: in.Limit,
	})
	if err != nil {
		return nil, err
	}

	pbOrders := &pb.ListOrder{}
	for _, order := range orders.Orders {
		pbOrders.Orders = append(pbOrders.Orders, &pb.Order{
			Id:        order.Id,
			ProductId: order.ProductID,
			UserId:    order.UserID,
			Status:    order.Status,
		})
	}

	return &pb.ListOrder{
		Orders:     pbOrders.Orders,
		TotalCount: orders.TotalCount,
	}, nil
}

func (p *productRPC) GetDiscountProducts(ctx context.Context, in *pb.ListRequest) (*pb.ListProduct, error) {
	orders, err := p.productUsecase.GetDiscountProducts(ctx, &entity.ListRequest{
		Page:  in.Page,
		Limit: in.Limit,
	})
	if err != nil {
		return nil, err
	}

	pbOrders := &pb.ListProduct{}
	for _, product := range orders.Products {
		pbOrders.Products = append(pbOrders.Products, &pb.Product{
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

	return &pb.ListProduct{
		Products:   pbOrders.Products,
		TotalCount: orders.TotalCount,
	}, nil
}
