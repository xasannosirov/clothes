package app_version

// import (
// 	"context"
// 	"time"

// 	"api-agteway/internal/entity"
// 	"api-agteway/internal/infrastructure/repository/postgresql/repo"
// 	// "evrone_api-agteway/internal/pkg/otlp"
// )

// type appVersionService struct {
// 	ctxTimeout time.Duration
// 	repo       repo.AppVersionRepo
// }

// func NewAppVersionService(ctxTimeout time.Duration, repo repo.AppVersionRepo) AppVersion {
// 	return &appVersionService{
// 		ctxTimeout: ctxTimeout,
// 		repo:       repo,
// 	}
// }

// func (r *appVersionService) beforeCreate(m *entity.AppVersion) {
// 	m.CreatedAt = time.Now().UTC()
// 	m.UpdatedAt = time.Now().UTC()
// }

// func (r *appVersionService) beforeUpdate(m *entity.AppVersion) {
// 	m.UpdatedAt = time.Now().UTC()
// }

// func (r *appVersionService) Get(ctx context.Context) (*entity.AppVersion, error) {
// 	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
// 	defer cancel()

// 	return r.repo.Get(ctx)
// }

// func (r *appVersionService) Create(ctx context.Context, m *entity.AppVersion) error {
// 	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
// 	defer cancel()

// 	r.beforeCreate(m)
// 	return r.repo.Create(ctx, m)
// }

// func (r *appVersionService) Update(ctx context.Context, m *entity.AppVersion) error {
// 	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
// 	defer cancel()

//     r.beforeUpdate(m)
// 	return r.repo.Update(ctx, m)
// }
