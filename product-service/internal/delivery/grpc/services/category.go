package services

import (
	"context"
	"errors"
	productproto "product-service/genproto/product_service"
	"product-service/internal/entity"
)

func (p *productRPC) CreateCategory(ctx context.Context, in *productproto.Category) (*productproto.Category, error) {
	category, err := p.productUsecase.CreateCategory(ctx, &entity.Category{
		ID:   in.Id,
		Name: in.Name,
	})

	if err != nil {
		return nil, err
	}

	return &productproto.Category{
		Id:   category.ID,
		Name: category.Name,
	}, nil
}

func (p *productRPC) DeleteCategory(ctx context.Context, in *productproto.GetWithID) (*productproto.MoveResponse, error) {
	err := p.productUsecase.DeleteCategory(ctx, in.Id)

	if err != nil {
		return nil, err
	}

	return &productproto.MoveResponse{
		Status: true,
	}, nil
}

func (p *productRPC) UpdateCategory(ctx context.Context, in *productproto.Category) (*productproto.Category, error) {
	response, err := p.productUsecase.UpdateCategory(ctx, &entity.Category{
		ID:   in.Id,
		Name: in.Name,
	})

	if err != nil {
		return nil, err
	}

	return &productproto.Category{
		Id:   response.ID,
		Name: response.Name,
	}, nil
}

func (p *productRPC) GetCategory(ctx context.Context, in *productproto.GetWithID) (*productproto.Category, error) {
	category, err := p.productUsecase.GetCategory(ctx, &entity.GetWithID{
		ID: in.Id,
	})

	if err != nil {
		return nil, err
	}

	return &productproto.Category{
		Id:   category.ID,
		Name: category.Name,
	}, nil
}

func (p *productRPC) ListCategories(ctx context.Context, in *productproto.ListRequest) (*productproto.ListCategory, error) {
	listCategory, err := p.productUsecase.ListCategories(ctx, &entity.ListRequest{
		Page:  in.Page,
		Limit: in.Limit,
	})

	if err != nil {
		return nil, err
	}

	var response productproto.ListCategory
	for _, category := range listCategory.Categories {
		response.Categories = append(response.Categories, &productproto.Category{
			Id:   category.ID,
			Name: category.Name,
		})
	}

	response.TotalCount = listCategory.TotalCount

	return &response, nil
}

func (p *productRPC) SearchCategory(ctx context.Context, in *productproto.SearchRequest) (*productproto.ListProduct, error) {
	products, err := p.productUsecase.SearchCategory(ctx, &entity.SearchRequest{
		Page:   in.Page,
		Limit:  in.Limit,
		Params: in.Params,
	})

	if err != nil {
		return nil, err
	}

	var response productproto.ListProduct
	for _, product := range products.Products {
		response.Products = append(response.Products, &productproto.Product{
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

func (p *productRPC) UniqueCategory(ctx context.Context, in *productproto.Params) (*productproto.MoveResponse, error) {
	exist, err := p.productUsecase.UniqueCategory(ctx, &entity.Params{
		Filter: map[string]string{
			"category_name": in.Filter["category_name"],
		},
	})
	if err != nil {
		return nil, err
	}
	if exist.Status {
		return nil, errors.New("category already created")
	}

	return &productproto.MoveResponse{
		Status: false,
	}, nil
}
