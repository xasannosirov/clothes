package postgres

import (
	"context"
	"fmt"
	"media-service/internal/entity"
	"media-service/internal/pkg/postgres"

	"github.com/Masterminds/squirrel"
)

const (
	mediaServiceTableName   = "media"
	mediaServiceServiceName = "mediaService"
)

type mediaRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewMediaRepo(db *postgres.PostgresDB) *mediaRepo {
	return &mediaRepo{
		tableName: mediaServiceTableName,
		db:        db,
	}
}

func (m *mediaRepo) mediaSelectQueryPrefix() squirrel.SelectBuilder {
	return m.db.Sq.Builder.Select(
		"id",
		"product_id",
		"image_url",
		"file_name",
		"created_at",
		"updated_at",
	).From(m.tableName)
}

// CreateMedia creates media consist of image_url, product_id and id
func (m mediaRepo) CreateMedia(ctx context.Context, media *entity.Media) (*entity.Media, error) {
	data := map[string]any{
		"id":         media.Id,
		"product_id": media.ProductID,
		"image_url":  media.ImageUrl,
		"file_name":  media.FileName,
		"created_at": media.CreatedAt,
		"updated_at": media.UpdatedAt,
	}

	query, args, err := m.db.Sq.Builder.Insert(m.tableName).SetMap(data).ToSql()
	if err != nil {
		return nil, m.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", m.tableName, "create"))
	}

	_, err = m.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, m.db.Error(err)
	}

	return media, nil
}

// GetMediaWithProductId returns list media by product id
func (m mediaRepo) GetMediaWithProductId(ctx context.Context, filter map[string]string) ([]*entity.Media, error) {
	var (
		ListMedia []*entity.Media
	)
	queryBuilder := m.mediaSelectQueryPrefix()
	queryBuilder = queryBuilder.Where(m.db.Sq.Equal("product_id", filter["product_id"])).Where("deleted_at IS NULL").OrderBy("created_at")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, m.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", m.tableName, "list"))
	}

	rows, err := m.db.Query(ctx, query, args...)
	if err != nil {
		return nil, m.db.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		var media entity.Media
		if err = rows.Scan(
			&media.Id,
			&media.ProductID,
			&media.ImageUrl,
			&media.FileName,
			&media.CreatedAt,
			&media.UpdatedAt,
		); err != nil {
			return nil, m.db.Error(err)
		}

		ListMedia = append(ListMedia, &media)
	}

	return ListMedia, nil
}

// DeleteMedia delete all media by product id
func (m mediaRepo) DeleteMedia(ctx context.Context, params map[string]any) error {
	query, args, err := m.db.Sq.Builder.
		Update(m.tableName).
		SetMap(params).
		Where(m.db.Sq.Equal("product_id", params["product_id"])).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return m.db.ErrSQLBuild(err, m.tableName+" delete")
	}

	commandTag, err := m.db.Exec(ctx, query, args...)
	if err != nil {
		return m.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return m.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}
