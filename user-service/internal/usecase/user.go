package usecase

import (
	"context"
	"time"
	"user-service/internal/entity"
	"user-service/internal/infrastructure/repository"
	"user-service/internal/pkg/otlp"
)

type User interface {
	Create(ctx context.Context, article *entity.User) (*entity.User, error)
	Update(ctx context.Context, article *entity.User) (*entity.User, error)
	Delete(ctx context.Context, guid string) error
	Get(ctx context.Context, params map[string]string) (*entity.User, error)
	List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.User, error)
	UniqueEmail(ctx context.Context, request *entity.IsUnique) (*entity.Response, error)
	UpdateRefresh(ctx context.Context, request *entity.UpdateRefresh) (*entity.Response, error)
	UpdatePassword(ctx context.Context, request *entity.UpdatePassword) (*entity.Response, error)
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

func (u userService) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, "user_grpc-usercase", "CreateUser")
	defer span.End()

	u.beforeRequest(nil, &user.CreatedAt, &user.UpdatedAt)

	return u.repo.Create(ctx, user)
}

func (u userService) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, "user_grpc-usercase", "UpdateUser")
	defer span.End()

	u.beforeRequest(nil, nil, &user.UpdatedAt)

	return u.repo.Update(ctx, user)
}

func (u userService) Delete(ctx context.Context, guid string) error {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, "user_grpc-usercase", "DeleteUser")
	defer span.End()

	return u.repo.Delete(ctx, guid)
}

func (u userService) Get(ctx context.Context, params map[string]string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, "user_grpc-usercase", "GetUser")
	defer span.End()

	return u.repo.Get(ctx, params)
}

func (u userService) List(ctx context.Context, limit, offset uint64, filter map[string]string) ([]*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, "user_grpc-usercase", "ListUsers")
	defer span.End()

	return u.repo.List(ctx, limit, offset, filter)
}

func (u userService) UniqueEmail(ctx context.Context, request *entity.IsUnique) (*entity.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, "user_grpc-usercase", "UniqueEmail")
	defer span.End()

	return u.repo.UniqueEmail(ctx, request)
}

func (u userService) UpdateRefresh(ctx context.Context, request *entity.UpdateRefresh) (*entity.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, "user_grpc-usercase", "UpdateRefresh")
	defer span.End()

	return u.repo.UpdateRefresh(ctx, request)
}

func (u userService) UpdatePassword(ctx context.Context, request *entity.UpdatePassword) (*entity.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, u.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, "user_grpc-usercase", "UpdatePassword")
	defer span.End()

	return u.repo.UpdatePassword(ctx, request)
}
