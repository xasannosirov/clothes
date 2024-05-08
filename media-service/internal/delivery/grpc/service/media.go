package service

import (
	"context"
	mediaproto "media-service/genproto/media_service"
	"media-service/internal/entity"
	grpc_service_clients "media-service/internal/infrastructure/grpc_service_client"
	"media-service/internal/pkg/otlp"
	"media-service/internal/usecase"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"go.uber.org/zap"
)

type mediaRPC struct {
	logger *zap.Logger
	media  usecase.Media
	client grpc_service_clients.ServiceClients
}

func NewRPC(logger *zap.Logger, media usecase.Media, clients grpc_service_clients.ServiceClients) mediaRPC {
	return mediaRPC{
		logger: logger,
		media:  media,
		client: clients,
	}
}

// Create method creates media
func (m mediaRPC) Create(ctx context.Context, req *mediaproto.Media) (*mediaproto.MediaWithID, error) {
	ctx, span := otlp.Start(ctx, "media_grpc-delivery", "CreateMedia")
	span.SetAttributes(
		attribute.Key("guid").String(req.ProductId),
	)
	defer span.End()

	media := &entity.Media{
		Id:        req.Id,
		ProductID: req.ProductId,
		ImageUrl:  req.ImageUrl,
		FileName:  req.FileName,
	}

	respMedia, err := m.media.CreateMedia(ctx, media)
	if err != nil {
		m.logger.Error("media.Create", zap.Error(err))
		return nil, err
	}

	return &mediaproto.MediaWithID{
		Id: respMedia.Id,
	}, nil
}

// Get method returns media
func (m mediaRPC) Get(ctx context.Context, req *mediaproto.MediaWithProductID) (*mediaproto.ProductImages, error) {
	ctx, span := otlp.Start(ctx, "media_grpc-delivery", "GetMedia")
	span.SetAttributes(
		attribute.Key("guid").String(req.ProductId),
	)
	defer span.End()

	params := make(map[string]string)
	params["product_id"] = req.ProductId

	response, err := m.media.GetMediaWithProductId(ctx, params)
	if err != nil {
		m.logger.Error("getMediaWithProductId", zap.Error(err))
	}

	respMedia := []*mediaproto.Media{}
	for _, media := range response {
		resp := &mediaproto.Media{
			Id:        media.Id,
			ProductId: media.ProductID,
			ImageUrl:  media.ImageUrl,
			FileName:  media.FileName,
			CreatedAt: media.CreatedAt.Format(time.RFC3339),
			UpdatedAt: media.UpdatedAt.Format(time.RFC3339),
		}

		respMedia = append(respMedia, resp)
	}

	return &mediaproto.ProductImages{
		Images: respMedia,
	}, nil

}

// Delete method delete media
func (m mediaRPC) Delete(ctx context.Context, req *mediaproto.MediaWithProductID) (*mediaproto.DeleteMediaResponse, error) {
	ctx, span := otlp.Start(ctx, "media_grpc-delivery", "DeleteMedia")
	span.SetAttributes(
		attribute.Key("guid").String(req.ProductId),
	)
	defer span.End()

	params := make(map[string]any)
	params["product_id"] = req.ProductId

	err := m.media.DeleteMedia(ctx, params)
	if err != nil {
		m.logger.Error("error", zap.Error(err))
		return nil, err
	}

	return &mediaproto.DeleteMediaResponse{
		Status: true,
	}, nil
}
