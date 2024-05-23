package postgresql

import (
	"context"
	"errors"
	"product-service/internal/entity"
	"time"

	"github.com/Masterminds/squirrel"
)

func (u *productRepo) CreateCategory(ctx context.Context, req *entity.Category) (*entity.Category, error) {
	data := map[string]any{
		"id":   req.ID,
		"name": req.Name,
	}

	query, args, err := u.db.Sq.Builder.Insert(u.categoryTable).SetMap(data).ToSql()
	if err != nil {
		return nil, err
	}

	_, err = u.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (u *productRepo) UpdateCategory(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	data := map[string]any{
		"name": category.Name,
	}

	sqlStr, args, err := u.db.Sq.Builder.
		Update(u.categoryTable).
		SetMap(data).
		Where(squirrel.Eq{"id": category.ID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	commandTag, err := u.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	if commandTag.RowsAffected() == 0 {
		return nil, errors.New("no sql rows")
	}

	return category, nil
}

func (u *productRepo) DeleteCategory(ctx context.Context, id string) error {
	caluses := map[string]interface{}{
		"deleted_at": time.Now().Format(time.RFC3339),
	}

	sqlStr, args, err := u.db.Sq.Builder.
		Update(u.categoryTable).
		SetMap(caluses).
		Where(u.db.Sq.Equal("id", id)).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return err
	}

	commandTag, err := u.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return errors.New("no sql rows")
	}

	return nil
}

func (u *productRepo) GetCategory(ctx context.Context, id string) (*entity.Category, error) {
	queryBuilder := u.categorySelectQueryPrefix()

	queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	queryBuilder = queryBuilder.Where(u.db.Sq.Equal("id", id))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	row := u.db.QueryRow(ctx, query, args...)

	var category entity.Category
	if err = row.Scan(
		&category.ID,
		&category.Name,
	); err != nil {
		return nil, err
	}

	return &category, nil
}

func (u *productRepo) ListCategories(ctx context.Context, req *entity.ListRequest) (*entity.LiestCategory, error) {
	categories := &entity.LiestCategory{}
	queryBuilder := u.categorySelectQueryPrefix()

	offset := (req.Page - 1) * req.Limit

	queryBuilder = queryBuilder.Limit(uint64(req.Limit)).Offset(uint64(offset))

	queryBuilder = queryBuilder.Where("deleted_at IS NULL").OrderBy("created_at")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category entity.Category
		if err = rows.Scan(
			&category.ID,
			&category.Name,
		); err != nil {
			return nil, err
		}

		categories.Categories = append(categories.Categories, &category)
	}

	var count uint64
	totol := `SELECT COUNT(*) FROM category WHERE deleted_at IS NULL`
	if err := u.db.QueryRow(ctx, totol).Scan(&count); err != nil {
		categories.TotalCount = 0
	}
	categories.TotalCount = count

	return categories, nil
}

func (u *productRepo) SearchCategory(ctx context.Context, req *entity.SearchRequest) (*entity.ListProduct, error) {
	products := entity.ListProduct{}
	offset := req.Limit * (req.Page - 1)
	queryBuilder := u.categorySelectQueryPrefix()

	queryBuilder = queryBuilder.Where(squirrel.Eq{"name": req.Params["name"]})
	queryBuilder = queryBuilder.Where("deleted_at IS NULL")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var category entity.Category
	err = u.db.QueryRow(ctx, query, args...).Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, err
	}

	builder := u.productsSelectQueryPrefix()
	builder = builder.Where("deleted_at IS NULL")
	builder = builder.Where(squirrel.Eq{"category_id": category.ID})
	builder = builder.Limit(req.Limit)
	builder = builder.Offset(offset)

	queryP, argsP, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := u.db.Query(ctx, queryP, argsP...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Category,
			&product.MadeIn,
			&product.Color,
			&product.Count,
			&product.Cost,
			&product.Discount,
			&product.AgeMin,
			&product.AgeMax,
			&product.ForGender,
			&product.Size,
		)
		if err != nil {
			return nil, err
		}
		products.Products = append(products.Products, &product)
	}

	var count uint64
	totol := `SELECT COUNT(category_id) FROM products WHERE deleted_at IS NULL AND category_id = $1`
	if err := u.db.QueryRow(ctx, totol, category.ID).Scan(&count); err != nil {
		return nil, err
	}
	products.TotalCount = count

	return &products, nil
}

func (u *productRepo) UniqueCategory(ctx context.Context, in *entity.Params) (*entity.MoveResponse, error) {
	query := `SELECT COUNT(category_name) FROM category WHERE deleted_at IS NULL AND category_name = $1`

	var count uint64
	err := u.db.QueryRow(ctx, query, in.Filter["category_name"]).Scan(&count)
	if err != nil {
		return nil, err
	}

	if count == 0 {
		return &entity.MoveResponse{Status: false}, nil
	}

	return &entity.MoveResponse{Status: true}, nil
}
