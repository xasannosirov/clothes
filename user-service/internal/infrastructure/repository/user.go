package repository

import (
	"context"
	"user-service/internal/entity"
)

type Users interface {
	Create(ctx context.Context, kyc *entity.User) (*entity.User, error)
	Update(ctx context.Context, kyc *entity.User) (*entity.User, error)
	Delete(ctx context.Context, guid string) error
	GetDelete(ctx context.Context, params map[string]string) (*entity.User, error)
	Get(ctx context.Context, params map[string]string) (*entity.User, error)
	List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.User, error)
	UniqueEmail(ctx context.Context, request *entity.IsUnique) (*entity.Response, error)
	UpdateRefresh(ctx context.Context, request *entity.UpdateRefresh) (*entity.Response, error)
	UpdatePassword(ctx context.Context, request *entity.UpdatePassword) (*entity.Response, error)
	Total(ctx context.Context, role string) uint64
}
