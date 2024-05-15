package postgresql

import (
	"context"
	"fmt"
	"product-service/internal/entity"
	"product-service/internal/pkg/otlp"
	"time"

	"github.com/Masterminds/squirrel"
)

func (u *productRepo) basketsSelectQueryPrefix() squirrel.SelectBuilder {
	return u.db.Sq.Builder.Select(
		"id",
		"product_id",
		"user_id",
		"created_at",
		"updated_at",
	).From(u.basketTable)
}

func (u *productRepo) CreateBasket(ctx context.Context, req *entity.Basket) (*entity.Basket, error) {
	ctx, span := otlp.Start(ctx, "basket-grpc-reposiroty", "CreateBasket")
	defer span.End()

	data := map[string]any{
		"id":         req.Id,
		"user_id":    req.UserId,
		"product_id": req.ProductId,
		"created_at": req.Created_at,
		"updated_at": req.Updated_at,
	}

	query, args, err := u.db.Sq.Builder.Insert(u.basketTable).SetMap(data).ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.productTable, "createProduct"))
	}

	_, err = u.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, u.db.Error(err)
	}

	return req, nil
}
func (u *productRepo) GetBasket(ctx context.Context, params map[string]string) (*entity.Basket, error) {
	var (
		basket entity.Basket
	)

	queryBuilder := u.basketsSelectQueryPrefix()

	for key, value := range params {
		if key == "id" {
			queryBuilder = queryBuilder.Where(squirrel.Eq{key: value})
		}
	}

	queryBuilder = queryBuilder.Where("deleted_at IS NULL")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.productTable, "getProduct"))
	}

	if err = u.db.QueryRow(ctx, query, args...).Scan(
		&basket.Id,
		&basket.ProductId,
		&basket.UserId,
		&basket.Created_at,
		&basket.Updated_at,
	); err != nil {
		return nil, u.db.Error(err)
	}
	return &basket, nil
}

func (u *productRepo) GetBaskets(ctx context.Context, req *entity.ListBasketReq) (*entity.ListBasketRes, error) {
	products := &entity.ListBasketRes{}
	queryBuilder := u.basketsSelectQueryPrefix()

	offset := (req.Page - 1) * req.Limit

	queryBuilder = queryBuilder.Where("deleted_at IS NULL LIMIT $1 OFFSET $2").OrderBy("created_at")

	query, _, err := queryBuilder.ToSql()
	if err != nil {
		return nil, u.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", u.productTable, "getProducts"))
	}

	rows, err := u.db.Query(ctx, query, req.Limit, offset)
	if err != nil {
		return nil, u.db.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		var product entity.Basket
		if err = rows.Scan(
			&product.Id,
			&product.ProductId,
			&product.UserId,
			&product.Created_at,
			&product.Updated_at,
		); err != nil {
			return nil, err
		}
		products.Basket = append(products.Basket, &product)

	}
	var count int
	total := "SELECT COUNT(*) FROM products WHERE deleted_at IS NULL"
	if err := u.db.QueryRow(ctx, total).Scan(&count); err != nil {
		products.TotalCount = 0
	}
	products.TotalCount = count
	return products, nil
}

func (u *productRepo) DeleteBasket(ctx context.Context, id string) error {
	caluses := map[string]interface{}{
		"deleted_at": time.Now().Format(time.RFC3339),
	}

	sqlStr, args, err := u.db.Sq.Builder.
		Update(u.basketTable).
		SetMap(caluses).
		Where(u.db.Sq.Equal("id", id)).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return u.db.ErrSQLBuild(err, u.basketTable+" deleteProduct")
	}

	commandTag, err := u.db.Exec(ctx, sqlStr, args...)
	if err != nil {
		return u.db.Error(err)
	}

	if commandTag.RowsAffected() == 0 {
		return u.db.Error(fmt.Errorf("no sql rows"))
	}

	return nil
}
