package usecase

import (
	"context"
	"time"
	"user-service/internal/entity"
	"user-service/internal/infrastructure/repository"
)

type User interface {
	Create(ctx context.Context, article *entity.User) (string, error)
	Update(ctx context.Context, article *entity.User) error
	Delete(ctx context.Context, guid string) error
	Get(ctx context.Context, params map[string]string) (*entity.User, error)
	List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.User, error)
}

type userService struct {
	BaseUseCase
	repo       repository.Users
	ctxTimeout time.Duration
}

func NewUserService(ctxTimeout time.Duration, repo repository.Users) User {
	return userService{
		ctxTimeout: ctxTimeout,
		repo:       repo,
	}
}

func (u userService) Create(ctx context.Context, article *entity.User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	u.beforeRequest(&article.GUID, &article.CreatedAt, &article.UpdatedAt)

	return article.GUID, u.repo.Create(ctx, article)
}

func (u userService) Get(ctx context.Context, params map[string]string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.Get(ctx, params)
}

func (u userService) List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.List(ctx, limit, offset, filter)
}

func (u userService) Update(ctx context.Context, article *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	u.beforeRequest(nil, nil, &article.UpdatedAt)

	return u.repo.Update(ctx, article)
}

func (u userService) Delete(ctx context.Context, guid string) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	return u.repo.Delete(ctx, guid)
}
