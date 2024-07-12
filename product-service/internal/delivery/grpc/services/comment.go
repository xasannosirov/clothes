package services

import (
	"context"
	pb "product-service/genproto/product_service"
	"product-service/internal/entity"
)

func (p *productRPC) CreateComment(ctx context.Context, comment *pb.Comment) (*pb.Comment, error) {

	resComment, err := p.productUsecase.CreateComment(ctx, &entity.Comment{
		Id:    comment.Id,
		OwnerId: comment.OwnerId,
		ProductId: comment.ProductId,
		Message:   comment.Message,
	})
	if err != nil {
		return nil, err
	}
	return &pb.Comment{
		Id: resComment.Id,
		OwnerId:    resComment.OwnerId,
		ProductId: resComment.ProductId,
		Message:   resComment.Message,
	}, nil
}
func (p *productRPC) UpdateComment(ctx context.Context, comment *pb.CommentUpdateRequst) (*pb.Comment, error) {
	resComment, err := p.productUsecase.UpdateComment(ctx, &entity.CommentUpdateRequest{
		Id:    comment.Id,
		Message:   comment.Message,
	})
	if err != nil {
		return nil, err
	}

	return &pb.Comment{
		Id: resComment.Id,
		OwnerId:    resComment.OwnerId,
		ProductId: resComment.ProductId,
		Message:   resComment.Message,
	}, nil
}
func (p *productRPC) GetComment(ctx context.Context, req *pb.CommentGetRequst) (*pb.Comment, error) {
		
	resComments, err := p.productUsecase.GetComment(ctx, &entity.GetRequest{
		Filter: req.Filter,
	})
	if err != nil {
		return nil, err
	}
	

	return &pb.Comment{
		Id: resComments.Id,
		OwnerId: resComments.OwnerId,
		ProductId: resComments.ProductId,
		Message: resComments.Message,
	}, nil
}
func (p *productRPC) ListComment(ctx context.Context, req *pb.CommentListRequest) (*pb.CommentListResponse, error) {
	resComments, err := p.productUsecase.ListComment(ctx, &entity.ListRequest{
		Page:   req.Page,
		Limit:  req.Limit,
		Filter: req.Filter,
	})
	if err != nil {
		return nil, err
	}
	var comments pb.CommentListResponse
	for _, resComment := range resComments.Comment {
		comment := pb.Comment{
			Id: resComment.Id,
			OwnerId:    resComment.OwnerId,
			ProductId: resComment.ProductId,
			Message:   resComment.Message,
		}
		comments.Comments = append(comments.Comments, &comment)
	}
	comments.TotalCount = int64(resComments.TotalCount)

	return &comments, nil
}

func (p *productRPC) DeleteComment(ctx context.Context, req *pb.CommentDeleteRequest)(*pb.MoveResponse, error){
	err := p.productUsecase.DeleteComment(ctx, &entity.DeleteRequest{
		Id: req.Id,
	})
	if err != nil{
		return &pb.MoveResponse{
			Status: false,
		}, err
	}

	return &pb.MoveResponse{
		Status: true,
	}, nil
}