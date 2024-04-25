package postgres

import (
	"clothes-store/media-service/internal/entity"
	"clothes-store/media-service/internal/pkg/postgres"
	"context"
	"fmt"

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
		"created_at",
		"updated_at",
	).From(m.tableName)
}

func (m mediaRepo) CreateMedia(ctx context.Context, media *entity.Media) (*entity.Media, error) {
	data := map[string]any{
		"id":         media.Id,
		"product_id": media.Product_Id,
		"image_url":  media.Image_Url,
		"created_at": media.Created_at,
		"updated_at": media.Updated_at,
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

func (m mediaRepo) GetMediaWithProductId(ctx context.Context, filter map[string]string) ([]*entity.Media, error) {
	var (
		ListMedia []*entity.Media
	)
	queryBuilder := m.mediaSelectQueryPrefix()
	queryBuilder = queryBuilder.Where(m.db.Sq.Equal("product_id", filter["product_id"])).Where("deleted_at IS NULL")

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
			&media.Product_Id,
			&media.Image_Url,
			&media.Created_at,
			&media.Updated_at,
		); err != nil {
			return nil, m.db.Error(err)
		}
		ListMedia = append(ListMedia, &media)
	}
	return ListMedia, nil
}

func (m mediaRepo) DeleteMedia(ctx context.Context, params map[string]any) error {
	sqlStr, args, err := m.db.Sq.Builder.
		Update(m.tableName).
		SetMap(params).
		Where(m.db.Sq.Equal("product_id", params["product_id"])).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return m.db.ErrSQLBuild(err, m.tableName+" delete")
	}

	commandTag, err := m.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return m.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return m.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}
