package services

import (
	"context"
	pb "product-service/genproto/product_service"
	grpc "product-service/internal/delivery"
	"product-service/internal/entity"
	"product-service/internal/usecase"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type productRPC struct {
	logger         *zap.Logger
	productUsecase usecase.Product
}

func NewRPC(logger *zap.Logger, productUsecase usecase.Product) pb.ProductServiceServer {
	return &productRPC{
		logger:         logger,
		productUsecase: productUsecase,
	}
}

func (d *productRPC) CreateProduct(ctx context.Context, in *pb.Product) (*pb.GetWithID, error) {
	respProduct, err := d.productUsecase.CreateProduct(ctx, &entity.Product{
		Name:           in.Name,
		Description:    in.Description,
		Category:       in.Category,
		MadeIn:         in.MadeIn,
		Color:          in.Color,
		Count:          in.Count,
		Cost:           in.Cost,
		Discount:       in.Discount,
		AgeMin:         in.AgeMin,
		AgeMax:         in.AgeMax,
		TemperatureMin: in.TemperatureMin,
		TemperatureMax: in.TemperatureMax,
		Size:           int64(in.Size()),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		d.logger.Error("productUseCase.CreateProduct", zap.Error(err))
		return &pb.GetWithID{}, grpc.Error(ctx, err)
	}
	return &pb.GetWithID{Id: respProduct.Id}, nil
}

func (d *productRPC) UpdateProduct(ctx context.Context, in *pb.Product) (*pb.Product, error) {
	err := d.productUsecase.UpdateProduct(ctx, &entity.Product{
		Id:             in.Id,
		Name:           in.Name,
		Description:    in.Description,
		Category:       in.Category,
		MadeIn:         in.MadeIn,
		Color:          in.Color,
		Count:          in.Count,
		Cost:           in.Cost,
		Discount:       in.Discount,
		AgeMin:         in.AgeMin,
		AgeMax:         in.AgeMax,
		TemperatureMin: in.TemperatureMin,
		TemperatureMax: in.TemperatureMax,
		ForGender:      in.ForGender,
		Size:           int64(in.Size()),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		d.logger.Error("productUseCase.UpdateProduct", zap.Error(err))
		return &pb.Product{}, grpc.Error(ctx, err)
	}

	return in, nil
}

func (d *productRPC) GetProductByID(ctx context.Context, in *pb.GetWithID) (*pb.Product, error) {
	product, err := d.productUsecase.GetProductByID(ctx, map[string]string{"id": in.Id})
	if err != nil {
		d.logger.Error("productUseCase.GetProductByID", zap.Error(err))
		return &pb.Product{}, grpc.Error(ctx, err)
	}

	return &pb.Product{
		Id:             product.Id,
		Name:           product.Name,
		Description:    product.Description,
		Category:       product.Category,
		MadeIn:         product.MadeIn,
		Color:          product.Color,
		Count:          product.Count,
		Cost:           product.Cost,
		Discount:       product.Discount,
		AgeMin:         product.AgeMin,
		AgeMax:         product.AgeMax,
		TemperatureMin: product.TemperatureMin,
		TemperatureMax: product.TemperatureMax,
		ForGender:      product.ForGender,
		Size_:          product.Size,
		CreatedAt:      product.CreatedAt.String(),
		UpdatedAt:      product.UpdatedAt.String(),
	}, nil
}

func (d *productRPC) DeleteProduct(ctx context.Context, in *pb.GetWithID) (*pb.DeleteResponse, error) {
	err := d.productUsecase.DeleteProduct(ctx, in.Id)
	if err != nil {
		d.logger.Error("productUseCase.DeleteProduct", zap.Error(err))
		return &pb.DeleteResponse{Status: false}, grpc.Error(ctx, err)
	}

	return &pb.DeleteResponse{Status: true}, nil
}

func (d *productRPC) GetAllProducts(ctx context.Context, in *pb.ListRequest) (*pb.ListProductResponse, error) {
	filter := &entity.ListRequest{
		Limit: in.Limit,
		Page:  in.Page,
	}

	products, err := d.productUsecase.GetAllProducts(ctx, filter)
	if err != nil {
		d.logger.Error("productUseCase.List", zap.Error(err))
		return &pb.ListProductResponse{}, grpc.Error(ctx, err)
	}

	var pbProducts []*pb.Product
	for _, product := range products {
		pbProducts = append(pbProducts, &pb.Product{
			Id:             product.Id,
			Name:           product.Name,
			Description:    product.Description,
			Category:       product.Category,
			MadeIn:         product.MadeIn,
			Color:          product.Color,
			Count:          int64(product.Count),
			Cost:           float32(product.Cost),
			Discount:       float32(product.Discount),
			AgeMin:         product.AgeMin,
			AgeMax:         product.AgeMax,
			TemperatureMin: product.TemperatureMin,
			TemperatureMax: product.TemperatureMax,
			ForGender:      product.ForGender,
			Size_:          product.Size,
			CreatedAt:      product.CreatedAt.String(),
			UpdatedAt:      product.UpdatedAt.String(),
		})
	}

	return &pb.ListProductResponse{Products: pbProducts}, nil
}

func (d *productRPC) CreateOrder(ctx context.Context, in *pb.Order) (*pb.GetWithID, error) {
	id := uuid.New().String()
	_, err := d.productUsecase.CreateOrder(ctx, &entity.Order{
		Id:        id,
		ProductID: in.ProductId,
		UserID:    in.UserId,
		Status:    in.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		d.logger.Error("productUsecase.CreateOrder", zap.Error(err))
		return &pb.GetWithID{}, grpc.Error(ctx, err)
	}

	in.Id = id
	return &pb.GetWithID{Id: in.Id}, nil
}

func (d *productRPC) CancelOrder(ctx context.Context, in *pb.GetWithID) (*pb.DeleteResponse, error) {
	err := d.productUsecase.CancelOrder(ctx, in.Id)
	if err != nil {
		d.logger.Error("productUsecase.CancelOrder", zap.Error(err))
		return &pb.DeleteResponse{Status: false}, grpc.Error(ctx, err)
	}

	return &pb.DeleteResponse{Status: true}, nil
}

func (d *productRPC) GetOrderByID(ctx context.Context, in *pb.GetWithID) (*pb.Order, error) {
	order, err := d.productUsecase.GetOrderByID(ctx, map[string]string{"id": in.Id})
	if err != nil {
		d.logger.Error("productUsecase.GetOrderByID", zap.Error(err))
		return &pb.Order{}, grpc.Error(ctx, err)
	}
	return &pb.Order{
		Id:        order.Id,
		ProductId: order.ProductID,
		UserId:    order.UserID,
		Status:    order.Status,
		CreatedAt: order.CreatedAt.String(),
		UpdatedAt: order.UpdatedAt.String(),
	}, nil
}

func (d *productRPC) GetAllOrders(ctx context.Context, in *pb.ListRequest) (*pb.ListOrderResponse, error) {
	filter := &entity.ListRequest{
		Limit: in.Limit,
		Page:  in.Page,
	}

	orders, err := d.productUsecase.GetAllOrders(ctx, filter)
	if err != nil {
		d.logger.Error("productUsecase.GetAllOrders", zap.Error(err))
		return &pb.ListOrderResponse{}, grpc.Error(ctx, err)
	}

	var pbOrders []*pb.Order
	for _, order := range orders {
		pbOrders = append(pbOrders, &pb.Order{
			Id:        order.Id,
			ProductId: order.ProductID,
			UserId:    order.UserID,
			Status:    order.Status,
			CreatedAt: order.CreatedAt.String(),
			UpdatedAt: order.UpdatedAt.String(),
		})
	}
	return &pb.ListOrderResponse{Orders: pbOrders}, nil
}

func (d *productRPC) SearchProduct(context.Context, *pb.Filter) (*pb.ListProductResponse, error) {
	return nil, nil

}
func (d *productRPC) Recommendation(context.Context, *pb.Recom) (*pb.ListProductResponse, error) {
	return nil, nil

}
func (d *productRPC) GetSavedProductsByUserID(context.Context, *pb.GetWithUserID) (*pb.ListProductResponse, error) {
	return nil, nil

}
func (d *productRPC) GetWishlistByUserID(context.Context, *pb.GetWithUserID) (*pb.ListProductResponse, error) {
	return nil, nil

}
func (d *productRPC) GetOrderedProductsByUserID(context.Context, *pb.GetWithUserID) (*pb.ListProductResponse, error) {
	return nil, nil

}
func (d *productRPC) LikeProduct(ctx context.Context, req *pb.Like) (*pb.MoveResponse, error) {
	status, err := d.productUsecase.LikeProduct(ctx, &entity.LikeProduct{
		Product_id: req.ProductId,
		User_id: req.UserId,
	})

	if err != nil {
		d.logger.Error("productUseCase.CreateProduct", zap.Error(err))
		return nil, grpc.Error(ctx, err)
	}

	return &pb.MoveResponse{
		Status: status,
	},nil
}

func (d *productRPC) SaveProduct(ctx context.Context, req *pb.Save) (*pb.MoveResponse, error) {
	status, err := d.productUsecase.SaveProduct(ctx, &entity.SaveProduct{
		Product_id: req.ProductId,
		User_id: req.UserId,
	})

	if err != nil {
		d.logger.Error("productUseCase.CreateProduct", zap.Error(err))
		return nil, grpc.Error(ctx, err)
	}

	return &pb.MoveResponse{
		Status: status,
	},nil
}

func (d *productRPC) StarProduct(context.Context, *pb.Star) (*pb.MoveResponse, error) {
	return nil, nil

}
func (d *productRPC) CommentToProduct(ctx context.Context, req *pb.Comment) (*pb.MoveResponse, error) {
	status, err := d.productUsecase.CommentToProduct(ctx, &entity.CommentToProduct{
		UserId: req.UserId,
		Product_Id: req.ProductId,
		Comment: req.Comment,
	})
	if err != nil{
		return nil, err
	}
	return &pb.MoveResponse{
		Status: status,
	}, nil
}
func (d *productRPC) GetDisableProducts(context.Context, *pb.ListRequest) (*pb.ListOrderResponse, error) {
	return nil, nil

}
func (d *productRPC) GetProductOrders(context.Context, *pb.GetWithID) (*pb.ListOrderResponse, error) {
	return nil, nil

}
func (d *productRPC) GetProductComments(context.Context, *pb.GetWithID) (*pb.ListCommentResponse, error) {
	return nil, nil

}
func (d *productRPC) GetProductLikes(context.Context, *pb.GetWithID) (*pb.ListWishlistResponse, error) {
	return nil, nil

}
func (d *productRPC) GetProductStars(context.Context, *pb.GetWithID) (*pb.ListStarsResponse, error) {
	return nil, nil

}
func (d *productRPC) GetAllComments(context.Context, *pb.ListRequest) (*pb.ListCommentResponse, error) {
	return nil, nil

}
func (d *productRPC) GetAllStars(context.Context, *pb.ListRequest) (*pb.ListStarsResponse, error) {
	return nil, nil

}
