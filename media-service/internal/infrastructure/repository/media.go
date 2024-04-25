package repository

import(
	"context"
	"clothes-store/media-service/internal/entity"
)

type MediaStorageI interface{
	CreateMedia(ctx context.Context, media *entity.Media)(*entity.Media, error)
	GetMediaWithProductId(ctx context.Context,  params map[string]string)([]*entity.Media, error)
	DeleteMedia(ctx context.Context, params map[string]any)error
}

