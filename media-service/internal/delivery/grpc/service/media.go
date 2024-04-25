package service

import (
	pm "clothes-store/media-service/genproto/media_service"
	"clothes-store/media-service/internal/entity"
	"clothes-store/media-service/internal/usecase"
	"context"

	"go.uber.org/zap"
)

type mediaRPC struct {
	logger *zap.Logger
	media  usecase.Media
}

func NewRPC(logger *zap.Logger, media usecase.Media) mediaRPC {
	return mediaRPC{
		logger: logger,
		media:  media,
	}
}

func (m mediaRPC) Create(ctx context.Context, req *pm.Media) (*pm.MediaWithID, error) {
	media := &entity.Media{
		Id:         req.Id,
		Product_Id: req.ProductId,
		Image_Url:  req.ImageUrl,
	}

	respMedia, err := m.media.CreateMedia(ctx, media)
	if err != nil {
		m.logger.Error("user.Create", zap.Error(err))
		return nil, err
	}
	
	return &pm.MediaWithID{
		Id: respMedia.Id,
	}, nil
}

func (m mediaRPC) Get(ctx context.Context, req *pm.MediaWithProductID) (*pm.ProductImages, error) {
	params := make(map[string]string)
	params["product_id"] = req.ProductId

	response, err := m.media.GetMediaWithProductId(ctx, params)
	if err != nil {
		m.logger.Error("getMediaWithproductId", zap.Error(err))
	}

	respMedia := []*pm.Media{}
	for _, media := range response {
		resp := &pm.Media{
			Id:        media.Id,
			ProductId: media.Product_Id,
			ImageUrl:  media.Image_Url,
			CreatedAt: media.Created_at.String(),
			UpdatedAt: media.Updated_at.String(),
		}

		respMedia = append(respMedia, resp)
	}

	return &pm.ProductImages{
		Images: respMedia,
	}, nil

}

func (m mediaRPC) Delete(ctx context.Context, req *pm.MediaWithProductID) (*pm.DeleteMediaResponse, error) {
	params := make(map[string]any)
	params["product_id"] = req.ProductId

	err := m.media.DeleteMedia(ctx, params)
	if err != nil {
		m.logger.Error("error", zap.Error(err))
	}

	return &pm.DeleteMediaResponse{
		Status: true,
	}, nil
}
