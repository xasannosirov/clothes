package services

import (
	"context"
	productproto "product-service/genproto/product_service"
	"product-service/internal/entity"
)

func (d *productRPC) CreateCategory(ctx context.Context, in *productproto.Category) (*productproto.Category, error) {
	category, err := d.productUsecase.CreateCategory(ctx, &entity.Category{
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

func (d *productRPC) DeleteCategory(ctx context.Context, in *productproto.GetWithID) (*productproto.DeleteResponse, error) {
	err := d.productUsecase.DeleteCategory(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &productproto.DeleteResponse{
		Status: true,
	}, nil
}

func (d *productRPC) GetAllCategory(ctx context.Context, in *productproto.ListRequest) (*productproto.ListCategory, error) {
	listCategory, err := d.productUsecase.ListCategory(ctx, &entity.ListRequest{
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

func (d *productRPC) UpdateCategory(ctx context.Context, in *productproto.Category) (*productproto.Category, error) {
	updatedCategory, err := d.productUsecase.UpdateCategory(ctx, &entity.Category{
		ID:   in.Id,
		Name: in.Name,
	})
	if err != nil {
		return nil, err
	}

	return &productproto.Category{
		Id:   updatedCategory.ID,
		Name: updatedCategory.Name,
	}, nil
}

func (d *productRPC) GetCategory(ctx context.Context, in *productproto.GetWithID) (*productproto.Category, error) {
	category, err := d.productUsecase.GetCategory(ctx, &entity.GetWithID{
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
