package services

import (
	"context"
	pb "product-service/genproto/product_service"
	grpc "product-service/internal/delivery"
	"product-service/internal/entity"

	"go.uber.org/zap"
)

func (d *productRPC) GetAllComments(ctx context.Context, in *pb.ListRequest) (*pb.ListCommentResponse, error) {
	filter := &entity.ListRequest{
		Limit: in.Limit,
		Page:  in.Page,
	}

	comments, err := d.productUsecase.GetAllComments(ctx, filter)
	if err != nil {
		d.logger.Error("commentUseCase.List", zap.Error(err))
		return &pb.ListCommentResponse{}, grpc.Error(ctx, err)
	}

	pbComments := &pb.ListCommentResponse{}
	for _, comment := range comments.Comments {
		pbComments.Comments = append(pbComments.Comments, &pb.Comment{
			Id:        comment.Id,
			ProductId: comment.ProductID,
			UserId:    comment.UserID,
			Comment:   comment.Comment,
			CreatedAt: comment.CreatedAt.String(),
			UpdatedAt: comment.UpdatedAt.String(),
		})
	}

	return &pb.ListCommentResponse{
		Comments:   pbComments.Comments,
		TotalCount: comments.TotalCount,
	}, nil
}

func (d *productRPC) GetAllStars(ctx context.Context, in *pb.ListRequest) (*pb.ListStarsResponse, error) {
	filter := &entity.ListRequest{
		Limit: in.Limit,
		Page:  in.Page,
	}

	stars, err := d.productUsecase.GetAllStars(ctx, filter)
	if err != nil {
		return &pb.ListStarsResponse{}, grpc.Error(ctx, err)
	}
	pbStars := &pb.ListStarsResponse{}
	for _, star := range stars.Stars {
		pbStars.Stars = append(pbStars.Stars, &pb.Star{
			Id:        star.Id,
			ProductId: star.ProductID,
			UserId:    star.UserID,
			Star:      star.Stars,
			CreatedAt: star.CreatedAt.String(),
			UpdatedAt: star.UpdatedAt.String(),
		})
	}

	return &pb.ListStarsResponse{
		Stars:      pbStars.Stars,
		TotalCount: stars.TotalCount,
	}, nil
}

func (d *productRPC) GetProductOrders(ctx context.Context, in *pb.GetWithID) (*pb.ListOrderResponse, error) {
	orders, err := d.productUsecase.GetProductOrders(ctx, &entity.GetWithID{ID: in.Id})
	if err != nil {
		d.logger.Error("productUsecase.GetProductOrders", zap.Error(err))
		return &pb.ListOrderResponse{}, grpc.Error(ctx, err)
	}

	pbOrders := &pb.ListOrderResponse{}
	for _, order := range orders.Orders {
		pbOrders.Orders = append(pbOrders.Orders, &pb.Order{
			Id:        order.Id,
			ProductId: order.ProductID,
			UserId:    order.UserID,
			Status:    order.Status,
			CreatedAt: order.CreatedAt.String(),
			UpdatedAt: order.UpdatedAt.String(),
		})
	}
	return &pb.ListOrderResponse{
		Orders:     pbOrders.Orders,
		TotalCount: orders.TotalCount,
	}, nil
}

func (d *productRPC) GetProductComments(ctx context.Context, in *pb.GetWithID) (*pb.ListCommentResponse, error) {
	comments, err := d.productUsecase.GetProductComments(ctx, &entity.GetWithID{ID: in.Id})
	if err != nil {
		d.logger.Error("productUsecase.GetProductComments", zap.Error(err))
		return &pb.ListCommentResponse{}, grpc.Error(ctx, err)
	}

	pbComments := &pb.ListCommentResponse{}
	for _, comment := range comments.Comments {
		pbComments.Comments = append(pbComments.Comments, &pb.Comment{
			Id:        comment.Id,
			ProductId: comment.ProductID,
			UserId:    comment.UserID,
			Comment:   comment.Comment,
			CreatedAt: comment.CreatedAt.String(),
			UpdatedAt: comment.UpdatedAt.String(),
		})
	}
	return &pb.ListCommentResponse{
		Comments:   pbComments.Comments,
		TotalCount: comments.TotalCount,
	}, nil
}

func (d *productRPC) GetProductLikes(ctx context.Context, in *pb.GetWithID) (*pb.ListWishlistResponse, error) {
	likes, err := d.productUsecase.GetProductLikes(ctx, &entity.GetWithID{ID: in.Id})
	if err != nil {
		d.logger.Error("productUsecase.GetProductLikes", zap.Error(err))
		return &pb.ListWishlistResponse{}, grpc.Error(ctx, err)
	}

	pbLikes := &pb.ListWishlistResponse{}
	for _, like := range likes.Likes {
		pbLikes.Likes = append(pbLikes.Likes, &pb.Like{
			Id:        like.Id,
			ProductId: like.ProductID,
			UserId:    like.UserID,
			CreatedAt: like.CreatedAt.String(),
			UpdatedAt: like.UpdatedAt.String(),
		})
	}
	return &pb.ListWishlistResponse{
		Likes:      pbLikes.Likes,
		TotalCount: likes.TotalCount,
	}, nil
}

func (d *productRPC) GetProductStars(ctx context.Context, in *pb.GetWithID) (*pb.ListStarsResponse, error) {
	stars, err := d.productUsecase.GetProductStars(ctx, &entity.GetWithID{ID: in.Id})
	if err != nil {
		d.logger.Error("productUsecase.GetProductStars", zap.Error(err))
		return &pb.ListStarsResponse{}, grpc.Error(ctx, err)
	}

	pbStars := &pb.ListStarsResponse{}
	for _, star := range stars.Stars {
		pbStars.Stars = append(pbStars.Stars, &pb.Star{
			Id:        star.Id,
			ProductId: star.ProductID,
			UserId:    star.UserID,
			Star:      star.Stars,
			CreatedAt: star.CreatedAt.String(),
			UpdatedAt: star.UpdatedAt.String(),
		})
	}
	return &pb.ListStarsResponse{
		Stars:      pbStars.Stars,
		TotalCount: stars.TotalCount,
	}, nil
}
