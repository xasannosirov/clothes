package repository

import (
	"context"
	"user-service/internal/entity"
)

type Users interface {
	Create(ctx context.Context, kyc *entity.User) error
	Update(ctx context.Context, kyc *entity.User) error
	Delete(ctx context.Context, guid string) error
	Get(ctx context.Context, params map[string]string) (*entity.User, error)
	List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.User, error)
}
