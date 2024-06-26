package usecase

import (
	"context"
	"media-service/internal/entity"
	"media-service/internal/infrastructure/repository"
	"media-service/internal/pkg/otlp"
	"time"
)

type Media interface {
	CreateMedia(ctx context.Context, media *entity.Media) (*entity.Media, error)
	GetMediaWithProductId(ctx context.Context, params map[string]string) ([]*entity.Media, error)
	DeleteMedia(ctx context.Context, params map[string]any) error
}

type mediaService struct {
	BaseUseCase
	repo       repository.MediaStorageI
	ctxTimeout time.Duration
}

func NewMediaService(ctxTimout time.Duration, repo repository.MediaStorageI) mediaService {
	return mediaService{
		ctxTimeout: ctxTimout,
		repo:       repo,
	}
}

func (m mediaService) CreateMedia(ctx context.Context, media *entity.Media) (*entity.Media, error) {
	ctx, cancel := context.WithTimeout(ctx, m.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, "media_grpc-usercase", "CreateMedia")
	defer span.End()

	m.beforeRequest(nil, &media.CreatedAt, &media.UpdatedAt)
	return m.repo.CreateMedia(ctx, media)
}

func (m mediaService) GetMediaWithProductId(ctx context.Context, params map[string]string) ([]*entity.Media, error) {
	ctx, cancel := context.WithTimeout(ctx, m.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, "product_grpc-usercase", "GetProduct")
	defer span.End()

	return m.repo.GetMediaWithProductId(ctx, params)
}

func (m mediaService) DeleteMedia(ctx context.Context, params map[string]any) error {
	ctx, cancel := context.WithTimeout(ctx, m.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, "product_grpc-usercase", "DeleteProduct")
	defer span.End()

	params["deleted_at"] = time.Now().UTC()

	return m.repo.DeleteMedia(ctx, params)
}
